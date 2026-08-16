package repositories

import (
	"database/sql"
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// QueryExecutor matches both *sql.DB and *sql.Tx
type QueryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// === USER REPOSITORY ===

func CreateUser(exec QueryExecutor, u *models.User) (int64, error) {
	query := `INSERT INTO users (name, email, phone_number, password_hash, status, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, u.Name, u.Email, u.PhoneNumber, u.PasswordHash, u.Status, u.CreatedAt, u.CreatedBy, u.UpdatedAt, u.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	u.ID = id
	return id, nil
}

func GetUserByID(exec QueryExecutor, id int64) (*models.User, error) {
	query := `SELECT id, name, email, phone_number, password_hash, status, created_at, created_by, updated_at, updated_by FROM users WHERE id = $1`
	var u models.User
	err := exec.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.PhoneNumber, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmailOrPhone(exec QueryExecutor, identifier string) (*models.User, error) {
	query := `SELECT id, name, email, phone_number, password_hash, status, created_at, created_by, updated_at, updated_by FROM users WHERE email = $1 OR phone_number = $1`
	var u models.User
	err := exec.QueryRow(query, identifier).Scan(&u.ID, &u.Name, &u.Email, &u.PhoneNumber, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(exec QueryExecutor, u *models.User) error {
	query := `UPDATE users SET name = ?, email = ?, phone_number = ?, password_hash = ?, status = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, u.Name, u.Email, u.PhoneNumber, u.PasswordHash, u.Status, u.UpdatedAt, u.UpdatedBy, u.ID)
	return err
}

func UpdateUserPassword(exec QueryExecutor, userID int64, passwordHash, updatedAt, updatedBy string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, passwordHash, updatedAt, updatedBy, userID)
	return err
}

func DeleteUser(exec QueryExecutor, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetUsersList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.UserFilters) ([]models.User, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.Name != "" {
		whr += " AND name ILIKE '%" + filters.Name + "%'"
	}
	if filters.Email != "" {
		whr += " AND email ILIKE '%" + filters.Email + "%'"
	}
	if filters.PhoneNumber != "" {
		whr += " AND phone_number ILIKE '%" + filters.PhoneNumber + "%'"
	}
	if filters.Status != "" {
		whr += " AND status = '" + filters.Status + "'"
	}
	if search != "" {
		whr += " AND (name ILIKE '%" + search + "%' OR email ILIKE '%" + search + "%' OR phone_number ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM users WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, email, phone_number, password_hash, status, created_at, created_by, updated_at, updated_by FROM users WHERE true` + whr
	if order != "" {
		// sanitization simple
		order = strings.ReplaceAll(order, ";", "")
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY id DESC"
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

	var list []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.PhoneNumber, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.CreatedBy, &u.UpdatedAt, &u.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, u)
	}
	return list, count, nil
}

// === ROLE REPOSITORY ===

func CreateRole(exec QueryExecutor, r *models.Role) (int64, error) {
	query := `INSERT INTO roles (role_code, role_name, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, r.RoleCode, r.RoleName, r.CreatedAt, r.CreatedBy, r.UpdatedAt, r.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	r.ID = id
	return id, nil
}

func GetRoleByID(exec QueryExecutor, id int64) (*models.Role, error) {
	query := `SELECT id, role_code, role_name, created_at, created_by, updated_at, updated_by FROM roles WHERE id = $1`
	var r models.Role
	err := exec.QueryRow(query, id).Scan(&r.ID, &r.RoleCode, &r.RoleName, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRoleByCode(exec QueryExecutor, code string) (*models.Role, error) {
	query := `SELECT id, role_code, role_name, created_at, created_by, updated_at, updated_by FROM roles WHERE role_code = $1`
	var r models.Role
	err := exec.QueryRow(query, code).Scan(&r.ID, &r.RoleCode, &r.RoleName, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func UpdateRole(exec QueryExecutor, r *models.Role) error {
	query := `UPDATE roles SET role_code = ?, role_name = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, r.RoleCode, r.RoleName, r.UpdatedAt, r.UpdatedBy, r.ID)
	return err
}

