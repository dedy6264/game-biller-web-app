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
	// 1.cek email dan phone number sudah ada atau belum,
	// 2.jika refferal == "", gunakan default agent dan segment yaitu agen_id = 1 dan segment_id = 2
	//3. jika referal!="", cari agent berdasarkan kode referal dan segment pertama yang dimiliki oleh agent
	//4. jika tidak ada agent yang ditemukan, gunakan default agent dan segment
	//5. lanjut pembuatan user dan binding dengan role sebagai merchant
	//6. jika segment tidak ada, gunakan default segment
	//7. lanjut pembuatan saving account
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

	// Tentukan Agent & Segment berdasarkan referral code atau default
	var (
		assignedAgentID   int64 = 1 // Default Agent
		assignedSegmentID int64 = 2 // Default Public_Retail Segment
	)

	if req.ReferralCode != "" {
		// Jika ada referral code: cari agent aktif berdasarkan kode tersebut
		ag, agErr := repositories.GetAgentByReferralCode(db, req.ReferralCode)
		if agErr == nil && ag != nil && ag.Status == "active" {
			assignedAgentID = ag.ID
			// Ambil segment pertama milik agent tersebut
			segs, _, segErr := repositories.GetSegmentsList(db, "", 0, 1, "", "", models.SegmentFilters{AgentID: ag.ID})
			if segErr == nil && len(segs) > 0 {
				assignedSegmentID = segs[0].ID
			}
			// Jika segment tidak ditemukan, tetap gunakan default segment (ID 2)
		} else {
			// Referral code tidak valid atau agent tidak aktif — fallback ke default
			helpers.ProcessLogger(c, svc, "Referral code tidak valid atau agent tidak aktif, menggunakan default", "Referral warning")
		}
	} else {
		// Tidak ada referral: gunakan default agent dan segment
		// Verifikasi bahwa default segment (ID 2) masih ada
		if seg, segErr := repositories.GetSegmentByID(db, 2); segErr == nil {
			assignedSegmentID = seg.ID
			assignedAgentID = seg.AgentID
		}
		// Jika default segment tidak ada, tetap gunakan hardcoded ID 1 & 2 sebagai fallback terakhir
	}

	now := time.Now().Format("2006-01-02T15:04:05Z07:00")

	var userID int64
	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		// 1. Buat User
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

		// 2. Bind role merchant (role_id = 4), actor_id diisi 0 dulu — diupdate setelah merchant dibuat
		mhr := models.ModelHasRole{
			UserID:    uid,
			RoleID:    4,
			ActorID:   0,
			CreatedAt: now,
			CreatedBy: "system",
		}
		mhrID, err := repositories.CreateModelHasRole(tx, &mhr)
		if err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to map user role")
			return err
		}

		// 3. Buat Merchant — auto-bind AgentID & SegmentID dari hasil referral / default
		merch := models.Merchant{
			UserID:       uid,
			AgentID:      assignedAgentID,
			SegmentID:    assignedSegmentID,
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

		// 4. Update actor_id di model_has_roles dengan merchant_id yang baru dibuat
		if err := repositories.UpdateModelHasRoleActorID(tx, mhrID, mid); err != nil {
			helpers.ProcessLogger(c, svc, err.Error(), "Failed to update actor_id on model_has_roles")
			return err
		}

		// 5. Buat Saving Account dengan PIN default "123456" dan saldo awal 1.000.000 (untuk testing)
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
		"user_id":  userID,
		"agent_id": assignedAgentID,
	}))
}
