package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goth/internal/handlers"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HelloHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
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

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		// If graceful shutdown fails, force close
		server.Close()
	}

	log.Println("Server stopped")
}
