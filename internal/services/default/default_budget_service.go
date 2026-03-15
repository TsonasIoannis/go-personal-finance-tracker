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

	if budget.StartDate.After(budget.EndDate) {
		return errors.New("start date cannot be after end date")
	}

	return s.budgetRepo.CreateBudget(budget)
}

// UpdateBudget modifies an existing budget
func (s *DefaultBudgetService) UpdateBudget(budget *models.Budget) error {
	if budget.Limit <= 0 {
		return errors.New("budget limit must be positive")
	}

	if budget.StartDate.After(budget.EndDate) {
		return errors.New("start date cannot be after end date")
	}

	return s.budgetRepo.UpdateBudget(budget)
}

// GetBudgetsByUser retrieves budgets for a user
func (s *DefaultBudgetService) GetBudgetsByUser(userID uint) ([]models.Budget, error) {
	return s.budgetRepo.GetBudgetsByUserID(userID)
}

// DeleteBudgetForUser removes a budget that belongs to the authenticated user.
func (s *DefaultBudgetService) DeleteBudgetForUser(userID, budgetID uint) error {
	budget, err := s.budgetRepo.GetBudgetByID(budgetID)
	if err != nil {
		return errors.New("budget not found")
	}

	if budget.UserID != userID {
		return errors.New("budget not found")
	}

	return s.budgetRepo.DeleteBudget(budgetID)
}
