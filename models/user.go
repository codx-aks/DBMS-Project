package models

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Email         string  `json:"email" db:"email"`
	Password      string  `json:"password" db:"password"`
	WalletPin     string  `json:"wallet_pin" db:"wallet_pin"`
	IsApproved    bool    `json:"is_approved" db:"is_approved"`
	WalletBalance float64 `json:"wallet_balance" db:"wallet_balance"`
}

func CreateUser(tx pgx.Tx, user User) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO users (name, email, password, wallet_pin, is_approved, wallet_balance) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Name, user.Email, user.Password, user.WalletPin, user.IsApproved, user.WalletBalance,
	)
	return err
}

func GetUserByEmailAndPassword(tx pgx.Tx, email, password string) (User, error) {
	var user User

	err := tx.QueryRow(
		context.Background(),
		"SELECT id, name, email, password, is_approved, wallet_balance FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsApproved, &user.WalletBalance)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	if user.Password != password {
		log.Printf("Password mismatch: Stored=%s, Input=%s", user.Password, password)
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}
