package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// === PAYMENT METHODS REPOSITORY ===

func CreatePaymentMethod(exec QueryExecutor, pm *models.PaymentMethod) (int64, error) {
	query := `INSERT INTO payment_methods (method_code, method_name, is_active, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pm.MethodCode, pm.MethodName, pm.IsActive, pm.CreatedAt, pm.CreatedBy, pm.UpdatedAt, pm.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	pm.ID = id
	return id, nil
}

func GetPaymentMethodByID(exec QueryExecutor, id int64) (*models.PaymentMethod, error) {
	query := `SELECT id, method_code, method_name, is_active, created_at, created_by, updated_at, updated_by FROM payment_methods WHERE id = $1`
	var pm models.PaymentMethod
	err := exec.QueryRow(query, id).Scan(&pm.ID, &pm.MethodCode, &pm.MethodName, &pm.IsActive, &pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func GetPaymentMethodByCode(exec QueryExecutor, code string) (*models.PaymentMethod, error) {
	query := `SELECT id, method_code, method_name, is_active, created_at, created_by, updated_at, updated_by FROM payment_methods WHERE method_code = $1`
	var pm models.PaymentMethod
	err := exec.QueryRow(query, code).Scan(&pm.ID, &pm.MethodCode, &pm.MethodName, &pm.IsActive, &pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func UpdatePaymentMethod(exec QueryExecutor, pm *models.PaymentMethod) error {
	query := `UPDATE payment_methods SET method_code = ?, method_name = ?, is_active = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pm.MethodCode, pm.MethodName, pm.IsActive, pm.UpdatedAt, pm.UpdatedBy, pm.ID)
	return err
}

func DeletePaymentMethod(exec QueryExecutor, id int64) error {
	query := `DELETE FROM payment_methods WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetPaymentMethodsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.PaymentMethodFilters) ([]models.PaymentMethod, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.MethodCode != "" {
		whr += " AND method_code = '" + filters.MethodCode + "'"
	}
	if filters.MethodName != "" {
		whr += " AND method_name ILIKE '%" + filters.MethodName + "%'"
	}
	if search != "" {
		whr += " AND (method_name ILIKE '%" + search + "%' OR method_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM payment_methods WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, method_code, method_name, is_active, created_at, created_by, updated_at, updated_by FROM payment_methods WHERE true` + whr
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

	var list []models.PaymentMethod
	for rows.Next() {
		var pm models.PaymentMethod
		err = rows.Scan(&pm.ID, &pm.MethodCode, &pm.MethodName, &pm.IsActive, &pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pm)
	}
	return list, count, nil
}

// === PAYMENT CHANNELS REPOSITORY ===

func CreatePaymentChannel(exec QueryExecutor, pc *models.PaymentChannel) (int64, error) {
	query := `INSERT INTO payment_channels (payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pc.PaymentMethodID, pc.ChannelCode, pc.ChannelName, pc.FeeType, pc.FeeValue, pc.IsActive, pc.CreatedAt, pc.CreatedBy, pc.UpdatedAt, pc.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	pc.ID = id
	return id, nil
}

func GetPaymentChannelByID(exec QueryExecutor, id int64) (*models.PaymentChannel, error) {
	query := `SELECT id, payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, created_by, updated_at, updated_by FROM payment_channels WHERE id = $1`
	var pc models.PaymentChannel
	err := exec.QueryRow(query, id).Scan(&pc.ID, &pc.PaymentMethodID, &pc.ChannelCode, &pc.ChannelName, &pc.FeeType, &pc.FeeValue, &pc.IsActive, &pc.CreatedAt, &pc.CreatedBy, &pc.UpdatedAt, &pc.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func GetPaymentChannelByCode(exec QueryExecutor, code string) (*models.PaymentChannel, error) {
	return GetPaymentChannelByCodeSafe(exec, code)
}

func GetPaymentChannelByCodeSafe(exec QueryExecutor, code string) (*models.PaymentChannel, error) {
	query := `SELECT id, payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, created_by, updated_at, updated_by FROM payment_channels WHERE channel_code = $1`
	var pc models.PaymentChannel
	err := exec.QueryRow(query, code).Scan(&pc.ID, &pc.PaymentMethodID, &pc.ChannelCode, &pc.ChannelName, &pc.FeeType, &pc.FeeValue, &pc.IsActive, &pc.CreatedAt, &pc.CreatedBy, &pc.UpdatedAt, &pc.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func UpdatePaymentChannel(exec QueryExecutor, pc *models.PaymentChannel) error {
	query := `UPDATE payment_channels SET payment_method_id = ?, channel_code = ?, channel_name = ?, fee_type = ?, fee_value = ?, is_active = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pc.PaymentMethodID, pc.ChannelCode, pc.ChannelName, pc.FeeType, pc.FeeValue, pc.IsActive, pc.UpdatedAt, pc.UpdatedBy, pc.ID)
	return err
}

