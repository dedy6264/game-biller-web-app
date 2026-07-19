package repositories

import (
	"database/sql"
	"fmt"
	"gamebiller/helpers"
	"gamebiller/models"
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

func GetSavingAccountsList(exec QueryExecutor, search string, start, length int, order, sort string) ([]models.SavingAccount, int64, error) {
	var count int64
	countQuery := `SELECT COUNT(*) FROM saving_accounts WHERE account_number ILIKE $1`
	likeSearch := "%" + search + "%"
	err := exec.QueryRow(countQuery, likeSearch).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, merchant_id, account_number, balance, account_pin_hash, status, created_at, created_by, updated_at, updated_by FROM saving_accounts WHERE account_number ILIKE ?`
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

	rows, err := exec.Query(query, likeSearch)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.SavingAccount
	for rows.Next() {
		var sa models.SavingAccount
		err = rows.Scan(&sa.ID, &sa.MerchantID, &sa.AccountNumber, &sa.Balance, &sa.AccountPinHash, &sa.Status, &sa.CreatedAt, &sa.CreatedBy, &sa.UpdatedAt, &sa.UpdatedBy)
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

func GetSavingTransactionsList(exec QueryExecutor, accountID int64, start, length int) ([]models.SavingTransaction, int64, error) {
	var count int64
	var countQuery string
	var err error
	if accountID > 0 {
		countQuery = `SELECT COUNT(*) FROM saving_transactions WHERE saving_account_id = $1`
		err = exec.QueryRow(countQuery, accountID).Scan(&count)
	} else {
		countQuery = `SELECT COUNT(*) FROM saving_transactions`
		err = exec.QueryRow(countQuery).Scan(&count)
	}
	if err != nil {
		return nil, 0, err
	}

	var query string
	var rows *sql.Rows
	if accountID > 0 {
		query = `SELECT id, saving_account_id, type_dc, amount, last_balance, reference_number, transaction_code, description, created_at, created_by, created_by_user FROM saving_transactions WHERE saving_account_id = ? ORDER BY id DESC`
		if length > 0 {
			query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
		}
		query = helpers.QuerySupport(query)
		rows, err = exec.Query(query, accountID)
	} else {
		query = `SELECT id, saving_account_id, type_dc, amount, last_balance, reference_number, transaction_code, description, created_at, created_by, created_by_user FROM saving_transactions ORDER BY id DESC`
		if length > 0 {
			query += fmt.Sprintf(" LIMIT %d OFFSET %d", length, start)
		}
		query = helpers.QuerySupport(query)
		rows, err = exec.Query(query)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []models.SavingTransaction
	for rows.Next() {
		var st models.SavingTransaction
		err = rows.Scan(&st.ID, &st.SavingAccountID, &st.TypeDC, &st.Amount, &st.LastBalance, &st.ReferenceNumber, &st.TransactionCode, &st.Description, &st.CreatedAt, &st.CreatedBy, &st.CreatedByUser)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, st)
	}
	return list, count, nil
}
