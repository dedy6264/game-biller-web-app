package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

// resolveSegmentName resolves the merchant segment ID from JWT claims.
// Returns the segment ID and true if authenticated, or 0 and false if not.
func resolveSegmentName(c echo.Context, db interface{}) (int64, bool) {
	claims, ok := helpers.GetClaims(c)
	if !ok || claims == nil || claims.MerchantID == 0 {
		return 0, false
	}
	merchant, err := repositories.GetMerchantByID(connections.DBconn(), claims.MerchantID)
	if err != nil || merchant == nil || merchant.Status != "active" {
		return 0, false
	}

	return merchant.SegmentID, true
}

func GetProductSegment(c echo.Context) error {
	var (
		svc     = "GetProductSegment"
		refCode = c.QueryParam("reference_code")
	)
	if refCode == "" {
		var body struct {
			ReferenceCode string `json:"reference_code"`
		}
		if err := c.Bind(&body); err == nil {
			refCode = body.ReferenceCode
		}
	}

	db := connections.DBconn()

	// Coba resolve segment dari JWT (opsional)
	segmentID, authenticated := resolveSegmentName(c, nil)

	var (
		list interface{}
		err  error
	)
	if authenticated {
		// User sudah login → tampilkan produk sesuai segment merchantnya saja
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, refCode, segmentID)
	} else {
		// Guest / tidak login → tampilkan default segment (Public_Retail)
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, refCode, 1)
	}

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get product segments")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, list))
}
