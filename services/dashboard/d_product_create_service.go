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

func AdminCreateProduct(c echo.Context) error {
	var (
		svc = "AdminCreateProduct"
		p   models.Product
	)
	if err := c.Bind(&p); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := repositories.CreateProduct(connections.DBconn(), &p)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create product")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", p))
}
