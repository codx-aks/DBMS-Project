package main

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"wallet-system/config"
	"wallet-system/routes"
	"runtime/debug"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	runtime.GC()
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(200)

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	if err := config.InitDBConnection(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer config.CloseDBConnection()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		log.Println(http.ListenAndServe(":6060", mux))
	}()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "",
		Output: os.Stdout,
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, 
		DisablePrintStack: true,   
	}))

	routes.RegisterRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}
	server := &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second, 
		Handler:           e,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := e.Shutdown(context.Background()); err != nil {
			log.Fatal("Server shutdown failed:", err)
		}
	}()

	log.Printf("Starting server on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}
}
