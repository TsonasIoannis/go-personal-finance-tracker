package controllers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

type BudgetController struct {
	budgetService *services.BudgetService
}

func NewBudgetController(budgetService *services.BudgetService) *BudgetController {
	return &BudgetController{budgetService: budgetService}
}

// CreateBudget adds a new budget
func (bc *BudgetController) CreateBudget(c *gin.Context) {
	var budget models.Budget

	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := bc.budgetService.CreateBudget(&budget)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Budget created"})
}

// GetBudgets fetches all budgets for a user
func (bc *BudgetController) GetBudgets(c *gin.Context) {
	userID := uint(1) // Extract from JWT in real scenario

	budgets, err := bc.budgetService.GetBudgetsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve budgets"})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// DeleteBudget removes a budget
func (bc *BudgetController) DeleteBudget(c *gin.Context) {
	budgetID := uint(1) // Extract from URL params in real implementation

	err := bc.budgetService.DeleteBudget(budgetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted"})
}
