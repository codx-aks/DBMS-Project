package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	ID              int       `json:"id" db:"id"`
	RequestID       string    `json:"request_id" db:"request_id"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	TransactionType int       `json:"transaction_type" db:"transaction_type"`
	Amount          float64   `json:"amount" db:"amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	FromUserID      int       `json:"from_user_id" db:"from_user_id"`
	ToVendorID      int       `json:"to_vendor_id" db:"to_vendor_id"`
}

func GetTransactions(db *pgx.Conn) ([]Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID, &t.RequestID, &t.TransactionID, &t.TransactionType, &t.Amount,
			&t.CreatedAt, &t.UpdatedAt, &t.FromUserID, &t.ToVendorID,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
