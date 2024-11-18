package helper

import (
	"context"
	"errors"
	"log"
	models "wallet-system/models"
	"github.com/jackc/pgx/v5"
)

func CreateUser(tx pgx.Tx, user models.User) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO users (roll_no, name, email, pin,is_verified, is_approved, wallet_balance) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.RollNo,user.Name, user.Email, user.Pin,user.IsVerified, user.IsApproved, user.WalletBalance,
	)
	return err
}

func GetUserByEmailAndPassword(tx pgx.Tx, email string, pin string) (models.User, error) {
	var user models.User

	err := tx.QueryRow(
		context.Background(),
		"SELECT roll_no, name, email, pin,is_verified, is_approved, wallet_balance FROM users WHERE email = $1",
		email,
	).Scan( &user.RollNo, &user.Name, &user.Email, &user.Pin, &user.IsVerified, &user.IsApproved, &user.WalletBalance)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	if user.Pin != pin {
		log.Printf("Password mismatch", user.Pin, pin)
		return models.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func GetTransactionsByUserID(ctx context.Context, tx pgx.Tx, rollNo string) ([]models.Transaction, error) {
	rows, err := tx.Query(context.Background(), "SELECT * FROM transactions WHERE sender = $1", rollNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID, &t.TransactionID, &t.Amount, &t.CreatedAt,
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

