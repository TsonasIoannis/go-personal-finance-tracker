package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
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
func (s *DefaultBudgetService) CreateBudget(ctx context.Context, budget *models.Budget) error {
	if budget.Limit <= 0 {
		return apperrors.Validation("invalid_budget_limit", "budget limit must be greater than zero")
	}

	if budget.StartDate.After(budget.EndDate) {
		return apperrors.Validation("invalid_budget_date_range", "start date cannot be after end date")
	}

	if err := s.budgetRepo.CreateBudget(ctx, budget); err != nil {
		return apperrors.Internal("budget_create_failed", "failed to create budget", err)
	}

	return nil
}

// UpdateBudget modifies an existing budget
func (s *DefaultBudgetService) UpdateBudget(ctx context.Context, budget *models.Budget) error {
	if budget.Limit <= 0 {
		return apperrors.Validation("invalid_budget_limit", "budget limit must be positive")
	}

	if budget.StartDate.After(budget.EndDate) {
		return apperrors.Validation("invalid_budget_date_range", "start date cannot be after end date")
	}

	if err := s.budgetRepo.UpdateBudget(ctx, budget); err != nil {
		return apperrors.Internal("budget_update_failed", "failed to update budget", err)
	}

	return nil
}

// GetBudgetsByUser retrieves budgets for a user
func (s *DefaultBudgetService) GetBudgetsByUser(ctx context.Context, userID uint) ([]models.Budget, error) {
	budgets, err := s.budgetRepo.GetBudgetsByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal("budgets_fetch_failed", "failed to retrieve budgets", err)
	}

	return budgets, nil
}

// DeleteBudgetForUser removes a budget that belongs to the authenticated user.
func (s *DefaultBudgetService) DeleteBudgetForUser(ctx context.Context, userID, budgetID uint) error {
	budget, err := s.budgetRepo.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return apperrors.NotFound("budget_not_found", "budget not found")
	}

	if budget.UserID != userID {
		return apperrors.NotFound("budget_not_found", "budget not found")
	}

	if err := s.budgetRepo.DeleteBudget(ctx, budgetID); err != nil {
		return apperrors.Internal("budget_delete_failed", "failed to delete budget", err)
	}

	return nil
}
