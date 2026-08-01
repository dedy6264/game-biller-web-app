package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// === TRANSACTIONS REPOSITORY ===

func CreateTransaction(exec QueryExecutor, t *models.Transaction) (int64, error) {
	query := `INSERT INTO transactions (merchant_id, product_id, product_segment_id, product_provider_id, provider_id, product_type_id, product_reference_id, payment_channel_id, 
	                                  product_code, merchant_name, product_name, product_segment_name, product_provider_code, product_provider_name, provider_name, product_type_name, payment_channel_name,
	                                  product_provider_price, product_price, product_admin_fee, product_merchant_fee, product_provider_admin_fee, product_provider_merchant_fee, payment_admin_fee, total_amount, 
	                                  customer_id, other_customer_id, customer_phone, reference_number_internal, reference_number_merchant, reference_number_provider, serial_number, 
	                                  status_code, status_message, retry_count, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, t.MerchantID, t.ProductID, t.ProductSegmentID, t.ProductProviderID, t.ProviderID, t.ProductTypeID, t.ProductReferenceID, t.PaymentChannelID,
		t.ProductCode, t.MerchantName, t.ProductName, t.ProductSegmentName, t.ProductProviderCode, t.ProductProviderName, t.ProviderName, t.ProductTypeName, t.PaymentChannelName,
		t.ProductProviderPrice, t.ProductPrice, t.ProductAdminFee, t.ProductMerchantFee, t.ProductProviderAdminFee, t.ProductProviderMerchantFee, t.PaymentAdminFee, t.TotalAmount,
		t.CustomerID, t.OtherCustomerID, t.CustomerPhone, t.ReferenceNumberInternal, t.ReferenceNumberMerchant, t.ReferenceNumberProvider, t.SerialNumber,
		t.StatusCode, t.StatusMessage, t.RetryCount, t.CreatedAt, t.CreatedBy, t.UpdatedAt, t.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	t.ID = id
	return id, nil
}

func GetTransactionByID(exec QueryExecutor, id int64) (*models.Transaction, error) {
	query := `SELECT id, merchant_id, product_id, product_segment_id, product_provider_id, provider_id, product_type_id, product_reference_id, payment_channel_id, 
	                 COALESCE(product_code, ''),
	                 COALESCE(merchant_name, ''), COALESCE(product_name, ''), COALESCE(product_segment_name, ''), COALESCE(product_provider_code, ''), COALESCE(product_provider_name, ''), COALESCE(provider_name, ''), COALESCE(product_type_name, ''), COALESCE(payment_channel_name, ''),
	                 product_provider_price, product_price, product_admin_fee, COALESCE(product_merchant_fee, 0), COALESCE(product_provider_admin_fee, 0), COALESCE(product_provider_merchant_fee, 0), payment_admin_fee, total_amount, 
	                 COALESCE(customer_id, ''), COALESCE(other_customer_id, ''), COALESCE(customer_phone, ''), reference_number_internal, reference_number_merchant, reference_number_provider, serial_number, 
	                 status_code, status_message, retry_count, created_at, created_by, updated_at, updated_by 
	          FROM transactions WHERE id = $1`
	var t models.Transaction
	err := exec.QueryRow(query, id).Scan(&t.ID, &t.MerchantID, &t.ProductID, &t.ProductSegmentID, &t.ProductProviderID, &t.ProviderID, &t.ProductTypeID, &t.ProductReferenceID, &t.PaymentChannelID,
		&t.ProductCode,
		&t.MerchantName, &t.ProductName, &t.ProductSegmentName, &t.ProductProviderCode, &t.ProductProviderName, &t.ProviderName, &t.ProductTypeName, &t.PaymentChannelName,
		&t.ProductProviderPrice, &t.ProductPrice, &t.ProductAdminFee, &t.ProductMerchantFee, &t.ProductProviderAdminFee, &t.ProductProviderMerchantFee, &t.PaymentAdminFee, &t.TotalAmount,
		&t.CustomerID, &t.OtherCustomerID, &t.CustomerPhone, &t.ReferenceNumberInternal, &t.ReferenceNumberMerchant, &t.ReferenceNumberProvider, &t.SerialNumber,
		&t.StatusCode, &t.StatusMessage, &t.RetryCount, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTransactionByRefInternal(exec QueryExecutor, refInternal string) (*models.Transaction, error) {
	query := `SELECT id, merchant_id, product_id, product_segment_id, product_provider_id, provider_id, product_type_id, product_reference_id, payment_channel_id, 
	                 COALESCE(product_code, ''), 
	                 COALESCE(merchant_name, ''), COALESCE(product_name, ''), COALESCE(product_segment_name, ''), COALESCE(product_provider_code, ''), COALESCE(product_provider_name, ''), COALESCE(provider_name, ''), COALESCE(product_type_name, ''), COALESCE(payment_channel_name, ''),
	                 product_provider_price, product_price, product_admin_fee, COALESCE(product_merchant_fee, 0), COALESCE(product_provider_admin_fee, 0), COALESCE(product_provider_merchant_fee, 0), payment_admin_fee, total_amount, 
	                 COALESCE(customer_id, ''), COALESCE(other_customer_id, ''), COALESCE(customer_phone, ''), reference_number_internal, reference_number_merchant, reference_number_provider, serial_number, 
	                 status_code, status_message, retry_count, created_at, created_by, updated_at, updated_by 
	          FROM transactions WHERE reference_number_internal = $1`
	var t models.Transaction
	err := exec.QueryRow(query, refInternal).Scan(&t.ID, &t.MerchantID, &t.ProductID, &t.ProductSegmentID, &t.ProductProviderID, &t.ProviderID, &t.ProductTypeID, &t.ProductReferenceID, &t.PaymentChannelID,
		&t.ProductCode,
		&t.MerchantName, &t.ProductName, &t.ProductSegmentName, &t.ProductProviderCode, &t.ProductProviderName, &t.ProviderName, &t.ProductTypeName, &t.PaymentChannelName,
		&t.ProductProviderPrice, &t.ProductPrice, &t.ProductAdminFee, &t.ProductMerchantFee, &t.ProductProviderAdminFee, &t.ProductProviderMerchantFee, &t.PaymentAdminFee, &t.TotalAmount,
		&t.CustomerID, &t.OtherCustomerID, &t.CustomerPhone, &t.ReferenceNumberInternal, &t.ReferenceNumberMerchant, &t.ReferenceNumberProvider, &t.SerialNumber,
		&t.StatusCode, &t.StatusMessage, &t.RetryCount, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func GetTransactionByRefProvider(exec QueryExecutor, refProvider string) (*models.Transaction, error) {
	query := `SELECT id, merchant_id, product_id, product_segment_id, product_provider_id, provider_id, product_type_id, product_reference_id, payment_channel_id, 
	                 COALESCE(product_code, ''),  
	                 COALESCE(merchant_name, ''), COALESCE(product_name, ''), COALESCE(product_segment_name, ''), COALESCE(product_provider_code, ''), COALESCE(product_provider_name, ''), COALESCE(provider_name, ''), COALESCE(product_type_name, ''), COALESCE(payment_channel_name, ''),
	                 product_provider_price, product_price, product_admin_fee, COALESCE(product_merchant_fee, 0), COALESCE(product_provider_admin_fee, 0), COALESCE(product_provider_merchant_fee, 0), payment_admin_fee, total_amount, 
	                 COALESCE(customer_id, ''), COALESCE(other_customer_id, ''), COALESCE(customer_phone, ''), reference_number_internal, reference_number_merchant, reference_number_provider, serial_number, 
	                 status_code, status_message, retry_count, created_at, created_by, updated_at, updated_by 
	          FROM transactions WHERE reference_number_provider = $1`
	var t models.Transaction
	err := exec.QueryRow(query, refProvider).Scan(&t.ID, &t.MerchantID, &t.ProductID, &t.ProductSegmentID, &t.ProductProviderID, &t.ProviderID, &t.ProductTypeID, &t.ProductReferenceID, &t.PaymentChannelID,
		&t.ProductCode,
		&t.MerchantName, &t.ProductName, &t.ProductSegmentName, &t.ProductProviderCode, &t.ProductProviderName, &t.ProviderName, &t.ProductTypeName, &t.PaymentChannelName,
		&t.ProductProviderPrice, &t.ProductPrice, &t.ProductAdminFee, &t.ProductMerchantFee, &t.ProductProviderAdminFee, &t.ProductProviderMerchantFee, &t.PaymentAdminFee, &t.TotalAmount,
		&t.CustomerID, &t.OtherCustomerID, &t.CustomerPhone, &t.ReferenceNumberInternal, &t.ReferenceNumberMerchant, &t.ReferenceNumberProvider, &t.SerialNumber,
		&t.StatusCode, &t.StatusMessage, &t.RetryCount, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTransaction(exec QueryExecutor, t *models.Transaction) error {
	query := `UPDATE transactions SET merchant_id = ?, product_id = ?, product_segment_id = ?, product_provider_id = ?, provider_id = ?, product_type_id = ?, product_reference_id = ?, payment_channel_id = ?, 
	                                 product_code = ?, merchant_name = ?, product_name = ?, product_segment_name = ?, product_provider_code = ?, product_provider_name = ?, provider_name = ?, product_type_name = ?, payment_channel_name = ?,
	                                 product_provider_price = ?, product_price = ?, product_admin_fee = ?, product_merchant_fee = ?, product_provider_admin_fee = ?, product_provider_merchant_fee = ?, payment_admin_fee = ?, total_amount = ?, 
	                                 customer_id = ?, other_customer_id = ?, customer_phone = ?, reference_number_merchant = ?, reference_number_provider = ?, serial_number = ?, 
	                                 status_code = ?, status_message = ?, retry_count = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, t.MerchantID, t.ProductID, t.ProductSegmentID, t.ProductProviderID, t.ProviderID, t.ProductTypeID, t.ProductReferenceID, t.PaymentChannelID,
		t.ProductCode, t.MerchantName, t.ProductName, t.ProductSegmentName, t.ProductProviderCode, t.ProductProviderName, t.ProviderName, t.ProductTypeName, t.PaymentChannelName,
		t.ProductProviderPrice, t.ProductPrice, t.ProductAdminFee, t.ProductMerchantFee, t.ProductProviderAdminFee, t.ProductProviderMerchantFee, t.PaymentAdminFee, t.TotalAmount,
		t.CustomerID, t.OtherCustomerID, t.CustomerPhone, t.ReferenceNumberMerchant, t.ReferenceNumberProvider, t.SerialNumber,
		t.StatusCode, t.StatusMessage, t.RetryCount, t.UpdatedAt, t.UpdatedBy, t.ID)
	return err
}

