package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
)

// === PRODUCT MASTER REPOSITORY ===

func CreateProductMaster(exec QueryExecutor, pm *models.ProductMaster) (int64, error) {
	pm.ProductSegmentIndex = strconv.Itoa(int(pm.ProductID)) + strconv.Itoa(int(pm.ProductProviderID))
	query := `INSERT INTO product_masters (
		provider_id, product_provider_id, product_name, product_segment_index, provider_name,
		product_provider_code, product_provider_name, product_id,
		product_price, admin_fee, merchant_fee,
		product_provider_price, product_provider_admin_fee, product_provider_merchant_fee,
		created_at, created_by, updated_at, updated_by
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query,
		pm.ProviderID, pm.ProductProviderID, pm.ProductName, pm.ProductSegmentIndex, pm.ProviderName,
		pm.ProductProviderCode, pm.ProductProviderName, pm.ProductID,
		pm.ProductPrice, pm.AdminFee, pm.MerchantFee,
		pm.ProductProviderPrice, pm.ProductProviderAdminFee, pm.ProductProviderMerchantFee,
		pm.CreatedAt, pm.CreatedBy, pm.UpdatedAt, pm.UpdatedBy,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	pm.ID = id
	return id, nil
}

func GetProductMasterByID(exec QueryExecutor, id int64) (*models.ProductMaster, error) {
	query := `SELECT
		pm.id, COALESCE(pm.provider_id, 0), pm.product_provider_id, pm.product_name, pm.product_segment_index,
		pm.provider_name, pm.product_provider_code, pm.product_provider_name,
		pm.product_id, pm.product_price, pm.admin_fee, pm.merchant_fee,
		pm.product_provider_price, pm.product_provider_admin_fee, pm.product_provider_merchant_fee,
		pm.created_at, pm.created_by, pm.updated_at, pm.updated_by,
		COALESCE(p.product_code, ''), COALESCE(p.product_name, '')
	FROM product_masters pm
	LEFT JOIN products p ON p.id = pm.product_id
	WHERE pm.id = $1`
	var pm models.ProductMaster
	err := exec.QueryRow(query, id).Scan(
		&pm.ID, &pm.ProviderID, &pm.ProductProviderID, &pm.ProductName, &pm.ProductSegmentIndex,
		&pm.ProviderName, &pm.ProductProviderCode, &pm.ProductProviderName,
		&pm.ProductID, &pm.ProductPrice, &pm.AdminFee, &pm.MerchantFee,
		&pm.ProductProviderPrice, &pm.ProductProviderAdminFee, &pm.ProductProviderMerchantFee,
		&pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy,
		&pm.ProductCode, &pm.ProductNameRef,
	)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func GetProductMasterByIndex(exec QueryExecutor, index string) (*models.ProductMaster, error) {
	query := `SELECT
		pm.id, COALESCE(pm.provider_id, 0), pm.product_provider_id, pm.product_name, pm.product_segment_index,
		pm.provider_name, pm.product_provider_code, pm.product_provider_name,
		pm.product_id, pm.product_price, pm.admin_fee, pm.merchant_fee,
		pm.product_provider_price, pm.product_provider_admin_fee, pm.product_provider_merchant_fee,
		pm.created_at, pm.created_by, pm.updated_at, pm.updated_by,
		COALESCE(p.product_code, ''), COALESCE(p.product_name, '')
	FROM product_masters pm
	LEFT JOIN products p ON p.id = pm.product_id
	WHERE pm.product_segment_index = $1`
	var pm models.ProductMaster
	err := exec.QueryRow(query, index).Scan(
		&pm.ID, &pm.ProviderID, &pm.ProductProviderID, &pm.ProductName, &pm.ProductSegmentIndex,
		&pm.ProviderName, &pm.ProductProviderCode, &pm.ProductProviderName,
		&pm.ProductID, &pm.ProductPrice, &pm.AdminFee, &pm.MerchantFee,
		&pm.ProductProviderPrice, &pm.ProductProviderAdminFee, &pm.ProductProviderMerchantFee,
		&pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy,
		&pm.ProductCode, &pm.ProductNameRef,
	)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func UpdateProductMaster(exec QueryExecutor, pm *models.ProductMaster) error {
	query := `UPDATE product_masters SET
		provider_id = ?, product_provider_id = ?, product_name = ?, product_segment_index = ?,
		provider_name = ?, product_provider_code = ?, product_provider_name = ?,
		product_id = ?, product_price = ?, admin_fee = ?, merchant_fee = ?,
		product_provider_price = ?, product_provider_admin_fee = ?, product_provider_merchant_fee = ?,
		updated_at = ?, updated_by = ?
	WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query,
		pm.ProviderID, pm.ProductProviderID, pm.ProductName, pm.ProductSegmentIndex,
		pm.ProviderName, pm.ProductProviderCode, pm.ProductProviderName,
		pm.ProductID, pm.ProductPrice, pm.AdminFee, pm.MerchantFee,
		pm.ProductProviderPrice, pm.ProductProviderAdminFee, pm.ProductProviderMerchantFee,
		pm.UpdatedAt, pm.UpdatedBy, pm.ID,
	)
	return err
}

