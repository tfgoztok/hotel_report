package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/service"
)

// HotelHandler handles hotel-related HTTP requests
type HotelHandler struct {
	service *service.HotelService // Service for hotel operations
}

// NewHotelHandler creates a new HotelHandler with the given service
func NewHotelHandler(service *service.HotelService) *HotelHandler {
	return &HotelHandler{service: service}
}

// CreateHotel handles the creation of a new hotel
func (h *HotelHandler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var hotel models.Hotel
	// Decode the incoming JSON request body into the hotel struct
	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Return error if decoding fails
		return
	}

	// Call the service to create the hotel
	if err := h.service.CreateHotel(r.Context(), &hotel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return error if creation fails
		return
	}

	w.WriteHeader(http.StatusCreated) // Respond with 201 Created
	json.NewEncoder(w).Encode(hotel)  // Return the created hotel as JSON
}

// DeleteHotel handles the deletion of a hotel by ID
func (h *HotelHandler) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)               // Get URL parameters
	id, err := uuid.Parse(params["id"]) // Parse the hotel ID from parameters
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest) // Return error if ID is invalid
		return
	}

	// Call the service to delete the hotel
	if err := h.service.DeleteHotel(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return error if deletion fails
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
}

// GetHotelDetails retrieves the details of a hotel by ID
func (h *HotelHandler) GetHotelDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	hotel, err := h.service.GetHotelDetails(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(hotel)
}
