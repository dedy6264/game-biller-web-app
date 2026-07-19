package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// === MERCHANT REPOSITORY ===

func CreateMerchant(exec QueryExecutor, m *models.Merchant) (int64, error) {
	query := `INSERT INTO merchants (user_id, segment_id, merchant_name, merchant_type, status, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT (user_id) DO UPDATE SET segment_id = EXCLUDED.segment_id, merchant_name = EXCLUDED.merchant_name, merchant_type = EXCLUDED.merchant_type, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at, updated_by = EXCLUDED.updated_by
	          RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, m.UserID, m.SegmentID, m.MerchantName, m.MerchantType, m.Status, m.CreatedAt, m.CreatedBy, m.UpdatedAt, m.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	m.ID = id
	return id, nil
}

func GetMerchantByID(exec QueryExecutor, id int64) (*models.Merchant, error) {
	query := `SELECT m.id, m.user_id, m.segment_id, m.merchant_name, m.merchant_type, m.status, m.created_at, m.created_by, m.updated_at, m.updated_by, COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(s.segment_name, ''), COALESCE(mac.client_key, ''), COALESCE(mac.whitelist_ips, ''), COALESCE(mac.is_active, false)
	          FROM merchants m
	          LEFT JOIN users u ON u.id = m.user_id
	          LEFT JOIN segments s ON s.id = m.segment_id
	          LEFT JOIN merchant_api_credentials mac ON mac.merchant_id = m.id
	          WHERE m.id = $1`
	var m models.Merchant
	err := exec.QueryRow(query, id).Scan(&m.ID, &m.UserID, &m.SegmentID, &m.MerchantName, &m.MerchantType, &m.Status, &m.CreatedAt, &m.CreatedBy, &m.UpdatedAt, &m.UpdatedBy, &m.UserName, &m.UserEmail, &m.SegmentName, &m.ClientKey, &m.WhitelistIPs, &m.ApiIsActive)
	if err != nil {
		return nil, err
	}
	if m.ClientKey != "" {
		m.ApiCredential = &models.MerchantApiCredential{
			MerchantID:   m.ID,
			ClientKey:    m.ClientKey,
			WhitelistIPs: m.WhitelistIPs,
			IsActive:     m.ApiIsActive,
		}
	}
	return &m, nil
}

func GetMerchantByUserID(exec QueryExecutor, userID int64) (*models.Merchant, error) {
	query := `SELECT m.id, m.user_id, m.segment_id, m.merchant_name, m.merchant_type, m.status, m.created_at, m.created_by, m.updated_at, m.updated_by, COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(s.segment_name, ''), COALESCE(mac.client_key, ''), COALESCE(mac.whitelist_ips, ''), COALESCE(mac.is_active, false)
	          FROM merchants m
	          LEFT JOIN users u ON u.id = m.user_id
	          LEFT JOIN segments s ON s.id = m.segment_id
	          LEFT JOIN merchant_api_credentials mac ON mac.merchant_id = m.id
	          WHERE m.user_id = $1`
	var m models.Merchant
	err := exec.QueryRow(query, userID).Scan(&m.ID, &m.UserID, &m.SegmentID, &m.MerchantName, &m.MerchantType, &m.Status, &m.CreatedAt, &m.CreatedBy, &m.UpdatedAt, &m.UpdatedBy, &m.UserName, &m.UserEmail, &m.SegmentName, &m.ClientKey, &m.WhitelistIPs, &m.ApiIsActive)
	if err != nil {
		return nil, err
	}
	if m.ClientKey != "" {
		m.ApiCredential = &models.MerchantApiCredential{
			MerchantID:   m.ID,
			ClientKey:    m.ClientKey,
			WhitelistIPs: m.WhitelistIPs,
			IsActive:     m.ApiIsActive,
		}
	}
	return &m, nil
}

func UpdateMerchant(exec QueryExecutor, m *models.Merchant) error {
	query := `UPDATE merchants SET segment_id = ?, merchant_name = ?, merchant_type = ?, status = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, m.SegmentID, m.MerchantName, m.MerchantType, m.Status, m.UpdatedAt, m.UpdatedBy, m.ID)
	return err
}

