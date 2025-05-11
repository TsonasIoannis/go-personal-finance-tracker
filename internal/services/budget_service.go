package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

type BudgetService interface {
	CreateBudget(budget *models.Budget) error
	UpdateBudget(budget *models.Budget) error
	GetBudgetsByUser(userID uint) ([]models.Budget, error)
	DeleteBudget(budgetID uint) error
}
