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

func AdminGetPaymentChannels(c echo.Context) error {
	var (
		svc = "AdminGetPaymentChannels"
		req models.RequestPaymentChannels
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetPaymentChannelsList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve payment channels list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreatePaymentChannel(c echo.Context) error {
	var (
		svc = "AdminCreatePaymentChannel"
		pc  models.PaymentChannel
	)
	if err := c.Bind(&pc); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	pc.CreatedAt = now
	pc.UpdatedAt = now
	pc.CreatedBy = "admin"
	pc.UpdatedBy = "admin"

	_, err := repositories.CreatePaymentChannel(connections.DBconn(), &pc)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pc))
}

func AdminUpdatePaymentChannel(c echo.Context) error {
	var (
		svc   = "AdminUpdatePaymentChannel"
		input models.PaymentChannel
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
	pc, err := repositories.GetPaymentChannelByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Payment channel not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	pc.PaymentMethodID = input.PaymentMethodID
	pc.ChannelCode = input.ChannelCode
	pc.ChannelName = input.ChannelName
	pc.FeeType = input.FeeType
	pc.FeeValue = input.FeeValue
	pc.IsActive = input.IsActive
	pc.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	pc.UpdatedBy = "admin"

	err = repositories.UpdatePaymentChannel(db, pc)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, pc))
}

func AdminDeletePaymentChannel(c echo.Context) error {
	var (
		svc = "AdminDeletePaymentChannel"
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

	err := repositories.DeletePaymentChannel(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete payment channel")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
