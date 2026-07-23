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

func AdminUpdateMerchant(c echo.Context) error {
	var (
		svc   = "AdminUpdateMerchant"
		input models.Merchant
	)
	if err := c.Bind(&input); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}
	if input.ApiCredential == nil && (input.ClientKey != "" || input.SecretKey != "" || input.WhitelistIPs != "") {
		input.ApiCredential = &models.MerchantApiCredential{
			ClientKey:    input.ClientKey,
			SecretKey:    input.SecretKey,
			WhitelistIPs: input.WhitelistIPs,
			IsActive:     input.ApiIsActive,
		}
	}
	if input.ID == 0 {
		helpers.ProcessLogger(c, svc, "ID cannot be zero", "Validation error")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeInvalidCustId, nil))
	}

	db := connections.DBconn()
	m, err := repositories.GetMerchantByID(db, input.ID)
	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Merchant not found")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrUser404, nil))
	}

	m.MerchantName = input.MerchantName
	m.MerchantType = input.MerchantType
	m.Status = input.Status
	m.SegmentID = input.SegmentID
	now := time.Now().Format("2006-01-02T15:04:05Z07:00")
	m.UpdatedAt = now

	err = helpers.DBTransaction(db, func(tx *sql.Tx) error {
		err = repositories.UpdateMerchant(tx, m)
		if err != nil {
			return err
		}

		if input.ApiCredential != nil {
			existing, err := repositories.GetMerchantApiCredentialByMerchantID(tx, input.ID)
			if err == nil {
				existing.ClientKey = input.ApiCredential.ClientKey
				existing.WhitelistIPs = input.ApiCredential.WhitelistIPs
				existing.IsActive = input.ApiCredential.IsActive
				existing.UpdatedAt = now
				existing.UpdatedBy = "admin"
				if input.ApiCredential.SecretKey != "" {
					hash, err := helpers.HashPassword(input.ApiCredential.SecretKey)
					if err != nil {
						return err
					}
					existing.SecretKeyHash = hash
				}
				err = repositories.UpdateMerchantApiCredential(tx, existing)
				if err != nil {
					return err
				}
				m.ApiCredential = existing
			} else {
				input.ApiCredential.MerchantID = input.ID
				input.ApiCredential.CreatedAt = now
				input.ApiCredential.UpdatedAt = now
				input.ApiCredential.CreatedBy = "admin"
				input.ApiCredential.UpdatedBy = "admin"
				if input.ApiCredential.SecretKey != "" {
					hash, err := helpers.HashPassword(input.ApiCredential.SecretKey)
					if err != nil {
						return err
					}
					input.ApiCredential.SecretKeyHash = hash
				}
				_, err = repositories.CreateMerchantApiCredential(tx, input.ApiCredential)
				if err != nil {
					return err
				}
				m.ApiCredential = input.ApiCredential
			}
		} else {
			existing, err := repositories.GetMerchantApiCredentialByMerchantID(tx, input.ID)
			if err == nil {
				m.ApiCredential = existing
			}
		}
		return nil
	})

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to update merchant")
		return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeErrSys500, nil))
	}
	return c.JSON(http.StatusOK, helpers.BuildResponse(helpers.CodeSucUser200, m))
}
