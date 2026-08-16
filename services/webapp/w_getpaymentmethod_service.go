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

	claims, ok := helpers.GetClaims(c)
	if !ok {
		helpers.ProcessLogger(c, svc, "Failed to get claims", "Authorization error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth419, nil))
	}
	// Resolve segment from JWT claims if user/merchant is logged in
	merchant, err := repositories.GetMerchantByID(db, claims.MerchantID)
	if err != nil || merchant.Status != "active" {
		helpers.ProcessLogger(c, svc, "Failed to get merchant or merchant inactive", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrInt201, nil))
	}
	if merchant.SegmentID != 0 {
		segmentID = merchant.SegmentID
	} else {
		segmentID = 1
	}

	list, err := repositories.GetPaymentMethodsWithChannelsBySegmentID(db, segmentID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get payment methods")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, list))
}
