package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/service"
)

type ContactHandler struct {
	service *service.ContactService // Service for handling contact operations
}

func NewContactHandler(service *service.ContactService) *ContactHandler {
	return &ContactHandler{service: service} // Constructor for ContactHandler
}

func (h *ContactHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)                    // Get URL parameters
	hotelID, err := uuid.Parse(params["id"]) // Parse hotel ID from parameters
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest) // Handle invalid hotel ID
		return
	}

	var contact models.Contact // Create a new contact instance
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Handle JSON decoding errors
		return
	}
	contact.HotelID = hotelID // Associate the contact with the hotel ID

	if err := h.service.AddContact(r.Context(), &contact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Handle service errors
		return
	}

	w.WriteHeader(http.StatusCreated)  // Respond with 201 Created
	json.NewEncoder(w).Encode(contact) // Return the created contact
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)                             // Get URL parameters
	contactID, err := uuid.Parse(params["contactId"]) // Parse contact ID from parameters
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest) // Handle invalid contact ID
		return
	}

	if err := h.service.DeleteContact(r.Context(), contactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Handle service errors
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
}