func DeleteMerchant(exec QueryExecutor, id int64) error {
	query := `DELETE FROM merchants WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetMerchantsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.MerchantFilters) ([]models.Merchant, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.UserID != 0 {
		whr += " AND user_id = " + strconv.FormatInt(filters.UserID, 10)
	}
	if filters.SegmentID != nil && *filters.SegmentID != 0 {
		whr += " AND segment_id = " + strconv.FormatInt(*filters.SegmentID, 10)
	}
	if filters.MerchantType != "" {
		whr += " AND merchant_type = '" + filters.MerchantType + "'"
	}
	if filters.Status != "" {
		whr += " AND status = '" + filters.Status + "'"
	}
	fmt.Println(":::", whr)
	countQuery := `SELECT COUNT(*) FROM merchants WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT m.id, m.user_id, m.segment_id, m.merchant_name, m.merchant_type, m.status, m.created_at, m.created_by, m.updated_at, m.updated_by, COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(s.segment_name, ''), COALESCE(mac.client_key, ''), COALESCE(mac.whitelist_ips, ''), COALESCE(mac.is_active, false)
	          FROM merchants m
	          LEFT JOIN users u ON u.id = m.user_id
	          LEFT JOIN segments s ON s.id = m.segment_id
	          LEFT JOIN merchant_api_credentials mac ON mac.merchant_id = m.id
	          WHERE true ` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		// Map order field prefix to avoid ambiguity
		if order == "id" {
			order = "m.id"
		} else if order == "created_at" {
			order = "m.created_at"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY m.id DESC"
	}
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Merchant
	for rows.Next() {
		var m models.Merchant
		err = rows.Scan(&m.ID, &m.UserID, &m.SegmentID, &m.MerchantName, &m.MerchantType, &m.Status, &m.CreatedAt, &m.CreatedBy, &m.UpdatedAt, &m.UpdatedBy, &m.UserName, &m.UserEmail, &m.SegmentName, &m.ClientKey, &m.WhitelistIPs, &m.ApiIsActive)
		if err != nil {
			return nil, 0, err
		}
		if m.ClientKey != "" {
			m.ApiCredential = &models.MerchantApiCredential{
				MerchantID:   m.ID,
				ClientKey:    m.ClientKey,
				WhitelistIPs: m.WhitelistIPs,
				IsActive:     m.ApiIsActive,
			}
		}
		list = append(list, m)
	}
	return list, count, nil
}

// === MERCHANT API CREDENTIALS REPOSITORY ===

func CreateMerchantApiCredential(exec QueryExecutor, mac *models.MerchantApiCredential) (int64, error) {
	query := `INSERT INTO merchant_api_credentials (merchant_id, client_key, secret_key_hash, whitelist_ips, is_active, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, mac.MerchantID, mac.ClientKey, mac.SecretKeyHash, mac.WhitelistIPs, mac.IsActive, mac.CreatedAt, mac.CreatedBy, mac.UpdatedAt, mac.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	mac.ID = id
	return id, nil
}

func GetMerchantApiCredentialByID(exec QueryExecutor, id int64) (*models.MerchantApiCredential, error) {
	query := `SELECT id, merchant_id, client_key, secret_key_hash, whitelist_ips, is_active, created_at, created_by, updated_at, updated_by FROM merchant_api_credentials WHERE id = $1`
	var mac models.MerchantApiCredential
	err := exec.QueryRow(query, id).Scan(&mac.ID, &mac.MerchantID, &mac.ClientKey, &mac.SecretKeyHash, &mac.WhitelistIPs, &mac.IsActive, &mac.CreatedAt, &mac.CreatedBy, &mac.UpdatedAt, &mac.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &mac, nil
}

func GetMerchantApiCredentialByMerchantID(exec QueryExecutor, merchantID int64) (*models.MerchantApiCredential, error) {
	query := `SELECT id, merchant_id, client_key, secret_key_hash, whitelist_ips, is_active, created_at, created_by, updated_at, updated_by FROM merchant_api_credentials WHERE merchant_id = $1`
	var mac models.MerchantApiCredential
	err := exec.QueryRow(query, merchantID).Scan(&mac.ID, &mac.MerchantID, &mac.ClientKey, &mac.SecretKeyHash, &mac.WhitelistIPs, &mac.IsActive, &mac.CreatedAt, &mac.CreatedBy, &mac.UpdatedAt, &mac.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &mac, nil
}

func UpdateMerchantApiCredential(exec QueryExecutor, mac *models.MerchantApiCredential) error {
	query := `UPDATE merchant_api_credentials SET client_key = ?, secret_key_hash = ?, whitelist_ips = ?, is_active = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, mac.ClientKey, mac.SecretKeyHash, mac.WhitelistIPs, mac.IsActive, mac.UpdatedAt, mac.UpdatedBy, mac.ID)
	return err
}

func DeleteMerchantApiCredential(exec QueryExecutor, id int64) error {
	query := `DELETE FROM merchant_api_credentials WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetMerchantApiCredentialsList(exec QueryExecutor, start, length int) ([]models.MerchantApiCredential, int64, error) {
	var count int64
	err := exec.QueryRow(`SELECT COUNT(*) FROM merchant_api_credentials`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, merchant_id, client_key, whitelist_ips, is_active, created_at, created_by, updated_at, updated_by FROM merchant_api_credentials ORDER BY id DESC`
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.MerchantApiCredential
	for rows.Next() {
		var mac models.MerchantApiCredential
		err = rows.Scan(&mac.ID, &mac.MerchantID, &mac.ClientKey, &mac.WhitelistIPs, &mac.IsActive, &mac.CreatedAt, &mac.CreatedBy, &mac.UpdatedAt, &mac.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, mac)
	}
	return list, count, nil
}
