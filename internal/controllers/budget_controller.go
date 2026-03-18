package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	budgetService services.BudgetService
}

// NewBudgetController initializes a BudgetController with an interface dependency
func NewBudgetController(budgetService services.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

type createBudgetRequest struct {
	CategoryID uint      `json:"category_id" binding:"required"`
	Limit      float64   `json:"limit"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
}

// CreateBudget adds a new budget
// @Summary Create a budget
// @Description Create a budget for the authenticated user.
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body createBudgetRequest true "Budget payload"
// @Success 201 {object} messageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/budgets [post]
func (bc *BudgetController) CreateBudget(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req createBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))
		return
	}

	budget := models.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Limit:      req.Limit,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}

	if err := bc.budgetService.CreateBudget(ctx, &budget); err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Budget created"})
}

// GetBudgets fetches all budgets for a user
func (bc *BudgetController) GetBudgets(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	budgets, err := bc.budgetService.GetBudgetsByUser(ctx, userID)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, newBudgetResponses(budgets))
}

// GetBudgetsPage fetches a paginated budget list for a user.
// @Summary List budgets
// @Description List the authenticated user's budgets with pagination.
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" minimum(1)
// @Param page_size query int false "Items per page" minimum(1) maximum(100)
// @Success 200 {object} budgetPageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/budgets [get]
func (bc *BudgetController) GetBudgetsPage(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	params, err := parsePaginationParams(c)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	budgets, total, err := bc.budgetService.GetBudgetsPageByUser(ctx, userID, params)
	if err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, paginatedResponse[budgetResponse]{
		Data:       newBudgetResponses(budgets),
		Pagination: newPaginationResponse(params, total),
	})
}

// DeleteBudget removes a budget
// @Summary Delete a budget
// @Description Delete one of the authenticated user's budgets.
// @Tags budgets
// @Produce json
// @Security BearerAuth
// @Param id path int true "Budget ID" minimum(1)
// @Success 200 {object} messageResponse
// @Failure 400 {object} httpapi.ErrorResponse
// @Failure 401 {object} httpapi.ErrorResponse
// @Failure 404 {object} httpapi.ErrorResponse
// @Failure 500 {object} httpapi.ErrorResponse
// @Router /api/v1/budgets/{id} [delete]
func (bc *BudgetController) DeleteBudget(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.WriteError(c, apperrors.Validation("invalid_budget_id", "invalid budget id"))
		return
	}

	if err := bc.budgetService.DeleteBudgetForUser(ctx, userID, uint(budgetID)); err != nil {
		httpapi.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted"})
}
