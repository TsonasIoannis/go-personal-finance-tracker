package services

import (
	"context"
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

func (m *MockCategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Create valid category", func(t *testing.T) {
		category := &models.Category{Name: "Groceries", Description: "Food and drinks"}
		mockRepo.On("CreateCategory", ctx, category).Return(nil)

		err := service.CreateCategory(ctx, category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCategories(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Retrieve multiple categories", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		categories := []models.Category{
			{ID: 1, Name: "Health", Description: "Medical expenses"},
			{ID: 2, Name: "Entertainment", Description: "Movies and fun"},
		}

		mockRepo.On("GetAllCategories", ctx).Return(categories, nil)

		result, err := service.GetCategories(ctx)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retrieve categories when none exist", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		mockRepo.On("GetAllCategories", ctx).Return([]models.Category{}, nil)

		result, err := service.GetCategories(ctx)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Update existing category", func(t *testing.T) {
		category := &models.Category{ID: 1, Name: "Travel", Description: "Flights & hotels"}
		mockRepo.On("UpdateCategory", ctx, category).Return(nil)

		err := service.UpdateCategory(ctx, category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to update non-existent category", func(t *testing.T) {
		category := &models.Category{ID: 9999, Name: "Luxury", Description: "Expensive items"}
		mockRepo.On("UpdateCategory", ctx, category).Return(errors.New("category not found"))

		err := service.UpdateCategory(ctx, category)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Delete existing category", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, uint(1)).Return(nil)

		err := service.DeleteCategory(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent category", func(t *testing.T) {
		mockRepo.On("DeleteCategory", ctx, uint(9999)).Return(errors.New("category not found"))

		err := service.DeleteCategory(ctx, 9999)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
