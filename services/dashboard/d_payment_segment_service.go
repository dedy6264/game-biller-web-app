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

func AdminGetPaymentSegments(c echo.Context) error {
	var (
		svc = "AdminGetPaymentSegments"
		req models.RequestPaymentSegments
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	list, total, err := repositories.GetPaymentSegmentsList(db, req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to fetch payment segments list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status_code":    helpers.CodeSuccess,
		"status_message": "SUCCESS",
		"status_desc":    "Payment segments retrieved successfully",
		"ui_message":     "Data payment segments berhasil dimuat.",
		"result": map[string]any{
			"draw":            req.Draw,
			"recordsTotal":    total,
			"recordsFiltered": total,
			"data":            list,
		},
	})
}

func AdminCreatePaymentSegment(c echo.Context) error {
	var (
		svc = "AdminCreatePaymentSegment"
		ps  models.PaymentSegment
	)
	if err := c.Bind(&ps); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	// Fill segment info if provided
	if ps.SegmentID > 0 && ps.SegmentName == "" {
		seg, err := repositories.GetSegmentByID(db, ps.SegmentID)
		if err == nil && seg != nil {
			ps.SegmentName = seg.SegmentName
		}
	}

	// Fill payment channel info
	if ps.PaymentChannelID > 0 {
		pc, err := repositories.GetPaymentChannelByID(db, ps.PaymentChannelID)
		if err == nil && pc != nil {
			ps.PaymentMethodID = pc.PaymentMethodID
			ps.ChannelCode = pc.ChannelCode
			ps.ChannelName = pc.ChannelName

			pm, errPm := repositories.GetPaymentMethodByID(db, pc.PaymentMethodID)
			if errPm == nil && pm != nil {
				ps.MethodCode = pm.MethodCode
			}
		}
	}

	now := time.Now().Format(time.RFC3339)
	ps.CreatedAt = now
	ps.UpdatedAt = now
	ps.CreatedBy = "admin"
	ps.UpdatedBy = "admin"

	id, err := repositories.CreatePaymentSegment(db, &ps)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create payment segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"id": id,
	}))
}

func AdminUpdatePaymentSegment(c echo.Context) error {
	var (
		svc = "AdminUpdatePaymentSegment"
		ps  models.PaymentSegment
	)
	if err := c.Bind(&ps); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	existing, err := repositories.GetPaymentSegmentByID(db, ps.ID)
	if err != nil || existing == nil {
		helpers.ProcessLogger(c, svc, "Payment segment not found", "Not Found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys404, nil))
	}

	// Fill segment info if provided
	if ps.SegmentID > 0 {
		seg, err := repositories.GetSegmentByID(db, ps.SegmentID)
		if err == nil && seg != nil {
			ps.SegmentName = seg.SegmentName
		}
	} else {
		ps.SegmentID = existing.SegmentID
		ps.SegmentName = existing.SegmentName
	}

	// Fill payment channel info if changed
	if ps.PaymentChannelID > 0 {
		pc, err := repositories.GetPaymentChannelByID(db, ps.PaymentChannelID)
		if err == nil && pc != nil {
			ps.PaymentMethodID = pc.PaymentMethodID
			ps.ChannelCode = pc.ChannelCode
			ps.ChannelName = pc.ChannelName

			pm, errPm := repositories.GetPaymentMethodByID(db, pc.PaymentMethodID)
			if errPm == nil && pm != nil {
				ps.MethodCode = pm.MethodCode
			}
		}
	} else {
		ps.PaymentChannelID = existing.PaymentChannelID
		ps.PaymentMethodID = existing.PaymentMethodID
		ps.ChannelCode = existing.ChannelCode
		ps.ChannelName = existing.ChannelName
		ps.MethodCode = existing.MethodCode
	}

	ps.UpdatedAt = time.Now().Format(time.RFC3339)
	ps.UpdatedBy = "admin"

	if err := repositories.UpdatePaymentSegment(db, &ps); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update payment segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"id": ps.ID,
	}))
}

func AdminDeletePaymentSegment(c echo.Context) error {
	var (
		svc = "AdminDeletePaymentSegment"
		req struct {
			ID int64 `json:"id"`
		}
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	if err := repositories.DeletePaymentSegment(db, req.ID); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete payment segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"id": req.ID,
	}))
}
