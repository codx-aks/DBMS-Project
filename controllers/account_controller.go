package controllers

import (
	"context"
	"net/http"

	"wallet-system/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func InitTableHandler(c echo.Context) error {
	err := models.InitTable(context.Background())
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

	err := models.InsertRows(context.Background(), accounts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error inserting rows")
	}
	return c.JSON(http.StatusOK, "Rows inserted.")
}

func PrintBalancesHandler(c echo.Context) error {
	balances, err := models.PrintBalances(context.Background())
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

	err := models.TransferFunds(context.Background(), req.From, req.To, req.Amount)
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

	err := models.DeleteRows(context.Background(), req.ID1, req.ID2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error deleting rows")
	}
	return c.JSON(http.StatusOK, "Rows deleted.")
}
