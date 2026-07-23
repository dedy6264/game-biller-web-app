package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

// === 9. TRANSACTION HISTORY ===
func TransactionHistory(c echo.Context) error {
	var (
		svc = "Transaction"
	)
	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth419, nil))
	}

	var req models.RequestTransactions
	if err := c.Bind(&req); err != nil {
		// Allow defaults
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")

	}

	if req.Length <= 0 {
		req.Length = 10
	}

	db := connections.DBconn()
	list, total, err := repositories.GetTransactionsListByMerchantID(db, claims.MerchantID, req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get transaction list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}
