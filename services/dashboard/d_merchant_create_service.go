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

func AdminCreateMerchant(c echo.Context) error {
	var (
		svc = "AdminCreateMerchant"
		m   models.Merchant
	)
	if err := c.Bind(&m); err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to bind request")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-VAL-104", nil))
	}
	if m.ApiCredential == nil && (m.ClientKey != "" || m.SecretKey != "" || m.WhitelistIPs != "") {
		m.ApiCredential = &models.MerchantApiCredential{
			ClientKey:    m.ClientKey,
			SecretKey:    m.SecretKey,
			WhitelistIPs: m.WhitelistIPs,
			IsActive:     m.ApiIsActive,
		}
	}
	now := time.Now().Format(time.RFC3339)
	m.CreatedAt = now
	m.UpdatedAt = now

	db := connections.DBconn()
	err := helpers.DBTransaction(db, func(tx *sql.Tx) error {
		mid, err := repositories.CreateMerchant(tx, &m)
		if err != nil {
			return err
		}
		m.ID = mid

		if m.ApiCredential != nil {
			m.ApiCredential.MerchantID = mid
			m.ApiCredential.CreatedAt = now
			m.ApiCredential.UpdatedAt = now
			if m.ApiCredential.CreatedBy == "" {
				m.ApiCredential.CreatedBy = "admin"
			}
			if m.ApiCredential.UpdatedBy == "" {
				m.ApiCredential.UpdatedBy = "admin"
			}
			if m.ApiCredential.SecretKey != "" {
				hash, err := helpers.HashPassword(m.ApiCredential.SecretKey)
				if err != nil {
					return err
				}
				m.ApiCredential.SecretKeyHash = hash
			}
			_, err = repositories.CreateMerchantApiCredential(tx, m.ApiCredential)
			return err
		}
		return nil
	})

	if err != nil {
		helpers.ProcessLogger(c, svc, err.Error(), "Failed to create merchant")
		return c.JSON(http.StatusOK, helpers.BuildResponse("ERR-SYS-500", nil))
	}

	return c.JSON(http.StatusOK, helpers.BuildResponse("SUC-INT-000", m))
}
