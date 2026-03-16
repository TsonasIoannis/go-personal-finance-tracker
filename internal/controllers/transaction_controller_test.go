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

func (m *MockTransactionService) DeleteTransactionForUser(userID, transactionID uint) error {
	args := m.Called(userID, transactionID)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		now := time.Now().UTC().Round(time.Second)
		payload := map[string]interface{}{
			"amount":      50.0,
			"category_id": 2,
			"type":        "expense",
			"date":        now.Format(time.RFC3339),
			"note":        "Lunch",
		}

		mockService.On("AddTransaction", mock.MatchedBy(func(t *models.Transaction) bool {
			return t.UserID == 1 &&
				t.Amount == 50.0 &&
				t.CategoryID == 2 &&
				t.Type == "expense" &&
				t.Note == "Lunch" &&
				t.Date.Equal(now)
		})).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(payload)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateTransaction(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction added")
	})

	t.Run("Exceeds Budget", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		now := time.Now().UTC()
		payload := map[string]interface{}{
			"amount":      5000.0,
			"category_id": 2,
			"type":        "expense",
			"date":        now.Format(time.RFC3339),
		}

		mockService.On("AddTransaction", mock.AnythingOfType("*models.Transaction")).Return(errors.New("transaction exceeds budget limit")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(payload)
		assert.NoError(t, err)
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
		c.Set("userID", uint(1))

		badJSON := `{"amount":50.0,"category_id":2`
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
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/transactions", nil)

		controller.GetTransactions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"amount":20`)
		assert.Contains(t, w.Body.String(), `"amount":100`)
		assert.Contains(t, w.Body.String(), `"user_id":1`)
		assert.NotContains(t, w.Body.String(), `"UserID"`)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("GetTransactionsByUser", uint(1)).Return([]models.Transaction(nil), errors.New("fetch failed")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
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

		mockService.On("DeleteTransactionForUser", uint(1), uint(1)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction deleted")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/abc", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid transaction id")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("DeleteTransactionForUser", uint(1), uint(1)).Return(errors.New("not found")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Transaction not found")
	})
}
