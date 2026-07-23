package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// === PROVIDERS REPOSITORY ===

func CreateProvider(exec QueryExecutor, p *models.Provider) (int64, error) {
	query := `INSERT INTO providers (provider_name, is_active, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, p.ProviderName, p.IsActive, p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	p.ID = id
	return id, nil
}

func GetProviderByID(exec QueryExecutor, id int64) (*models.Provider, error) {
	query := `SELECT id, provider_name, is_active, created_at, created_by, updated_at, updated_by FROM providers WHERE id = $1`
	var p models.Provider
	err := exec.QueryRow(query, id).Scan(&p.ID, &p.ProviderName, &p.IsActive, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProvider(exec QueryExecutor, p *models.Provider) error {
	query := `UPDATE providers SET provider_name = ?, is_active = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, p.ProviderName, p.IsActive, p.UpdatedAt, p.UpdatedBy, p.ID)
	return err
}

func DeleteProvider(exec QueryExecutor, id int64) error {
	query := `DELETE FROM providers WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProvidersList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.ProviderFilters) ([]models.Provider, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProviderName != "" {
		whr += " AND provider_name ILIKE '%" + filters.ProviderName + "%'"
	}
	if filters.IsActive {
		whr += " AND is_active = true"
	}
	if search != "" {
		whr += " AND provider_name ILIKE '%" + search + "%'"
	}

	countQuery := `SELECT COUNT(*) FROM providers WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, provider_name, is_active, created_at, created_by, updated_at, updated_by FROM providers WHERE true` + whr
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

	var list []models.Provider
	for rows.Next() {
		var p models.Provider
		err = rows.Scan(&p.ID, &p.ProviderName, &p.IsActive, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	return list, count, nil
}

// === PRODUCT TYPES REPOSITORY ===

func CreateProductType(exec QueryExecutor, pt *models.ProductType) (int64, error) {
	query := `INSERT INTO product_types (product_type_name) VALUES (?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pt.ProductTypeName).Scan(&id)
	if err != nil {
		return 0, err
	}
	pt.ID = id
	return id, nil
}

func GetProductTypeByID(exec QueryExecutor, id int64) (*models.ProductType, error) {
	query := `SELECT id, product_type_name FROM product_types WHERE id = $1`
	var pt models.ProductType
	err := exec.QueryRow(query, id).Scan(&pt.ID, &pt.ProductTypeName)
	if err != nil {
		return nil, err
	}
	return &pt, nil
}

func UpdateProductType(exec QueryExecutor, pt *models.ProductType) error {
	query := `UPDATE product_types SET product_type_name = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pt.ProductTypeName, pt.ID)
	return err
}

func DeleteProductType(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_types WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductTypesList(exec QueryExecutor, start, length int) ([]models.ProductType, int64, error) {
	var count int64
	err := exec.QueryRow(`SELECT COUNT(*) FROM product_types`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, product_type_name FROM product_types ORDER BY id DESC`
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ProductType
	for rows.Next() {
		var pt models.ProductType
		err = rows.Scan(&pt.ID, &pt.ProductTypeName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pt)
	}
	return list, count, nil
}

// === PRODUCT CATEGORIES REPOSITORY ===

func CreateProductCategory(exec QueryExecutor, pc *models.ProductCategory) (int64, error) {
	query := `INSERT INTO product_categories (product_category_name) VALUES (?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pc.ProductCategoryName).Scan(&id)
	if err != nil {
		return 0, err
	}
	pc.ID = id
	return id, nil
}

func GetProductCategoryByID(exec QueryExecutor, id int64) (*models.ProductCategory, error) {
	query := `SELECT id, product_category_name FROM product_categories WHERE id = $1`
	var pc models.ProductCategory
	err := exec.QueryRow(query, id).Scan(&pc.ID, &pc.ProductCategoryName)
	if err != nil {
		return nil, err
	}
	return &pc, nil
}

func UpdateProductCategory(exec QueryExecutor, pc *models.ProductCategory) error {
	query := `UPDATE product_categories SET product_category_name = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pc.ProductCategoryName, pc.ID)
	return err
}

