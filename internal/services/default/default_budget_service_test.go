package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBudgetRepository is a mock implementation of the BudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) CreateBudget(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) GetBudgetByID(ctx context.Context, id uint) (*models.Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Budget), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBudgetRepository) GetBudgetsByUserID(ctx context.Context, userID uint) ([]models.Budget, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Budget), args.Error(1)
}

func (m *MockBudgetRepository) UpdateBudget(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) DeleteBudget(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	t.Run("Create valid budget", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      1000.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		mockRepo.On("CreateBudget", ctx, budget).Return(nil)

		err := service.CreateBudget(ctx, budget)
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

		err := service.CreateBudget(ctx, budget)
		assert.Error(t, err)
		assert.Equal(t, "budget limit must be greater than zero", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindValidation))
		mockRepo.AssertNotCalled(t, "CreateBudget")
	})

	t.Run("Fail to create budget with invalid date range", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      100.00,
			StartDate:  time.Now().AddDate(0, 1, 0),
			EndDate:    time.Now(),
		}

		err := service.CreateBudget(ctx, budget)
		assert.Error(t, err)
		assert.Equal(t, "start date cannot be after end date", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindValidation))
		mockRepo.AssertNotCalled(t, "CreateBudget")
	})
}

func TestUpdateBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	t.Run("Update existing budget", func(t *testing.T) {
		budget := &models.Budget{
			ID:         1,
			UserID:     1,
			CategoryID: 2,
			Limit:      500.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		mockRepo.On("UpdateBudget", ctx, budget).Return(nil)

		err := service.UpdateBudget(ctx, budget)
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

		err := service.UpdateBudget(ctx, budget)
		assert.Error(t, err)
		assert.Equal(t, "budget limit must be positive", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindValidation))
		mockRepo.AssertNotCalled(t, "UpdateBudget")
	})

	t.Run("Fail to update budget with invalid date range", func(t *testing.T) {
		budget := &models.Budget{
			ID:         1,
			UserID:     1,
			CategoryID: 2,
			Limit:      100,
			StartDate:  time.Now().AddDate(0, 1, 0),
			EndDate:    time.Now(),
		}

		err := service.UpdateBudget(ctx, budget)
		assert.Error(t, err)
		assert.Equal(t, "start date cannot be after end date", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindValidation))
		mockRepo.AssertNotCalled(t, "UpdateBudget")
	})
}

func TestGetBudgetsByUser(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	t.Run("Retrieve budgets for user", func(t *testing.T) {
		budgets := []models.Budget{
			{ID: 1, UserID: 1, CategoryID: 2, Limit: 1000},
			{ID: 2, UserID: 1, CategoryID: 3, Limit: 500},
		}

		mockRepo.On("GetBudgetsByUserID", ctx, uint(1)).Return(budgets, nil)

		result, err := service.GetBudgetsByUser(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Retrieve budgets for user with no budgets", func(t *testing.T) {
		mockRepo.On("GetBudgetsByUserID", ctx, uint(999)).Return([]models.Budget{}, nil)

		result, err := service.GetBudgetsByUser(ctx, 999)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBudget(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	t.Run("Delete existing budget", func(t *testing.T) {
		mockRepo.On("GetBudgetByID", ctx, uint(1)).Return(&models.Budget{ID: 1, UserID: 1}, nil)
		mockRepo.On("DeleteBudget", ctx, uint(1)).Return(nil)

		err := service.DeleteBudgetForUser(ctx, 1, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent budget", func(t *testing.T) {
		mockRepo.On("GetBudgetByID", ctx, uint(9999)).Return(nil, errors.New("budget not found"))

		err := service.DeleteBudgetForUser(ctx, 1, 9999)
		assert.Error(t, err)
		assert.Equal(t, "budget not found", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete another user's budget", func(t *testing.T) {
		mockRepo.On("GetBudgetByID", ctx, uint(2)).Return(&models.Budget{ID: 2, UserID: 99}, nil)

		err := service.DeleteBudgetForUser(ctx, 1, 2)
		assert.Error(t, err)
		assert.Equal(t, "budget not found", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindNotFound))
		mockRepo.AssertNotCalled(t, "DeleteBudget", uint(2))
	})
}
