package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olivere/elastic/v7"
	"github.com/tfgoztok/hotel-service/internal/api/graphql"
	"github.com/tfgoztok/hotel-service/internal/api/handlers"
	"github.com/tfgoztok/hotel-service/internal/api/middleware"
	"github.com/tfgoztok/hotel-service/internal/messaging"
	"github.com/tfgoztok/hotel-service/internal/repository"
	"github.com/tfgoztok/hotel-service/internal/service"
	"github.com/tfgoztok/hotel-service/pkg/logger"
)

func NewRouter(db *sql.DB, logger logger.Logger, rabbitMQ messaging.RabbitMQInterface, esClient *elastic.Client) http.Handler {
	// Create a new router instance
	r := mux.NewRouter()

	// Initialize repositories for hotels and contacts
	hotelRepo := repository.NewHotelRepository(db)
	contactRepo := repository.NewContactRepository(db)

	// Initialize services for hotels and contacts
	hotelService := service.NewHotelService(hotelRepo)
	contactService := service.NewContactService(contactRepo)

	// Initialize handlers for hotels and contacts
	hotelHandler := handlers.NewHotelHandler(hotelService)
	contactHandler := handlers.NewContactHandler(contactService)

	// Initialize handler for elk
	reportHandler := handlers.NewReportHandler(rabbitMQ, esClient)

	graphqlService := graphql.NewGraphQLService(hotelService)
	graphqlHandler, err := handlers.NewGraphQLHandler(graphqlService)
	if err != nil {
		logger.Fatal("Failed to create GraphQL handler", "error", err)
	}
	r.Handle("/graphql", graphqlHandler)

	// Middleware for logging requests
	r.Use(middleware.Logging(logger))

	// Define routes for hotel operations
	r.HandleFunc("/hotels", hotelHandler.CreateHotel).Methods("POST")                                 // Create a new hotel
	r.HandleFunc("/hotels/{id}", hotelHandler.DeleteHotel).Methods("DELETE")                          // Delete a hotel by ID
	r.HandleFunc("/hotels/{id}/contacts", contactHandler.AddContact).Methods("POST")                  // Add a contact to a hotel
	r.HandleFunc("/hotels/{id}/contacts/{contactId}", contactHandler.DeleteContact).Methods("DELETE") // Delete a contact by ID
	r.HandleFunc("/hotels/{id}/officials", hotelHandler.ListOfficials).Methods("GET")                 // List officials for a hotel
	r.HandleFunc("/hotels/{id}", hotelHandler.GetHotelDetails).Methods("GET")                         // Get details of a hotel
	r.HandleFunc("/reports/request", reportHandler.RequestReport).Methods("POST")                     // Request report from report-service

	return r // Return the configured router
}