func DeletePaymentChannel(exec QueryExecutor, id int64) error {
	query := `DELETE FROM payment_channels WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetPaymentChannelsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.PaymentChannelFilters) ([]models.PaymentChannel, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.PaymentMethodID != 0 {
		whr += " AND payment_method_id = " + strconv.FormatInt(filters.PaymentMethodID, 10)
	}
	if filters.ChannelCode != "" {
		whr += " AND channel_code = '" + filters.ChannelCode + "'"
	}
	if filters.ChannelName != "" {
		whr += " AND channel_name ILIKE '%" + filters.ChannelName + "%'"
	}
	if search != "" {
		whr += " AND (channel_name ILIKE '%" + search + "%' OR channel_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM payment_channels WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT pc.id, pc.payment_method_id, pc.channel_code, pc.channel_name, pc.fee_type, pc.fee_value, pc.is_active, pc.created_at, pc.created_by, pc.updated_at, pc.updated_by, COALESCE(pm.method_name, ''), COALESCE(pm.method_code, '')
	          FROM payment_channels pc
	          LEFT JOIN payment_methods pm ON pm.id = pc.payment_method_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "pc.id"
		} else if order == "created_at" {
			order = "pc.created_at"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY pc.id DESC"
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

	var list []models.PaymentChannel
	for rows.Next() {
		var pc models.PaymentChannel
		err = rows.Scan(&pc.ID, &pc.PaymentMethodID, &pc.ChannelCode, &pc.ChannelName, &pc.FeeType, &pc.FeeValue, &pc.IsActive, &pc.CreatedAt, &pc.CreatedBy, &pc.UpdatedAt, &pc.UpdatedBy, &pc.MethodName, &pc.MethodCode)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pc)
	}
	return list, count, nil
}

// Helpers for scan
func ppcUpdatedBy(pc models.PaymentChannel) *string {
	return &pc.UpdatedBy
}

type PaymentMethodWithChannels struct {
	MethodCode string                  `json:"method_code"`
	MethodName string                  `json:"method_name"`
	IsActive   bool                    `json:"is_active"`
	Channels   []models.PaymentChannel `json:"channels"`
}

// GetPaymentMethodsWithChannels returns payment methods nested with active channels
func GetPaymentMethodsWithChannels(exec QueryExecutor) ([]PaymentMethodWithChannels, error) {
	methodsQuery := `SELECT id, method_code, method_name, is_active FROM payment_methods WHERE is_active = true`
	mRows, err := exec.Query(methodsQuery)
	if err != nil {
		return nil, err
	}
	defer mRows.Close()

	var list []PaymentMethodWithChannels
	var methodIDs []int64
	methodMap := make(map[int64]int)

	for mRows.Next() {
		var id int64
		var m PaymentMethodWithChannels
		err = mRows.Scan(&id, &m.MethodCode, &m.MethodName, &m.IsActive)
		if err != nil {
			return nil, err
		}
		m.Channels = []models.PaymentChannel{}
		methodMap[id] = len(list)
		list = append(list, m)
		methodIDs = append(methodIDs, id)
	}

	if len(methodIDs) == 0 {
		return list, nil
	}

	channelsQuery := `SELECT id, payment_method_id, channel_code, channel_name, fee_type, fee_value, is_active, created_at, created_by, updated_at, updated_by 
	                  FROM payment_channels WHERE is_active = true ORDER BY id ASC`
	cRows, err := exec.Query(channelsQuery)
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	for cRows.Next() {
		var pc models.PaymentChannel
		err = cRows.Scan(&pc.ID, &pc.PaymentMethodID, &pc.ChannelCode, &pc.ChannelName, &pc.FeeType, &pc.FeeValue, &pc.IsActive, &pc.CreatedAt, &pc.CreatedBy, &pc.UpdatedAt, &pc.UpdatedBy)
		if err != nil {
			return nil, err
		}
		if idx, exists := methodMap[pc.PaymentMethodID]; exists {
			list[idx].Channels = append(list[idx].Channels, pc)
		}
	}

	return list, nil
}
