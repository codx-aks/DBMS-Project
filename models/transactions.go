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



