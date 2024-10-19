package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

// Ensure HotelRepositoryImpl implements HotelRepository interface
var _ HotelRepository = (*HotelRepositoryImpl)(nil)

// HotelRepositoryImpl is a struct that holds the database connection.
type HotelRepositoryImpl struct {
	db *sql.DB // Database connection
}

// NewHotelRepository initializes a new HotelRepositoryImpl with the provided database connection.
func NewHotelRepository(db *sql.DB) *HotelRepositoryImpl {
	return &HotelRepositoryImpl{db: db} // Return a new instance of HotelRepositoryImpl
}

// Create inserts a new hotel record into the database.
func (r *HotelRepositoryImpl) Create(ctx context.Context, hotel *models.Hotel) error {
	query := `
		INSERT INTO hotels (id, official_name, official_surname, company_title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	// Execute the insert query with hotel details
	_, err := r.db.ExecContext(ctx, query, hotel.ID, hotel.OfficialName, hotel.OfficialSurname, hotel.CompanyTitle, hotel.CreatedAt, hotel.UpdatedAt)
	return err // Return any error encountered
}

// Delete removes a hotel record from the database by its ID.
func (r *HotelRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM hotels WHERE id = $1`
	// Execute the delete query using the hotel ID
	_, err := r.db.ExecContext(ctx, query, id)
	return err // Return any error encountered
}

// GetByID retrieves a hotel record from the database by its ID.
func (r *HotelRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Hotel, error) {
	query := `
		SELECT id, official_name, official_surname, company_title, created_at, updated_at
		FROM hotels
		WHERE id = $1
	`
	var hotel models.Hotel // Variable to hold the retrieved hotel
	// Execute the select query and scan the result into the hotel variable
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID, &hotel.OfficialName, &hotel.OfficialSurname, &hotel.CompanyTitle, &hotel.CreatedAt, &hotel.UpdatedAt,
	)
	if err != nil {
		return nil, err // Return nil and the error if something went wrong
	}
	return &hotel, nil // Return the retrieved hotel
}
