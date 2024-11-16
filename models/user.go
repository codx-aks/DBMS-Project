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

func CreateUser(tx pgx.Tx, user User) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO users (roll_no, name, email, pin,is_verified, is_approved, wallet_balance) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.RollNo,user.Name, user.Email, user.Pin,user.IsVerfied, user.IsApproved, user.WalletBalance
	)
	return err
}

func GetUserByEmailAndPassword(tx pgx.Tx, email string, password string) (User, error) {
	var user User

	err := tx.QueryRow(
		context.Background(),
		"SELECT roll_no, name, email, pin,is_verified, is_approved, wallet_balance FROM users WHERE email = $1",
		email,
	).Scan( &user.RollNo, &user.Name, &user.Email, &user.Pin, &user.IsVerfied, &user.IsApproved, &user.WalletBalance)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	if user.Password != password {
		log.Printf("Password mismatch: Stored=%s, Input=%s", user.Password, password)
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}