func DeleteProductMaster(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_masters WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductMastersList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.ProductMasterFilters) ([]models.ProductMaster, int64, error) {
	var (
		count int64
		whr   string
	)

	if filters.ID != 0 {
		whr += " AND pm.id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProviderID != 0 {
		whr += " AND pm.provider_id = " + strconv.FormatInt(filters.ProviderID, 10)
	}
	if filters.ProductID != 0 {
		whr += " AND pm.product_id = " + strconv.FormatInt(filters.ProductID, 10)
	}
	if filters.ProductProviderID != 0 {
		whr += " AND pm.product_provider_id = " + strconv.FormatInt(filters.ProductProviderID, 10)
	}
	if filters.ProductSegmentIndex != "" {
		whr += " AND pm.product_segment_index ILIKE '%" + filters.ProductSegmentIndex + "%'"
	}
	if filters.ProductName != "" {
		whr += " AND pm.product_name ILIKE '%" + filters.ProductName + "%'"
	}
	if search != "" {
		whr += " AND (pm.product_name ILIKE '%" + search + "%' OR pm.product_segment_index ILIKE '%" + search + "%' OR pm.product_provider_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM product_masters pm LEFT JOIN products p ON p.id = pm.product_id WHERE true` + whr
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "pm.id"
	if order != "" {
		orderBy = order
	}
	sortOrder := "DESC"
	if sort != "" {
		sortOrder = sort
	}

	query := `SELECT
		pm.id, COALESCE(pm.provider_id, 0), pm.product_provider_id, pm.product_name, pm.product_segment_index,
		pm.provider_name, pm.product_provider_code, pm.product_provider_name,
		pm.product_id, pm.product_price, pm.admin_fee, pm.merchant_fee,
		pm.product_provider_price, pm.product_provider_admin_fee, pm.product_provider_merchant_fee,
		pm.created_at, pm.created_by, pm.updated_at, pm.updated_by,
		COALESCE(p.product_code, ''), COALESCE(p.product_name, '')
	FROM product_masters pm
	LEFT JOIN products p ON p.id = pm.product_id
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

	var list []models.ProductMaster
	for rows.Next() {
		var pm models.ProductMaster
		err = rows.Scan(
			&pm.ID, &pm.ProviderID, &pm.ProductProviderID, &pm.ProductName, &pm.ProductSegmentIndex,
			&pm.ProviderName, &pm.ProductProviderCode, &pm.ProductProviderName,
			&pm.ProductID, &pm.ProductPrice, &pm.AdminFee, &pm.MerchantFee,
			&pm.ProductProviderPrice, &pm.ProductProviderAdminFee, &pm.ProductProviderMerchantFee,
			&pm.CreatedAt, &pm.CreatedBy, &pm.UpdatedAt, &pm.UpdatedBy,
			&pm.ProductCode, &pm.ProductNameRef,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pm)
	}
	return list, count, nil
}
