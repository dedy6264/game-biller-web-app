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

func AdminGetSavingAccounts(c echo.Context) error {
	var (
		svc = "AdminGetSavingAccounts"
		req models.RequestSavingAccounts
	)
	_ = c.Bind(&req)
	db := connections.DBconn()
	list, total, err := repositories.GetSavingAccountsList(db, req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve saving accounts list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminGetSavingAccountByID(c echo.Context) error {
	var (
		svc = "AdminGetSavingAccountByID"
		req struct {
			ID int64 `json:"id"`
		}
	)
	if err := c.Bind(&req); err != nil || req.ID == 0 {
		helpers.ProcessLogger(c, svc, "Invalid request payload", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	db := connections.DBconn()
	sa, err := repositories.GetSavingAccountByID(db, req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Saving account not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, sa))
}

func AdminUpdateSavingAccount(c echo.Context) error {
	var (
		svc   = "AdminUpdateSavingAccount"
		input models.SavingAccount
	)
	if err := c.Bind(&input); err != nil || input.ID == 0 {
		helpers.ProcessLogger(c, svc, "Invalid request payload or zero ID", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	db := connections.DBconn()
	sa, err := repositories.GetSavingAccountByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Saving account not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	if input.Status != "" {
		sa.Status = input.Status
	}
	if input.Balance != 0 {
		sa.Balance = input.Balance
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	sa.UpdatedAt = now
	sa.UpdatedBy = "admin"

	err = repositories.UpdateSavingAccount(db, sa)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update saving account")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, sa))
}

func AdminGetSavingTransactions(c echo.Context) error {
	var (
		svc = "AdminGetSavingTransactions"
		req models.RequestSavingTransactions
	)
	_ = c.Bind(&req)
	db := connections.DBconn()
	list, total, err := repositories.GetSavingTransactionsList(db, req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve saving transactions list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}
