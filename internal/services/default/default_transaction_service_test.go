package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository implements the TransactionRepository interface
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionByID(ctx context.Context, id uint) (*models.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uint) ([]models.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsPageByUserID(ctx context.Context, userID uint, params pagination.Params) ([]models.Transaction, int64, error) {
	args := m.Called(ctx, userID, params)
	return args.Get(0).([]models.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) DeleteTransaction(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAddTransaction(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	mockBudgetRepo := new(MockBudgetRepository)
	service := NewTransactionService(mockTransactionRepo, mockBudgetRepo)
	ctx := context.Background()

	t.Run("Create valid transaction", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil // Reset expectations
		mockBudgetRepo.ExpectedCalls = nil      // Reset expectations

		transaction := &models.Transaction{
			UserID:     1,
			Type:       "income",
			Amount:     500.00,
			CategoryID: 2,
			Date:       time.Now(),
		}

		mockBudgetRepo.On("GetBudgetsByUserID", ctx, uint(1)).Return([]models.Budget{}, nil)
		mockTransactionRepo.On("CreateTransaction", ctx, transaction).Return(nil)

		err := service.AddTransaction(ctx, transaction)
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Fail when transaction exceeds budget", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil // Reset expectations
		mockBudgetRepo.ExpectedCalls = nil      // Reset expectations

		transaction := &models.Transaction{
			UserID:     1,
			Type:       "expense",
			Amount:     1200.00, // Over budget
			CategoryID: 2,
			Date:       time.Now(),
		}

		budgets := []models.Budget{
			{UserID: 1, CategoryID: 2, Limit: 1000.00},
		}

		mockBudgetRepo.On("GetBudgetsByUserID", ctx, uint(1)).Return(budgets, nil)

		err := service.AddTransaction(ctx, transaction)
		assert.Error(t, err)
		assert.Equal(t, "transaction exceeds budget limit", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindValidation))
		mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
	})
	t.Run("Fail when budget retrieval fails", func(t *testing.T) {
		mockTransactionRepo.ExpectedCalls = nil // Reset expectations
		mockBudgetRepo.ExpectedCalls = nil      // Reset expectations

		transaction := &models.Transaction{
			UserID:     1,
			Type:       "expense",
			Amount:     500.00,
			CategoryID: 2,
			Date:       time.Now(),
		}

		// Simulate an error when retrieving budgets
		mockBudgetRepo.On("GetBudgetsByUserID", ctx, uint(1)).Return([]models.Budget{}, errors.New("database error"))

		err := service.AddTransaction(ctx, transaction)

		// Ensure the error is returned
		assert.Error(t, err)
		assert.Equal(t, "failed to validate transaction budget", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindInternal))

		// Ensure transaction was NOT created
		mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
	})

}

func TestGetTransactionsByUser(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockTransactionRepo, nil)
	ctx := context.Background()

	t.Run("Retrieve transactions for user", func(t *testing.T) {
		transactions := []models.Transaction{
			{ID: 1, UserID: 1, Type: "expense", Amount: 50, CategoryID: 1},
			{ID: 2, UserID: 1, Type: "income", Amount: 200, CategoryID: 2},
		}

		mockTransactionRepo.On("GetTransactionsByUserID", ctx, uint(1)).Return(transactions, nil)

		result, err := service.GetTransactionsByUser(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetTransactionsPageByUser(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockTransactionRepo, nil)
	ctx := context.Background()
	params := pagination.New(2, 1)

	t.Run("Retrieve paginated transactions for user", func(t *testing.T) {
		transactions := []models.Transaction{
			{ID: 2, UserID: 1, Type: "income", Amount: 200, CategoryID: 2},
		}

		mockTransactionRepo.On("GetTransactionsPageByUserID", ctx, uint(1), params).Return(transactions, int64(3), nil)

		result, total, err := service.GetTransactionsPageByUser(ctx, 1, params)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(3), total)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestDeleteTransaction(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockTransactionRepo, nil)
	ctx := context.Background()

	t.Run("Delete existing transaction", func(t *testing.T) {
		mockTransactionRepo.On("GetTransactionByID", ctx, uint(1)).Return(&models.Transaction{ID: 1, UserID: 1}, nil)
		mockTransactionRepo.On("DeleteTransaction", ctx, uint(1)).Return(nil)

		err := service.DeleteTransactionForUser(ctx, 1, 1)
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent transaction", func(t *testing.T) {
		mockTransactionRepo.On("GetTransactionByID", ctx, uint(9999)).Return(nil, errors.New("transaction not found"))

		err := service.DeleteTransactionForUser(ctx, 1, 9999)
		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindNotFound))
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete another user's transaction", func(t *testing.T) {
		mockTransactionRepo.On("GetTransactionByID", ctx, uint(2)).Return(&models.Transaction{ID: 2, UserID: 99}, nil)

		err := service.DeleteTransactionForUser(ctx, 1, 2)
		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())
		assert.True(t, isAppErrorKind(err, apperrors.KindNotFound))
		mockTransactionRepo.AssertNotCalled(t, "DeleteTransaction", uint(2))
	})
}
