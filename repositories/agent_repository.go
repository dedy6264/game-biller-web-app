package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
)

// === AGENT REPOSITORY ===

func CreateAgent(exec QueryExecutor, a *models.Agent) (int64, error) {
	query := `INSERT INTO agents (agent_name, user_id, referral_code, status, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, a.AgentName, a.UserID, a.ReferralCode, a.Status, a.CreatedAt, a.CreatedBy, a.UpdatedAt, a.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	a.ID = id
	return id, nil
}

func GetAgentByID(exec QueryExecutor, id int64) (*models.Agent, error) {
	query := `SELECT a.id, a.agent_name, a.user_id, COALESCE(a.referral_code, ''), a.status, a.created_at, a.created_by, a.updated_at, a.updated_by,
	                 COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(u.phone_number, '')
	          FROM agents a
	          LEFT JOIN users u ON u.id = a.user_id
	          WHERE a.id = $1`
	var a models.Agent
	err := exec.QueryRow(query, id).Scan(
		&a.ID, &a.AgentName, &a.UserID, &a.ReferralCode, &a.Status, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy,
		&a.UserName, &a.UserEmail, &a.UserPhone,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetAgentByUserID(exec QueryExecutor, userID int64) (*models.Agent, error) {
	query := `SELECT a.id, a.agent_name, a.user_id, COALESCE(a.referral_code, ''), a.status, a.created_at, a.created_by, a.updated_at, a.updated_by,
	                 COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(u.phone_number, '')
	          FROM agents a
	          LEFT JOIN users u ON u.id = a.user_id
	          WHERE a.user_id = $1`
	var a models.Agent
	err := exec.QueryRow(query, userID).Scan(
		&a.ID, &a.AgentName, &a.UserID, &a.ReferralCode, &a.Status, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy,
		&a.UserName, &a.UserEmail, &a.UserPhone,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func GetAgentByReferralCode(exec QueryExecutor, code string) (*models.Agent, error) {
	query := `SELECT a.id, a.agent_name, a.user_id, COALESCE(a.referral_code, ''), a.status, a.created_at, a.created_by, a.updated_at, a.updated_by,
	                 COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(u.phone_number, '')
	          FROM agents a
	          LEFT JOIN users u ON u.id = a.user_id
	          WHERE UPPER(a.referral_code) = UPPER($1) LIMIT 1`
	var a models.Agent
	err := exec.QueryRow(query, code).Scan(
		&a.ID, &a.AgentName, &a.UserID, &a.ReferralCode, &a.Status, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy,
		&a.UserName, &a.UserEmail, &a.UserPhone,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func UpdateAgent(exec QueryExecutor, a *models.Agent) error {
	query := `UPDATE agents SET agent_name = ?, referral_code = ?, status = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, a.AgentName, a.ReferralCode, a.Status, a.UpdatedAt, a.UpdatedBy, a.ID)
	return err
}

func DeleteAgent(exec QueryExecutor, id int64) error {
	query := `DELETE FROM agents WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetAgentsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.AgentFilters) ([]models.Agent, int64, error) {
	var (
		count int64
		whr   string
	)

	if filters.ID != 0 {
		whr += " AND a.id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.UserID != 0 {
		whr += " AND a.user_id = " + strconv.FormatInt(filters.UserID, 10)
	}
	if filters.AgentName != "" {
		whr += " AND a.agent_name ILIKE '%" + filters.AgentName + "%'"
	}
	if filters.ReferralCode != "" {
		whr += " AND a.referral_code ILIKE '%" + filters.ReferralCode + "%'"
	}
	if filters.Status != "" {
		whr += " AND a.status = '" + filters.Status + "'"
	}
	if search != "" {
		whr += " AND (a.agent_name ILIKE '%" + search + "%' OR a.referral_code ILIKE '%" + search + "%' OR u.name ILIKE '%" + search + "%' OR u.email ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM agents a LEFT JOIN users u ON u.id = a.user_id WHERE true` + whr
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "a.id"
	if order != "" {
		orderBy = order
	}
	sortOrder := "DESC"
	if sort != "" {
		sortOrder = sort
	}

	query := `SELECT a.id, a.agent_name, a.user_id, COALESCE(a.referral_code, ''), a.status, a.created_at, a.created_by, a.updated_at, a.updated_by,
	                 COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(u.phone_number, '')
	          FROM agents a
	          LEFT JOIN users u ON u.id = a.user_id
	          WHERE true` + whr + " ORDER BY " + orderBy + " " + sortOrder
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.Agent
	for rows.Next() {
		var a models.Agent
		err = rows.Scan(
			&a.ID, &a.AgentName, &a.UserID, &a.ReferralCode, &a.Status, &a.CreatedAt, &a.CreatedBy, &a.UpdatedAt, &a.UpdatedBy,
			&a.UserName, &a.UserEmail, &a.UserPhone,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, a)
	}
	return list, count, nil
}
