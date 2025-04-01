package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"urlShortner/api"
	"urlShortner/config"
	"urlShortner/db"
)

func main() {
	// Loding configuration from Config/config.go
	cfg := config.Load()

	// Connect to MongoDB
	db.Connect(cfg)

	// Create a new router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/api/shorten", api.CreateShortURLHandler(cfg)).Methods("POST")
	router.HandleFunc("/api/stats/{shortCode}", api.URLStatsHandler()).Methods("GET")
	router.HandleFunc("/health", api.HealthCheckHandler()).Methods("GET")
	router.HandleFunc("/{shortCode}", api.RedirectHandler()).Methods("GET")
	// Set up the server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s\n", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Disconnect from MongoDB
	if err := db.Disconnect(ctx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}

	log.Println("Server exiting")
}