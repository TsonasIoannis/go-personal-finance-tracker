package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/pagination"
)

type BudgetService interface {
	CreateBudget(ctx context.Context, budget *models.Budget) error
	UpdateBudget(ctx context.Context, budget *models.Budget) error
	GetBudgetsByUser(ctx context.Context, userID uint) ([]models.Budget, error)
	GetBudgetsPageByUser(ctx context.Context, userID uint, params pagination.Params) ([]models.Budget, int64, error)
	DeleteBudgetForUser(ctx context.Context, userID, budgetID uint) error
}
