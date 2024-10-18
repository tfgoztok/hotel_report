package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tfgoztok/hotel-service/internal/api/handlers"
	"github.com/tfgoztok/hotel-service/internal/api/middleware"
	"github.com/tfgoztok/hotel-service/internal/repository"
	"github.com/tfgoztok/hotel-service/internal/service"
	"github.com/tfgoztok/hotel-service/pkg/logger"
)

func NewRouter(db *sql.DB, logger logger.Logger) http.Handler {
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

	// Middleware for logging requests
	r.Use(middleware.Logging(logger))

	// Define routes for hotel operations
	r.HandleFunc("/hotels", hotelHandler.CreateHotel).Methods("POST")                                 // Create a new hotel
	r.HandleFunc("/hotels/{id}", hotelHandler.DeleteHotel).Methods("DELETE")                          // Delete a hotel by ID
	r.HandleFunc("/hotels/{id}/contacts", contactHandler.AddContact).Methods("POST")                  // Add a contact to a hotel
	r.HandleFunc("/hotels/{id}/contacts/{contactId}", contactHandler.DeleteContact).Methods("DELETE") // Delete a contact by ID
	r.HandleFunc("/hotels/{id}/officials", hotelHandler.ListOfficials).Methods("GET")                 // List officials for a hotel
	r.HandleFunc("/hotels/{id}", hotelHandler.GetHotelDetails).Methods("GET")                         // Get details of a hotel

	return r // Return the configured router
}
