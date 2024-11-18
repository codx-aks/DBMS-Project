package helper

import (
	"context"
	models "wallet-system/models"
	"github.com/jackc/pgx/v5"
	"errors"
	"log"
)

func GetTransactionsByVendor(ctx context.Context, tx pgx.Tx, vendorID int) ([]models.Transaction, int, error) {
	rows, err := tx.Query(context.Background(), "SELECT * FROM transactions WHERE receiver = $1", vendorID)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}
		total += t.Amount
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func GetAllVendors(ctx context.Context, tx pgx.Tx) ([]models.RespVendor, error) {
	rows, err := tx.Query(context.Background(), "SELECT id,name,description,image_url FROM vendors WHERE is_active = $1",true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vendors []models.RespVendor
	for rows.Next() {
		var v models.RespVendor
		if err := rows.Scan(
			&v.ID, &v.Name, &v.Description, &v.ImageURL,
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

func GetVendorByIDAndPassword(tx pgx.Tx, id int, password string) (models.Vendor, error) {
	var user models.Vendor

	err := tx.QueryRow(
		context.Background(),
		"SELECT * FROM vendors WHERE id = $1",
		id,
	).Scan( &user.ID,&user.Name,&user.Description,&user.Password,&user.ImageURL,&user.IsActive)
	if err != nil {
		return models.Vendor{}, errors.New("user not found")
	}

	if user.Password != password {
		log.Printf("Password mismatch", user.Password, password)
		return models.Vendor{}, errors.New("invalid credentials")
	}

	return user, nil
}