func DeleteRole(exec QueryExecutor, id int64) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetRolesList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.RoleFilters) ([]models.Role, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.RoleCode != "" {
		whr += " AND role_code = '" + filters.RoleCode + "'"
	}
	if filters.RoleName != "" {
		whr += " AND role_name ILIKE '%" + filters.RoleName + "%'"
	}
	if search != "" {
		whr += " AND (role_name ILIKE '%" + search + "%' OR role_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM roles WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, role_code, role_name, created_at, created_by, updated_at, updated_by FROM roles WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY id DESC"
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

	var list []models.Role
	for rows.Next() {
		var r models.Role
		err = rows.Scan(&r.ID, &r.RoleCode, &r.RoleName, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, r)
	}
	return list, count, nil
}

// === MODEL HAS ROLES REPOSITORY ===

func CreateModelHasRole(exec QueryExecutor, mhr *models.ModelHasRole) (int64, error) {
	query := `INSERT INTO model_has_roles (user_id, role_id, actor_id, created_at, created_by) VALUES (?, ?, ?, ?, ?)
	          RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, mhr.UserID, mhr.RoleID, mhr.ActorID, mhr.CreatedAt, mhr.CreatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	mhr.ID = id
	return id, nil
}

func UpdateModelHasRoleActorID(exec QueryExecutor, mhrID int64, actorID int64) error {
	query := `UPDATE model_has_roles SET actor_id = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, actorID, mhrID)
	return err
}

