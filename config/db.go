package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDBConnection() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL is not set in the environment variables")
		return nil
	}

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return err
	}

	log.Println("Connected to CockroachDB.")
	return nil
}

func CloseDBConnection() {
	if Conn != nil {
		Conn.Close(context.Background())
	}
}
