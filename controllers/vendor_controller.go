package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"fmt"
	utils "wallet-system/utils"
	helper "wallet-system/helper"
	models "wallet-system/models"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func VendorLoginHandler(c echo.Context) error {
	var req struct {
		ID int `json:"id"`
		Password   string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding login data: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}
	var err error
	req.Password,err = utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	var user models.Vendor
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var innerErr error
		user, innerErr = helper.GetVendorByIDAndPassword(tx, req.ID, req.Password)
		return innerErr
	})
	if err != nil {
		log.Printf("Login error: %v", err)
		return c.JSON(http.StatusUnauthorized, "Invalid id or password")
	}

	c.Set("user", user)

	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		log.Printf("Error generating session token: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error creating session")
	}

	setSessionCookie(c, sessionToken)

	return c.JSON(http.StatusOK, user)
}

func setSessionCookieVendor(c echo.Context, sessionToken string) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)
}

func vendorTransactions(c echo.Context) error {
	user, ok := c.Get("user").(models.Vendor)
	if !ok {
		log.Println("Error retrieving user from context")
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	var transactions []models.Transaction
	var total int
	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var innerErr error
		transactions, total, innerErr = helper.GetTransactionsByVendor(context.Background(), tx, user.ID)
		return innerErr
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error fetching transactions")
	}

	var resp struct {
		Transactions []models.Transaction `json:"transactions"`
		Total   int `json:"total"`
	}

	resp.Transactions = transactions
	resp.Total = total

	return c.JSON(http.StatusOK, resp)
}