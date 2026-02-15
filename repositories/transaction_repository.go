package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/models"
	"strings"
	"unicode/utf8"
)


type TransactionRepository struct {
	db *sql.DB
}


func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}


func RemoveLastChar(str string) string {
	if len(str) == 0 {
		return str
	}
	_, size := utf8.DecodeLastRuneInString(str)
	return str[:len(str)-size]
}


func (repo *TransactionRepository) GetAll() ([]models.Transaction, error) {
	query := "SELECT created_at, total_amount FROM transactions"
	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := make([]models.Transaction, 0)

	for rows.Next() {
		var p models.Transaction
		err := rows.Scan(&p.ID, &p.TotalAmount, &p.CreatedAt)
		
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, p)
	}

	return transactions, nil
}


func (repo *TransactionRepository) GetByID(id int) (*models.Transaction, error) {
	query := "SELECT id, total_amount, created_at FROM transactions"

	var p models.Transaction
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.TotalAmount, &p.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}


func (repo *TransactionRepository) Delete(id int) error {
	query := "DELETE FROM transactions WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("transaksi tidak ditemukan")
	}

	return err
}


func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}
  
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int

	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	
	if err != nil {
		return nil, err
	}

	var sqlInsert = "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "

	var insertValues strings.Builder
	insertValues.WriteString("")

	for i := range details {
		details[i].TransactionID = transactionID

		var localValue = fmt.Sprintf("(%d, %d, %d, %d),", transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		insertValues.WriteString(localValue)

	}

	finalInsert := RemoveLastChar(insertValues.String())
	sqlInsert += finalInsert

	_, err = tx.Exec(sqlInsert)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction {
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}


