package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"wallet-system/models"
	"wallet-system/utils"
	"wallet-system/helper"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

const (
	SessionCookieName = "session_token"
)

func SignupHandler(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return models.CreateUser(tx, user)
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error creating user")
	}

	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		log.Printf("Error generating session token: %v", err)
		return c.JSON(http.StatusInternalServerError, "Error creating session")
	}

	setSessionCookie(c, sessionToken)

	return c.JSON(http.StatusOK, "User successfully created")
}

func LoginHandler(c echo.Context) error {
	var req struct {
		Email string `json:"email"`
		Pin   string `json:"pin"`
	}
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding login data: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	var user models.User
	err := crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var innerErr error
		user, innerErr = models.GetUserByEmailAndPassword(tx, req.Email, req.Pin)
		return innerErr
	})
	if err != nil {
		log.Printf("Login error: %v", err)
		return c.JSON(http.StatusUnauthorized, "Invalid email or password")
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

func setSessionCookie(c echo.Context, sessionToken string) {
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