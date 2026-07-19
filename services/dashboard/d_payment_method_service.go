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

func AdminGetPaymentMethods(c echo.Context) error {
	var (
		svc = "AdminGetPaymentMethods"
		req models.RequestPaymentMethods
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetPaymentMethodsList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve payment methods list")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreatePaymentMethod(c echo.Context) error {
	var (
		svc = "AdminCreatePaymentMethod"
		pm  models.PaymentMethod
	)
	if err := c.Bind(&pm); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pm.CreatedAt = now
	pm.UpdatedAt = now
	pm.CreatedBy = "admin"
	pm.UpdatedBy = "admin"

	_, err := repositories.CreatePaymentMethod(connections.DBconn(), &pm)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create payment method")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pm))
}

func AdminUpdatePaymentMethod(c echo.Context) error {
	var (
		svc   = "AdminUpdatePaymentMethod"
		input models.PaymentMethod
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
	pm, err := repositories.GetPaymentMethodByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Payment method not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	pm.MethodCode = input.MethodCode
	pm.MethodName = input.MethodName
	pm.IsActive = input.IsActive
	pm.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pm.UpdatedBy = "admin"

	err = repositories.UpdatePaymentMethod(db, pm)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update payment method")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", pm))
}

func AdminDeletePaymentMethod(c echo.Context) error {
	var (
		svc = "AdminDeletePaymentMethod"
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

	err := repositories.DeletePaymentMethod(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete payment method")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{"id": req.ID}))
}
