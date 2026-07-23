package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

// === 5. GET POPULAR PRODUCT ===
func GetPopularProduct(c echo.Context) error {
	var (
		svc = "GetPopularProduct"
	)
	db := connections.DBconn()
	list, err := repositories.GetPopularProducts(db)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get popular products")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, list))
}
