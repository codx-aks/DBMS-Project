package helper

import (
	"context"
	"time"
	"wallet-system/models"
	"errors"
	"github.com/jackc/pgx/v5"
)

func GetVendorItems(ctx context.Context, tx pgx.Tx, vendorID int) ([]models.Item, error) {
	rows, err := tx.Query(context.Background(), "SELECT *  FROM items WHERE vendor_id = $1",vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(
			&i.ID, &i.Name, &i.Description, &i.ImageURL,&i.Cost, &i.VendorID
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}