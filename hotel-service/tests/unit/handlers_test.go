package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tfgoztok/hotel-service/internal/api/handlers"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/repository"
)

// MockHotelService is a mock implementation of the HotelService
type MockHotelService struct {
	mock.Mock
}

// Ensure MockContactRepository implements ContactRepository interface
var _ repository.ContactRepository = (*MockContactRepository)(nil)

func (m *MockHotelService) CreateHotel(ctx context.Context, hotel *models.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelService) DeleteHotel(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHotelService) GetHotelDetails(ctx context.Context, id uuid.UUID) (*models.Hotel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Hotel), args.Error(1)
}

func (m *MockHotelService) ListOfficials(ctx context.Context, id uuid.UUID) (*models.HotelOfficials, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.HotelOfficials), args.Error(1)
}

// MockContactService is a mock implementation of the ContactService
type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) AddContact(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactService) DeleteContact(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateHotelHandler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := handlers.NewHotelHandler(mockService)

	hotel := &models.Hotel{
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
	}

	mockService.On("CreateHotel", mock.Anything, mock.AnythingOfType("*models.Hotel")).Return(nil)

	body, _ := json.Marshal(hotel)
	req, _ := http.NewRequest("POST", "/hotels", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateHotel(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteHotelHandler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := handlers.NewHotelHandler(mockService)

	hotelID := uuid.New()
	mockService.On("DeleteHotel", mock.Anything, hotelID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/hotels/"+hotelID.String(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/hotels/{id}", handler.DeleteHotel).Methods("DELETE")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetHotelDetailsHandler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := handlers.NewHotelHandler(mockService)

	hotelID := uuid.New()
	hotel := &models.Hotel{
		ID:              hotelID,
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
	}

	mockService.On("GetHotelDetails", mock.Anything, hotelID).Return(hotel, nil)

	req, _ := http.NewRequest("GET", "/hotels/"+hotelID.String(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/hotels/{id}", handler.GetHotelDetails).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)

	var responseHotel models.Hotel
	json.Unmarshal(rr.Body.Bytes(), &responseHotel)
	assert.Equal(t, hotel, &responseHotel)
}

func TestListOfficialsHandler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := handlers.NewHotelHandler(mockService)

	hotelID := uuid.New()
	officials := &models.HotelOfficials{
		HotelID:         hotelID,
		OfficialName:    "John",
		OfficialSurname: "Doe",
	}

	mockService.On("ListOfficials", mock.Anything, hotelID).Return(officials, nil)

	req, _ := http.NewRequest("GET", "/hotels/"+hotelID.String()+"/officials", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/hotels/{id}/officials", handler.ListOfficials).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)

	var responseOfficials models.HotelOfficials
	json.Unmarshal(rr.Body.Bytes(), &responseOfficials)
	assert.Equal(t, officials, &responseOfficials)
}

func TestAddContactHandler(t *testing.T) {
	mockService := new(MockContactService)
	handler := handlers.NewContactHandler(mockService)

	hotelID := uuid.New()
	contact := &models.Contact{
		HotelID: hotelID,
		Type:    "email",
		Content: "test@example.com",
	}

	mockService.On("AddContact", mock.Anything, mock.AnythingOfType("*models.Contact")).Return(nil)

	body, _ := json.Marshal(contact)
	req, _ := http.NewRequest("POST", "/hotels/"+hotelID.String()+"/contacts", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/hotels/{id}/contacts", handler.AddContact).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteContactHandler(t *testing.T) {
	mockService := new(MockContactService)
	handler := handlers.NewContactHandler(mockService)

	contactID := uuid.New()
	mockService.On("DeleteContact", mock.Anything, contactID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/contacts/"+contactID.String(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/contacts/{contactId}", handler.DeleteContact).Methods("DELETE")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}
