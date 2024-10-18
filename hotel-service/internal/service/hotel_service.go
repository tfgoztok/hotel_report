package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/repository"
)

// HotelService provides methods to manage hotels.
type HotelService struct {
	repo *repository.HotelRepository // Repository for hotel data
}

// NewHotelService creates a new instance of HotelService.
func NewHotelService(repo *repository.HotelRepository) *HotelService {
	return &HotelService{repo: repo} // Initialize HotelService with the provided repository
}

// CreateHotel creates a new hotel record in the repository.
func (s *HotelService) CreateHotel(ctx context.Context, hotel *models.Hotel) error {
	hotel.ID = uuid.New()            // Generate a new unique ID for the hotel
	hotel.CreatedAt = time.Now()     // Set the creation timestamp
	hotel.UpdatedAt = time.Now()     // Set the updated timestamp
	return s.repo.Create(ctx, hotel) // Save the hotel to the repository
}

// DeleteHotel removes a hotel record from the repository by its ID.
func (s *HotelService) DeleteHotel(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id) // Call the repository to delete the hotel
}

// GetHotelDetails retrieves hotel details by its ID.
func (s *HotelService) GetHotelDetails(ctx context.Context, id uuid.UUID) (*models.Hotel, error) {
	return s.repo.GetByID(ctx, id) // Fetch the hotel details from the repository
}

// ListOfficials retrieves the officials of a hotel by its ID.
func (s *HotelService) ListOfficials(ctx context.Context, id uuid.UUID) (*models.HotelOfficials, error) {
	hotel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.HotelOfficials{
		HotelID:         hotel.ID,
		OfficialName:    hotel.OfficialName,
		OfficialSurname: hotel.OfficialSurname,
	}, nil
}
