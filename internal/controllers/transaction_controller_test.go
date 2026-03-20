package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/filters"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService implements services.TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) AddTransaction(ctx context.Context, t *models.Transaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionsByUser(ctx context.Context, userID uint, transactionFilters filters.TransactionFilters) ([]models.Transaction, error) {
	args := m.Called(ctx, userID, transactionFilters)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransactionsPageByUser(ctx context.Context, userID uint, params pagination.Params, transactionFilters filters.TransactionFilters) ([]models.Transaction, int64, error) {
	args := m.Called(ctx, userID, params, transactionFilters)
	return args.Get(0).([]models.Transaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockTransactionService) DeleteTransactionForUser(ctx context.Context, userID, transactionID uint) error {
	args := m.Called(ctx, userID, transactionID)
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

		mockService.On("AddTransaction", mock.Anything, mock.MatchedBy(func(t *models.Transaction) bool {
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

		mockService.On("AddTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).
			Return(apperrors.Validation("budget_limit_exceeded", "transaction exceeds budget limit")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(payload)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"budget_limit_exceeded"`)
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
		assert.Contains(t, w.Body.String(), `"code":"invalid_request"`)
		assert.Contains(t, w.Body.String(), "invalid request payload")
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

		mockService.On("GetTransactionsByUser", mock.Anything, uint(1), filters.TransactionFilters{}).Return(transactions, nil).Once()

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

		mockService.On("GetTransactionsByUser", mock.Anything, uint(1), filters.TransactionFilters{}).
			Return([]models.Transaction(nil), apperrors.Internal("transactions_fetch_failed", "failed to retrieve transactions", nil)).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/transactions", nil)

		controller.GetTransactions(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"transactions_fetch_failed"`)
		assert.Contains(t, w.Body.String(), "failed to retrieve transactions")
	})
}

func TestGetTransactionsPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		now := time.Now().UTC()
		params := pagination.New(2, 1)
		transactionFilters := filters.TransactionFilters{Type: "income"}
		transactions := []models.Transaction{
			{UserID: 1, Amount: 100.0, CategoryID: 2, Type: "income", Date: now, Note: "Salary"},
		}

		mockService.On("GetTransactionsPageByUser", mock.Anything, uint(1), params, transactionFilters).Return(transactions, int64(3), nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/transactions?page=2&page_size=1&type=income", nil)

		controller.GetTransactionsPage(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[`)
		assert.Contains(t, w.Body.String(), `"page":2`)
		assert.Contains(t, w.Body.String(), `"page_size":1`)
		assert.Contains(t, w.Body.String(), `"total":3`)
		assert.Contains(t, w.Body.String(), `"total_pages":3`)
	})

	t.Run("Legacy Endpoint Accepts Filters", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		transactionFilters := filters.TransactionFilters{Type: "expense"}
		transactions := []models.Transaction{
			{UserID: 1, Amount: 20.0, CategoryID: 1, Type: "expense", Date: time.Now().UTC(), Note: "Groceries"},
		}

		mockService.On("GetTransactionsByUser", mock.Anything, uint(1), transactionFilters).Return(transactions, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/transactions?type=expense", nil)

		controller.GetTransactions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"type":"expense"`)
	})

	t.Run("Invalid Page", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/transactions?page=0", nil)

		controller.GetTransactionsPage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_page"`)
	})

	t.Run("Invalid Filter", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/transactions?type=transfer", nil)

		controller.GetTransactionsPage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_transaction_type"`)
	})
}

func TestDeleteTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("DeleteTransactionForUser", mock.Anything, uint(1), uint(1)).Return(nil).Once()

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
		assert.Contains(t, w.Body.String(), `"code":"invalid_transaction_id"`)
		assert.Contains(t, w.Body.String(), "invalid transaction id")
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService := new(MockTransactionService)
		controller := NewTransactionController(mockService)

		mockService.On("DeleteTransactionForUser", mock.Anything, uint(1), uint(1)).
			Return(apperrors.NotFound("transaction_not_found", "transaction not found")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)

		controller.DeleteTransaction(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"transaction_not_found"`)
		assert.Contains(t, w.Body.String(), "transaction not found")
	})
}
