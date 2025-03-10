package services

import (
	"errors"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// BudgetRepository defines the required repository methods
type BudgetRepository interface {
	CreateBudget(budget *models.Budget) error
	GetBudgetByID(id uint) (*models.Budget, error)
	GetBudgetsByUserID(userID uint) ([]models.Budget, error)
	UpdateBudget(budget *models.Budget) error
	DeleteBudget(id uint) error
}

type BudgetService struct {
	budgetRepo BudgetRepository
}

func NewBudgetService(budgetRepo BudgetRepository) *BudgetService {
	return &BudgetService{budgetRepo: budgetRepo}
}

// CreateBudget validates and adds a budget
func (s *BudgetService) CreateBudget(budget *models.Budget) error {
	if budget.Limit <= 0 {
		return errors.New("budget limit must be greater than zero")
	}
	return s.budgetRepo.CreateBudget(budget)
}

// UpdateBudget modifies an existing budget
func (s *BudgetService) UpdateBudget(budget *models.Budget) error {
	if budget.Limit <= 0 {
		return errors.New("budget limit must be positive")
	}
	return s.budgetRepo.UpdateBudget(budget)
}

// GetBudgetsByUser retrieves budgets for a user
func (s *BudgetService) GetBudgetsByUser(userID uint) ([]models.Budget, error) {
	return s.budgetRepo.GetBudgetsByUserID(userID)
}

// DeleteBudget removes a budget
func (s *BudgetService) DeleteBudget(budgetID uint) error {
	return s.budgetRepo.DeleteBudget(budgetID)
}
