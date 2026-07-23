package dashboard

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

func AdminLogin(c echo.Context) error {
	var (
		svc = "AdminLogin"
		req models.LoginRequest
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()

	user, err := repositories.GetUserByEmailOrPhone(db, req.Username)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to find user")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth401, nil))
	}

	if !helpers.CheckPasswordHash(req.Password, user.PasswordHash) {
		helpers.ProcessLogger(c, svc, "Password hash mismatch", "Authentication failed")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth401, nil))
	}

	// Fetch Role
	role, err := repositories.GetUserRole(db, user.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to get user role")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
	}

	// Only allow internal admin roles
	if role.RoleCode != "super_admin" && role.RoleCode != "finance" && role.RoleCode != "cs" {
		helpers.ProcessLogger(c, svc, "User is not an admin", "Access denied")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
	}

	// Generate JWT
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &models.JwtCustomClaims{
		UserID:     user.ID,
		MerchantID: 0,
		RoleCode:   role.RoleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(configs.APP_KEY))
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to sign JWT token")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSuccessAuth, map[string]any{
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": 43200,
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  role.RoleCode,
		},
	}))
}
