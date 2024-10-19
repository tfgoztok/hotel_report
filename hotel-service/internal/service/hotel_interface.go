package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

type HotelService interface {
	CreateHotel(ctx context.Context, hotel *models.Hotel) error
	DeleteHotel(ctx context.Context, id uuid.UUID) error
	GetHotelDetails(ctx context.Context, id uuid.UUID) (*models.Hotel, error)
	ListOfficials(ctx context.Context, id uuid.UUID) (*models.HotelOfficials, error)
}
