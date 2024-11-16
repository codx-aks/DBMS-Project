package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	ID              int       `json:"id" db:"id"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	Credit          int   `json:"credit" db:"credit"`
	Debit          int   `json:"debit" db:"debit"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	Sender      string       `json:"sender" db:"sender"`
	Receiver      int       `json:"receiver" db:"receiver"`
	Description      string       `json:"description" db:"description"`
}

func GetTransactionsByVendor(db *pgx.Conn, vendorID string) ([]Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM transactions WHERE receiver = $1", vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID, &t.TransactionID, &t.Credit, &t.Debit, &t.CreatedAt,
			&t.Sender, &t.Receiver, &t.Description,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTransactionsByRollNo(db *pgx.Conn, rollNo string) ([]Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM transactions WHERE sender = $1", rollNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID, &t.TransactionID, &t.Credit, &t.Debit, &t.CreatedAt,
			&t.Sender, &t.Receiver, &t.Description,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

