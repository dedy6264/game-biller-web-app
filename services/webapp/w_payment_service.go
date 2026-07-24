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
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth419, nil))
	}

	var req models.PaymentRequest
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if req.ReferenceNumberInternal == "" {
		helpers.ProcessLogger(c, svc, "reference_number_internal is empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// 2. Validasi transaksi berdasarkan noreff
	trx, err := repositories.GetTransactionByRefInternal(db, req.ReferenceNumberInternal)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Transaction not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidProductNotFound, nil))
	}

	// Pastikan transaksi milik merchant yang memanggil API ini
	if trx.MerchantID != claims.MerchantID {
		helpers.ProcessLogger(c, svc, "Unauthorized access to transaction", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
	}

	// 3. Validasi status transaksi
	if trx.StatusCode != helpers.CodeInqSuccess {
		helpers.ProcessLogger(c, svc, "Invalid transaction status: "+trx.StatusCode, "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransaction, nil))
	}

	// 4. Fetch dan validasi payment channel
	// if trx.PaymentChannelID == nil {
	// 	helpers.ProcessLogger(c, svc, "Payment channel ID is nil", "Internal error")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	// }
	// channel, err := repositories.GetPaymentChannelByID(db, *trx.PaymentChannelID)
	// if err != nil {
	// 	helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment channel")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	// }

	// // Validasi channel masih aktif
	// if !channel.IsActive {
	// 	helpers.ProcessLogger(c, svc, "Payment channel is inactive: "+channel.ChannelCode, "Validation error")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt204, nil))
	// }

	// // Validasi payment method induk channel masih aktif
	// method, err := repositories.GetPaymentMethodByID(db, channel.PaymentMethodID)
	// if err != nil {
	// 	helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment method")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	// }
	// if !method.IsActive {
	// 	helpers.ProcessLogger(c, svc, "Payment method is inactive: "+method.MethodCode, "Validation error")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt204, nil))
	// }

	now := time.Now().Format(time.RFC3339)

	// 5. Jika metode pembayaran BALANCE_INTERNAL — potong saldo & verifikasi PIN
	// if channel.ChannelCode == "BALANCE_INTERNAL" {
	// 	sa, err := repositories.GetSavingAccountByMerchantID(db, claims.MerchantID)
	// 	if err != nil || sa.Status != "active" {
	// 		helpers.ProcessLogger(c, svc, "Failed to get saving account or inactive", "Validation error")
	// 		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrIntBalance, nil))
	// 	}

	// 	// Verifikasi PIN
	// 	if sa.AccountPinHash != "" {
	// 		if req.PIN == "" || !helpers.CheckPasswordHash(req.PIN, sa.AccountPinHash) {
	// 			helpers.ProcessLogger(c, svc, "Invalid PIN", "Validation error")
	// 			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrMerch403, nil))
	// 		}
	// 	}

	// 	// Verifikasi saldo mencukupi
	// 	if sa.Balance < trx.TotalAmount {
	// 		helpers.ProcessLogger(c, svc, "Insufficient balance", "Validation error")
	// 		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrIntBalance, nil))
	// 	}

	// 	// Potong saldo dan buat saving transaction dalam satu DB transaction
	// 	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
	// 		newBalance := sa.Balance - trx.TotalAmount
	// 		if err := repositories.UpdateSavingAccountBalance(tx, sa.ID, newBalance, strconv.FormatInt(claims.UserID, 10), now); err != nil {
	// 			return err
	// 		}
	// 		st := models.SavingTransaction{
	// 			SavingAccountID: sa.ID,
	// 			TypeDC:          "D",
	// 			Amount:          trx.TotalAmount,
	// 			LastBalance:     sa.Balance,
	// 			ReferenceNumber: trx.ReferenceNumberInternal,
	// 			TransactionCode: "GAME_TOPUP",
	// 			Description:     "Game Top Up: " + trx.SnapshotProductName,
	// 			CreatedAt:       now,
	// 			CreatedBy:       strconv.FormatInt(claims.UserID, 10),
	// 		}
	// 		_, err = repositories.CreateSavingTransaction(tx, &st)
	// 		return err
	// 	})

	// 	if err != nil {
	// 		helpers.ProcessLogger(c, svc, err.Error(), "Failed to process balance deduction")
	// 		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	// 	}
	// }

	// 6. Validasi provider — jika provider_id = 1 (IAK), panggil worker IAK payment
	trx.UpdatedAt = now
	trx.UpdatedBy = strconv.FormatInt(claims.UserID, 10)

	if trx.ProviderID != nil && *trx.ProviderID == ProviderIAK {
		// Siapkan bill_desc dari data yang tersimpan di transaksi
		billDesc := trx.SnapshotProductName + " - " + trx.CustomerID

		provRef := ""
		if trx.ReferenceNumberProvider != nil {
			provRef = *trx.ReferenceNumberProvider
		}

		typeID := int64(0)
		if trx.ProductTypeID != nil {
			typeID = *trx.ProductTypeID
		}
		refID := int64(0)
		if trx.ProductReferenceID != nil {
			refID = *trx.ProductReferenceID
		}

		iakReq := models.RequestPayment{
			RefID:           trx.ReferenceNumberInternal,
			ProviderRefID:   provRef,
			BillDesc:        billDesc,
			CustomerID:      trx.CustomerID,
			OtherCustomerID: trx.OtherCustomerID,
			DataProduct: models.DataProduct{
				ProviderID:         *trx.ProviderID,
				ProductTypeID:      typeID,
				ProductCode:        trx.ProductProviderCode,
				ProductReferenceID: refID,
			},
		}

		iakResult, err := callIAKPayment(iakReq)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "IAK worker payment call failed")
			trx.StatusCode = helpers.CodeIntrPending
			trx.StatusMessage = "PENDING_UPSTREAM"
		} else {
			if iakResult.StatusCode == helpers.CodeSuccess {
				trx.StatusCode = helpers.CodeSuccess
				trx.StatusMessage = "PAYMENT_SUCCESS"
				// Update data hasil worker ke transaksi
				if iakResult.DataTransaction.SerialNumber != "" {
					sn := iakResult.DataTransaction.SerialNumber
					trx.SerialNumber = &sn
				}
				if iakResult.ProviderRefID != "" {
					pref := iakResult.ProviderRefID
					trx.ReferenceNumberProvider = &pref
				}
			} else {
				trx.StatusCode = helpers.CodeIntrPending
				trx.StatusMessage = iakResult.ProviderDetail.Message
			}
		}
	} else {
		// Provider bukan IAK — langsung set payment success
		serialNum := "SN-" + helpers.RandomDigits(16)
		providerRef := "PVD-" + helpers.RandomDigits(12)
		trx.StatusCode = helpers.CodeSuccess
		trx.StatusMessage = "PAYMENT_SUCCESS"
		trx.SerialNumber = &serialNum
		trx.ReferenceNumberProvider = &providerRef
	}

	// 7. Update status transaksi di DB
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(trx.StatusCode, map[string]any{
		"reference_number_internal": trx.ReferenceNumberInternal,
		"product_code":              trx.SnapshotProductCode,
		"product_name":              trx.SnapshotProductName,
		"total_amount":              trx.TotalAmount,
		"target_user_id":            trx.CustomerID,
		"serial_number":             trx.SerialNumber,
		"status":                    trx.StatusMessage,
	}))
}
