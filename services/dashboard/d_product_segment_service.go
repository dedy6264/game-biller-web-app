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

func AdminGetProductSegments(c echo.Context) error {
	var (
		svc = "AdminGetProductSegments"
		req models.RequestProductSegments
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductSegmentsList(connections.DBconn(), req.Start, req.Length, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product segments list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductSegment(c echo.Context) error {
	var (
		svc = "AdminCreateProductSegment"
		ps  models.ProductSegment
	)
	if err := c.Bind(&ps); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if ps.ProductMasterID == 0 {
		helpers.ProcessLogger(c, svc, "ProductMasterID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if ps.SegmentID == 0 {
		helpers.ProcessLogger(c, svc, "SegmentID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// Fetch product master for snapshot data
	pm, err := repositories.GetProductMasterByID(db, ps.ProductMasterID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product master not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	// Fetch segment for name
	seg, err := repositories.GetSegmentByID(db, ps.SegmentID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Segment not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	masterCode := pm.ProductProviderCode
	if ps.ProductMasterCode != "" {
		masterCode = ps.ProductMasterCode
	}

	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	newPS := models.ProductSegment{
		SegmentID:                ps.SegmentID,
		AgentID:                  seg.AgentID,
		ProductMasterID:          pm.ID,
		SegmentName:              seg.SegmentName,
		ProductMasterCode:        masterCode,
		ProductMasterName:        pm.ProductName,
		ProductMasterPrice:       pm.ProductPrice,
		ProductMasterAdminFee:    pm.AdminFee,
		ProductMasterMerchantFee: pm.MerchantFee,
		// Allow override; fall back to master values if zero
		ProductPrice:       ps.ProductPrice,
		ProductAdminFee:    ps.ProductAdminFee,
		ProductMerchantFee: ps.ProductMerchantFee,
		CreatedAt:          now,
		CreatedBy:          "admin",
		UpdatedAt:          now,
		UpdatedBy:          "admin",
	}
	if newPS.ProductPrice == 0 {
		newPS.ProductPrice = pm.ProductPrice
	}
	if newPS.ProductAdminFee == 0 {
		newPS.ProductAdminFee = pm.AdminFee
	}
	if newPS.ProductMerchantFee == 0 {
		newPS.ProductMerchantFee = pm.MerchantFee
	}

	_, err = repositories.CreateProductSegment(db, newPS)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, newPS))
}

func AdminUpdateProductSegment(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductSegment"
		input models.ProductSegment
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	ps, err := repositories.GetProductSegmentByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product segment not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	// Update product master if changed
	if input.ProductMasterID != 0 && input.ProductMasterID != ps.ProductMasterID {
		pm, err := repositories.GetProductMasterByID(db, input.ProductMasterID)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Product master not found")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
		}
		ps.ProductMasterID = pm.ID
		ps.ProductMasterCode = pm.ProductProviderCode
		ps.ProductMasterName = pm.ProductName
		ps.ProductMasterPrice = pm.ProductPrice
		ps.ProductMasterAdminFee = pm.AdminFee
		ps.ProductMasterMerchantFee = pm.MerchantFee
	}

	if input.ProductMasterCode != "" {
		ps.ProductMasterCode = input.ProductMasterCode
	}

	// Update segment if changed
	if input.SegmentID != 0 && input.SegmentID != ps.SegmentID {
		seg, err := repositories.GetSegmentByID(db, input.SegmentID)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Segment not found")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
		}
		ps.SegmentID = seg.ID
		ps.SegmentName = seg.SegmentName
		ps.AgentID = seg.AgentID
	}

	// Allow price overrides
	if input.ProductPrice != 0 {
		ps.ProductPrice = input.ProductPrice
	}
	if input.ProductAdminFee != 0 {
		ps.ProductAdminFee = input.ProductAdminFee
	}
	if input.ProductMerchantFee != 0 {
		ps.ProductMerchantFee = input.ProductMerchantFee
	}

	ps.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	ps.UpdatedBy = "admin"

	err = repositories.UpdateProductSegment(db, ps)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, ps))
}

func AdminDeleteProductSegment(c echo.Context) error {
	var (
		svc = "AdminDeleteProductSegment"
		req struct {
			ID int64 `json:"id"`
		}
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if req.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	err := repositories.DeleteProductSegment(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
