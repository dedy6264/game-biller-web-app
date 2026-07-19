package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetReferenceProduct(c echo.Context) error {
	var (
		svc = "GetReferenceProduct"
	)
	db := connections.DBconn()
	list, _, err := repositories.GetProductReferencesList(db, "", 0, 0, "", "", models.ProductReferenceFilters{})
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get product references")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", list))
}
