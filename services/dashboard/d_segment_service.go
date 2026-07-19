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

func AdminGetSegments(c echo.Context) error {
	var (
		svc = "AdminGetSegments"
		req models.RequestSegments
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetSegmentsList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve segments list")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateSegment(c echo.Context) error {
	var (
		svc = "AdminCreateSegment"
		s   models.Segment
	)
	if err := c.Bind(&s); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	if s.SegmentName == "" {
		helpers.ProcessLogger(c, svc, "Segment name is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	s.CreatedAt = now
	s.UpdatedAt = now
	s.CreatedBy = "admin"
	s.UpdatedBy = "admin"

	_, err := repositories.CreateSegment(connections.DBconn(), &s)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", s))
}

func AdminUpdateSegment(c echo.Context) error {
	var (
		svc   = "AdminUpdateSegment"
		input models.Segment
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	if input.SegmentName == "" {
		helpers.ProcessLogger(c, svc, "Segment name cannot be empty", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	db := connections.DBconn()
	s, err := repositories.GetSegmentByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Segment not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-USER-404", nil))
	}

	s.SegmentName = input.SegmentName
	s.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	s.UpdatedBy = "admin"

	err = repositories.UpdateSegment(db, s)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", s))
}

func AdminDeleteSegment(c echo.Context) error {
	var (
		svc = "AdminDeleteSegment"
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

	err := repositories.DeleteSegment(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete segment")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", map[string]any{"id": req.ID}))
}
