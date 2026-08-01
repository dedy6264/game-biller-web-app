package webapp

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetProductSegmentByCustomer(c echo.Context) error {
	var (
		svc        = "GetProductSegmentByCustomer"
		customerID = c.QueryParam("customer_id")
	)

	if customerID == "" {
		var body struct {
			CustomerID string `json:"customer_id"`
		}
		if err := c.Bind(&body); err == nil {
			customerID = body.CustomerID
		}
	}

	if customerID == "" {
		helpers.ProcessLogger(c, svc, "Customer ID is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if len(customerID) < 5 {
		helpers.ProcessLogger(c, svc, "Customer ID is too short", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// 1. Resolve product reference via prefix matching
	ref, err := repositories.MatchPhonePrefix(db, customerID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to match customer ID prefix")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys404, nil))
	}

	// 2. Coba resolve segment dari JWT (opsional)
	segmentID, authenticated := resolveSegmentName(c, nil)

	var list interface{}

	if authenticated {
		// User sudah login → tampilkan produk sesuai segment merchantnya saja
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, ref.ProductReferenceCode, segmentID)
	} else {
		// Guest / tidak login → tampilkan Public_Retail
		list, err = repositories.GetProductSegmentsByRefCodeAndSegment(db, ref.ProductReferenceCode, 1)
	}

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get product segments")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 3. Return resolved reference code beserta daftar produk di segment yang sesuai
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"reference_code": ref.ProductReferenceCode,
		// "segment":          list[0].SegmentName,
		"is_authenticated": authenticated,
		"product_segments": list,
	}))
}
