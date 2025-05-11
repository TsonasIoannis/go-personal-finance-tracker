package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService implements services.TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) AddTransaction(t *models.Transaction) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionsByUser(userID uint) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionService) DeleteTransaction(transactionID uint) error {
	args := m.Called(transactionID)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		transaction := models.Transaction{
			UserID:     1,
			Amount:     50.0,
			CategoryID: 2,
			Type:       "expense",
			Date:       time.Now().UTC(),
			Note:       "Lunch",
		}

		mockService.On("AddTransaction", mock.MatchedBy(func(t *models.Transaction) bool {
			return t.UserID == transaction.UserID &&
				t.Amount == transaction.Amount &&
				t.CategoryID == transaction.CategoryID &&
				t.Type == transaction.Type &&
				t.Note == transaction.Note
		})).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(transaction)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateTransaction(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction added")
	})

	t.Run("Exceeds Budget", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		transaction := models.Transaction{
			UserID:     1,
			Amount:     5000.0,
			CategoryID: 2,
			Type:       "expense",
			Date:       time.Now().UTC(),
		}

		mockService.On("AddTransaction", &transaction).Return(errors.New("transaction exceeds budget limit")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(transaction)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "transaction exceeds budget limit")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		badJSON := `{"UserID":1,"Amount":50.0,"CategoryID":2`
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(badJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})
}

func TestGetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		now := time.Now().UTC()
		transactions := []models.Transaction{
			{UserID: 1, Amount: 20.0, CategoryID: 1, Type: "expense", Date: now, Note: "Groceries"},
			{UserID: 1, Amount: 100.0, CategoryID: 2, Type: "income", Date: now, Note: "Salary"},
		}

		mockService.On("GetTransactionsByUser", uint(1)).Return(transactions, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/transactions", nil)

		controller.GetTransactions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"Amount":20`)
		assert.Contains(t, w.Body.String(), `"Amount":100`)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("GetTransactionsByUser", uint(1)).Return([]models.Transaction(nil), errors.New("fetch failed")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/transactions", nil)

		controller.GetTransactions(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve transactions")
	})
}

func TestDeleteTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("DeleteTransaction", uint(1)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction deleted")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("DeleteTransaction", uint(1)).Return(errors.New("not found")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction not found")
	})
}
