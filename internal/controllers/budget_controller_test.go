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

// MockBudgetService is a mock implementation of BudgetService
type MockBudgetService struct {
	mock.Mock
}

func (m *MockBudgetService) CreateBudget(budget *models.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudgetsByUser(userID uint) ([]models.Budget, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Budget), args.Error(1)
}

func (m *MockBudgetService) DeleteBudget(budgetID uint) error {
	args := m.Called(budgetID)
	return args.Error(0)
}

func (m *MockBudgetService) UpdateBudget(budget *models.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func TestCreateBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		budget := models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      1000,
			StartDate:  time.Now().UTC(),
			EndDate:    time.Now().AddDate(0, 1, 0).UTC(),
		}

		mockService.On("CreateBudget", mock.MatchedBy(func(b *models.Budget) bool {
			return b.UserID == budget.UserID &&
				b.CategoryID == budget.CategoryID &&
				b.Limit == budget.Limit &&
				b.StartDate.Equal(budget.StartDate) &&
				b.EndDate.Equal(budget.EndDate)
		})).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(budget)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Budget created")
	})
	t.Run("Invalid Limit", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		budget := models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      0, // Invalid
			StartDate:  time.Now().UTC(),
			EndDate:    time.Now().AddDate(0, 1, 0).UTC(),
		}

		mockService.On("CreateBudget", &budget).Return(errors.New("budget limit must be greater than zero")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Encode JSON
		jsonBody, _ := json.Marshal(budget)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "budget limit must be greater than zero")
	})

	t.Run("Start Date After End Date", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		budget := models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      500,
			StartDate:  time.Now().AddDate(0, 1, 0).UTC(), // Start after End
			EndDate:    time.Now().UTC(),
		}

		mockService.On("CreateBudget", &budget).Return(errors.New("start date cannot be after end date")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Encode JSON
		jsonBody, _ := json.Marshal(budget)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "start date cannot be after end date")
	})
	t.Run("Invalid JSON Payload", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Send malformed JSON (missing closing brace)
		invalidJSON := `{"UserID": 1, "CategoryID": 2, "Limit": 1000, "StartDate": "2025-03-12T01:06:59Z", "EndDate": "2025-04-12T01:06:59Z"`
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the controller method
		controller.CreateBudget(c)

		// Expect a 400 Bad Request response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})
}

func TestGetBudgets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		now := time.Now().UTC()
		expectedBudgets := []models.Budget{
			{UserID: 1, CategoryID: 2, Limit: 1000, StartDate: now, EndDate: now.AddDate(0, 1, 0)},
			{UserID: 1, CategoryID: 3, Limit: 500, StartDate: now, EndDate: now.AddDate(0, 2, 0)},
		}

		mockService.On("GetBudgetsByUser", uint(1)).Return(expectedBudgets, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/budgets", nil)

		controller.GetBudgets(c)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check that both budgets appear in response
		assert.Contains(t, w.Body.String(), `"Limit":1000`)
		assert.Contains(t, w.Body.String(), `"Limit":500`)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("GetBudgetsByUser", uint(1)).Return([]models.Budget(nil), errors.New("DB error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodGet, "/budgets", nil)

		controller.GetBudgets(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve budgets")
	})
}

func TestDeleteBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("DeleteBudget", uint(1)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)

		controller.DeleteBudget(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Budget deleted")
	})

	t.Run("Budget Not Found", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("DeleteBudget", uint(1)).Return(errors.New("not found")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)

		controller.DeleteBudget(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Budget not found")
	})
}
