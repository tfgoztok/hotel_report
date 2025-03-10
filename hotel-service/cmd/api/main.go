package main

import (
	"log"
	"net/http"

	"github.com/olivere/elastic/v7"
	"github.com/tfgoztok/hotel-service/internal/api"
	"github.com/tfgoztok/hotel-service/internal/config"
	"github.com/tfgoztok/hotel-service/internal/db"
	"github.com/tfgoztok/hotel-service/internal/messaging"
	"github.com/tfgoztok/hotel-service/pkg/logger"
)

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

	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ", "error", err)
	}
	defer rabbitMQ.Close()

	// Connect to Elasticsearch
	esClient, err := elastic.NewClient(
		elastic.SetURL(cfg.ElasticsearchURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		logger.Fatal("Failed to connect to Elasticsearch", "error", err)
	}

	// Run migrations
	if err := db.RunMigrations(database, "./internal/db/migrations"); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	router := api.NewRouter(database, logger, rabbitMQ, esClient)

	logger.Info("Starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Fatal("Server failed to start", "error", err)
	}
}