func GetModelHasRoleByID(exec QueryExecutor, id int64) (*models.ModelHasRole, error) {
	query := `SELECT id, user_id, role_id, COALESCE(actor_id, 0), created_at, created_by FROM model_has_roles WHERE id = $1`
	var mhr models.ModelHasRole
	err := exec.QueryRow(query, id).Scan(&mhr.ID, &mhr.UserID, &mhr.RoleID, &mhr.ActorID, &mhr.CreatedAt, &mhr.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &mhr, nil
}
func GetModelHasRoleByUserID(exec QueryExecutor, id int64) (*models.ModelHasRole, error) {
	query := `SELECT id, user_id, role_id, COALESCE(actor_id, 0), created_at, created_by FROM model_has_roles WHERE user_id = $1`
	var mhr models.ModelHasRole
	err := exec.QueryRow(query, id).Scan(&mhr.ID, &mhr.UserID, &mhr.RoleID, &mhr.ActorID, &mhr.CreatedAt, &mhr.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &mhr, nil
}

func GetUserRole(exec QueryExecutor, userID int64) (*models.Role, error) {
	query := `SELECT r.id, r.role_code, r.role_name, r.created_at, r.created_by, r.updated_at, r.updated_by 
	          FROM roles r 
	          JOIN model_has_roles mhr ON mhr.role_id = r.id 
	          WHERE mhr.user_id = $1`
	var r models.Role
	err := exec.QueryRow(query, userID).Scan(&r.ID, &r.RoleCode, &r.RoleName, &r.CreatedAt, &r.CreatedBy, &r.UpdatedAt, &r.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func DeleteModelHasRole(exec QueryExecutor, id int64) error {
	query := `DELETE FROM model_has_roles WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetModelHasRolesList(exec QueryExecutor, start, length int, filters models.ModelHasRoleFilters) ([]models.ModelHasRole, int64, error) {
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
	if filters.RoleID != 0 {
		whr += " AND role_id = " + strconv.FormatInt(filters.RoleID, 10)
	}
	if filters.ActorID != 0 {
		whr += " AND actor_id = " + strconv.FormatInt(filters.ActorID, 10)
	}

	countQuery := `SELECT COUNT(*) FROM model_has_roles WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT mhr.id, mhr.user_id, mhr.role_id, COALESCE(mhr.actor_id, 0), mhr.created_at, mhr.created_by, COALESCE(u.name, ''), COALESCE(u.email, ''), COALESCE(r.role_name, ''), COALESCE(r.role_code, '')
	          FROM model_has_roles mhr
	          LEFT JOIN users u ON u.id = mhr.user_id
	          LEFT JOIN roles r ON r.id = mhr.role_id
	          WHERE true` + whr + " ORDER BY mhr.id DESC"
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ModelHasRole
	for rows.Next() {
		var mhr models.ModelHasRole
		err = rows.Scan(&mhr.ID, &mhr.UserID, &mhr.RoleID, &mhr.ActorID, &mhr.CreatedAt, &mhr.CreatedBy, &mhr.UserName, &mhr.UserEmail, &mhr.RoleName, &mhr.RoleCode)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, mhr)
	}
	return list, count, nil
}

// === OTP CODES REPOSITORY ===

func CreateOtpCode(exec QueryExecutor, otp *models.OtpCode) (int64, error) {
	query := `INSERT INTO otp_codes (user_id, identifier, otp_code, otp_type, expired_at, is_used, attempt_count, max_attempt, created_at, created_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, otp.UserID, otp.Identifier, otp.OtpCode, otp.OtpType, otp.ExpiredAt, otp.IsUsed, otp.AttemptCount, otp.MaxAttempt, otp.CreatedAt, otp.CreatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	otp.ID = id
	return id, nil
}

func GetOtpCodeByID(exec QueryExecutor, id int64) (*models.OtpCode, error) {
	query := `SELECT id, user_id, identifier, otp_code, otp_type, expired_at, is_used, attempt_count, max_attempt, created_at, created_by FROM otp_codes WHERE id = $1`
	var otp models.OtpCode
	err := exec.QueryRow(query, id).Scan(&otp.ID, &otp.UserID, &otp.Identifier, &otp.OtpCode, &otp.OtpType, &otp.ExpiredAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempt, &otp.CreatedAt, &otp.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func GetActiveOtp(exec QueryExecutor, userID int64, identifier, otpType string) (*models.OtpCode, error) {
	query := `SELECT id, user_id, identifier, otp_code, otp_type, expired_at, is_used, attempt_count, max_attempt, created_at, created_by 
	          FROM otp_codes 
	          WHERE user_id = $1 AND identifier = $2 AND otp_type = $3 AND is_used = false 
	          ORDER BY id DESC LIMIT 1`
	var otp models.OtpCode
	err := exec.QueryRow(query, userID, identifier, otpType).Scan(&otp.ID, &otp.UserID, &otp.Identifier, &otp.OtpCode, &otp.OtpType, &otp.ExpiredAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempt, &otp.CreatedAt, &otp.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func UpdateOtpCode(exec QueryExecutor, otp *models.OtpCode) error {
	query := `UPDATE otp_codes SET is_used = ?, attempt_count = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, otp.IsUsed, otp.AttemptCount, otp.ID)
	return err
}

func DeleteOtpCode(exec QueryExecutor, id int64) error {
	query := `DELETE FROM otp_codes WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetOtpCodesList(exec QueryExecutor, start, length int) ([]models.OtpCode, int64, error) {
	var count int64
	err := exec.QueryRow(`SELECT COUNT(*) FROM otp_codes`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, user_id, identifier, otp_code, otp_type, expired_at, is_used, attempt_count, max_attempt, created_at, created_by FROM otp_codes ORDER BY id DESC`
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.OtpCode
	for rows.Next() {
		var otp models.OtpCode
		err = rows.Scan(&otp.ID, &otp.UserID, &otp.Identifier, &otp.OtpCode, &otp.OtpType, &otp.ExpiredAt, &otp.IsUsed, &otp.AttemptCount, &otp.MaxAttempt, &otp.CreatedAt, &otp.CreatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, otp)
	}
	return list, count, nil
}
