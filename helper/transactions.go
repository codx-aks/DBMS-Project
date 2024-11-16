package helper

import (
	"context"
	"time"
	"wallet-system/models"
	"github.com/jackc/pgx/v5"
)

func GetTransactionsByVendor(db *pgx.Conn, vendorID string) ([]models.Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM transactions WHERE receiver = $1", vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
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

func GetTransactionsByRollNo(db *pgx.Conn, rollNo string) ([]models.Transaction, error) {
	rows, err := db.Query(context.Background(), "SELECT * FROM transactions WHERE sender = $1", rollNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
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