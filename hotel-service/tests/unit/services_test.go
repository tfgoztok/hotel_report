package unit

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tfgoztok/hotel-service/internal/models"
	"github.com/tfgoztok/hotel-service/internal/repository"
	"github.com/tfgoztok/hotel-service/internal/service"
)

var _ repository.HotelRepository = (*MockHotelRepository)(nil)

// MockHotelRepository is a mock implementation of the HotelRepository
type MockHotelRepository struct {
	mock.Mock
}

func (m *MockHotelRepository) Create(ctx context.Context, hotel *models.Hotel) error {
	args := m.Called(ctx, hotel)
	return args.Error(0)
}

func (m *MockHotelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHotelRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Hotel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Hotel), args.Error(1)
}

// MockContactRepository is a mock implementation of the ContactRepository
type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContactRepository) GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]*models.Contact, error) {
	args := m.Called(ctx, hotelID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func TestHotelServiceCreate(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := service.NewHotelService(mockRepo)

	hotel := &models.Hotel{
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Hotel")).Return(nil)

	err := service.CreateHotel(context.Background(), hotel)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, hotel.ID)
	assert.NotZero(t, hotel.CreatedAt)
	assert.NotZero(t, hotel.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestHotelServiceDelete(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := service.NewHotelService(mockRepo)

	hotelID := uuid.New()

	mockRepo.On("Delete", mock.Anything, hotelID).Return(nil)

	err := service.DeleteHotel(context.Background(), hotelID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHotelServiceGetDetails(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := service.NewHotelService(mockRepo)

	hotelID := uuid.New()
	hotel := &models.Hotel{
		ID:              hotelID,
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
	}

	mockRepo.On("GetByID", mock.Anything, hotelID).Return(hotel, nil)

	result, err := service.GetHotelDetails(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.Equal(t, hotel, result)
	mockRepo.AssertExpectations(t)
}

func TestHotelServiceListOfficials(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := service.NewHotelService(mockRepo)

	hotelID := uuid.New()
	hotel := &models.Hotel{
		ID:              hotelID,
		OfficialName:    "John",
		OfficialSurname: "Doe",
		CompanyTitle:    "Test Hotel",
	}

	mockRepo.On("GetByID", mock.Anything, hotelID).Return(hotel, nil)

	officials, err := service.ListOfficials(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.Equal(t, hotelID, officials.HotelID)
	assert.Equal(t, hotel.OfficialName, officials.OfficialName)
	assert.Equal(t, hotel.OfficialSurname, officials.OfficialSurname)
	mockRepo.AssertExpectations(t)
}

func TestContactServiceAdd(t *testing.T) {
	mockRepo := new(MockContactRepository)
	service := service.NewContactService(mockRepo)

	contact := &models.Contact{
		HotelID: uuid.New(),
		Type:    "email",
		Content: "test@example.com",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Contact")).Return(nil)

	err := service.AddContact(context.Background(), contact)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, contact.ID)
	assert.NotZero(t, contact.CreatedAt)
	assert.NotZero(t, contact.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestContactServiceDelete(t *testing.T) {
	mockRepo := new(MockContactRepository)
	service := service.NewContactService(mockRepo)

	contactID := uuid.New()

	mockRepo.On("Delete", mock.Anything, contactID).Return(nil)

	err := service.DeleteContact(context.Background(), contactID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestContactServiceGetByHotelID(t *testing.T) {
	mockRepo := new(MockContactRepository)
	service := service.NewContactService(mockRepo)

	hotelID := uuid.New()
	contacts := []*models.Contact{
		{
			ID:      uuid.New(),
			HotelID: hotelID,
			Type:    "email",
			Content: "test1@example.com",
		},
		{
			ID:      uuid.New(),
			HotelID: hotelID,
			Type:    "phone",
			Content: "+1234567890",
		},
	}

	mockRepo.On("GetByHotelID", mock.Anything, hotelID).Return(contacts, nil)

	result, err := service.GetContactsByHotelID(context.Background(), hotelID)

	assert.NoError(t, err)
	assert.Equal(t, contacts, result)
	mockRepo.AssertExpectations(t)
}
