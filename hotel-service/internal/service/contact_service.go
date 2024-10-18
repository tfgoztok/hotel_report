package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/repository"
)

// ContactService provides methods to manage contacts.
type ContactService struct {
	repo *repository.ContactRepository // Repository for contact data
}

// NewContactService creates a new instance of ContactService.
func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo} // Initialize ContactService with the provided repository
}

// AddContact adds a new contact to the repository.
func (s *ContactService) AddContact(ctx context.Context, contact *models.Contact) error {
	contact.ID = uuid.New()            // Generate a new unique ID for the contact
	contact.CreatedAt = time.Now()     // Set the creation timestamp
	contact.UpdatedAt = time.Now()     // Set the updated timestamp
	return s.repo.Create(ctx, contact) // Save the contact to the repository
}

// DeleteContact removes a contact from the repository by its ID.
func (s *ContactService) DeleteContact(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id) // Call the repository to delete the contact
}

// GetContactsByHotelID retrieves all contacts associated with a specific hotel ID.
func (s *ContactService) GetContactsByHotelID(ctx context.Context, hotelID uuid.UUID) ([]*models.Contact, error) {
	return s.repo.GetByHotelID(ctx, hotelID) // Fetch contacts from the repository by hotel ID
}
