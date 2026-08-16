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

func AdminGetProductMasters(c echo.Context) error {
	var (
		svc = "AdminGetProductMasters"
		req models.RequestProductMasters
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductMastersList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product masters list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductMaster(c echo.Context) error {
	var (
		svc = "AdminCreateProductMaster"
		pm  models.ProductMaster
	)
	if err := c.Bind(&pm); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if pm.ProductName == "" {
		helpers.ProcessLogger(c, svc, "Product name is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if pm.ProductProviderID == 0 {
		helpers.ProcessLogger(c, svc, "Product provider ID is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if pm.ProductID == 0 {
		helpers.ProcessLogger(c, svc, "Product ID is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// Auto-lookup ProviderID from ProductProvider if zero
	if pm.ProviderID == 0 && pm.ProductProviderID != 0 {
		if pprov, err := repositories.GetProductProviderByID(db, pm.ProductProviderID); err == nil {
			pm.ProviderID = pprov.ProviderID
			if pm.ProviderName == "" {
				pm.ProviderName = pprov.ProviderName
			}
		}
	}

	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pm.CreatedAt = now
	pm.UpdatedAt = now
	pm.CreatedBy = "admin"
	pm.UpdatedBy = "admin"

	_, err := repositories.CreateProductMaster(db, &pm)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product master")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pm))
}

func AdminUpdateProductMaster(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductMaster"
		input models.ProductMaster
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
	pm, err := repositories.GetProductMasterByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product master not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	// Update only provided fields
	if input.ProviderID != 0 {
		pm.ProviderID = input.ProviderID
	}
	if input.ProductName != "" {
		pm.ProductName = input.ProductName
	}
	if input.ProductSegmentIndex != "" {
		pm.ProductSegmentIndex = input.ProductSegmentIndex
	}
	if input.ProductProviderID != 0 {
		pm.ProductProviderID = input.ProductProviderID
	}
	if input.ProductID != 0 {
		pm.ProductID = input.ProductID
	}
	if input.ProviderName != "" {
		pm.ProviderName = input.ProviderName
	}
	if input.ProductProviderCode != "" {
		pm.ProductProviderCode = input.ProductProviderCode
	}
	if input.ProductProviderName != "" {
		pm.ProductProviderName = input.ProductProviderName
	}
	if input.ProductPrice != 0 {
		pm.ProductPrice = input.ProductPrice
	}
	if input.AdminFee != 0 {
		pm.AdminFee = input.AdminFee
	}
	if input.MerchantFee != 0 {
		pm.MerchantFee = input.MerchantFee
	}
	if input.ProductProviderPrice != 0 {
		pm.ProductProviderPrice = input.ProductProviderPrice
	}
	if input.ProductProviderAdminFee != 0 {
		pm.ProductProviderAdminFee = input.ProductProviderAdminFee
	}
	if input.ProductProviderMerchantFee != 0 {
		pm.ProductProviderMerchantFee = input.ProductProviderMerchantFee
	}

	pm.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pm.UpdatedBy = "admin"

	err = repositories.UpdateProductMaster(db, pm)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product master")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pm))
}

func AdminDeleteProductMaster(c echo.Context) error {
	var (
		svc = "AdminDeleteProductMaster"
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

	err := repositories.DeleteProductMaster(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product master")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
