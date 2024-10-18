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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database, "./internal/db/migrations"); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	router := api.NewRouter(database, logger)

	logger.Info("Starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Fatal("Server failed to start", "error", err)
	}
}
