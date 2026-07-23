package dashboard

import (
	"database/sql"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func AdminGetModelHasRoles(c echo.Context) error {
	var (
		svc = "AdminGetModelHasRoles"
		req models.RequestModelHasRoles
	)
	_ = c.Bind(&req)
	list, total, err := repositories.GetModelHasRolesList(connections.DBconn(), req.Start, req.Length, req.Filters)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to retrieve model has roles list")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{
		"draw":            req.Draw,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"data":            list,
	}))
}

func AdminCreateModelHasRole(c echo.Context) error {
	var (
		svc = "AdminCreateModelHasRole"
		mhr models.ModelHasRole
	)
	if err := c.Bind(&mhr); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	mhr.CreatedAt = now
	mhr.CreatedBy = "admin"

	_, err := repositories.CreateModelHasRole(connections.DBconn(), &mhr)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create model has role mapping")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, mhr))
}

func AdminUpdateModelHasRole(c echo.Context) error {
	var (
		svc   = "AdminUpdateModelHasRole"
		input models.ModelHasRole
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
	existing, err := repositories.GetModelHasRoleByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "ModelHasRole mapping not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		err = repositories.DeleteModelHasRole(tx, input.ID)
		if err != nil {
			return err
		}
		existing.UserID = input.UserID
		existing.RoleID = input.RoleID
		existing.CreatedAt = time.Now().Format("2006-01-02T15:04:05Z07:00")
		existing.CreatedBy = "admin"
		_, err = repositories.CreateModelHasRole(tx, existing)
		return err
	})

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update model has role mapping")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, existing))
}

func AdminDeleteModelHasRole(c echo.Context) error {
	var (
		svc = "AdminDeleteModelHasRole"
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

	err := repositories.DeleteModelHasRole(connections.DBconn(), req.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to delete model has role mapping")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccess, map[string]any{"id": req.ID}))
}
