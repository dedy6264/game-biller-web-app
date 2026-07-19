package webapp

import (
	"database/sql"
	"fmt"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Status code constants
const (
	StatusInquirySuccess = "INQ-SYS-001" // inquiry success - awaiting payment
	StatusPaymentSuccess = "SUC-INT-000" // payment success
)

// === 7. INQUIRY ===
func Inquiry(c echo.Context) error {
	var (
		svc = "Inquiry"
	)

	// 1. Validasi kredensial pengguna (JWT)
	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-419", nil))
	}

	var req models.InquiryRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	// 2. Validasi customer_id (target_user_id)
	if strings.TrimSpace(req.TargetUserID) == "" {
		helpers.ProcessLogger(c, svc, "target_user_id is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	db := connections.DBconn()

	// 3. Validasi product — cari product berdasarkan product_code
	product, err := repositories.GetProductByCode(db, req.ProductCode)
	if err != nil || !product.IsActive {
		helpers.ProcessLogger(c, svc, "Failed to get product or product inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-101", nil))
	}

	// 4. Fetch Merchant Details untuk menentukan segment
	merchant, err := repositories.GetMerchantByID(db, claims.MerchantID)
	if err != nil || merchant.Status != "active" {
		helpers.ProcessLogger(c, svc, "Failed to get merchant or merchant inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-201", nil))
	}

	// Tentukan segment berdasarkan tipe merchant
	segmentName := "Public_Retail"
	switch merchant.MerchantType {
	case "guest_retail":
		segmentName = "Public_Retail"
	case "member_premium":
		segmentName = "Gold_Reseller"
	case "h2h_api":
		segmentName = "H2H_Partner"
	}

	// 5. Validasi product segment — cari harga segment sesuai product_code
	segmentPrice, err := repositories.GetProductSegmentByProductAndSegment(db, product.ID, segmentName)
	if err != nil {
		// Product tidak tersedia di segment merchant ini → inquiry gagal
		helpers.ProcessLogger(c, svc, err.Error(), "Product segment not found for this merchant type")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-105", nil))
	}

	// 6. Fetch Payment Channel
	channel, err := repositories.GetPaymentChannelByCodeSafe(db, req.PaymentChannelCode)
	if err != nil || !channel.IsActive {
		helpers.ProcessLogger(c, svc, "Payment channel not found or inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	// Hitung payment fee
	var paymentFee float64
	if channel.FeeType == "FIXED" {
		paymentFee = channel.FeeValue
	} else if channel.FeeType == "PERCENTAGE" {
		paymentFee = (channel.FeeValue / 100.0) * (segmentPrice.ProductPrice + segmentPrice.AdminFee)
	}

	totalAmount := segmentPrice.ProductPrice + segmentPrice.AdminFee + paymentFee

	// 7. Resolve harga provider
	providerPrice := segmentPrice.ProviderProductPrice
	providerMerchantFee := segmentPrice.ProviderProductMerchantFee
	var productProviderID *int64 = segmentPrice.ProductProviderID

	if providerPrice == 0 && segmentPrice.ProductProviderID != nil && *segmentPrice.ProductProviderID != 0 {
		pp, err := repositories.GetProductProviderByID(db, *segmentPrice.ProductProviderID)
		if err == nil {
			providerPrice = pp.ProviderPrice
			providerMerchantFee = pp.ProviderMerchantFee
		}
	}
	if providerPrice == 0 {
		providerPrice = segmentPrice.ProductPrice * 0.95 // default 5% margin
	}

	// 8. Generate reference number internal
	nowStr := time.Now().Format("20060102150405")
	refInternal := fmt.Sprintf("TRX-%s-%s", nowStr, helpers.RandomDigits(6))

	now := time.Now().Format(time.RFC3339)

	var refMerchant *string
	if req.ReferenceNumberMerchant != "" {
		refMerchant = &req.ReferenceNumberMerchant
	}

	trx := models.Transaction{
		MerchantID:              claims.MerchantID,
		ProductID:               &product.ID,
		ProductSegmentID:        &segmentPrice.ID,
		ProductProviderID:       productProviderID,
		PaymentChannelID:        &channel.ID,
		SnapshotProductCode:     product.ProductCode,
		SnapshotProductName:     product.ProductName,
		BuyPrice:                providerPrice + providerMerchantFee,
		SellPrice:               segmentPrice.ProductPrice,
		AdminFee:                segmentPrice.AdminFee,
		PaymentFee:              paymentFee,
		TotalAmount:             totalAmount,
		TargetUserID:            req.TargetUserID,
		ReferenceNumberInternal: refInternal,
		ReferenceNumberMerchant: refMerchant,
		StatusCode:              StatusInquirySuccess,
		StatusMessage:           "INQUIRY_SUCCESS",
		RetryCount:              0,
		CreatedAt:               now,
		CreatedBy:               strconv.FormatInt(claims.UserID, 10),
		UpdatedAt:               now,
		UpdatedBy:               strconv.FormatInt(claims.UserID, 10),
	}

	// 9. Simpan data transaksi menggunakan DB transaction
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		_, err := repositories.CreateTransaction(tx, &trx)
		return err
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create transaction")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse("INQ-SYS-001", map[string]any{
		"reference_number_internal": refInternal,
		"product_code":              product.ProductCode,
		"product_name":              product.ProductName,
		"sell_price":                trx.SellPrice,
		"admin_fee":                 trx.AdminFee,
		"payment_fee":               trx.PaymentFee,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.TargetUserID,
	}))
}
