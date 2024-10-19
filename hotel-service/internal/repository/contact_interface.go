package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

// ContactRepository interface defines the methods for contact operations
type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]*models.Contact, error)
}
