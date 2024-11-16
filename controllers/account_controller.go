// package controllers

// import (
// 	"context"
// 	"net/http"

// 	"wallet-system/models"

// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"
// )

// func InitTableHandler(c echo.Context) error {
// 	err := models.InitTable(context.Background())
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "error initializing table")
// 	}
// 	return c.JSON(http.StatusOK, "Table initialized.")
// }

// func InsertRowsHandler(c echo.Context) error {
// 	accounts := [4]uuid.UUID{}
// 	for i := 0; i < len(accounts); i++ {
// 		accounts[i] = uuid.New()
// 	}

// 	err := models.InsertRows(context.Background(), accounts)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "error inserting rows")
// 	}
// 	return c.JSON(http.StatusOK, "Rows inserted.")
// }

// func PrintBalancesHandler(c echo.Context) error {
// 	balances, err := models.PrintBalances(context.Background())
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "error fetching balances")
// 	}
// 	return c.JSON(http.StatusOK, balances)
// }

// func TransferFundsHandler(c echo.Context) error {
// 	var req models.TransferRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, "invalid request payload")
// 	}

// 	err := models.TransferFunds(context.Background(), req.From, req.To, req.Amount)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "error transferring funds")
// 	}
// 	return c.JSON(http.StatusOK, "Funds transferred.")
// }

// func DeleteRowsHandler(c echo.Context) error {
// 	var req models.DeleteRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, "invalid request payload")
// 	}

// 	err := models.DeleteRows(context.Background(), req.ID1, req.ID2)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, "error deleting rows")
// 	}
// 	return c.JSON(http.StatusOK, "Rows deleted.")
// }

package controllers

import (
	"context"
	"log"
	"net/http"
	"os"

	"wallet-system/models"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Balance struct {
	ID      uuid.UUID `json:"id"`
	Balance int       `json:"balance"`
}

var conn *pgx.Conn

func init() {
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error parsing connection configuration: %v", err)
	}
	config.RuntimeParams["application_name"] = "wallet_system"

	conn, err = pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
}

func InitTableHandler(c echo.Context) error {
	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return models.InitTable(context.Background(), tx)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error initializing table")
	}
	return c.JSON(http.StatusOK, "Table initialized.")
}

func InsertRowsHandler(c echo.Context) error {
	accounts := [4]uuid.UUID{}
	for i := 0; i < len(accounts); i++ {
		accounts[i] = uuid.New()
	}

	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return models.InsertRows(context.Background(), tx, accounts)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error inserting rows")
	}
	return c.JSON(http.StatusOK, "Rows inserted.")
}

func PrintBalancesHandler(c echo.Context) error {
	var balances []models.Balance
	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var innerErr error
		balances, innerErr = models.PrintBalances(context.Background(), tx)
		return innerErr
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error fetching balances")
	}
	return c.JSON(http.StatusOK, balances)
}

func TransferFundsHandler(c echo.Context) error {
	var req models.TransferRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request payload")
	}

	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return models.TransferFunds(context.Background(), req.From, req.To, req.Amount)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error transferring funds")
	}
	return c.JSON(http.StatusOK, "Funds transferred.")
}

func DeleteRowsHandler(c echo.Context) error {
	var req models.DeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request payload")
	}

	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return models.DeleteRows(context.Background(), req.ID)
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error deleting rows")
	}
	return c.JSON(http.StatusOK, "Rows deleted.")
}