func DeleteProductCategory(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_categories WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductCategoriesList(exec QueryExecutor, start, length int, filters models.ProductCategoryFilters) ([]models.ProductCategory, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProductCategoryName != "" {
		whr += " AND product_category_name = '" + filters.ProductCategoryName + "'"
	}

	countQuery := `SELECT COUNT(*) FROM product_categories WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, product_category_name FROM product_categories WHERE true` + whr + " ORDER BY id DESC"
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ProductCategory
	for rows.Next() {
		var pc models.ProductCategory
		err = rows.Scan(&pc.ID, &pc.ProductCategoryName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pc)
	}
	return list, count, nil
}

// === PRODUCT REFERENCES REPOSITORY ===

func CreateProductReference(exec QueryExecutor, pr *models.ProductReference) (int64, error) {
	query := `INSERT INTO product_references (product_reference_code, product_reference_name, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64

	err := exec.QueryRow(query, pr.ProductReferenceCode, pr.ProductReferenceName, pr.CreatedAt, pr.CreatedBy, pr.UpdatedAt, pr.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	pr.ID = id
	return id, nil
}

func GetProductReferenceByID(exec QueryExecutor, id int64) (*models.ProductReference, error) {
	query := `SELECT id, product_reference_code, product_reference_name, created_at, created_by, updated_at, updated_by FROM product_references WHERE id = $1`
	var pr models.ProductReference
	err := exec.QueryRow(query, id).Scan(&pr.ID, &pr.ProductReferenceCode, &pr.ProductReferenceName, &pr.CreatedAt, &pr.CreatedBy, &pr.UpdatedAt, &pr.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func GetProductReferenceByCode(exec QueryExecutor, code string) (*models.ProductReference, error) {
	query := `SELECT id, product_reference_code, product_reference_name, created_at, created_by, updated_at, updated_by FROM product_references WHERE product_reference_code = $1`
	var pr models.ProductReference
	err := exec.QueryRow(query, code).Scan(&pr.ID, &pr.ProductReferenceCode, &pr.ProductReferenceName, &pr.CreatedAt, &pr.CreatedBy, &pr.UpdatedAt, &pr.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func UpdateProductReference(exec QueryExecutor, pr *models.ProductReference) error {
	query := `UPDATE product_references SET product_reference_code = ?, product_reference_name = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pr.ProductReferenceCode, pr.ProductReferenceName, pr.UpdatedAt, pr.UpdatedBy, pr.ID)
	return err
}

