package dashboard

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminGetProductCategories(c echo.Context) error {
	var (
		svc = "AdminGetProductCategories"
		req models.RequestProductCategories
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProductCategoriesList(connections.DBconn(), req.Start, req.Length, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve product categories list")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProductCategory(c echo.Context) error {
	var (
		svc = "AdminCreateProductCategory"
		pc  models.ProductCategory
	)
	if err := c.Bind(&pc); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	_, err := repositories.CreateProductCategory(connections.DBconn(), &pc)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product category")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pc))
}

func AdminUpdateProductCategory(c echo.Context) error {
	var (
		svc   = "AdminUpdateProductCategory"
		input models.ProductCategory
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
	pc, err := repositories.GetProductCategoryByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product category not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	pc.ProductCategoryName = input.ProductCategoryName

	err = repositories.UpdateProductCategory(db, pc)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product category")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pc))
}

func AdminDeleteProductCategory(c echo.Context) error {
	var (
		svc = "AdminDeleteProductCategory"
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

	err := repositories.DeleteProductCategory(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete product category")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{"id": req.ID}))
}
