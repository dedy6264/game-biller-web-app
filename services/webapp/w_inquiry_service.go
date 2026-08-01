package webapp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gamebiller/configs"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"gamebiller/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// IAK provider ID
const ProviderIAK int64 = 1

// callIAKInquiry sends inquiry request to IAK worker using utils.WorkerRequestPOST
func callIAKInquiry(payload models.RequestInquiry) (*models.InquiryResult, error) {
	data, _, err := utils.WorkerRequestPOST("json", configs.HOST_WORKER+configs.ENDPOINT_IAK_INQUIRY, payload, models.ReqHeader{}, 30*time.Second)
	if err != nil {
		return nil, err
	}
	var result models.InquiryResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// === 7. INQUIRY ===
func Inquiry(c echo.Context) error {
	var (
		svc                     = "Inquiry"
		bytes                   []byte
		db                      = connections.DBconn()
		totalAmount, paymentFee float64 //, providerPrice, providerMerchantFee, providerAdminFee float64
		segmentID               int64
	)

	// 1. Validasi kredensial pengguna (JWT)
	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth419, nil))
	}

	var req models.InquiryRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	// 2. Validasi customer_id (target_user_id)
	// if strings.TrimSpace(req.TargetUserID) == "" {
	// 	helpers.ProcessLogger(c, svc, "target_user_id is empty", "Validation error")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	// }

	// 3. Get segment ID melalui repo GetMerchantByID
	merchant, err := repositories.GetMerchantByID(db, claims.MerchantID)
	if err != nil || merchant.Status != "active" {
		helpers.ProcessLogger(c, svc, "Failed to get merchant or merchant inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt201, nil))
	}
	if merchant.SegmentID != 0 {
		segmentID = merchant.SegmentID
	} else {
		segmentID = 1
	}

	// 4. Get product segment JOIN product provider by segment_id dan product_code
	productSegment, err := repositories.GetProductSegmentJoinProvider(db, segmentID, req.ProductCode)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product segment not found for this merchant/product")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductSegmnt, nil))
	}
	// Validasi ketersediaan produk & provider
	if !productSegment.ProductIsActive {
		helpers.ProcessLogger(c, svc, "Product is inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductNotFound, nil))
	}

	if !productSegment.ProviderIsAvailable {
		helpers.ProcessLogger(c, svc, "Product provider is not available", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrPvd2301, nil))
	}
	// 5. Fetch Payment Channel
	channel, err := repositories.GetPaymentChannelByCode(db, req.PaymentChannelCode)
	if err != nil || !channel.IsActive {
		helpers.ProcessLogger(c, svc, "Payment channel invalid or inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	paymentFee = channel.FeeValue

	// Resolve harga buy price dari provider
	// providerPrice = productSegment.ProductProviderPrice
	// providerMerchantFee = productSegment.ProductProviderMerchantFee
	// providerAdminFee = productSegment.ProductProviderAdminFee

	// 6. Generate reference number internal & construct Transaction
	nowStr := time.Now().Format("20060102150405")
	refInternal := fmt.Sprintf("TRX-%s-%s", nowStr, helpers.RandomDigits(6))
	now := time.Now().Format(time.RFC3339)

	if req.ZoneID != "" || req.ServerID != "" {
		otherCustId := models.OtherCustomerID{
			ZoneID:   req.ZoneID,
			ServerID: req.ServerID,
		}
		bytes, _ = json.Marshal(otherCustId)

	}
	trx := models.Transaction{
		ReferenceNumberInternal: refInternal,
		ReferenceNumberMerchant: req.ReferenceNumberMerchant,

		MerchantID:         claims.MerchantID,
		MerchantName:       merchant.MerchantName,
		ProductSegmentID:   productSegment.ProductSegmentID,
		ProviderID:         productSegment.ProviderID,
		ProviderName:       productSegment.ProviderName,
		ProductTypeID:      productSegment.ProductTypeID,
		ProductTypeName:    productSegment.ProductTypeName,
		ProductReferenceID: productSegment.ProductReferenceID,

		ProductProviderID:          productSegment.ProductProviderID,
		ProductProviderName:        productSegment.ProductProviderName,
		ProductProviderCode:        productSegment.ProductProviderCode,
		ProductProviderPrice:       productSegment.ProductProviderPrice,
		ProductProviderAdminFee:    productSegment.ProductProviderAdminFee,
		ProductProviderMerchantFee: productSegment.ProductProviderMerchantFee,

		ProductID:          productSegment.ProductID,
		ProductName:        productSegment.ProductName,
		ProductCode:        productSegment.ProductCode,
		ProductPrice:       productSegment.ProductPrice,
		ProductAdminFee:    productSegment.AdminFee,
		ProductMerchantFee: productSegment.MerchantFee,

		PaymentChannelID:   channel.ID,
		PaymentChannelName: channel.ChannelName,
		PaymentAdminFee:    paymentFee,

		ProductSegmentName: productSegment.SegmentName,

		TotalAmount: totalAmount,

		CustomerID:      req.TargetUserID,
		OtherCustomerID: string(bytes),
		CustomerPhone:   req.CustomerPhone,

		StatusCode:    "", //belum memiliki status atau inquiry proccess
		StatusMessage: "INQUIRY_PROCCESS",
		RetryCount:    0,
		CreatedAt:     now,
		CreatedBy:     "sys",
		UpdatedAt:     now,
		UpdatedBy:     "sys",
	}
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		_, err := repositories.CreateTransaction(tx, &trx)
		return err
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create transaction")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 8. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK inquiry
	if productSegment.ProviderID == ProviderIAK {
		var refProductID int64
		refProductID = productSegment.ProductReferenceID
		iakReq := models.RequestInquiry{
			RefID:      refInternal,
			CustomerID: req.TargetUserID,
			DataProduct: models.DataProduct{
				ProviderID:         ProviderIAK,
				ProductCategoryID:  productSegment.ProductCategoryID,
				ProductTypeID:      productSegment.ProductTypeID,
				ProductCode:        productSegment.ProductProviderCode,
				ProductReferenceID: refProductID,
			},
		}
		////bru sampai sini
		iakResult, err := callIAKInquiry(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker call failed, transaction left as PENDING")
		} else {
			trx.UpdatedAt = time.Now().Format(time.RFC3339)
			trx.UpdatedBy = "sys"

			if iakResult.StatusCode == helpers.CodeInqSuccess {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = "INQUIRY_SUCCESS"
				if iakResult.DataTransaction.SerialNumber != "" {
					sn := iakResult.DataTransaction.SerialNumber
					trx.SerialNumber = sn
				}
				if iakResult.ProviderRefID != "" {
					pref := iakResult.ProviderRefID
					trx.ReferenceNumberProvider = pref
				}
			} else {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = iakResult.ProviderDetail.Message
			}
			//sementara pembeda baku berdasar prepaid/postpaid
			if productSegment.ProductTypeID == 1 { //prepaid
				totalAmount = productSegment.ProductPrice + iakResult.DataTransaction.AdminFee + paymentFee
			} else {
				totalAmount = iakResult.DataTransaction.Price + iakResult.DataTransaction.AdminFee + paymentFee
			}
			if iakResult.DataTransaction.Price != 0 {
				trx.ProductProviderPrice = iakResult.DataTransaction.Price
			}
			if iakResult.DataTransaction.AdminFee != 0 {
				trx.ProductProviderAdminFee = iakResult.DataTransaction.AdminFee
			}
			if iakResult.DataTransaction.MerchantFee != 0 {
				trx.ProductProviderMerchantFee = iakResult.DataTransaction.MerchantFee
			}
			if totalAmount != 0 {
				trx.TotalAmount = totalAmount
			}
			trx.CustomerID = iakResult.DataTransaction.CustomerID
			// 		OtherCustomerID: string(bytes),
			// _ = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			// 	return repositories.UpdateTransaction(tx, &trx)
			// })
		}
	} else {
		// Provider bukan IAK — langsung set inquiry success
		trx.StatusCode = helpers.CodeInqSuccess
		trx.StatusMessage = "INQUIRY_SUCCESS"
		trx.UpdatedAt = time.Now().Format(time.RFC3339)
		trx.UpdatedBy = "sys"

	}
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, &trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction")
		trx.StatusCode = helpers.CodeErrSys500
		// return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": refInternal,
		"product_code":              productSegment.ProductCode,
		"product_name":              productSegment.ProductName,
		"sell_price":                trx.ProductPrice,
		"admin_fee":                 trx.ProductAdminFee,
		"payment_fee":               trx.PaymentAdminFee,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.CustomerID,
		"status":                    trx.StatusMessage,
	}))
}