func DeleteProductReference(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_references WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductReferencesList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.ProductReferenceFilters) ([]models.ProductReference, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProductReferenceCode != "" {
		whr += " AND product_reference_code = '" + filters.ProductReferenceCode + "'"
	}
	if filters.ProductReferenceName != "" {
		whr += " AND product_reference_name ILIKE '%" + filters.ProductReferenceName + "%'"
	}
	if search != "" {
		whr += " AND (product_reference_name ILIKE '%" + search + "%' OR product_reference_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM product_references WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, product_reference_code, product_reference_name, created_at, created_by, updated_at, updated_by FROM product_references WHERE true` + whr
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

	var list []models.ProductReference
	for rows.Next() {
		var pr models.ProductReference
		err = rows.Scan(&pr.ID, &pr.ProductReferenceCode, &pr.ProductReferenceName, &pr.CreatedAt, &pr.CreatedBy, &pr.UpdatedAt, &pr.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pr)
	}
	return list, count, nil
}

// === PRODUCT PREFIXES REPOSITORY ===

func CreateProductPrefix(exec QueryExecutor, pp *models.ProductPrefix) (int64, error) {
	query := `INSERT INTO product_prefixes (product_reference_id, prefix_number, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pp.ProductReferenceID, pp.PrefixNumber, pp.CreatedAt, pp.CreatedBy, pp.UpdatedAt, pp.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	pp.ID = id
	return id, nil
}

func GetProductPrefixByID(exec QueryExecutor, id int64) (*models.ProductPrefix, error) {
	query := `SELECT id, product_reference_id, prefix_number, created_at, created_by, updated_at, updated_by FROM product_prefixes WHERE id = $1`
	var pp models.ProductPrefix
	err := exec.QueryRow(query, id).Scan(&pp.ID, &pp.ProductReferenceID, &pp.PrefixNumber, &pp.CreatedAt, &pp.CreatedBy, &pp.UpdatedAt, &pp.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pp, nil
}

// MatchPhonePrefix finds the product reference for a phone number
func MatchPhonePrefix(exec QueryExecutor, phoneNumber string) (*models.ProductReference, error) {
	// Format to simple digits, check first 5 digits or first 4 digits
	// Wait, the seed prefixes are e.g. '08120', '08121', '08510', etc. (5 characters)
	if len(phoneNumber) < 5 {
		return nil, fmt.Errorf("phone number too short")
	}
	prefix5 := phoneNumber[:5]
	prefix4 := phoneNumber[:4]

	query := `SELECT pr.id, pr.product_reference_code, pr.product_reference_name, pr.created_at, pr.created_by, pr.updated_at, pr.updated_by 
	          FROM product_references pr 
	          JOIN product_prefixes pp ON pp.product_reference_id = pr.id 
	          WHERE pp.prefix_number = $1 OR pp.prefix_number = $2 
	          LIMIT 1`
	var pr models.ProductReference
	err := exec.QueryRow(query, prefix5, prefix4).Scan(&pr.ID, &pr.ProductReferenceCode, &pr.ProductReferenceName, &pr.CreatedAt, &pr.CreatedBy, &pr.UpdatedAt, &pr.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func UpdateProductPrefix(exec QueryExecutor, pp *models.ProductPrefix) error {
	query := `UPDATE product_prefixes SET product_reference_id = ?, prefix_number = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pp.ProductReferenceID, pp.PrefixNumber, pp.UpdatedAt, pp.UpdatedBy, pp.ID)
	return err
}

func DeleteProductPrefix(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_prefixes WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductPrefixesList(exec QueryExecutor, start, length int, filters models.ProductPrefixFilters) ([]models.ProductPrefix, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProductReferenceID != 0 {
		whr += " AND product_reference_id = " + strconv.FormatInt(filters.ProductReferenceID, 10)
	}
	if filters.PrefixNumber != "" {
		whr += " AND prefix_number = '" + filters.PrefixNumber + "'"
	}

	countQuery := `SELECT COUNT(*) FROM product_prefixes WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT pp.id, pp.product_reference_id, pp.prefix_number, pp.created_at, pp.created_by, pp.updated_at, pp.updated_by, COALESCE(pr.product_reference_name, '')
	          FROM product_prefixes pp
	          LEFT JOIN product_references pr ON pr.id = pp.product_reference_id
	          WHERE true` + whr + " ORDER BY pp.id DESC"
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ProductPrefix
	for rows.Next() {
		var pp models.ProductPrefix
		err = rows.Scan(&pp.ID, &pp.ProductReferenceID, &pp.PrefixNumber, &pp.CreatedAt, &pp.CreatedBy, &pp.UpdatedAt, &pp.UpdatedBy, &pp.ProductReferenceName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pp)
	}
	return list, count, nil
}

// === PRODUCTS REPOSITORY ===

func CreateProduct(exec QueryExecutor, p *models.Product) (int64, error) {
	query := `INSERT INTO products (product_reference_id, product_type_id, product_category_id, product_code, product_name, is_active, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, p.ProductReferenceID, p.ProductTypeID, p.ProductCategoryID, p.ProductCode, p.ProductName, p.IsActive, p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	p.ID = id
	return id, nil
}

func GetProductByID(exec QueryExecutor, id int64) (*models.Product, error) {
	query := `SELECT id, product_reference_id, product_type_id, product_category_id, product_code, product_name, is_active, created_at, created_by, updated_at, updated_by FROM products WHERE id = $1`
	var p models.Product
	err := exec.QueryRow(query, id).Scan(&p.ID, &p.ProductReferenceID, &p.ProductTypeID, &p.ProductCategoryID, &p.ProductCode, &p.ProductName, &p.IsActive, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetProductByCode(exec QueryExecutor, code string) (*models.Product, error) {
	query := `SELECT id, product_reference_id, product_type_id, product_category_id, product_code, product_name, is_active, created_at, created_by, updated_at, updated_by FROM products WHERE product_code = $1`
	var p models.Product
	err := exec.QueryRow(query, code).Scan(&p.ID, &p.ProductReferenceID, &p.ProductTypeID, &p.ProductCategoryID, &p.ProductCode, &p.ProductName, &p.IsActive, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProduct(exec QueryExecutor, p *models.Product) error {
	query := `UPDATE products SET product_reference_id = ?, product_type_id = ?, product_category_id = ?, product_code = ?, product_name = ?, is_active = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, p.ProductReferenceID, p.ProductTypeID, p.ProductCategoryID, p.ProductCode, p.ProductName, p.IsActive, p.UpdatedAt, p.UpdatedBy, p.ID)
	return err
}

