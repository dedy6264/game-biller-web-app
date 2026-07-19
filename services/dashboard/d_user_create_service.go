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

func AdminCreateUser(c echo.Context) error {
	var (
		svc = "AdminCreateUser"
		u   models.User
	)
	if err := c.Bind(&u); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	now := time.Now().Format(time.RFC3339)
	u.CreatedAt = now
	u.UpdatedAt = now
	
	// Check email/phone uniqueness
	db := connections.DBconn()
	if _, err := repositories.GetUserByEmailOrPhone(db, u.Email); err == nil {
		helpers.ProcessLogger(c, svc, "Email already exists", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("VAL-USER-422", nil))
	}
	if _, err := repositories.GetUserByEmailOrPhone(db, u.PhoneNumber); err == nil {
		helpers.ProcessLogger(c, svc, "Phone already exists", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("VAL-USER-423", nil))
	}

	// default password hash if not provided
	if u.PasswordHash == "" {
		hashed, _ := helpers.HashPassword("password123")
		u.PasswordHash = hashed
	}

	_, err := repositories.CreateUser(db, &u)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create user")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-AUTH-201", u))
}
