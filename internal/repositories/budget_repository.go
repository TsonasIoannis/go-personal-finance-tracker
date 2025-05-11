package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// BudgetRepository defines the required repository methods
// This ensures other services can use different implementations if needed.
type BudgetRepository interface {
	CreateBudget(budget *models.Budget) error
	GetBudgetByID(id uint) (*models.Budget, error)
	GetBudgetsByUserID(userID uint) ([]models.Budget, error)
	UpdateBudget(budget *models.Budget) error
	DeleteBudget(id uint) error
}