func DeleteProduct(exec QueryExecutor, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.ProductFilters) ([]models.Product, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProductReferenceID != nil && *filters.ProductReferenceID != 0 {
		whr += " AND product_reference_id = " + strconv.FormatInt(*filters.ProductReferenceID, 10)
	}
	if filters.ProductTypeID != 0 {
		whr += " AND product_type_id = " + strconv.FormatInt(filters.ProductTypeID, 10)
	}
	if filters.ProductCategoryID != 0 {
		whr += " AND product_category_id = " + strconv.FormatInt(filters.ProductCategoryID, 10)
	}
	if filters.ProductCode != "" {
		whr += " AND product_code = '" + filters.ProductCode + "'"
	}
	if filters.ProductName != "" {
		whr += " AND product_name ILIKE '%" + filters.ProductName + "%'"
	}
	if search != "" {
		whr += " AND (product_name ILIKE '%" + search + "%' OR product_code ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM products WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT p.id, p.product_reference_id, p.product_type_id, p.product_category_id, p.product_code, p.product_name, p.is_active, p.created_at, p.created_by, p.updated_at, p.updated_by, COALESCE(pr.product_reference_name, ''), COALESCE(pc.product_category_name, '')
	          FROM products p
	          LEFT JOIN product_references pr ON pr.id = p.product_reference_id
	          LEFT JOIN product_categories pc ON pc.id = p.product_category_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "p.id"
		} else if order == "created_at" {
			order = "p.created_at"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY p.id DESC"
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

	var list []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.ID, &p.ProductReferenceID, &p.ProductTypeID, &p.ProductCategoryID, &p.ProductCode, &p.ProductName, &p.IsActive, &p.CreatedAt, &p.CreatedBy, &p.UpdatedAt, &p.UpdatedBy, &p.ProductReferenceName, &p.ProductCategoryName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}
	return list, count, nil
}

// === PRODUCT PROVIDERS REPOSITORY ===

