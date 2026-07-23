package dashboard

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func AdminUpdateTransaction(c echo.Context) error {
	var (
		svc   = "AdminUpdateTransaction"
		input models.Transaction
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	t, err := repositories.GetTransactionByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Transaction not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	t.StatusCode = input.StatusCode
	t.StatusMessage = input.StatusMessage
	t.SerialNumber = input.SerialNumber
	t.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	err = repositories.UpdateTransaction(db, t)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update transaction")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucUser200, t))
}
