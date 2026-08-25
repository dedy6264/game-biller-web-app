package webapp

import (
	"database/sql"
	"encoding/json"
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

// callIAKPayment sends payment request to IAK worker using utils.WorkerRequestPOST
func callIAKPayment(payload models.RequestPayment) (*models.PaymentResult, error) {
	data, _, err := utils.WorkerRequestPOST("json", configs.HOST_WORKER+configs.ENDPOINT_IAK_PAYMENT, payload, models.ReqHeader{}, 30*time.Second)
	if err != nil {
		return nil, err
	}
	var result models.PaymentResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// === 8. PAYMENT ===
func Payment(c echo.Context) error {
	var (
		svc = "Payment"
	)

	// 1. Validasi kredensial pengguna (JWT)
	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusUnauthorized, helpers.BuildResponse(helpers.CodeErrAuth419, nil))
	}

	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidRequest, nil))
	}

	if req.ReferenceNumberInternal == "" {
		helpers.ProcessLogger(c, svc, "reference_number_internal is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidRequest, nil))
	}

	db := connections.DBconn()

	// 2. Validasi transaksi berdasarkan noreff
	trx, err := repositories.GetTransactionByRefInternal(db, req.ReferenceNumberInternal)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Transaction not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeServiceDisruption, nil))
	}

	// Pastikan transaksi milik merchant yang memanggil API ini
	if trx.MerchantID != claims.MerchantID {
		helpers.ProcessLogger(c, svc, "Unauthorized access to transaction", "Authorization error")
		return c.JSON(http.StatusUnauthorized, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
	}

	// 3. Validasi status transaksi
	if trx.StatusCode != helpers.CodeInqSuccess {
		helpers.ProcessLogger(c, svc, "Invalid transaction status: "+trx.StatusCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransactionNoOrStatus, nil))
	}

	channel, err := repositories.GetPaymentChannelByID(db, trx.PaymentChannelID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))
	}

	// Validasi channel masih aktif->harus update gagal
	if !channel.IsActive {
		trx.StatusCode = helpers.CodeInvalidPayment
		trx.StatusMessage = "PAYMENT_METHOD_UNAVAILABLE"
		err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			return repositories.UpdateTransaction(tx, trx)
		})
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeServiceDisruption, nil))
		}
		helpers.ProcessLogger(c, svc, "Payment channel is inactive: "+channel.ChannelCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))
	}

	now := time.Now().Format(time.RFC3339)

	// 5. Jika metode pembayaran BALANCE_INTERNAL — potong saldo & verifikasi PIN
	if channel.ChannelCode == "BALANCE_INTERNAL" {
		sa, err := repositories.GetSavingAccountByMerchantID(db, claims.MerchantID)
		if err != nil || sa.Status != "active" {
			helpers.ProcessLogger(c, svc, "Failed to get saving account or inactive", "Validation error")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))
		}

		// 	// Verifikasi PIN
		if sa.AccountPinHash != "" {
			if req.PIN == "" || !helpers.CheckPasswordHash(req.PIN, sa.AccountPinHash) {
				helpers.ProcessLogger(c, svc, "Invalid PIN", "Validation error")
				return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPin, nil))
			}
		}

		// 	// Verifikasi saldo mencukupi
		if sa.Balance < trx.TotalAmount {
			helpers.ProcessLogger(c, svc, "Insufficient balance", "Validation error")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeBalanceLimit, nil))
		}

		// 	// Potong saldo dan buat saving transaction dalam satu DB transaction
		err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			newBalance := sa.Balance - trx.TotalAmount
			if err := repositories.UpdateSavingAccountBalance(tx, sa.ID, newBalance, strconv.FormatInt(claims.UserID, 10), now); err != nil {
				return err
			}
			st := models.SavingTransaction{
				SavingAccountID: sa.ID,
				TypeDC:          "D",
				Amount:          trx.TotalAmount,
				LastBalance:     sa.Balance,
				ReferenceNumber: trx.ReferenceNumberInternal,
				TransactionCode: "GAME_TOPUP",
				Description:     "Game Top Up: " + trx.ProductName,
				CreatedAt:       now,
				CreatedBy:       strconv.FormatInt(claims.UserID, 10),
			}
			_, err = repositories.CreateSavingTransaction(tx, &st)
			return err
		})

		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to process balance deduction")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeServiceDisruption, nil))
		}
	}

	// 6. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK payment
	trx.UpdatedAt = now
	trx.UpdatedBy = strconv.FormatInt(claims.UserID, 10)

	if trx.ProviderID == ProviderIAK {
		iakReq := models.RequestPayment{
			RefID:           trx.ReferenceNumberInternal,
			ProviderRefID:   trx.ReferenceNumberProvider,
			BillDesc:        trx.OtherCustomerID,
			CustomerID:      trx.CustomerID,
			OtherCustomerID: trx.OtherCustomerID,
			DataProduct: models.DataProduct{
				ProviderID:         trx.ProviderID,
				ProductTypeID:      trx.ProductTypeID,
				ProductCode:        trx.ProductProviderCode,
				ProductReferenceID: trx.ProductReferenceID,
			},
		}

		iakResult, err := callIAKPayment(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker payment call failed")
			trx.StatusCode = helpers.CodePending
			trx.StatusMessage = "PENDING_UPSTREAM"
		} else {
			trx.StatusCodeDetail = iakResult.ProviderDetail.Code
			trx.StatusMessageDetail = iakResult.ProviderDetail.Message
			switch iakResult.StatusCode {
			case helpers.CodeSuccess:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = "PAYMENT_SUCCESS"
				// Update data hasil worker ke transaksi
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			case helpers.CodePending:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = iakResult.ProviderDetail.Message
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			default:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = iakResult.ProviderDetail.Message
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			}
		}
	} else {
		// Provider bukan IAK — langsung set payment success
		trx.StatusCode = helpers.CodeSuccess
		trx.StatusMessage = "PAYMENT_SUCCESS"
		trx.SerialNumber = "SN-" + helpers.RandomDigits(16)
		trx.ReferenceNumberProvider = "PVD-" + helpers.RandomDigits(12)
	}

	// 7. Update status transaksi di DB
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
		trx.StatusCode = helpers.CodeServiceDisruption
		// return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": trx.ReferenceNumberInternal,
		"product_code":              trx.ProductCode,
		"product_name":              trx.ProductName,
		"total_amount":              trx.TotalAmount,
		"customer_id":               trx.CustomerID,
		"serial_number":             trx.SerialNumber,
		"status":                    trx.StatusMessage,
	}))
}
func PaymentUnSubscribe(c echo.Context) error {
	var (
		svc = "PaymentUnSubscribe"
	)

	// 1. Validasi kredensial pengguna (JWT)

	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidRequest, nil))
	}

	if req.ReferenceNumberInternal == "" {
		helpers.ProcessLogger(c, svc, "reference_number_internal is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidRequest, nil))
	}

	db := connections.DBconn()

	// 2. Validasi transaksi berdasarkan noreff
	trx, err := repositories.GetTransactionByRefInternal(db, req.ReferenceNumberInternal)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Transaction not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransactionNoOrStatus, nil))
	}

	// 3. Validasi status transaksi
	if trx.StatusCode != helpers.CodeInqSuccess {
		helpers.ProcessLogger(c, svc, "Invalid transaction status: "+trx.StatusCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransactionNoOrStatus, nil))
	}

	channel, err := repositories.GetPaymentChannelByID(db, trx.PaymentChannelID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))
	}

	// Validasi channel masih aktif->harus update gagal
	if !channel.IsActive {
		trx.StatusCode = helpers.CodeInvalidPayment
		trx.StatusMessage = "PAYMENT_METHOD_UNAVAILABLE"
		err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
			return repositories.UpdateTransaction(tx, trx)
		})
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeServiceDisruption, nil))
		}
		helpers.ProcessLogger(c, svc, "Payment channel is inactive: "+channel.ChannelCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))
	}

	now := time.Now().Format(time.RFC3339)

	// 5. Jika metode pembayaran BALANCE_INTERNAL — potong saldo & verifikasi PIN
	if channel.ChannelCode == "BALANCE_INTERNAL" {

		helpers.ProcessLogger(c, svc, err.Error(), "Invalid payment method")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidPayment, nil))

	}

	// 6. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK payment
	trx.UpdatedAt = now
	trx.UpdatedBy = "sys"

	if trx.ProviderID == ProviderIAK {
		iakReq := models.RequestPayment{
			RefID:           trx.ReferenceNumberInternal,
			ProviderRefID:   trx.ReferenceNumberProvider,
			BillDesc:        trx.OtherCustomerID,
			CustomerID:      trx.CustomerID,
			OtherCustomerID: trx.OtherCustomerID,
			DataProduct: models.DataProduct{
				ProviderID:         trx.ProviderID,
				ProductTypeID:      trx.ProductTypeID,
				ProductCode:        trx.ProductProviderCode,
				ProductReferenceID: trx.ProductReferenceID,
			},
		}

		iakResult, err := callIAKPayment(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker payment call failed")
			trx.StatusCode = helpers.CodePending
			trx.StatusMessage = "PENDING_UPSTREAM"
		} else {
			trx.StatusCodeDetail = iakResult.ProviderDetail.Code
			trx.StatusMessageDetail = iakResult.ProviderDetail.Message
			switch iakResult.StatusCode {
			case helpers.CodeSuccess:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = "PAYMENT_SUCCESS"
				// Update data hasil worker ke transaksi
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			case helpers.CodePending:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = iakResult.ProviderDetail.Message
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			default:
				trx.StatusCode = iakResult.StatusCode
				trx.StatusMessage = iakResult.ProviderDetail.Message
				if iakResult.DataTransaction.SerialNumber != "" {
					trx.SerialNumber = iakResult.DataTransaction.SerialNumber
				}
				if iakResult.ProviderRefID != "" {
					trx.ReferenceNumberProvider = iakResult.ProviderRefID
				}
			}
		}
	} else {
		// Provider bukan IAK — langsung set payment success
		trx.StatusCode = helpers.CodeSuccess
		trx.StatusMessage = "PAYMENT_SUCCESS"
		trx.SerialNumber = "SN-" + helpers.RandomDigits(16)
		trx.ReferenceNumberProvider = "PVD-" + helpers.RandomDigits(12)
	}

	// 7. Update status transaksi di DB
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
		trx.StatusCode = helpers.CodeServiceDisruption
		// return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": trx.ReferenceNumberInternal,
		"product_code":              trx.ProductCode,
		"product_name":              trx.ProductName,
		"total_amount":              trx.TotalAmount,
		"customer_id":               trx.CustomerID,
		"serial_number":             trx.SerialNumber,
		"status":                    trx.StatusMessage,
	}))
}
