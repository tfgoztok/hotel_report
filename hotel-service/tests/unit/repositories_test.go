package unit

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/repository"
)

func TestHotelRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewHotelRepository(db)

	hotel := &models.Hotel{
		ID:              uuid.New(),
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
		Location:        "New York",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mock.ExpectExec("INSERT INTO hotels").
		WithArgs(hotel.ID, hotel.OfficialName, hotel.OfficialSurname, hotel.CompanyTitle, hotel.Location, hotel.CreatedAt, hotel.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), hotel)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHotelRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewHotelRepository(db)

	hotelID := uuid.New()

	mock.ExpectExec("DELETE FROM hotels").
		WithArgs(hotelID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHotelRepositoryGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewHotelRepository(db)

	hotelID := uuid.New()
	expectedHotel := &models.Hotel{
		ID:              hotelID,
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "official_name", "official_surname", "company_title", "created_at", "updated_at"}).
		AddRow(expectedHotel.ID, expectedHotel.OfficialName, expectedHotel.OfficialSurname, expectedHotel.CompanyTitle, expectedHotel.CreatedAt, expectedHotel.UpdatedAt)

	mock.ExpectQuery("SELECT (.+) FROM hotels").
		WithArgs(hotelID).
		WillReturnRows(rows)

	hotel, err := repo.GetByID(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.Equal(t, expectedHotel, hotel)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContactRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewContactRepository(db)

	contact := &models.Contact{
		ID:        uuid.New(),
		HotelID:   uuid.New(),
		Type:      "email",
		Content:   "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO contacts").
		WithArgs(contact.ID, contact.HotelID, contact.Type, contact.Content, contact.CreatedAt, contact.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), contact)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContactRepositoryDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewContactRepository(db)

	contactID := uuid.New()

	mock.ExpectExec("DELETE FROM contacts").
		WithArgs(contactID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), contactID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestContactRepositoryGetByHotelID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewContactRepository(db)

	hotelID := uuid.New()
	expectedContacts := []*models.Contact{
		{
			ID:        uuid.New(),
			HotelID:   hotelID,
			Type:      "email",
			Content:   "test1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			HotelID:   hotelID,
			Type:      "phone",
			Content:   "+1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "hotel_id", "type", "content", "created_at", "updated_at"})
	for _, contact := range expectedContacts {
		rows.AddRow(contact.ID, contact.HotelID, contact.Type, contact.Content, contact.CreatedAt, contact.UpdatedAt)
	}

	mock.ExpectQuery("SELECT (.+) FROM contacts").
		WithArgs(hotelID).
		WillReturnRows(rows)

	contacts, err := repo.GetByHotelID(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.Equal(t, expectedContacts, contacts)
	assert.NoError(t, mock.ExpectationsWereMet())
}
