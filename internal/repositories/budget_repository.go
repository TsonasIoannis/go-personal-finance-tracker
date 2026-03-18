package repositories

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// BudgetRepository defines the required repository methods
// This ensures other services can use different implementations if needed.
type BudgetRepository interface {
	CreateBudget(ctx context.Context, budget *models.Budget) error
	GetBudgetByID(ctx context.Context, id uint) (*models.Budget, error)
	GetBudgetsByUserID(ctx context.Context, userID uint) ([]models.Budget, error)
	UpdateBudget(ctx context.Context, budget *models.Budget) error
	DeleteBudget(ctx context.Context, id uint) error
}
