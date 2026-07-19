package webapp

import (
	"database/sql"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// === 8. PAYMENT ===
func Payment(c echo.Context) error {
	var (
		svc = "Payment"
	)

	// 1. Validasi kredensial pengguna (JWT)
	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-419", nil))
	}

	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	if req.ReferenceNumberInternal == "" {
		helpers.ProcessLogger(c, svc, "reference_number_internal is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	db := connections.DBconn()

	// 2. Validasi transaksi berdasarkan noreff
	trx, err := repositories.GetTransactionByRefInternal(db, req.ReferenceNumberInternal)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Transaction not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-101", nil))
	}

	// Pastikan transaksi milik merchant yang memanggil API ini
	if trx.MerchantID != claims.MerchantID {
		helpers.ProcessLogger(c, svc, "Unauthorized access to transaction", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-403", nil))
	}

	// 3. Validasi status transaksi — harus dalam status inquiry success
	if trx.StatusCode != StatusInquirySuccess {
		// Status selain inquiry success dinyatakan invalid transaction
		helpers.ProcessLogger(c, svc, "Invalid transaction status: "+trx.StatusCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-102", nil))
	}

	// 4. Fetch dan validasi payment channel
	if trx.PaymentChannelID == nil {
		helpers.ProcessLogger(c, svc, "Payment channel ID is nil", "Internal error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	channel, err := repositories.GetPaymentChannelByID(db, *trx.PaymentChannelID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}

	// Validasi channel masih aktif
	if !channel.IsActive {
		helpers.ProcessLogger(c, svc, "Payment channel is inactive: "+channel.ChannelCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-204", nil))
	}

	// Validasi payment method induk channel masih aktif
	method, err := repositories.GetPaymentMethodByID(db, channel.PaymentMethodID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment method")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	if !method.IsActive {
		helpers.ProcessLogger(c, svc, "Payment method is inactive: "+method.MethodCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-204", nil))
	}

	now := time.Now().Format(time.RFC3339)

	// 5. Jika metode pembayaran BALANCE_INTERNAL — potong saldo & verifikasi PIN
	if channel.ChannelCode == "BALANCE_INTERNAL" {
		sa, err := repositories.GetSavingAccountByMerchantID(db, claims.MerchantID)
		if err != nil || sa.Status != "active" {
			helpers.ProcessLogger(c, svc, "Failed to get saving account or inactive", "Validation error")
			return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-103", nil))
		}

		// Verifikasi PIN
		if sa.AccountPinHash != "" {
			if req.PIN == "" || !helpers.CheckPasswordHash(req.PIN, sa.AccountPinHash) {
				helpers.ProcessLogger(c, svc, "Invalid PIN", "Validation error")
				return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-MERCH-403", nil))
			}
		}

		// Verifikasi saldo mencukupi
		if sa.Balance < trx.TotalAmount {
			helpers.ProcessLogger(c, svc, "Insufficient balance", "Validation error")
			return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-INT-103", nil))
		}

		// Potong saldo dan buat saving transaction dalam satu DB transaction
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
				Description:     "Game Top Up: " + trx.SnapshotProductName,
				CreatedAt:       now,
				CreatedBy:       strconv.FormatInt(claims.UserID, 10),
			}
			_, err = repositories.CreateSavingTransaction(tx, &st)
			return err
		})

		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to process balance deduction")
			return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
		}
	}

	// 6. Update status transaksi menjadi payment success (dalam DB transaction)
	providerRef := "PVD-" + helpers.RandomDigits(12)
	serialNum := "SN-" + helpers.RandomDigits(16)

	trx.StatusCode = StatusPaymentSuccess
	trx.StatusMessage = "PAYMENT_SUCCESS"
	trx.ReferenceNumberProvider = &providerRef
	trx.SerialNumber = &serialNum
	trx.UpdatedAt = now
	trx.UpdatedBy = strconv.FormatInt(claims.UserID, 10)

	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"reference_number_internal": trx.ReferenceNumberInternal,
		"product_code":              trx.SnapshotProductCode,
		"product_name":              trx.SnapshotProductName,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.TargetUserID,
		"serial_number":             serialNum,
		"status":                    "PAYMENT_SUCCESS",
	}))
}
