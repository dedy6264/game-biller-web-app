package dashboard

import (
	"encoding/json"
	"fmt"
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
		svc   = "AdminCreateProductSegment"
		ps    models.ProductSegment
		pprov models.ProductProvider
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

	db := connections.DBconn()
	if ps.ProductName == "" && ps.ProductID != 0 {
		if prod, err := repositories.GetProductByID(db, ps.ProductID); err == nil {
			ps.ProductName = prod.ProductName
		}
	}
	if ps.ProductProviderID == 0 {
		//return invalid product provider
		helpers.ProcessLogger(c, svc, "ProductProviderID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if ps.ProductID == 0 {
		//return invalid product
		helpers.ProcessLogger(c, svc, "ProductID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	p, err := repositories.GetProductByID(db, ps.ProductID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}
	pprov, err = repositories.GetProductProviderByID(db, ps.ProductProviderID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product provider not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}
	ps = models.ProductSegment{
		ID:                         ps.ID,
		SegmentName:                ps.SegmentName,
		SegmentID:                  ps.SegmentID,
		ProductProviderID:          pprov.ID,
		ProductID:                  p.ID,
		ProductName:                p.ProductName,
		ProviderName:               pprov.ProviderName,
		ProductPrice:               ps.ProductPrice,
		AdminFee:                   ps.AdminFee,
		MerchantFee:                ps.MerchantFee,
		ProductProviderCode:        pprov.ProductProviderCode,
		ProductProviderName:        pprov.ProductProviderName,
		ProductProviderPrice:       pprov.ProductProviderPrice,
		ProductProviderAdminFee:    pprov.ProductProviderAdminFee,
		ProductProviderMerchantFee: pprov.ProductProviderMerchantFee,
		CreatedAt:                  ps.CreatedAt,
		CreatedBy:                  ps.CreatedBy,
		UpdatedAt:                  ps.UpdatedAt,
		UpdatedBy:                  ps.UpdatedBy,
	}
	_, err = repositories.CreateProductSegment(db, ps)
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
	if input.ProductProviderID == 0 {
		//return invalid product provider
		helpers.ProcessLogger(c, svc, "ProductProviderID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if input.ProductID == 0 {
		//return invalid product
		helpers.ProcessLogger(c, svc, "ProductID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	p, err := repositories.GetProductByID(db, input.ProductID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}
	pprov, err := repositories.GetProductProviderByID(db, input.ProductProviderID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product provider not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}
	ss, _ := json.Marshal(pprov)
	fmt.Println("====", string(ss))
	ps = models.ProductSegment{
		ID:                         ps.ID,
		SegmentName:                ps.SegmentName,
		SegmentID:                  ps.SegmentID,
		ProductProviderID:          pprov.ID,
		ProductID:                  p.ID,
		ProductName:                p.ProductName,
		ProviderName:               pprov.ProviderName,
		ProductPrice:               input.ProductPrice,
		AdminFee:                   input.AdminFee,
		MerchantFee:                input.MerchantFee,
		ProductProviderCode:        pprov.ProductProviderCode,
		ProductProviderName:        pprov.ProductProviderName,
		ProductProviderPrice:       pprov.ProductProviderPrice,
		ProductProviderAdminFee:    pprov.ProductProviderAdminFee,
		ProductProviderMerchantFee: pprov.ProductProviderMerchantFee,
		CreatedAt:                  ps.CreatedAt,
		CreatedBy:                  ps.CreatedBy,
		UpdatedAt:                  time.Now().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedBy:                  "admin",
	}
	ss, _ = json.Marshal(ps)
	fmt.Println("====", string(ss))
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
