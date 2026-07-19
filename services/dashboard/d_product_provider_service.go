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

func AdminGetProductProviders(c echo.Context) error {
	var (
		svc = "AdminGetProductProviders"
		req models.RequestProductProviders
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductProvidersList(connections.DBconn(), req.Start, req.Length, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product providers list")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductProvider(c echo.Context) error {
	var (
		svc = "AdminCreateProductProvider"
		pp  models.ProductProvider
	)
	if err := c.Bind(&pp); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pp.CreatedAt = now
	pp.UpdatedAt = now
	pp.CreatedBy = "admin"
	pp.UpdatedBy = "admin"

	_, err := repositories.CreateProductProvider(connections.DBconn(), &pp)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pp))
}

func AdminUpdateProductProvider(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductProvider"
		input models.ProductProvider
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	db := connections.DBconn()
	pp, err := repositories.GetProductProviderByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product provider not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	pp.ProviderID = input.ProviderID
	pp.ProviderProductCode = input.ProviderProductCode
	pp.ProviderPrice = input.ProviderPrice
	pp.ProviderAdminFee = input.ProviderAdminFee
	pp.ProviderMerchantFee = input.ProviderMerchantFee
	pp.ProviderIndex = input.ProviderIndex
	pp.IsAvailable = input.IsAvailable
	pp.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pp.UpdatedBy = "admin"

	err = repositories.UpdateProductProvider(db, pp)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pp))
}

func AdminDeleteProductProvider(c echo.Context) error {
	var (
		svc = "AdminDeleteProductProvider"
		req struct {
			ID int64 `json:"id"`
		}
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	if req.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	err := repositories.DeleteProductProvider(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{"id": req.ID}))
}