func DeleteTransaction(exec QueryExecutor, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetTransactionsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.TransactionFilters) ([]models.Transaction, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.MerchantID != 0 {
		whr += " AND merchant_id = " + strconv.FormatInt(filters.MerchantID, 10)
	}
	if filters.ProductID != nil && *filters.ProductID != 0 {
		whr += " AND product_id = " + strconv.FormatInt(*filters.ProductID, 10)
	}
	if filters.ProductSegmentID != nil && *filters.ProductSegmentID != 0 {
		whr += " AND product_segment_id = " + strconv.FormatInt(*filters.ProductSegmentID, 10)
	}
	if filters.ProductProviderID != nil && *filters.ProductProviderID != 0 {
		whr += " AND product_provider_id = " + strconv.FormatInt(*filters.ProductProviderID, 10)
	}
	if filters.ProductProviderCode != "" {
		whr += " AND product_provider_code = '" + filters.ProductProviderCode + "'"
	}
	if filters.ProviderID != nil && *filters.ProviderID != 0 {
		whr += " AND provider_id = " + strconv.FormatInt(*filters.ProviderID, 10)
	}
	if filters.ProductTypeID != nil && *filters.ProductTypeID != 0 {
		whr += " AND product_type_id = " + strconv.FormatInt(*filters.ProductTypeID, 10)
	}
	if filters.ProductReferenceID != nil && *filters.ProductReferenceID != 0 {
		whr += " AND product_reference_id = " + strconv.FormatInt(*filters.ProductReferenceID, 10)
	}
	if filters.PaymentChannelID != nil && *filters.PaymentChannelID != 0 {
		whr += " AND payment_channel_id = " + strconv.FormatInt(*filters.PaymentChannelID, 10)
	}
	if filters.ProductCode != "" {
		whr += " AND product_code = '" + filters.ProductCode + "'"
	}
	if filters.CustomerID != "" {
		whr += " AND customer_id = '" + filters.CustomerID + "'"
	}
	if filters.OtherCustomerID != "" {
		whr += " AND other_customer_id = '" + filters.OtherCustomerID + "'"
	}
	if filters.CustomerPhone != "" {
		whr += " AND customer_phone = '" + filters.CustomerPhone + "'"
	}
	if filters.StatusCode != "" {
		whr += " AND status_code = '" + filters.StatusCode + "'"
	}
	if filters.ReferenceNumberInternal != "" {
		whr += " AND reference_number_internal = '" + filters.ReferenceNumberInternal + "'"
	}
	if filters.ReferenceNumberMerchant != nil && *filters.ReferenceNumberMerchant != "" {
		whr += " AND reference_number_merchant = '" + *filters.ReferenceNumberMerchant + "'"
	}
	if filters.StartDate != "" {
		whr += " AND created_at >= '" + filters.StartDate + "'"
	}
	if filters.EndDate != "" {
		whr += " AND created_at <= '" + filters.EndDate + "'"
	}
	if search != "" {
		whr += " AND ( product_code ILIKE '%" + search + "%' OR product_provider_code ILIKE '%" + search + "%' OR reference_number_internal ILIKE '%" + search + "%' OR customer_id ILIKE '%" + search + "%' OR other_customer_id ILIKE '%" + search + "%' OR customer_phone ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM transactions WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT t.id, t.merchant_id, t.product_id, t.product_segment_id, t.product_provider_id, t.provider_id, t.product_type_id, t.product_reference_id, t.payment_channel_id, 
	                 COALESCE(t.product_code, ''),  
	                 COALESCE(t.merchant_name, m.merchant_name, ''), COALESCE(t.product_name, p.product_name, ''), COALESCE(t.product_segment_name, ps.segment_name, ''), COALESCE(t.product_provider_code, pprov.product_provider_code, ''), COALESCE(t.product_provider_name, pprov.product_provider_name, pprov.product_provider_code, ''), COALESCE(t.provider_name, prov.provider_name, ''), COALESCE(t.product_type_name, ''), COALESCE(t.payment_channel_name, pc.channel_name, ''),
	                 t.product_provider_price, t.product_price, t.product_admin_fee, COALESCE(t.product_merchant_fee, 0), COALESCE(t.product_provider_admin_fee, 0), COALESCE(t.product_provider_merchant_fee, 0), t.payment_admin_fee, t.total_amount, 
	                 COALESCE(t.customer_id, ''), COALESCE(t.other_customer_id, ''), COALESCE(t.customer_phone, ''), t.reference_number_internal, t.reference_number_merchant, t.reference_number_provider, t.serial_number, 
	                 t.status_code, t.status_message, t.retry_count, t.created_at, t.created_by, t.updated_at, t.updated_by
	          FROM transactions t
	          LEFT JOIN merchants m ON m.id = t.merchant_id
	          LEFT JOIN products p ON p.id = t.product_id
	          LEFT JOIN product_segments ps ON ps.id = t.product_segment_id
	          LEFT JOIN product_providers pprov ON pprov.id = t.product_provider_id
	          LEFT JOIN providers prov ON prov.id = COALESCE(t.provider_id, pprov.provider_id)
	          LEFT JOIN payment_channels pc ON pc.id = t.payment_channel_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "t.id"
		} else if order == "created_at" {
			order = "t.created_at"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY t.id DESC"
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

	var list []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err = rows.Scan(&t.ID, &t.MerchantID, &t.ProductID, &t.ProductSegmentID, &t.ProductProviderID, &t.ProviderID, &t.ProductTypeID, &t.ProductReferenceID, &t.PaymentChannelID,
			&t.ProductCode,
			&t.MerchantName, &t.ProductName, &t.ProductSegmentName, &t.ProductProviderCode, &t.ProductProviderName, &t.ProviderName, &t.ProductTypeName, &t.PaymentChannelName,
			&t.ProductProviderPrice, &t.ProductPrice, &t.ProductAdminFee, &t.ProductMerchantFee, &t.ProductProviderAdminFee, &t.ProductProviderMerchantFee, &t.PaymentAdminFee, &t.TotalAmount,
			&t.CustomerID, &t.OtherCustomerID, &t.CustomerPhone, &t.ReferenceNumberInternal, &t.ReferenceNumberMerchant, &t.ReferenceNumberProvider, &t.SerialNumber,
			&t.StatusCode, &t.StatusMessage, &t.RetryCount, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	return list, count, nil
}

func GetTransactionsListByMerchantID(exec QueryExecutor, merchantID int64, search string, start, length int, order, sort string, filters models.TransactionFilters) ([]models.Transaction, int64, error) {
	var (
		count int64
		whr   string
	)
	if merchantID != 0 {
		whr += " AND merchant_id = " + strconv.FormatInt(merchantID, 10)
	}
	if filters.ID != 0 {
		whr += " AND id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.ProductID != nil && *filters.ProductID != 0 {
		whr += " AND product_id = " + strconv.FormatInt(*filters.ProductID, 10)
	}
	if filters.ProductSegmentID != nil && *filters.ProductSegmentID != 0 {
		whr += " AND product_segment_id = " + strconv.FormatInt(*filters.ProductSegmentID, 10)
	}
	if filters.ProductProviderID != nil && *filters.ProductProviderID != 0 {
		whr += " AND product_provider_id = " + strconv.FormatInt(*filters.ProductProviderID, 10)
	}
	if filters.ProductProviderCode != "" {
		whr += " AND product_provider_code = '" + filters.ProductProviderCode + "'"
	}
	if filters.ProviderID != nil && *filters.ProviderID != 0 {
		whr += " AND provider_id = " + strconv.FormatInt(*filters.ProviderID, 10)
	}
	if filters.ProductTypeID != nil && *filters.ProductTypeID != 0 {
		whr += " AND product_type_id = " + strconv.FormatInt(*filters.ProductTypeID, 10)
	}
	if filters.ProductReferenceID != nil && *filters.ProductReferenceID != 0 {
		whr += " AND product_reference_id = " + strconv.FormatInt(*filters.ProductReferenceID, 10)
	}
	if filters.PaymentChannelID != nil && *filters.PaymentChannelID != 0 {
		whr += " AND payment_channel_id = " + strconv.FormatInt(*filters.PaymentChannelID, 10)
	}
	if filters.ProductCode != "" {
		whr += " AND (product_code = '" + filters.ProductCode + "' )"
	}
	if filters.CustomerID != "" {
		whr += " AND customer_id = '" + filters.CustomerID + "'"
	}
	if filters.OtherCustomerID != "" {
		whr += " AND other_customer_id = '" + filters.OtherCustomerID + "'"
	}
	if filters.CustomerPhone != "" {
		whr += " AND customer_phone = '" + filters.CustomerPhone + "'"
	}
	if filters.StatusCode != "" {
		whr += " AND status_code = '" + filters.StatusCode + "'"
	}
	if filters.ReferenceNumberInternal != "" {
		whr += " AND reference_number_internal = '" + filters.ReferenceNumberInternal + "'"
	}
	if filters.ReferenceNumberMerchant != nil && *filters.ReferenceNumberMerchant != "" {
		whr += " AND reference_number_merchant = '" + *filters.ReferenceNumberMerchant + "'"
	}
	if filters.StartDate != "" {
		whr += " AND created_at >= '" + filters.StartDate + "'"
	}
	if filters.EndDate != "" {
		whr += " AND created_at <= '" + filters.EndDate + "'"
	}
	if search != "" {
		whr += " AND ( product_code ILIKE '%" + search + "%' OR product_provider_code ILIKE '%" + search + "%' OR reference_number_internal ILIKE '%" + search + "%' OR customer_id ILIKE '%" + search + "%' OR other_customer_id ILIKE '%" + search + "%' OR customer_phone ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM transactions WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT t.id, t.merchant_id, t.product_id, t.product_segment_id, t.product_provider_id, t.provider_id, t.product_type_id, t.product_reference_id, t.payment_channel_id, 
	                 COALESCE(t.product_code, ''),  
	                 COALESCE(t.merchant_name, m.merchant_name, ''), COALESCE(t.product_name, p.product_name, ''), COALESCE(t.product_segment_name, ps.segment_name, ''), COALESCE(t.product_provider_code, pprov.product_provider_code, ''), COALESCE(t.product_provider_name, pprov.product_provider_name, pprov.product_provider_code, ''), COALESCE(t.provider_name, prov.provider_name, ''), COALESCE(t.product_type_name, ''), COALESCE(t.payment_channel_name, pc.channel_name, ''),
	                 t.product_provider_price, t.product_price, t.product_admin_fee, COALESCE(t.product_merchant_fee, 0), COALESCE(t.product_provider_admin_fee, 0), COALESCE(t.product_provider_merchant_fee, 0), t.payment_admin_fee, t.total_amount, 
	                 COALESCE(t.customer_id, ''), COALESCE(t.other_customer_id, ''), COALESCE(t.customer_phone, ''), t.reference_number_internal, t.reference_number_merchant, t.reference_number_provider, t.serial_number, 
	                 t.status_code, t.status_message, t.retry_count, t.created_at, t.created_by, t.updated_at, t.updated_by
	          FROM transactions t
	          LEFT JOIN merchants m ON m.id = t.merchant_id
	          LEFT JOIN products p ON p.id = t.product_id
	          LEFT JOIN product_segments ps ON ps.id = t.product_segment_id
	          LEFT JOIN product_providers pprov ON pprov.id = t.product_provider_id
	          LEFT JOIN providers prov ON prov.id = COALESCE(t.provider_id, pprov.provider_id)
	          LEFT JOIN payment_channels pc ON pc.id = t.payment_channel_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "t.id"
		} else if order == "created_at" {
			order = "t.created_at"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY t.id DESC"
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

	var list []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err = rows.Scan(&t.ID, &t.MerchantID, &t.ProductID, &t.ProductSegmentID, &t.ProductProviderID, &t.ProviderID, &t.ProductTypeID, &t.ProductReferenceID, &t.PaymentChannelID,
			&t.ProductCode,
			&t.MerchantName, &t.ProductName, &t.ProductSegmentName, &t.ProductProviderCode, &t.ProductProviderName, &t.ProviderName, &t.ProductTypeName, &t.PaymentChannelName,
			&t.ProductProviderPrice, &t.ProductPrice, &t.ProductAdminFee, &t.ProductMerchantFee, &t.ProductProviderAdminFee, &t.ProductProviderMerchantFee, &t.PaymentAdminFee, &t.TotalAmount,
			&t.CustomerID, &t.OtherCustomerID, &t.CustomerPhone, &t.ReferenceNumberInternal, &t.ReferenceNumberMerchant, &t.ReferenceNumberProvider, &t.SerialNumber,
			&t.StatusCode, &t.StatusMessage, &t.RetryCount, &t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	return list, count, nil
}

type PopularProductResult struct {
	ProductReferenceCode string `json:"product_reference_code"`
	ProductReferenceName string `json:"product_reference_name"`
	TransactionCount     int    `json:"transaction_count"`
}

// GetPopularProducts fetches references ordered by transaction counts
func GetPopularProducts(exec QueryExecutor) ([]PopularProductResult, error) {
	// Aggregate transactions grouped by reference code
	query := `SELECT pr.product_reference_code, pr.product_reference_name, COUNT(t.id) as trx_count 
	          FROM product_references pr 
	          JOIN products p ON p.product_reference_id = pr.id 
	          LEFT JOIN transactions t ON t.product_id = p.id AND t.status_code = 'SUC-INT-000' 
	          GROUP BY pr.id, pr.product_reference_code, pr.product_reference_name 
	          ORDER BY trx_count DESC, pr.product_reference_name ASC`
	rows, err := exec.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []PopularProductResult
	for rows.Next() {
		var r PopularProductResult
		err = rows.Scan(&r.ProductReferenceCode, &r.ProductReferenceName, &r.TransactionCount)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

// === TRANSACTION PAYLOAD LOGS REPOSITORY ===

func CreateTransactionPayloadLog(exec QueryExecutor, tpl *models.TransactionPayloadLog) (int64, error) {
	query := `INSERT INTO transaction_payload_logs (transaction_id, request_payload, response_payload, created_at, created_by)
	          VALUES (?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, tpl.TransactionID, tpl.RequestPayload, tpl.ResponsePayload, tpl.CreatedAt, tpl.CreatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	tpl.ID = id
	return id, nil
}

func GetTransactionPayloadLogByID(exec QueryExecutor, id int64) (*models.TransactionPayloadLog, error) {
	query := `SELECT id, transaction_id, request_payload, response_payload, created_at, created_by FROM transaction_payload_logs WHERE id = $1`
	var tpl models.TransactionPayloadLog
	err := exec.QueryRow(query, id).Scan(&tpl.ID, &tpl.TransactionID, &tpl.RequestPayload, &tpl.ResponsePayload, &tpl.CreatedAt, &tpl.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}