func CreateProductProvider(exec QueryExecutor, pp *models.ProductProvider) (int64, error) {
	query := `INSERT INTO product_providers (provider_id, provider_product_code, provider_price, provider_admin_fee, provider_merchant_fee, provider_index, is_available, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, pp.ProviderID, pp.ProviderProductCode, pp.ProviderPrice, pp.ProviderAdminFee, pp.ProviderMerchantFee, pp.ProviderIndex, pp.IsAvailable, pp.CreatedAt, pp.CreatedBy, pp.UpdatedAt, pp.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	pp.ID = id
	return id, nil
}

func GetProductProviderByID(exec QueryExecutor, id int64) (*models.ProductProvider, error) {
	query := `SELECT id, provider_id, provider_product_code, provider_price, provider_admin_fee, provider_merchant_fee, provider_index, is_available, created_at, created_by, updated_at, updated_by FROM product_providers WHERE id = $1`
	var pp models.ProductProvider
	err := exec.QueryRow(query, id).Scan(&pp.ID, &pp.ProviderID, &pp.ProviderProductCode, &pp.ProviderPrice, &pp.ProviderAdminFee, &pp.ProviderMerchantFee, &pp.ProviderIndex, &pp.IsAvailable, &pp.CreatedAt, &pp.CreatedBy, &pp.UpdatedAt, &pp.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &pp, nil
}

func UpdateProductProvider(exec QueryExecutor, pp *models.ProductProvider) error {
	query := `UPDATE product_providers SET provider_id = ?, provider_product_code = ?, provider_price = ?, provider_admin_fee = ?, provider_merchant_fee = ?, provider_index = ?, is_available = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, pp.ProviderID, pp.ProviderProductCode, pp.ProviderPrice, pp.ProviderAdminFee, pp.ProviderMerchantFee, pp.ProviderIndex, pp.IsAvailable, pp.UpdatedAt, pp.UpdatedBy, pp.ID)
	return err
}

func DeleteProductProvider(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_providers WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductProvidersList(exec QueryExecutor, start, length int, filters models.ProductProviderFilters) ([]models.ProductProvider, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProviderID != 0 {
		whr += " AND provider_id = " + strconv.FormatInt(filters.ProviderID, 10)
	}
	if filters.ProviderProductCode != "" {
		whr += " AND provider_product_code = '" + filters.ProviderProductCode + "'"
	}

	countQuery := `SELECT COUNT(*) FROM product_providers WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT pp.id, pp.provider_id, pp.provider_product_code, pp.provider_price, pp.provider_admin_fee, pp.provider_merchant_fee, pp.provider_index, pp.is_available, pp.created_at, pp.created_by, pp.updated_at, pp.updated_by, COALESCE(pr.provider_name, '')
	          FROM product_providers pp
	          LEFT JOIN providers pr ON pr.id = pp.provider_id
	          WHERE true` + whr + " ORDER BY pp.id DESC"
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ProductProvider
	for rows.Next() {
		var pp models.ProductProvider
		err = rows.Scan(&pp.ID, &pp.ProviderID, &pp.ProviderProductCode, &pp.ProviderPrice, &pp.ProviderAdminFee, &pp.ProviderMerchantFee, &pp.ProviderIndex, &pp.IsAvailable, &pp.CreatedAt, &pp.CreatedBy, &pp.UpdatedAt, &pp.UpdatedBy, &pp.ProviderName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, pp)
	}
	return list, count, nil
}

// === PRODUCT SEGMENTS REPOSITORY ===

func CreateProductSegment(exec QueryExecutor, ps *models.ProductSegment) (int64, error) {
	query := `INSERT INTO product_segments (segment_id, product_provider_id, segment_name, product_provider_code, product_provider_name, product_id, product_price, admin_fee, merchant_fee, provider_product_price, provider_product_admin_fee, provider_product_merchant_fee, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, ps.SegmentID, ps.ProductProviderID, ps.SegmentName, ps.ProductProviderCode, ps.ProductProviderName, ps.ProductID, ps.ProductPrice, ps.AdminFee, ps.MerchantFee, ps.ProviderProductPrice, ps.ProviderProductAdminFee, ps.ProviderProductMerchantFee, ps.CreatedAt, ps.CreatedBy, ps.UpdatedAt, ps.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	ps.ID = id
	return id, nil
}

func GetProductSegmentByID(exec QueryExecutor, id int64) (*models.ProductSegment, error) {
	query := `SELECT id, segment_id, product_provider_id, segment_name, COALESCE(product_provider_code, ''), COALESCE(product_provider_name, ''), product_id, product_price, admin_fee, merchant_fee, provider_product_price, provider_product_admin_fee, provider_product_merchant_fee, created_at, created_by, updated_at, updated_by FROM product_segments WHERE id = $1`
	var ps models.ProductSegment
	err := exec.QueryRow(query, id).Scan(&ps.ID, &ps.SegmentID, &ps.ProductProviderID, &ps.SegmentName, &ps.ProductProviderCode, &ps.ProductProviderName, &ps.ProductID, &ps.ProductPrice, &ps.AdminFee, &ps.MerchantFee, &ps.ProviderProductPrice, &ps.ProviderProductAdminFee, &ps.ProviderProductMerchantFee, &ps.CreatedAt, &ps.CreatedBy, &ps.UpdatedAt, &ps.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func GetProductSegmentByProductAndSegment(exec QueryExecutor, productID int64, segmentName string) (*models.ProductSegment, error) {
	query := `SELECT id, segment_id, product_provider_id, segment_name, COALESCE(product_provider_code, ''), COALESCE(product_provider_name, ''), product_id, product_price, admin_fee, merchant_fee, provider_product_price, provider_product_admin_fee, provider_product_merchant_fee, created_at, created_by, updated_at, updated_by 
	          FROM product_segments 
	          WHERE product_id = $1 AND (segment_name = $2 OR segment_id = (SELECT id FROM segments WHERE segment_name = $2 LIMIT 1))
	          LIMIT 1`
	var ps models.ProductSegment
	err := exec.QueryRow(query, productID, segmentName).Scan(&ps.ID, &ps.SegmentID, &ps.ProductProviderID, &ps.SegmentName, &ps.ProductProviderCode, &ps.ProductProviderName, &ps.ProductID, &ps.ProductPrice, &ps.AdminFee, &ps.MerchantFee, &ps.ProviderProductPrice, &ps.ProviderProductAdminFee, &ps.ProviderProductMerchantFee, &ps.CreatedAt, &ps.CreatedBy, &ps.UpdatedAt, &ps.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func UpdateProductSegment(exec QueryExecutor, ps *models.ProductSegment) error {
	query := `UPDATE product_segments SET segment_id = ?, product_provider_id = ?, segment_name = ?, product_provider_code = ?, product_provider_name = ?, product_id = ?, product_price = ?, admin_fee = ?, merchant_fee = ?, provider_product_price = ?, provider_product_admin_fee = ?, provider_product_merchant_fee = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, ps.SegmentID, ps.ProductProviderID, ps.SegmentName, ps.ProductProviderCode, ps.ProductProviderName, ps.ProductID, ps.ProductPrice, ps.AdminFee, ps.MerchantFee, ps.ProviderProductPrice, ps.ProviderProductAdminFee, ps.ProviderProductMerchantFee, ps.UpdatedAt, ps.UpdatedBy, ps.ID)
	return err
}

