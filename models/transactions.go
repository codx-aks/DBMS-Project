package models

import (
	"time"

)

type Transaction struct {
	ID              int       `json:"id" db:"id"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	Amount          int   `json:"amount" db:"amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	Sender      	string       `json:"sender" db:"sender"`
	Receiver      	int       `json:"receiver" db:"receiver"`
	Description     string       `json:"description" db:"description"`
}



