package services

import (
	"errors"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
)

type DefaultBudgetService struct {
	budgetRepo repositories.BudgetRepository
}

func NewBudgetService(budgetRepo repositories.BudgetRepository) *DefaultBudgetService {
	return &DefaultBudgetService{budgetRepo: budgetRepo}
}

// CreateBudget validates and adds a budget
func (s *DefaultBudgetService) CreateBudget(budget *models.Budget) error {
	if budget.Limit <= 0 {
		return errors.New("budget limit must be greater than zero")
	}
	return s.budgetRepo.CreateBudget(budget)
}

// UpdateBudget modifies an existing budget
func (s *DefaultBudgetService) UpdateBudget(budget *models.Budget) error {
	if budget.Limit <= 0 {
		return errors.New("budget limit must be positive")
	}
	return s.budgetRepo.UpdateBudget(budget)
}

// GetBudgetsByUser retrieves budgets for a user
func (s *DefaultBudgetService) GetBudgetsByUser(userID uint) ([]models.Budget, error) {
	return s.budgetRepo.GetBudgetsByUserID(userID)
}

// DeleteBudget removes a budget
func (s *DefaultBudgetService) DeleteBudget(budgetID uint) error {
	return s.budgetRepo.DeleteBudget(budgetID)
}
