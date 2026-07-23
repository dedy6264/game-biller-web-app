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
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	ps.CreatedAt = now
	ps.UpdatedAt = now
	ps.CreatedBy = "admin"
	ps.UpdatedBy = "admin"

	_, err := repositories.CreateProductSegment(connections.DBconn(), &ps)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, ps))
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

	ps.SegmentName = input.SegmentName
	ps.SegmentID = input.SegmentID
	ps.ProductProviderID = input.ProductProviderID
	ps.ProductID = input.ProductID
	ps.ProductPrice = input.ProductPrice
	ps.AdminFee = input.AdminFee
	ps.MerchantFee = input.MerchantFee
	ps.ProviderProductPrice = input.ProviderProductPrice
	ps.ProviderProductAdminFee = input.ProviderProductAdminFee
	ps.ProviderProductMerchantFee = input.ProviderProductMerchantFee
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
