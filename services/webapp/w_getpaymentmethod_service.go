package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

// === 6. GET PAYMENT METHOD ===
func GetPaymentMethod(c echo.Context) error {
	var (
		svc = "GetPaymentMethod"
	)
	db := connections.DBconn()
	list, err := repositories.GetPaymentMethodsWithChannels(db)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment methods")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", list))
}
