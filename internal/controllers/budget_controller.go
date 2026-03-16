package controllers

import (
	"net/http"
	"strconv"
	"time"

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
func (bc *BudgetController) CreateBudget(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req createBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	budget := models.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Limit:      req.Limit,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	}

	if err := bc.budgetService.CreateBudget(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Budget created"})
}

// GetBudgets fetches all budgets for a user
func (bc *BudgetController) GetBudgets(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	budgets, err := bc.budgetService.GetBudgetsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budgets"})
		return
	}

	c.JSON(http.StatusOK, newBudgetResponses(budgets))
}

// DeleteBudget removes a budget
func (bc *BudgetController) DeleteBudget(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	budgetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget id"})
		return
	}

	if err := bc.budgetService.DeleteBudgetForUser(userID, uint(budgetID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted"})
}
