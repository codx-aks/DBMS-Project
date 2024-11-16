package models

import (
	"context"
	"wallet-system/config"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TransferRequest struct {
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Amount int       `json:"amount"`
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

type Balance struct {
	ID      uuid.UUID `json:"id"`
	Balance int       `json:"balance"`
}

func InitTable(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS accounts (id UUID PRIMARY KEY, balance INT8)")
	return err
}

func InsertRows(ctx context.Context, tx pgx.Tx, accounts [4]uuid.UUID) error {
	_, err := tx.Exec(ctx,
		"INSERT INTO accounts (id, balance) VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8)",
		accounts[0], 250, accounts[1], 100, accounts[2], 500, accounts[3], 300,
	)
	return err
}

func PrintBalances(ctx context.Context, tx pgx.Tx) ([]Balance, error) {
	rows, err := tx.Query(ctx, "SELECT id, balance FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []Balance
	for rows.Next() {
		var bal Balance
		if err := rows.Scan(&bal.ID, &bal.Balance); err != nil {
			return nil, err
		}
		balances = append(balances, bal)
	}
	return balances, nil
}

func TransferFunds(ctx context.Context, from uuid.UUID, to uuid.UUID, amount int) error {
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

func DeleteRows(ctx context.Context, id uuid.UUID) error {
	_, err := config.Conn.Exec(ctx, "DELETE FROM accounts WHERE id=$1", id)
	return err
}
