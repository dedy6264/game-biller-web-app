package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

// resolveSegmentName resolves the merchant segment name from JWT claims.
// Returns the segment name and true if authenticated, or empty string and false if not.
func resolveSegmentName(c echo.Context, db interface { /* QueryExecutor */
}) (string, bool) {
	claims, ok := helpers.GetClaims(c)
	if !ok || claims.MerchantID == 0 {
		return "", false
	}
	merchant, err := repositories.GetMerchantByID(connections.DBconn(), claims.MerchantID)
	if err != nil || merchant.Status != "active" {
		return "", false
	}
	switch merchant.MerchantType {
	case "member_premium":
		return "Gold_Reseller", true
	case "h2h_api":
		return "H2H_Partner", true
	default: // guest_retail dan semua tipe lainnya
		return "Public_Retail", true
	}
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

	// if refCode == "" {
	// 	helpers.ProcessLogger(c, svc, "Reference code is required", "Validation error")
	// 	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	// }

	db := connections.DBconn()

	// Coba resolve segment dari JWT (opsional)
	segmentName, authenticated := resolveSegmentName(c, nil)

	var (
		list interface{}
		err  error
	)

	if authenticated {
		// User sudah login → tampilkan produk sesuai segment merchantnya saja
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, refCode, segmentName)
	} else {
		// Guest / tidak login → tampilkan semua segment (Public_Retail)
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, refCode, "Public_Retail")
	}

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get product segments")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, list))
}
