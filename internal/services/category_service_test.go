package services

import (
	"errors"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository implements the CategoryRepository interface
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) CreateCategory(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) GetAllCategories() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategory(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) DeleteCategory(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	t.Run("Create valid category", func(t *testing.T) {
		category := &models.Category{Name: "Groceries", Description: "Food and drinks"}
		mockRepo.On("CreateCategory", category).Return(nil)

		err := service.CreateCategory(category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCategories(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	t.Run("Retrieve multiple categories", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		categories := []models.Category{
			{ID: 1, Name: "Health", Description: "Medical expenses"},
			{ID: 2, Name: "Entertainment", Description: "Movies and fun"},
		}

		mockRepo.On("GetAllCategories").Return(categories, nil)

		result, err := service.GetCategories()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retrieve categories when none exist", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		mockRepo.On("GetAllCategories").Return([]models.Category{}, nil)

		result, err := service.GetCategories()
		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	t.Run("Update existing category", func(t *testing.T) {
		category := &models.Category{ID: 1, Name: "Travel", Description: "Flights & hotels"}
		mockRepo.On("UpdateCategory", category).Return(nil)

		err := service.UpdateCategory(category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to update non-existent category", func(t *testing.T) {
		category := &models.Category{ID: 9999, Name: "Luxury", Description: "Expensive items"}
		mockRepo.On("UpdateCategory", category).Return(errors.New("category not found"))

		err := service.UpdateCategory(category)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	t.Run("Delete existing category", func(t *testing.T) {
		mockRepo.On("DeleteCategory", uint(1)).Return(nil)

		err := service.DeleteCategory(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent category", func(t *testing.T) {
		mockRepo.On("DeleteCategory", uint(9999)).Return(errors.New("category not found"))

		err := service.DeleteCategory(9999)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
