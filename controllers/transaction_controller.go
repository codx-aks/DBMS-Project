package controllers

import (
	"context"
	"log"
	"net/http"
	helper "wallet-system/helper"
	"strconv"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"fmt"
	"github.com/labstack/echo/v4"
)

func TransactionApprovalHandler(c echo.Context) error {
	// user, ok := c.Get("user").(models.User)
	// if !ok {
	// 	log.Println("Error retrieving user from context")
	// 	fmt.Println(user.RollNo)
	// 	return c.JSON(http.StatusUnauthorized, "User not authenticated")
	// }

	// var user = models.User{
	// 	Name:         "John Doe",
	// 	RollNo:       "123456789",
	// 	Email:        "johndoe@example.com",
	// 	Pin:          "$2a$10$Y3HQLfjU/q2twM9oA8buaO8Fe1XTeoMlEp84kS4FK5cdJYYcVzPae", 
	// 	IsVerified:   true,
	// 	IsApproved:   true,
	// 	WalletBalance: 1000,
	// }

	var req struct {
		RollNo string `json:"roll_no"`
		Pin    string `json:"pin"`
		Amount int    `json:"amount"`
		VendorID string  `json:"vendor_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	vID, err := strconv.Atoi(req.VendorID)
	if err != nil {
		fmt.Println("Error:", err)
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		transactionID, err := helper.TransactionApproval(context.Background(),req.RollNo, req.Pin, req.Amount,vID)
		if err != nil {
			log.Printf("Error approving transaction for user %v: %v", req.RollNo, err)
			return err
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Transaction approved",
			"transaction_id": transactionID ,
		})
	})

	if err != nil {
		log.Printf("Error processing transaction for user %v: %v", req.RollNo, err)
		return c.JSON(http.StatusInternalServerError, "Transaction failed")
	}

	return nil
}
