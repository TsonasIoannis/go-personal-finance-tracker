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
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBudgetService is a mock implementation of BudgetService
type MockBudgetService struct {
	mock.Mock
}

func (m *MockBudgetService) CreateBudget(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudgetsByUser(ctx context.Context, userID uint) ([]models.Budget, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Budget), args.Error(1)
}

func (m *MockBudgetService) GetBudgetsPageByUser(ctx context.Context, userID uint, params pagination.Params) ([]models.Budget, int64, error) {
	args := m.Called(ctx, userID, params)
	return args.Get(0).([]models.Budget), args.Get(1).(int64), args.Error(2)
}

func (m *MockBudgetService) DeleteBudgetForUser(ctx context.Context, userID, budgetID uint) error {
	args := m.Called(ctx, userID, budgetID)
	return args.Error(0)
}

func (m *MockBudgetService) UpdateBudget(ctx context.Context, budget *models.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func TestCreateBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		startDate := time.Now().UTC().Round(time.Second)
		endDate := startDate.AddDate(0, 1, 0)
		budgetPayload := map[string]interface{}{
			"category_id": 2,
			"limit":       1000.0,
			"start_date":  startDate.Format(time.RFC3339),
			"end_date":    endDate.Format(time.RFC3339),
		}

		mockService.On("CreateBudget", mock.Anything, mock.MatchedBy(func(b *models.Budget) bool {
			return b.UserID == 1 &&
				b.CategoryID == 2 &&
				b.Limit == 1000 &&
				b.StartDate.Equal(startDate) &&
				b.EndDate.Equal(endDate)
		})).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(budgetPayload)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Budget created")
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		now := time.Now().UTC()
		budgetPayload := map[string]interface{}{
			"category_id": 2,
			"limit":       0.0,
			"start_date":  now.Format(time.RFC3339),
			"end_date":    now.AddDate(0, 1, 0).Format(time.RFC3339),
		}

		mockService.On("CreateBudget", mock.Anything, mock.AnythingOfType("*models.Budget")).
			Return(apperrors.Validation("invalid_budget_limit", "budget limit must be greater than zero")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(budgetPayload)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_budget_limit"`)
		assert.Contains(t, w.Body.String(), "budget limit must be greater than zero")
	})

	t.Run("Start Date After End Date", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		startDate := time.Now().AddDate(0, 1, 0).UTC()
		endDate := time.Now().UTC()
		budgetPayload := map[string]interface{}{
			"category_id": 2,
			"limit":       500.0,
			"start_date":  startDate.Format(time.RFC3339),
			"end_date":    endDate.Format(time.RFC3339),
		}

		mockService.On("CreateBudget", mock.Anything, mock.AnythingOfType("*models.Budget")).
			Return(apperrors.Validation("invalid_budget_date_range", "start date cannot be after end date")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		jsonBody, err := json.Marshal(budgetPayload)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_budget_date_range"`)
		assert.Contains(t, w.Body.String(), "start date cannot be after end date")
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		invalidJSON := `{"category_id": 2, "limit": 1000, "start_date": "2025-03-12T01:06:59Z", "end_date": "2025-04-12T01:06:59Z"`
		c.Request = httptest.NewRequest(http.MethodPost, "/budgets", bytes.NewBufferString(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.CreateBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_request"`)
		assert.Contains(t, w.Body.String(), "invalid request payload")
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

		mockService.On("GetBudgetsByUser", mock.Anything, uint(1)).Return(expectedBudgets, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/budgets", nil)

		controller.GetBudgets(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"limit":1000`)
		assert.Contains(t, w.Body.String(), `"limit":500`)
		assert.Contains(t, w.Body.String(), `"user_id":1`)
		assert.NotContains(t, w.Body.String(), `"UserID"`)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("GetBudgetsByUser", mock.Anything, uint(1)).
			Return([]models.Budget(nil), apperrors.Internal("budgets_fetch_failed", "failed to retrieve budgets", nil)).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/budgets", nil)

		controller.GetBudgets(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"budgets_fetch_failed"`)
		assert.Contains(t, w.Body.String(), "failed to retrieve budgets")
	})
}

func TestGetBudgetsPage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		now := time.Now().UTC()
		params := pagination.New(1, 2)
		expectedBudgets := []models.Budget{
			{UserID: 1, CategoryID: 2, Limit: 1000, StartDate: now, EndDate: now.AddDate(0, 1, 0)},
			{UserID: 1, CategoryID: 3, Limit: 500, StartDate: now, EndDate: now.AddDate(0, 2, 0)},
		}

		mockService.On("GetBudgetsPageByUser", mock.Anything, uint(1), params).Return(expectedBudgets, int64(4), nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/budgets?page=1&page_size=2", nil)

		controller.GetBudgetsPage(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[`)
		assert.Contains(t, w.Body.String(), `"page":1`)
		assert.Contains(t, w.Body.String(), `"page_size":2`)
		assert.Contains(t, w.Body.String(), `"total":4`)
		assert.Contains(t, w.Body.String(), `"total_pages":2`)
	})

	t.Run("Invalid Page Size", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/budgets?page_size=0", nil)

		controller.GetBudgetsPage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_page_size"`)
	})
}

func TestDeleteBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("DeleteBudgetForUser", mock.Anything, uint(1), uint(1)).Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)

		controller.DeleteBudget(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Budget deleted")
	})

	t.Run("Invalid ID", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "oops"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/budgets/oops", nil)

		controller.DeleteBudget(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_budget_id"`)
		assert.Contains(t, w.Body.String(), "invalid budget id")
	})

	t.Run("Budget Not Found", func(t *testing.T) {
		mockService := new(MockBudgetService)
		controller := NewBudgetController(mockService)

		mockService.On("DeleteBudgetForUser", mock.Anything, uint(1), uint(1)).
			Return(apperrors.NotFound("budget_not_found", "budget not found")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/budgets/1", nil)

		controller.DeleteBudget(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"budget_not_found"`)
		assert.Contains(t, w.Body.String(), "budget not found")
	})
}
