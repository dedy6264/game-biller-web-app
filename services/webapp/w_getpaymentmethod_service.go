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
		req struct {
			SegmentID int64 `json:"segment_id"`
		}
	)
	_ = c.Bind(&req)

	db := connections.DBconn()
	var segmentID int64 = req.SegmentID

	// Default to Open_Biller segment (ID 1) if no segment specified
	if segmentID == 0 {
		segmentID = 1
	}

	list, err := repositories.GetPaymentMethodsWithChannelsBySegmentID(db, segmentID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment methods")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, list))
}
