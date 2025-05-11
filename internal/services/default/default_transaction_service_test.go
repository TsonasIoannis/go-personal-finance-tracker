package services

import (
	"errors"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository implements the TransactionRepository interface
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionByID(id uint) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) DeleteTransaction(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAddTransaction(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	mockBudgetRepo := new(MockBudgetRepository)
	service := NewTransactionService(mockTransactionRepo, mockBudgetRepo)

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

		mockBudgetRepo.On("GetBudgetsByUserID", uint(1)).Return([]models.Budget{}, nil)
		mockTransactionRepo.On("CreateTransaction", transaction).Return(nil)

		err := service.AddTransaction(transaction)
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

		mockBudgetRepo.On("GetBudgetsByUserID", uint(1)).Return(budgets, nil)

		err := service.AddTransaction(transaction)
		assert.Error(t, err)
		assert.Equal(t, "transaction exceeds budget limit", err.Error())
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
		mockBudgetRepo.On("GetBudgetsByUserID", uint(1)).Return([]models.Budget{}, errors.New("database error"))

		err := service.AddTransaction(transaction)

		// Ensure the error is returned
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		// Ensure transaction was NOT created
		mockTransactionRepo.AssertNotCalled(t, "CreateTransaction")
	})

}

func TestGetTransactionsByUser(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockTransactionRepo, nil)

	t.Run("Retrieve transactions for user", func(t *testing.T) {
		transactions := []models.Transaction{
			{ID: 1, UserID: 1, Type: "expense", Amount: 50, CategoryID: 1},
			{ID: 2, UserID: 1, Type: "income", Amount: 200, CategoryID: 2},
		}

		mockTransactionRepo.On("GetTransactionsByUserID", uint(1)).Return(transactions, nil)

		result, err := service.GetTransactionsByUser(1)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestDeleteTransaction(t *testing.T) {
	mockTransactionRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockTransactionRepo, nil)

	t.Run("Delete existing transaction", func(t *testing.T) {
		mockTransactionRepo.On("DeleteTransaction", uint(1)).Return(nil)

		err := service.DeleteTransaction(1)
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
	})

	t.Run("Fail to delete non-existent transaction", func(t *testing.T) {
		mockTransactionRepo.On("DeleteTransaction", uint(9999)).Return(errors.New("transaction not found"))

		err := service.DeleteTransaction(9999)
		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
}
