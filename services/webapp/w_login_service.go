package webapp

import (
	"gamebiller/configs"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var (
		req models.LoginRequest
		svc = "Login"
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}

	db := connections.DBconn()

	user, err := repositories.GetUserByEmailOrPhone(db, req.Username)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get user")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-401", nil))
	}

	if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
		helpers.ProcessLogger(c, svc, "Invalid password", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-401", nil))
	}

	// Fetch Role
	role, err := repositories.GetUserRole(db, user.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get user role")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-403", nil))
	}

	// Fetch Merchant
	merchant, err := repositories.GetMerchantByUserID(db, user.ID)
	var merchantID int64
	if err == nil {
		merchantID = merchant.ID
	} else {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get merchant")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-AUTH-403", nil))
	}

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.JwtCustomClaims{
		UserID:     user.ID,
		MerchantID: merchantID,
		RoleCode:   role.RoleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.APP_KEY))
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to sign JWT")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-AUTH-200", map[string]any{
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": 86400,
		"user": map[string]any{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"phone_number": user.PhoneNumber,
			"role":         role.RoleCode,
		},
	}))
}
