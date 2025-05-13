package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goth/internal/components"
	"goth/internal/handlers"
	"goth/internal/middleware"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "dashboard" // From docker-compose.yml
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "changeme" // From docker-compose.yml
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "dashboard" // From docker-compose.yml
	}

	err := components.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		// Continue execution even if database connection fails
	} else {
		log.Println("Successfully connected to the database")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.HandleFunc("/healthz", handlers.HealthzHandler)

	stack := middleware.CreateStack(
		middleware.LoggingMiddleware,
	)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: stack(mux),
	}

	// static server
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	// Handle SIGINT (Ctrl+C), SIGTERM, and SIGQUIT
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Run server in a goroutine so it doesn't block shutdown handling
	go func() {
		log.Println("Starting server on port " + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-quit
	log.Printf("Shutdown signal received (%v), initiating graceful shutdown...", sig)

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close database connection before shutting down
	if components.DB != nil {
		log.Println("Closing database connection...")
		components.DB.Close()
	}

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		// If graceful shutdown fails, force close
		server.Close()
	}

	log.Println("Server stopped")
}
