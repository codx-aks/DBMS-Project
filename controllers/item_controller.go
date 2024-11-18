package controllers

import (
	"context"
	"log"
	"net/http"
	helper "wallet-system/helper"
	models "wallet-system/models"
	"strconv"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func VendorItems(c echo.Context) error {
	ID := c.QueryParam("id") 
	
	if ID == "" {
		log.Printf("Error: Vendor ID not provided")
		return c.JSON(http.StatusBadRequest, "Vendor ID is required")
	}
	
	vendorID, err := strconv.Atoi(ID)
	if err != nil {
		log.Printf("Invalid vendorId format: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid vendor ID format")
	}

	log.Printf("Received vendorId: %s", ID)
	
	var items []models.Item
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var innerErr error
		items, innerErr = helper.GetVendorItems(context.Background(), tx, vendorID)
		return innerErr
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error fetching items")
	}
	return c.JSON(http.StatusOK, items)
}