func InquiryUnSubscribe(c echo.Context) error {
	var (
		svc                     = "InquiryUnSubscribe"
		bytes                   []byte
		db                      = connections.DBconn()
		totalAmount, paymentFee float64 //, providerPrice, providerMerchantFee, providerAdminFee float64
		segmentID               int64
	)

	var req models.InquiryRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	// 2. Validasi customer_id (target_user_id)
	if req.CustomerPhone == "" {
		helpers.ProcessLogger(c, svc, "Customer Phone is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	mid, _ := strconv.Atoi(configs.DEFAULT_MID)
	// 3. Get segment ID melalui repo GetMerchantByID
	merchant, err := repositories.GetMerchantByID(db, int64(mid))
	if err != nil || merchant.Status != "active" {
		helpers.ProcessLogger(c, svc, "Failed to get merchant or merchant inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt201, nil))
	}
	if merchant.SegmentID != 0 {
		segmentID = merchant.SegmentID
	} else {
		segmentID = 1
	}

	// 4. Get product segment JOIN product provider by segment_id dan product_code
	productSegment, err := repositories.GetProductSegmentJoinProvider(db, segmentID, req.ProductCode)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product segment not found for this merchant/product")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductSegmnt, nil))
	}

	// Validasi ketersediaan produk & provider
	if !productSegment.ProductIsActive {
		helpers.ProcessLogger(c, svc, "Product is inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductNotFound, nil))
	}

	if !productSegment.ProviderIsAvailable {
		helpers.ProcessLogger(c, svc, "Product provider is not available", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrPvd2301, nil))
	}
	// 5. Fetch Payment Channel
	channel, err := repositories.GetPaymentChannelByCode(db, "QRIS_GATEWAY") //default QRIS_GATEWAY
	if err != nil || !channel.IsActive {
		helpers.ProcessLogger(c, svc, "Payment channel invalid or inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	paymentFee = channel.FeeValue

	// Resolve harga buy price dari provider
	// providerPrice = productSegment.ProductProviderPrice
	// providerMerchantFee = productSegment.ProductProviderMerchantFee
	// providerAdminFee = productSegment.ProductProviderAdminFee

	// 6. Generate reference number internal & construct Transaction
	nowStr := time.Now().Format("20060102150405")
	refInternal := fmt.Sprintf("TRX-%s-%s", nowStr, helpers.RandomDigits(6))
	now := time.Now().Format(time.RFC3339)

	if req.ZoneID != "" || req.ServerID != "" {
		otherCustId := models.OtherCustomerID{
			ZoneID:   req.ZoneID,
			ServerID: req.ServerID,
		}
		bytes, _ = json.Marshal(otherCustId)

	}
	trx := models.Transaction{
		ReferenceNumberInternal: refInternal,
		ReferenceNumberMerchant: req.ReferenceNumberMerchant,

		MerchantID:         int64(mid),
		MerchantName:       merchant.MerchantName,
		ProductSegmentID:   productSegment.ProductSegmentID,
		ProviderID:         productSegment.ProviderID,
		ProviderName:       productSegment.ProviderName,
		ProductTypeID:      productSegment.ProductTypeID,
		ProductTypeName:    productSegment.ProductTypeName,
		ProductReferenceID: productSegment.ProductReferenceID,

		ProductProviderID:          productSegment.ProductProviderID,
		ProductProviderName:        productSegment.ProductProviderName,
		ProductProviderCode:        productSegment.ProductProviderCode,
		ProductProviderPrice:       productSegment.ProductProviderPrice,
		ProductProviderAdminFee:    productSegment.ProductProviderAdminFee,
		ProductProviderMerchantFee: productSegment.ProductProviderMerchantFee,

		ProductID:          productSegment.ProductID,
		ProductName:        productSegment.ProductName,
		ProductCode:        productSegment.ProductCode,
		ProductPrice:       productSegment.ProductPrice,
		ProductAdminFee:    productSegment.AdminFee,
		ProductMerchantFee: productSegment.MerchantFee,

		PaymentChannelID:   channel.ID,
		PaymentChannelName: channel.ChannelName,
		PaymentAdminFee:    paymentFee,

		ProductSegmentName: productSegment.SegmentName,

		TotalAmount: totalAmount,

		CustomerID:      req.TargetUserID,
		OtherCustomerID: string(bytes),

		StatusCode:    "", //belum memiliki status atau inquiry proccess
		StatusMessage: "INQUIRY_PROCCESS",
		RetryCount:    0,
		CreatedAt:     now,
		CreatedBy:     "sys",
		UpdatedAt:     now,
		UpdatedBy:     "sys",
	}

	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		_, err := repositories.CreateTransaction(tx, &trx)
		return err
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create transaction")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 8. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK inquiry
	if productSegment.ProviderID == ProviderIAK {
		var refProductID int64
		refProductID = productSegment.ProductReferenceID
		iakReq := models.RequestInquiry{
			RefID:      refInternal,
			CustomerID: req.TargetUserID,
			DataProduct: models.DataProduct{
				ProviderID:         ProviderIAK,
				ProductCategoryID:  productSegment.ProductCategoryID,
				ProductTypeID:      productSegment.ProductTypeID,
				ProductCode:        productSegment.ProductProviderCode,
				ProductReferenceID: refProductID,
			},
		}
		////bru sampai sini
		iakResult, err := callIAKInquiry(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker call failed, transaction left as PENDING")
		} else {
			trx.UpdatedAt = time.Now().Format(time.RFC3339)
			trx.UpdatedBy = "sys"

			if iakResult.StatusCode == helpers.CodeInqSuccess {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = "INQUIRY_SUCCESS"
				if iakResult.DataTransaction.SerialNumber != "" {
					sn := iakResult.DataTransaction.SerialNumber
					trx.SerialNumber = sn
				}
				if iakResult.ProviderRefID != "" {
					pref := iakResult.ProviderRefID
					trx.ReferenceNumberProvider = pref
				}
			} else {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = iakResult.ProviderDetail.Message
			}
			//sementara pembeda baku berdasar prepaid/postpaid
			if productSegment.ProductTypeID == 1 { //prepaid
				totalAmount = productSegment.ProductPrice + iakResult.DataTransaction.AdminFee + paymentFee
			} else {
				totalAmount = iakResult.DataTransaction.Price + iakResult.DataTransaction.AdminFee + paymentFee
			}
			trx.ProductProviderPrice = iakResult.DataTransaction.Price
			// 			ProductProviderPrice:       productSegment.ProductProviderPrice,
			trx.ProductProviderAdminFee = iakResult.DataTransaction.AdminFee
			// 		ProductProviderAdminFee:    productSegment.ProductProviderAdminFee,
			trx.ProductProviderMerchantFee = iakResult.DataTransaction.MerchantFee
			// 		ProductProviderMerchantFee: productSegment.ProductProviderMerchantFee,
			trx.TotalAmount = totalAmount
			// TotalAmount: totalAmount,

			trx.CustomerID = iakResult.DataTransaction.CustomerID
			// 		OtherCustomerID: string(bytes),
			// _ = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			// 	return repositories.UpdateTransaction(tx, &trx)
			// })
		}
	} else {
		// Provider bukan IAK — langsung set inquiry success
		trx.StatusCode = helpers.CodeInqSuccess
		trx.StatusMessage = "INQUIRY_SUCCESS"
		trx.UpdatedAt = time.Now().Format(time.RFC3339)
		trx.UpdatedBy = "sys"

	}
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, &trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction")
		trx.StatusCode = helpers.CodeErrSys500
		// return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": refInternal,
		"product_code":              productSegment.ProductCode,
		"product_name":              productSegment.ProductName,
		"sell_price":                trx.ProductPrice,
		"admin_fee":                 trx.ProductAdminFee,
		"payment_fee":               trx.PaymentAdminFee,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.CustomerID,
		"status":                    trx.StatusMessage,
	}))
}
