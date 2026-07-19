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

func AdminUpdateProduct(c echo.Context) error {
	var (
		svc   = "AdminUpdateProduct"
		input models.Product
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
	p, err := repositories.GetProductByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Product not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	p.ProductName = input.ProductName
	p.ProductCode = input.ProductCode
	p.IsActive = input.IsActive
	p.ProductReferenceID = input.ProductReferenceID
	p.ProductTypeID = input.ProductTypeID
	p.ProductCategoryID = input.ProductCategoryID
	p.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	err = repositories.UpdateProduct(db, p)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update product")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-USER-200", p))
}
