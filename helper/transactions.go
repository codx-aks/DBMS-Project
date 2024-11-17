package helper

import (
	"context"
	"time"
	"wallet-system/models"
	"errors"
	"github.com/jackc/pgx/v5"
)

const (
    maxRetries     = 5
    initialBackoff = 10 * time.Millisecond 
)

func GetTransactionsByVendor(ctx context.Context, tx pgx.Tx, vendorID int) ([]models.Transaction, error) {
	rows, err := tx.Query(context.Background(), "SELECT * FROM transactions WHERE receiver = $1", vendorID)
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

func GetTransactionsByRollNo(ctx context.Context, tx pgx.Tx, rollNo string) ([]models.Transaction, error) {
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

func TransactionApproval(ctx context.Context,userID string, pin string, amount int, vendorID int) (string, error) {
	tx, err := config.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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

	err = retryUpdateWallet(tx, user.RollNo, amount)
	if err != nil {
		return "", fmt.Errorf("failed to update user wallet: %w", err)
	}

	transactionID, err := insertTransactionWithRetry(tx, user.RollNo, vendorID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to insert transaction: %w", err)
	}

	if err != nil {
		return "", err
	}

	err = tx.Commit()
		if err != nil {
			return "", fmt.Errorf("failed to commit transaction: %w", err)
		}

	return "Transaction Approved", nil
}

func retryUpdateWallet(tx pgx.Tx, userID string, amount int) error {
	var retries int
	for retries < maxRetries {
		_, err := tx.Exec(
			context.Background(),
			"UPDATE users SET wallet_balance = wallet_balance - $1 WHERE roll_no = $2",
			amount, userID,
		)
		if err == nil {
			return nil
		}
		if isNetworkOrTransientError(err) {
            backoffDuration := initialBackoff * (1 << (retries - 1))
            time.Sleep(backoffDuration)
			retries++
            continue 
        }

        return "", fmt.Errorf("failed to log transaction after multiple attempts: %w", err)
	}

	return fmt.Errorf("failed to update wallet balance after %d retries: %w", maxRetries , err)
}

func insertTransactionWithRetry(tx pgx.Tx, senderID string, receiverID int, amount int) (string, error) {
    var transactionID string

    var retries int
	for retries < maxRetries {
        err := tx.QueryRow(
            context.Background(),
            "INSERT INTO transactions (sender, receiver, amount, description, created_at) "+
                "VALUES ($1, $2, $3, $4, NOW()) RETURNING transaction_id",
            senderID, receiverID, amount, "Payment to Vendor", 
        ).Scan(&transactionID)

        if err == nil {
            return transactionID, nil
        }

        if isNetworkOrTransientError(err) {
            backoffDuration := initialBackoff * (1 << (retries - 1)) 
            time.Sleep(backoffDuration)
			retries++
            continue 
        }
        return "", fmt.Errorf("failed to log transaction after multiple attempts: %w", err)
		
    }

    return "", errors.New("failed to insert transaction after  %d retries %w",maxRetries,err)
}

func isNetworkOrTransientError(err error) bool {
    if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
        return true
    }

    if pgxErr, ok := err.(*pgx.PgError); ok {
		//PostgreSQL deadlock or timeout error codes
        if pgxErr.Code == "40001" || pgxErr.Code == "57P03" { 
            return true
        }
    }

    return false
}
