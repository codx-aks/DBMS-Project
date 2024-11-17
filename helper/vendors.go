package helper

import (
	"context"
	"time"
	"wallet-system/models"
	"errors"
	"github.com/jackc/pgx/v5"
)

func GetTransactionsByVendor(ctx context.Context, tx pgx.Tx, vendorID int) ([]models.Transaction, int, error) {
	rows, err := tx.Query(context.Background(), "SELECT * FROM transactions WHERE receiver = $1", vendorID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var total int
	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID, &t.TransactionID, &t.Amount, &t.CreatedAt,
			&t.Sender, &t.Receiver, &t.Description,
		); err != nil {
			return nil, nil, err
		}
		total += t.Amount
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return transactions, total, nil
}

func GetAllVendors(ctx context.Context, tx pgx.Tx) ([]models.RespVendor, error) {
	rows, err := tx.Query(context.Background(), "SELECT id,name,description,image_url FROM vendors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vendors []models.RespVendor
	for rows.Next() {
		var v models.RespVendor
		if err := rows.Scan(
			&v.ID, &v.Name, &v.Description, &v.ImageURL
		); err != nil {
			return nil, err
		}
		vendors = append(vendors, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vendors, nil
}

