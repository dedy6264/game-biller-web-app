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

func AdminGetProviders(c echo.Context) error {
	var (
		svc = "AdminGetProviders"
		req models.RequestProviders
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetProvidersList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve providers list")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateProvider(c echo.Context) error {
	var (
		svc = "AdminCreateProvider"
		p   models.Provider
	)
	if err := c.Bind(&p); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	p.CreatedAt = now
	p.UpdatedAt = now
	p.CreatedBy = "admin"
	p.UpdatedBy = "admin"

	_, err := repositories.CreateProvider(connections.DBconn(), &p)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", p))
}

func AdminUpdateProvider(c echo.Context) error {
	var (
		svc   = "AdminUpdateProvider"
		input models.Provider
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
	p, err := repositories.GetProviderByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Provider not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	p.ProviderName = input.ProviderName
	p.IsActive = input.IsActive
	p.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	p.UpdatedBy = "admin"

	err = repositories.UpdateProvider(db, p)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", p))
}

func AdminDeleteProvider(c echo.Context) error {
	var (
		svc = "AdminDeleteProvider"
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

	err := repositories.DeleteProvider(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete provider")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{"id": req.ID}))
}
