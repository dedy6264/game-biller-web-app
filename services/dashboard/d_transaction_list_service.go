package dashboard

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminGetTransactions(c echo.Context) error {
	var (
		svc = "AdminGetTransactions"
		req models.RequestTransactions
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetTransactionsList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve transactions list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}
