package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tfgoztok/hotel-service/internal/models"
)

// ContactRepository is a struct that holds the database connection.
type ContactRepository struct {
	db *sql.DB // Database connection
}

// NewContactRepository initializes a new ContactRepository with the provided database connection.
func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create inserts a new contact into the database.
// It takes a context for managing request-scoped values and a pointer to a Contact model.
func (r *ContactRepository) Create(ctx context.Context, contact *models.Contact) error {
	query := `
		INSERT INTO contacts (id, hotel_id, type, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	// Execute the insert query with the contact's details.
	_, err := r.db.ExecContext(ctx, query, contact.ID, contact.HotelID, contact.Type, contact.Content, contact.CreatedAt, contact.UpdatedAt)
	return err // Return any error encountered during execution.
}

// Delete removes a contact from the database by its ID.
// It takes a context and the UUID of the contact to be deleted.
func (r *ContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM contacts WHERE id = $1`
	// Execute the delete query using the provided contact ID.
	_, err := r.db.ExecContext(ctx, query, id)
	return err // Return any error encountered during execution.
}

// GetByHotelID retrieves all contacts associated with a specific hotel ID.
// It returns a slice of pointers to Contact models and an error if any occurs.
func (r *ContactRepository) GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]*models.Contact, error) {
	query := `
		SELECT id, hotel_id, type, content, created_at, updated_at
		FROM contacts
		WHERE hotel_id = $1
	`
	// Execute the query to fetch contacts for the specified hotel ID.
	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err // Return nil and the error if the query fails.
	}
	defer rows.Close() // Ensure rows are closed after processing.

	var contacts []*models.Contact // Slice to hold the retrieved contacts.
	for rows.Next() {
		var contact models.Contact // Temporary variable to hold each contact.
		// Scan the row into the contact variable.
		err := rows.Scan(&contact.ID, &contact.HotelID, &contact.Type, &contact.Content, &contact.CreatedAt, &contact.UpdatedAt)
		if err != nil {
			return nil, err // Return nil and the error if scanning fails.
		}
		contacts = append(contacts, &contact) // Append the contact to the slice.
	}
	return contacts, nil // Return the slice of contacts and nil for no error.
}
