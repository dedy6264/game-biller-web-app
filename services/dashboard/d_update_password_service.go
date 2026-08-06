package dashboard

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

func AdminUpdatePassword(c echo.Context) error {
	var (
		svc = "AdminUpdatePassword"
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

	// 2. Determine target user ID
	targetUserID := claims.UserID
	isSelfUpdate := true

	if req.UserID > 0 && req.UserID != claims.UserID {
		// Only super_admin can update another user's password
		if claims.RoleCode != "super_admin" {
			helpers.ProcessLogger(c, svc, "Permission denied to update another user password", "Forbidden")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth403, nil))
		}
		targetUserID = req.UserID
		isSelfUpdate = false
	}

	// 3. Fetch target user
	targetUser, err := repositories.GetUserByID(db, targetUserID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Target user not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	// 4. If self update or old_password provided, verify old_password
	if isSelfUpdate || req.OldPassword != "" {
		if req.OldPassword == "" || !helpers.CheckPasswordHash(req.OldPassword, targetUser.PasswordHash) {
			helpers.ProcessLogger(c, svc, "Old password mismatch", "Authentication failed")
			return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrAuth401, nil))
		}
	}

	// 5. Validate new_password length (min 8 characters)
	if len(req.NewPassword) < 8 {
		helpers.ProcessLogger(c, svc, "Password too short", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser424, nil))
	}

	// 6. Validate confirm_password
	if req.NewPassword != req.ConfirmPassword {
		helpers.ProcessLogger(c, svc, "Password confirmation mismatch", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser424, nil))
	}

	// 7. Hash new password
	newHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to hash new password")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// 8. Update in DB
	now := time.Now().Format(time.RFC3339)
	updater := fmt.Sprintf("user_%d", claims.UserID)
	if err := repositories.UpdateUserPassword(db, targetUser.ID, newHash, now, updater); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update password in DB")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucUser200, map[string]any{
		"user_id": targetUser.ID,
		"message": "Password successfully updated",
	}))
}
