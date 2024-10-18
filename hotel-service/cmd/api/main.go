package main

import (
	"log"
	"net/http"

	"github.com/tfgoztok/hotel-service/internal/api"
	"github.com/tfgoztok/hotel-service/internal/config"
	"github.com/tfgoztok/hotel-service/internal/db"
	"github.com/tfgoztok/hotel-service/pkg/logger"
)

// main function is the entry point of the application
func main() {
	// Load the application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err) // Log fatal error if config loading fails
	}

	// Initialize the logger with the specified log level
	logger := logger.New(cfg.LogLevel)

	// Connect to the database using the provided database URL
	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err) // Log fatal error if database connection fails
	}
	defer db.Close() // Ensure the database connection is closed when main exits

	// Create a new router with the database and logger
	router := api.NewRouter(db, logger)

	// Start the HTTP server on the specified port
	logger.Info("Starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Fatal("Server failed to start", "error", err) // Log fatal error if server fails to start
	}
}
