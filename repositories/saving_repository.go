package repositories

import (
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
	"strconv"
	"strings"
)

// === SAVING ACCOUNT REPOSITORY ===

func CreateSavingAccount(exec QueryExecutor, sa *models.SavingAccount) (int64, error) {
	query := `INSERT INTO saving_accounts (merchant_id, account_number, balance, account_pin_hash, status, created_at, created_by, updated_at, updated_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, sa.MerchantID, sa.AccountNumber, sa.Balance, sa.AccountPinHash, sa.Status, sa.CreatedAt, sa.CreatedBy, sa.UpdatedAt, sa.UpdatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}
	sa.ID = id
	return id, nil
}

func GetSavingAccountByID(exec QueryExecutor, id int64) (*models.SavingAccount, error) {
	query := `SELECT id, merchant_id, account_number, balance, account_pin_hash, status, created_at, created_by, updated_at, updated_by FROM saving_accounts WHERE id = $1`
	var sa models.SavingAccount
	err := exec.QueryRow(query, id).Scan(&sa.ID, &sa.MerchantID, &sa.AccountNumber, &sa.Balance, &sa.AccountPinHash, &sa.Status, &sa.CreatedAt, &sa.CreatedBy, &sa.UpdatedAt, &sa.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

func GetSavingAccountByMerchantID(exec QueryExecutor, merchantID int64) (*models.SavingAccount, error) {
	query := `SELECT id, merchant_id, account_number, balance, account_pin_hash, status, created_at, created_by, updated_at, updated_by FROM saving_accounts WHERE merchant_id = $1`
	var sa models.SavingAccount
	err := exec.QueryRow(query, merchantID).Scan(&sa.ID, &sa.MerchantID, &sa.AccountNumber, &sa.Balance, &sa.AccountPinHash, &sa.Status, &sa.CreatedAt, &sa.CreatedBy, &sa.UpdatedAt, &sa.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

func UpdateSavingAccount(exec QueryExecutor, sa *models.SavingAccount) error {
	query := `UPDATE saving_accounts SET balance = ?, account_pin_hash = ?, status = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, sa.Balance, sa.AccountPinHash, sa.Status, sa.UpdatedAt, sa.UpdatedBy, sa.ID)
	return err
}

func UpdateSavingAccountBalance(exec QueryExecutor, id int64, newBalance float64, updatedBy, updatedAt string) error {
	query := `UPDATE saving_accounts SET balance = ?, updated_by = ?, updated_at = ? WHERE id = ?`
	query = helpers.QuerySupport(query)
	_, err := exec.Exec(query, newBalance, updatedBy, updatedAt, id)
	return err
}

func DeleteSavingAccount(exec QueryExecutor, id int64) error {
	query := `DELETE FROM saving_accounts WHERE id = $1`
	_, err := exec.Exec(query, id)
	return err
}

func GetSavingAccountsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.SavingAccountFilters) ([]models.SavingAccount, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND sa.id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.MerchantID != 0 {
		whr += " AND sa.merchant_id = " + strconv.FormatInt(filters.MerchantID, 10)
	}
	if filters.AccountNumber != "" {
		whr += " AND sa.account_number = '" + filters.AccountNumber + "'"
	}
	if filters.Status != "" {
		whr += " AND sa.status = '" + filters.Status + "'"
	}
	if search != "" {
		whr += " AND (sa.account_number ILIKE '%" + search + "%' OR m.merchant_name ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM saving_accounts sa LEFT JOIN merchants m ON m.id = sa.merchant_id WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT sa.id, sa.merchant_id, sa.account_number, sa.balance, sa.account_pin_hash, sa.status, sa.created_at, sa.created_by, sa.updated_at, sa.updated_by, COALESCE(m.merchant_name, '')
	          FROM saving_accounts sa
	          LEFT JOIN merchants m ON m.id = sa.merchant_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "sa.id"
		} else if order == "balance" {
			order = "sa.balance"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY sa.id DESC"
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

	var list []models.SavingAccount
	for rows.Next() {
		var sa models.SavingAccount
		err = rows.Scan(&sa.ID, &sa.MerchantID, &sa.AccountNumber, &sa.Balance, &sa.AccountPinHash, &sa.Status, &sa.CreatedAt, &sa.CreatedBy, &sa.UpdatedAt, &sa.UpdatedBy, &sa.MerchantName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, sa)
	}
	return list, count, nil
}

// === SAVING TRANSACTION REPOSITORY ===

func CreateSavingTransaction(exec QueryExecutor, st *models.SavingTransaction) (int64, error) {
	query := `INSERT INTO saving_transactions (saving_account_id, type_dc, amount, last_balance, reference_number, transaction_code, description, created_at, created_by, created_by_user)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`
	query = helpers.QuerySupport(query)
	var id int64
	err := exec.QueryRow(query, st.SavingAccountID, st.TypeDC, st.Amount, st.LastBalance, st.ReferenceNumber, st.TransactionCode, st.Description, st.CreatedAt, st.CreatedBy, st.CreatedByUser).Scan(&id)
	if err != nil {
		return 0, err
	}
	st.ID = id
	return id, nil
}

func GetSavingTransactionByID(exec QueryExecutor, id int64) (*models.SavingTransaction, error) {
	query := `SELECT id, saving_account_id, type_dc, amount, last_balance, reference_number, transaction_code, description, created_at, created_by, created_by_user FROM saving_transactions WHERE id = $1`
	var st models.SavingTransaction
	err := exec.QueryRow(query, id).Scan(&st.ID, &st.SavingAccountID, &st.TypeDC, &st.Amount, &st.LastBalance, &st.ReferenceNumber, &st.TransactionCode, &st.Description, &st.CreatedAt, &st.CreatedBy, &st.CreatedByUser)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func GetSavingTransactionsList(exec QueryExecutor, search string, start, length int, order, sort string, filters models.SavingTransactionFilters) ([]models.SavingTransaction, int64, error) {
	var (
		count int64
		whr   string
	)
	if filters.ID != 0 {
		whr += " AND st.id = " + strconv.FormatInt(filters.ID, 10)
	}
	if filters.SavingAccountID != 0 {
		whr += " AND st.saving_account_id = " + strconv.FormatInt(filters.SavingAccountID, 10)
	}
	if filters.MerchantID != 0 {
		whr += " AND sa.merchant_id = " + strconv.FormatInt(filters.MerchantID, 10)
	}
	if filters.TypeDC != "" {
		whr += " AND st.type_dc = '" + filters.TypeDC + "'"
	}
	if filters.TransactionCode != "" {
		whr += " AND st.transaction_code = '" + filters.TransactionCode + "'"
	}
	if filters.ReferenceNumber != "" {
		whr += " AND st.reference_number = '" + filters.ReferenceNumber + "'"
	}
	if search != "" {
		whr += " AND (st.reference_number ILIKE '%" + search + "%' OR sa.account_number ILIKE '%" + search + "%' OR m.merchant_name ILIKE '%" + search + "%')"
	}

	countQuery := `SELECT COUNT(*) FROM saving_transactions st LEFT JOIN saving_accounts sa ON sa.id = st.saving_account_id LEFT JOIN merchants m ON m.id = sa.merchant_id WHERE true` + whr
	countQuery = helpers.QuerySupport(countQuery)
	err := exec.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT st.id, st.saving_account_id, st.type_dc, st.amount, st.last_balance, st.reference_number, st.transaction_code, st.description, st.created_at, st.created_by, st.created_by_user, COALESCE(sa.account_number, ''), COALESCE(m.merchant_name, '')
	          FROM saving_transactions st
	          LEFT JOIN saving_accounts sa ON sa.id = st.saving_account_id
	          LEFT JOIN merchants m ON m.id = sa.merchant_id
	          WHERE true` + whr
	if order != "" {
		order = strings.ReplaceAll(order, ";", "")
		if order == "id" {
			order = "st.id"
		}
		if strings.ToLower(sort) != "desc" {
			sort = "ASC"
		} else {
			sort = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", order, sort)
	} else {
		query += " ORDER BY st.id DESC"
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

	var list []models.SavingTransaction
	for rows.Next() {
		var st models.SavingTransaction
		err = rows.Scan(&st.ID, &st.SavingAccountID, &st.TypeDC, &st.Amount, &st.LastBalance, &st.ReferenceNumber, &st.TransactionCode, &st.Description, &st.CreatedAt, &st.CreatedBy, &st.CreatedByUser, &st.AccountNumber, &st.MerchantName)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, st)
	}
	return list, count, nil
}
