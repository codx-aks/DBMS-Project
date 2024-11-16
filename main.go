package main

import (
	"log"
	"os"
	"wallet-system/config"
	"wallet-system/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	if err := config.InitDBConnection(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer config.CloseDBConnection()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}
	log.Fatal(e.Start(":" + port))
}
