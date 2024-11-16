package models

import (
	"context"
	"wallet-system/config"

	"github.com/google/uuid"
)

type TransferRequest struct {
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount int       `json:"amount"`
}

type DeleteRequest struct {
	ID1 uuid.UUID `json:"id1"`
	ID2 uuid.UUID `json:"id2"`
}

func InitTable(ctx context.Context) error {
	_, err := config.Conn.Exec(ctx, "DROP TABLE IF EXISTS accounts")
	if err != nil {
		return err
	}
	_, err = config.Conn.Exec(ctx, "CREATE TABLE accounts (id UUID PRIMARY KEY, balance INT8)")
	return err
}

func InsertRows(ctx context.Context, accounts [4]uuid.UUID) error {
	_, err := config.Conn.Exec(ctx, "INSERT INTO accounts (id, balance) VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8)",
		accounts[0], 250, accounts[1], 100, accounts[2], 500, accounts[3], 300)
	return err
}

func PrintBalances(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := config.Conn.Query(ctx, "SELECT id, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var balance int
		if err := rows.Scan(&id, &balance); err != nil {
			return nil, err
		}
		balances = append(balances, map[string]interface{}{
			"id":      id.String(),
			"balance": balance,
		})
	}
	return balances, nil
}

func TransferFunds(ctx context.Context, from, to uuid.UUID, amount int) error {
	tx, err := config.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var fromBalance int
	err = tx.QueryRow(ctx, "SELECT balance FROM accounts WHERE id = $1", from).Scan(&fromBalance)
	if err != nil || fromBalance < amount {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, from)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, to)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func DeleteRows(ctx context.Context, id1, id2 uuid.UUID) error {
	_, err := config.Conn.Exec(ctx, "DELETE FROM accounts WHERE id IN ($1, $2)", id1, id2)
	return err
}
