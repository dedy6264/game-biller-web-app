package webapp

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

func Register(c echo.Context) error {
	var (
		svc = "Register"
		req models.RegisterRequest
	)
	if err := c.Bind(&req); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	if len(req.Password) < 8 {
		helpers.ProcessLogger(c, svc, "Password too short", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser424, nil))
	}

	db := connections.DBconn()

	// Check email exist
	_, err := repositories.GetUserByEmailOrPhone(db, req.Email)
	if err == nil {
		helpers.ProcessLogger(c, svc, "Email already exists", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser422, nil))
	}

	// Check phone exist
	_, err = repositories.GetUserByEmailOrPhone(db, req.PhoneNumber)
	if err == nil {
		helpers.ProcessLogger(c, svc, "Phone number already exists", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeValUser423, nil))
	}

	// Hash password
	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to hash password")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	// Lookup segment Public_Retail untuk auto-binding saat register
	var publicRetailSegmentID int64
	seg, err := repositories.GetSegmentByID(db, 2)
	if err == nil {
		publicRetailSegmentID = seg.ID
	} else {
		helpers.ProcessLogger(c, svc, err.Error(), "Retail Biller segment not found, proceeding without segment binding")
	}

	now := time.Now().Format(time.RFC3339)

	var userID int64
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		// 1. Create User
		user := models.User{
			Name:         req.Name,
			Email:        req.Email,
			PhoneNumber:  req.PhoneNumber,
			PasswordHash: hashed,
			Status:       "active",
			CreatedAt:    now,
			CreatedBy:    "system",
			UpdatedAt:    now,
			UpdatedBy:    "system",
		}
		uid, err := repositories.CreateUser(tx, &user)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to create user")
			return err
		}
		userID = uid

		// 2. Map Role (member_reseller = ID 4)
		mhr := models.ModelHasRole{
			UserID:    uid,
			RoleID:    4,
			CreatedAt: now,
			CreatedBy: "system",
		}
		_, err = repositories.CreateModelHasRole(tx, &mhr)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to map user role")
			return err
		}

		// 3. Create Merchant (guest_retail) — auto-bind segment Public_Retail
		merch := models.Merchant{
			UserID:       uid,
			SegmentID:    publicRetailSegmentID,
			MerchantName: req.Name + " Store",
			MerchantType: "Retail Biller",
			Status:       "active",
			CreatedAt:    now,
			CreatedBy:    "system",
			UpdatedAt:    now,
			UpdatedBy:    "system",
		}
		mid, err := repositories.CreateMerchant(tx, &merch)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to create merchant")
			return err
		}

		// 4. Create Saving Account with PIN (default PIN "123456") and 1,000,000 balance for testing
		pinHash, _ := helpers.HashPassword("123456")
		sa := models.SavingAccount{
			MerchantID:     mid,
			AccountNumber:  "SA-" + helpers.RandomDigits(8),
			Balance:        1000000.00,
			AccountPinHash: pinHash,
			Status:         "active",
			CreatedAt:      now,
			CreatedBy:      "system",
			UpdatedAt:      now,
			UpdatedBy:      "system",
		}
		_, err = repositories.CreateSavingAccount(tx, &sa)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to create saving account")
			return err
		}
		return nil
	})

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Registration transaction failed")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucAuth201, map[string]any{
		"user_id": userID,
	}))
}
