package models

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

type User struct {
	RollNo        string  `json:"roll_no" db:"roll_no"`
	Name          string  `json:"name" db:"name"`
	Email         string  `json:"email" db:"email"`
	Pin      string  `json:"pin" db:"pin"`
	IsApproved    bool    `json:"is_approved" db:"is_approved"`
	IsVerfied    bool    `json:"is_verified" db:"is_verified"`
	WalletBalance int `json:"wallet_balance" db:"wallet_balance"`
}




