package helper

import (
	"context"
	"time"
	"wallet-system/models"
	"errors"
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


const maxRetries = 5 

func TransactionApproval(tx pgx.Tx, userID string, pin string, amount int, vendorID int) (string, error) {
	var user models.User

	err := tx.QueryRow(
		context.Background(),
		"SELECT roll_no, name, email, pin, is_verified, is_approved, wallet_balance FROM users WHERE roll_no = $1 FOR UPDATE",
		userID,
	).Scan(
		&user.RollNo, &user.Name, &user.Email, &user.Pin, &user.IsVerified, &user.IsApproved, &user.WalletBalance,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", errors.New("error retrieving user data")
	}

	if !user.IsVerified {
		return "", errors.New("account not verified")
	}

	if !user.IsApproved {
		return "", errors.New("account not approved by admin")
	}

	if user.Pin != pin {
		return "", errors.New("invalid credentials")
	}

	if user.WalletBalance < amount {
		return "", errors.New("insufficient funds")
	}

	err = tx.BeginFunc(context.Background(), func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), "SAVEPOINT save_transaction")
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			context.Background(),
			"UPDATE users SET wallet_balance = $1 WHERE roll_no = $2",
			user.WalletBalance-amount, user.RollNo,
		)
		if err != nil {
			return err
		}

		transactionID, err := insertTransactionWithRetry(tx, user.RollNo, vendorID, amount)
		if err != nil {
			_, rollbackErr := tx.Exec(context.Background(), "ROLLBACK TO SAVEPOINT save_transaction")
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return "", nil
}

func insertTransactionWithRetry(tx pgx.Tx, senderID string, receiverID int, amount int) (string, error) {
	var transactionID string

	for i := 0; i < maxRetries; i++ {
		err := tx.QueryRow(
			context.Background(),
			"INSERT INTO transactions (sender, receiver, amount, description, created_at) "+
				"VALUES ($1, $2, $3, $4, NOW()) RETURNING transaction_id",
			senderID, receiverID, amount, "Payment to Vendor", 
		).Scan(&transactionID)

		if err == nil {
			return transactionID, nil
		}
 
	}

	return "", errors.New("failed to log transaction after multiple attempts")
}
