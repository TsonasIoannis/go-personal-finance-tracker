package services

import (
	"errors"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBudgetRepository is a mock implementation of the BudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) CreateBudget(budget *models.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) GetBudgetByID(id uint) (*models.Budget, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Budget), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBudgetRepository) GetBudgetsByUserID(userID uint) ([]models.Budget, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Budget), args.Error(1)
}

func (m *MockBudgetRepository) UpdateBudget(budget *models.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) DeleteBudget(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)

	t.Run("Create valid budget", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      1000.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		mockRepo.On("CreateBudget", budget).Return(nil)

		err := service.CreateBudget(budget)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to create budget with negative limit", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      -100.00, // Invalid limit
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		err := service.CreateBudget(budget)
		assert.Error(t, err)
		assert.Equal(t, "budget limit must be greater than zero", err.Error())
		mockRepo.AssertNotCalled(t, "CreateBudget")
	})
}

func TestUpdateBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)

	t.Run("Update existing budget", func(t *testing.T) {
		budget := &models.Budget{
			ID:         1,
			UserID:     1,
			CategoryID: 2,
			Limit:      500.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		mockRepo.On("UpdateBudget", budget).Return(nil)

		err := service.UpdateBudget(budget)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to update budget with zero limit", func(t *testing.T) {
		budget := &models.Budget{
			ID:         1,
			UserID:     1,
			CategoryID: 2,
			Limit:      0, // Invalid limit
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		err := service.UpdateBudget(budget)
		assert.Error(t, err)
		assert.Equal(t, "budget limit must be positive", err.Error())
		mockRepo.AssertNotCalled(t, "UpdateBudget")
	})
}

func TestGetBudgetsByUser(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)

	t.Run("Retrieve budgets for user", func(t *testing.T) {
		budgets := []models.Budget{
			{ID: 1, UserID: 1, CategoryID: 2, Limit: 1000},
			{ID: 2, UserID: 1, CategoryID: 3, Limit: 500},
		}

		mockRepo.On("GetBudgetsByUserID", uint(1)).Return(budgets, nil)

		result, err := service.GetBudgetsByUser(1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retrieve budgets for user with no budgets", func(t *testing.T) {
		mockRepo.On("GetBudgetsByUserID", uint(999)).Return([]models.Budget{}, nil)

		result, err := service.GetBudgetsByUser(999)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)

	t.Run("Delete existing budget", func(t *testing.T) {
		mockRepo.On("DeleteBudget", uint(1)).Return(nil)

		err := service.DeleteBudget(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent budget", func(t *testing.T) {
		mockRepo.On("DeleteBudget", uint(9999)).Return(errors.New("budget not found"))

		err := service.DeleteBudget(9999)
		assert.Error(t, err)
		assert.Equal(t, "budget not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
