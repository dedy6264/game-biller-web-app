package dashboard

import (
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func AdminGetAgents(c echo.Context) error {
	var (
		svc = "AdminGetAgents"
		req models.RequestAgents
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetAgentsList(connections.DBconn(), req.Search, req.Start, req.Length, req.Order, req.Sort, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve agents list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateAgent(c echo.Context) error {
	var (
		svc = "AdminCreateAgent"
		a   models.Agent
	)
	if err := c.Bind(&a); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if a.AgentName == "" {
		helpers.ProcessLogger(c, svc, "Agent name is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if a.UserID == 0 {
		helpers.ProcessLogger(c, svc, "User ID is required", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if a.Status == "" {
		a.Status = "active"
	}

	// Auto-generate referral_code if empty
	if strings.TrimSpace(a.ReferralCode) == "" {
		a.ReferralCode = "AG-" + strings.ToUpper(helpers.RandomDigits(6))
	} else {
		a.ReferralCode = strings.ToUpper(strings.TrimSpace(a.ReferralCode))
	}

	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	a.CreatedAt = now
	a.UpdatedAt = now
	a.CreatedBy = "admin"
	a.UpdatedBy = "admin"

	_, err := repositories.CreateAgent(connections.DBconn(), &a)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create agent")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, a))
}

func AdminUpdateAgent(c echo.Context) error {
	var (
		svc   = "AdminUpdateAgent"
		input models.Agent
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
	a, err := repositories.GetAgentByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Agent not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	if input.AgentName != "" {
		a.AgentName = input.AgentName
	}
	if input.ReferralCode != "" {
		a.ReferralCode = strings.ToUpper(strings.TrimSpace(input.ReferralCode))
	}
	if input.Status != "" {
		a.Status = input.Status
	}
	a.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
	a.UpdatedBy = "admin"

	err = repositories.UpdateAgent(db, a)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update agent")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, a))
}

func AdminDeleteAgent(c echo.Context) error {
	var (
		svc = "AdminDeleteAgent"
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

	err := repositories.DeleteAgent(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete agent")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
