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

func AdminUpdateUser(c echo.Context) error {
	var (
		svc   = "AdminUpdateUser"
		input models.User
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
	u, err := repositories.GetUserByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "User not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	u.Name = input.Name
	u.Email = input.Email
	u.PhoneNumber = input.PhoneNumber
	u.Status = input.Status
	u.UpdatedAt = time.Now().Format(time.RFC3339)
	u.UpdatedBy = "admin"

	err = repositories.UpdateUser(db, u)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update user")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucUser200, u))
}