func DeleteProductSegment(exec QueryExecutor, id int64) error {
	query := `DELETE FROM product_segments WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetProductSegmentsList(exec QueryExecutor, start, length int, filters models.ProductSegmentFilters) ([]models.ProductSegment, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.SegmentID != nil && *filters.SegmentID != 0 {
		whr += " AND segment_id = " + strconv.FormatInt(*filters.SegmentID, 10)
	}
	if filters.ProductProviderID != nil && *filters.ProductProviderID != 0 {
		whr += " AND product_provider_id = " + strconv.FormatInt(*filters.ProductProviderID, 10)
	}
	if filters.ProductID != 0 {
		whr += " AND product_id = " + strconv.FormatInt(filters.ProductID, 10)
	}
	if filters.SegmentName != "" {
		whr += " AND segment_name = '" + filters.SegmentName + "'"
	}

	countQuery := `SELECT COUNT(*) FROM product_segments WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT ps.id, ps.segment_id, ps.product_provider_id, ps.segment_name, COALESCE(ps.product_provider_code, ''), COALESCE(ps.product_provider_name, ''), ps.product_id, ps.product_price, ps.admin_fee, ps.merchant_fee, ps.provider_product_price, ps.provider_product_admin_fee, ps.provider_product_merchant_fee, ps.created_at, ps.created_by, ps.updated_at, ps.updated_by, COALESCE(p.product_name, ''), COALESCE(p.product_code, ''), COALESCE(pprov.provider_product_code, ''), COALESCE(pprov.provider_price, 0.00), COALESCE(pprov.provider_admin_fee, 0.00), COALESCE(pprov.provider_merchant_fee, 0.00)
	          FROM product_segments ps
	          LEFT JOIN products p ON p.id = ps.product_id
	          LEFT JOIN product_providers pprov ON pprov.id = ps.product_provider_id
	          WHERE true` + whr + " ORDER BY ps.id DESC"
	if length > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
	}
	query = helpers.QuerySupport(query)

	rows, err := exec.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.ProductSegment
	for rows.Next() {
		var ps models.ProductSegment
		err = rows.Scan(&ps.ID, &ps.SegmentID, &ps.ProductProviderID, &ps.SegmentName, &ps.ProductProviderCode, &ps.ProductProviderName, &ps.ProductID, &ps.ProductPrice, &ps.AdminFee, &ps.MerchantFee, &ps.ProviderProductPrice, &ps.ProviderProductAdminFee, &ps.ProviderProductMerchantFee, &ps.CreatedAt, &ps.CreatedBy, &ps.UpdatedAt, &ps.UpdatedBy, &ps.ProductName, &ps.ProductCode, &ps.ProviderProductCode, &ps.ProviderPrice, &ps.ProviderAdminFee, &ps.ProviderMerchantFee)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, ps)
	}
	return list, count, nil
}

// === HIGH LEVEL BUSINESS METHODS ===

type ReferenceCodeSegmentResult struct {
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	IsActive    bool    `json:"is_active"`
	SegmentName string  `json:"segment_name"`
	SellPrice   float64 `json:"sell_price"`
	AdminFee    float64 `json:"admin_fee"`
	MerchantFee float64 `json:"merchant_fee"`
}

// GetProductSegmentsByRefCode fetches segment pricing for all products under a specific reference code
func GetProductSegmentsByRefCode(exec QueryExecutor, refCode string) ([]ReferenceCodeSegmentResult, error) {
	query := `SELECT p.product_code, p.product_name, p.is_active, ps.segment_name, ps.product_price, ps.admin_fee, ps.merchant_fee 
	          FROM products p 
	          JOIN product_references pr ON p.product_reference_id = pr.id 
	          JOIN product_segments ps ON ps.product_id = p.id 
	          WHERE pr.product_reference_code = $1 AND p.is_active = true 
	          ORDER BY p.product_code ASC, ps.segment_name ASC`
	rows, err := exec.Query(query, refCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ReferenceCodeSegmentResult
	for rows.Next() {
		var r ReferenceCodeSegmentResult
		err = rows.Scan(&r.ProductCode, &r.ProductName, &r.IsActive, &r.SegmentName, &r.SellPrice, &r.AdminFee, &r.MerchantFee)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// GetProductSegmentsByRefCodeAndSegment fetches segment pricing filtered by a specific segment name.
// Used when the caller is an authenticated merchant — returns only products available in their segment.
func GetProductSegmentsByRefCodeAndSegment(exec QueryExecutor, refCode string, segmentName string) ([]ReferenceCodeSegmentResult, error) {
	query := `SELECT p.product_code, p.product_name, p.is_active, ps.segment_name, ps.product_price, ps.admin_fee, ps.merchant_fee 
	          FROM products p 
	          JOIN product_references pr ON p.product_reference_id = pr.id 
	          JOIN product_segments ps ON ps.product_id = p.id 
	          WHERE p.is_active = true 
	          `
	if refCode != "" {
		query += ` and pr.product_reference_code = '` + refCode + `'`
	}
	if segmentName != "" {
		query += ` AND ps.segment_name = '` + segmentName + `'`
	}
	query += ` ORDER BY p.product_code ASC`
	rows, err := exec.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ReferenceCodeSegmentResult
	for rows.Next() {
		var r ReferenceCodeSegmentResult
		err = rows.Scan(&r.ProductCode, &r.ProductName, &r.IsActive, &r.SegmentName, &r.SellPrice, &r.AdminFee, &r.MerchantFee)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

type ProductSegmentDetailResult struct {
	ProductSegmentID           int64
	SegmentID                  *int64
	ProductID                  int64
	ProductCode                string
	ProductName                string
	ProductReferenceID         *int64
	ProductTypeID              int64
	ProductCategoryID          int64
	ProductIsActive            bool
	ProductPrice               float64
	AdminFee                   float64
	MerchantFee                float64
	ProviderProductPrice       float64
	ProviderProductAdminFee    float64
	ProviderProductMerchantFee float64
	ProductProviderID          *int64
	ProductProviderCode        string
	ProductProviderName        string
	ProviderID                 int64
	ProviderProductCode        string
	ProviderPrice              float64
	ProviderAdminFee           float64
	ProviderMerchantFee        float64
	ProviderIsAvailable        bool
	SegmentName                string
	ProviderName               string
	ProductTypeName            string
}

// GetProductSegmentJoinProvider fetches product segment pricing joined with products and product_providers
// by segmentID and productCode.
func GetProductSegmentJoinProvider(exec QueryExecutor, segmentID int64, productCode string) (*ProductSegmentDetailResult, error) {
	query := `SELECT 
		ps.id, 
		ps.segment_id, 
		ps.product_id, 
		ps.product_provider_id, 
		ps.product_price, 
		ps.admin_fee, 
		ps.merchant_fee,
		ps.provider_product_price, 
		ps.provider_product_admin_fee, 
		ps.provider_product_merchant_fee,
		p.product_code, 
		p.product_name, 
		p.is_active, 
		p.product_reference_id, 
		p.product_type_id, 
		p.product_category_id,
		COALESCE(pp.provider_id, 0), 
		COALESCE(ps.product_provider_code, pp.provider_product_code, ''),
		COALESCE(ps.product_provider_name, ''),
		COALESCE(pp.provider_product_code, ''),
		COALESCE(pp.provider_price, 0), 
		COALESCE(pp.provider_admin_fee, 0), 
		COALESCE(pp.provider_merchant_fee, 0), 
		COALESCE(pp.is_available, true),
		COALESCE(psg.segment_name, ''),
		COALESCE(prov.provider_name, ''),
		COALESCE(pt.product_type_name, '')
	FROM product_segments ps
	JOIN products p ON ps.product_id = p.id
	LEFT JOIN product_providers pp ON ps.product_provider_id = pp.id
	LEFT JOIN segments psg ON psg.id = ps.segment_id
	LEFT JOIN providers prov ON prov.id = pp.provider_id
	LEFT JOIN product_types pt ON pt.id = p.product_type_id
	WHERE ps.segment_id = $1 AND p.product_code = $2
	LIMIT 1`
	var res ProductSegmentDetailResult
	err := exec.QueryRow(query, segmentID, productCode).Scan(
		&res.ProductSegmentID,
		&res.SegmentID,
		&res.ProductID,
		&res.ProductProviderID,
		&res.ProductPrice,
		&res.AdminFee,
		&res.MerchantFee,
		&res.ProviderProductPrice,
		&res.ProviderProductAdminFee,
		&res.ProviderProductMerchantFee,
		&res.ProductCode,
		&res.ProductName,
		&res.ProductIsActive,
		&res.ProductReferenceID,
		&res.ProductTypeID,
		&res.ProductCategoryID,
		&res.ProviderID,
		&res.ProductProviderCode,
		&res.ProductProviderName,
		&res.ProviderProductCode,
		&res.ProviderPrice,
		&res.ProviderAdminFee,
		&res.ProviderMerchantFee,
		&res.ProviderIsAvailable,
		&res.SegmentName,
		&res.ProviderName,
		&res.ProductTypeName,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}


