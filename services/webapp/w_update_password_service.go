package webapp

import (
	"fmt"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/models"
	"gamebiller/repositories"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func UpdatePassword(c echo.Context) error {
	var (
		svc = "UpdatePassword"
		req models.UpdatePasswordRequest
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	// 1. Get logged-in user from JWT claims
	claimsVal := c.Get("user")
	if claimsVal == nil {
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
	}
	token := claimsVal.(*jwt.Token)
	claims := token.Claims.(*models.JwtCustomClaims)

	db := connections.DBconn()

	// 2. Fetch logged-in user from DB
	user, err := repositories.GetUserByID(db, claims.UserID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "User not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	// 3. Verify old_password
	if req.OldPassword == "" || !helpers.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		helpers.ProcessLogger(c, svc, "Old password mismatch", "Authentication failed")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth401, nil))
	}

	// 4. Validate new_password length (min 8 characters)
	if len(req.NewPassword) < 8 {
		helpers.ProcessLogger(c, svc, "Password too short", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser424, nil))
	}

	// 5. Validate confirm_password
	if req.NewPassword != req.ConfirmPassword {
		helpers.ProcessLogger(c, svc, "Password confirmation mismatch", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser424, nil))
	}

	// 6. Hash new password
	newHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to hash new password")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 7. Update in DB
	now := time.Now().Format(time.RFC3339)
	updater := fmt.Sprintf("user_%d", user.ID)
	if err := repositories.UpdateUserPassword(db, user.ID, newHash, now, updater); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update password in DB")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucUser200, map[string]any{
		"user_id": user.ID,
		"message": "Password successfully updated",
	}))
}
