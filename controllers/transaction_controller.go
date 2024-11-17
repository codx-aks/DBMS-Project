package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	models "wallet-system/models"
	utils "wallet-system/utils"
	helper "wallet-system/helper"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func TransactionApprovalHandler(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		log.Println("Error retrieving user from context")
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	var req struct {
		Pin    string `json:"pin"`
		Amount int    `json:"amount"`
		VendorID int  `json:"vendor_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		transactionID, err := models.TransactionApproval(context.Background(),user.RollNo, req.Pin, req.Amount,req.VendorID)
		if err != nil {
			log.Printf("Error approving transaction for user %v: %v", user.RollNo, err)
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Transaction approved",
			"transaction_id": transactionID ,
		})
	})

	if err != nil {
		log.Printf("Error processing transaction for user %v: %v", user.RollNo, err)
		return c.JSON(http.StatusInternalServerError, "Transaction failed")
	}

	return nil
}
