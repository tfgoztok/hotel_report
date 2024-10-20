package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

// HotelRepository is a struct that holds the database connection.
type HotelRepository struct {
	db *sql.DB // Database connection
}

// NewHotelRepository initializes a new HotelRepository with the provided database connection.
func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{db: db} // Return a new instance of HotelRepository
}

// Create inserts a new hotel record into the database.
func (r *HotelRepository) Create(ctx context.Context, hotel *models.Hotel) error {
	query := `
		INSERT INTO hotels (id, official_name, official_surname, company_title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	// Execute the insert query with hotel details
	_, err := r.db.ExecContext(ctx, query, hotel.ID, hotel.OfficialName, hotel.OfficialSurname, hotel.CompanyTitle, hotel.CreatedAt, hotel.UpdatedAt)
	return err // Return any error encountered
}

// Delete removes a hotel record from the database by its ID.
func (r *HotelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM hotels WHERE id = $1`
	// Execute the delete query using the hotel ID
	_, err := r.db.ExecContext(ctx, query, id)
	return err // Return any error encountered
}

// GetByID retrieves a hotel record from the database by its ID.
func (r *HotelRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Hotel, error) {
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

// GetByLocation retrieves a list of hotels based on the provided location.
func (r *HotelRepository) GetByLocation(ctx context.Context, location string) ([]*models.Hotel, error) {
	query := `
		SELECT id, official_name, official_surname, company_title, location
		FROM hotels
		WHERE location = $1
	`
	rows, err := r.db.QueryContext(ctx, query, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []*models.Hotel
	for rows.Next() {
		var hotel models.Hotel
		err := rows.Scan(&hotel.ID, &hotel.OfficialName, &hotel.OfficialSurname, &hotel.CompanyTitle, &hotel.Location)
		if err != nil {
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}
	return hotels, nil
}

// GetContactsByLocation retrieves a list of contacts associated with hotels based on the provided location.
func (r *HotelRepository) GetContactsByLocation(ctx context.Context, location string) ([]*models.Contact, error) {
	query := `
		SELECT c.id, c.hotel_id, c.type, c.content
		FROM contacts c
		JOIN hotels h ON c.hotel_id = h.id
		WHERE h.location = $1
	`
	rows, err := r.db.QueryContext(ctx, query, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(&contact.ID, &contact.HotelID, &contact.Type, &contact.Content)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &contact)
	}
	return contacts, nil
}