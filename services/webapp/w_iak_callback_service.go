package webapp

import (
	"database/sql"
	"fmt"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// === IAK CALLBACK ===
func IAKCallback(c echo.Context) error {
	var (
		svc     = "IAKCallback"
		payload models.IAKCallbackPayload
	)

	if err := c.Bind(&payload); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind callback request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransaction, nil))
	}

	data := payload.Data
	if strings.TrimSpace(data.RefID) == "" && strings.TrimSpace(data.TrID) == "" {
		helpers.ProcessLogger(c, svc, "ref_id and tr_id are both empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransaction, nil))
	}

	db := connections.DBconn()

	// 2. Get transaction by referenceNumberInternal using ref_id, or referenceNumberProvider using tr_id
	var (
		trx *models.Transaction
		err error
	)
	if data.RefID != "" {
		trx, err = repositories.GetTransactionByRefInternal(db, data.RefID)
	}
	if (err != nil || trx == nil) && data.TrID != "" {
		trx, err = repositories.GetTransactionByRefProvider(db, data.TrID)
	}
	if err != nil || trx == nil {
		helpers.ProcessLogger(c, svc, fmt.Sprintf("Transaction not found for ref_id=%s tr_id=%s", data.RefID, data.TrID), "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransaction, nil))
	}

	// 3. Validasi status data transaksi, jika status bukan payment pending return invalid transaksi
	if trx.StatusCode != helpers.CodeIntrPending {
		helpers.ProcessLogger(c, svc, fmt.Sprintf("Transaction status is not pending: %s", trx.StatusCode), "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidTransaction, nil))
	}

	// 4. Jika status pending, update sesuai status callback, sukses atau gagal berikut dengan kelengkapan data seperti sn dan data lainnya
	mainCode, providerMsg, _ := helpers.ConvertIAKPaymentResponse(data.RC)
	trx.StatusCode = mainCode
	if data.Message != "" {
		trx.StatusMessage = data.Message
	} else {
		trx.StatusMessage = providerMsg
	}

	if data.TrID != "" {
		trx.ReferenceNumberProvider = data.TrID
	}
	if data.SN != "" {
		trx.SerialNumber = data.SN
	}
	now := time.Now().Format(time.RFC3339)
	trx.UpdatedAt = now
	trx.UpdatedBy = "CALLBACK_IAK"

	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		return repositories.UpdateTransaction(tx, trx)
	})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction status")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"reference_number_internal": trx.ReferenceNumberInternal,
		"reference_number_provider": trx.ReferenceNumberProvider,
		"status_code":               trx.StatusCode,
		"status_message":            trx.StatusMessage,
		"serial_number":             trx.SerialNumber,
	}))
}
