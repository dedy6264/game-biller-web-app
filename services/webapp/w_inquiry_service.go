package webapp

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// IAK provider ID
const ProviderIAK int64 = 1

// IAK worker base URL
const IakWorkerHost = "http://localhost:10003"

// callIAKInquiry sends inquiry request to IAK worker
func callIAKInquiry(payload models.RequestInquiry) (*models.InquiryResult, error) {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(IakWorkerHost+"/api/iak/inquiry", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	var result models.InquiryResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// === 7. INQUIRY ===
func Inquiry(c echo.Context) error {
	var (
		svc = "Inquiry"
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
	if strings.TrimSpace(req.TargetUserID) == "" {
		helpers.ProcessLogger(c, svc, "target_user_id is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// 3. Get segment ID melalui repo GetMerchantByID
	merchant, err := repositories.GetMerchantByID(db, claims.MerchantID)
	if err != nil || merchant.Status != "active" {
		helpers.ProcessLogger(c, svc, "Failed to get merchant or merchant inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt201, nil))
	}

	var segmentID int64
	if merchant.SegmentID != nil && *merchant.SegmentID != 0 {
		segmentID = *merchant.SegmentID
	} else {
		segmentID = 1
	}

	// 4. Get product segment JOIN product provider by segment_id dan product_code
	psDetail, err := repositories.GetProductSegmentJoinProvider(db, segmentID, req.ProductCode)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product segment not found for this merchant/product")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductSegmnt, nil))
	}

	// Validasi ketersediaan produk & provider
	if !psDetail.ProductIsActive {
		helpers.ProcessLogger(c, svc, "Product is inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductNotFound, nil))
	}

	if !psDetail.ProviderIsAvailable {
		helpers.ProcessLogger(c, svc, "Product provider is not available", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrPvd2301, nil))
	}

	// 5. Fetch Payment Channel
	channel, err := repositories.GetPaymentChannelByCodeSafe(db, req.PaymentChannelCode)
	if err != nil || !channel.IsActive {
		helpers.ProcessLogger(c, svc, "Payment channel not found or inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	// Hitung payment fee
	var paymentFee float64
	if channel.FeeType == "FIXED" {
		paymentFee = channel.FeeValue
	} else if channel.FeeType == "PERCENTAGE" {
		paymentFee = (channel.FeeValue / 100.0) * (psDetail.ProductPrice + psDetail.AdminFee)
	}

	totalAmount := psDetail.ProductPrice + psDetail.AdminFee + paymentFee

	// Resolve harga buy price dari provider
	providerPrice := psDetail.ProviderPrice
	providerMerchantFee := psDetail.ProviderMerchantFee
	if providerPrice == 0 {
		providerPrice = psDetail.ProviderProductPrice
		providerMerchantFee = psDetail.ProviderProductMerchantFee
	}
	if providerPrice == 0 {
		providerPrice = psDetail.ProductPrice * 0.95
	}
	buyPrice := providerPrice + providerMerchantFee

	// 6. Generate reference number internal & construct Transaction
	nowStr := time.Now().Format("20060102150405")
	refInternal := fmt.Sprintf("TRX-%s-%s", nowStr, helpers.RandomDigits(6))
	now := time.Now().Format(time.RFC3339)

	var refMerchant *string
	if req.ReferenceNumberMerchant != "" {
		refMerchant = &req.ReferenceNumberMerchant
	}

	trx := models.Transaction{
		MerchantID:              claims.MerchantID,
		ProductID:               &psDetail.ProductID,
		ProductSegmentID:        &psDetail.ProductSegmentID,
		ProductProviderID:       psDetail.ProductProviderID,
		PaymentChannelID:        &channel.ID,
		SnapshotProductCode:     psDetail.ProductCode,
		SnapshotProductName:     psDetail.ProductName,
		BuyPrice:                buyPrice,
		SellPrice:               psDetail.ProductPrice,
		AdminFee:                psDetail.AdminFee,
		PaymentFee:              paymentFee,
		TotalAmount:             totalAmount,
		TargetUserID:            req.TargetUserID,
		ReferenceNumberInternal: refInternal,
		ReferenceNumberMerchant: refMerchant,
		StatusCode:              helpers.CodeIntrPending,
		StatusMessage:           "PENDING_UPSTREAM",
		RetryCount:              0,
		CreatedAt:               now,
		CreatedBy:               strconv.FormatInt(claims.UserID, 10),
		UpdatedAt:               now,
		UpdatedBy:               strconv.FormatInt(claims.UserID, 10),
	}

	// 7. Simpan data transaksi (status: PENDING) menggunakan DB transaction
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		_, err := repositories.CreateTransaction(tx, &trx)
		return err
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create transaction")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 8. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK inquiry
	if psDetail.ProviderID == ProviderIAK {
		var refProductID int64
		if psDetail.ProductReferenceID != nil {
			refProductID = *psDetail.ProductReferenceID
		}
		iakReq := models.RequestInquiry{
			RefID:      refInternal,
			CustomerID: req.TargetUserID,
			DataProduct: models.DataProduct{
				ProviderID:         ProviderIAK,
				ProductCategoryID:  psDetail.ProductCategoryID,
				ProductTypeID:      psDetail.ProductTypeID,
				ProductCode:        psDetail.ProductCode,
				ProductReferenceID: refProductID,
			},
		}

		iakResult, err := callIAKInquiry(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker call failed, transaction left as PENDING")
		} else {
			trx.UpdatedAt = time.Now().Format(time.RFC3339)
			trx.UpdatedBy = strconv.FormatInt(claims.UserID, 10)

			if iakResult.StatusCode == helpers.CodeInqSuccess {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = "INQUIRY_SUCCESS"
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = &iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = &iakResult.ProviderRefID
				}
			} else {
				trx.StatusCode = helpers.CodeInqSuccess
				trx.StatusMessage = iakResult.ProviderDetail.Message
			}

			_ = helpers.DBTransaction(db, func(tx *sql.Tx) error {
				return repositories.UpdateTransaction(tx, &trx)
			})
		}
	} else {
		// Provider bukan IAK — langsung set inquiry success
		trx.StatusCode = helpers.CodeInqSuccess
		trx.StatusMessage = "INQUIRY_SUCCESS"
		trx.UpdatedAt = time.Now().Format(time.RFC3339)
		trx.UpdatedBy = strconv.FormatInt(claims.UserID, 10)
		_ = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			return repositories.UpdateTransaction(tx, &trx)
		})
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": refInternal,
		"product_code":              psDetail.ProductCode,
		"product_name":              psDetail.ProductName,
		"sell_price":                trx.SellPrice,
		"admin_fee":                 trx.AdminFee,
		"payment_fee":               trx.PaymentFee,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.TargetUserID,
		"status":                    trx.StatusMessage,
	}))
}
