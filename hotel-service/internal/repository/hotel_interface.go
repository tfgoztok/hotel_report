package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

type HotelRepository interface {
	Create(ctx context.Context, hotel *models.Hotel) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Hotel, error)
}
