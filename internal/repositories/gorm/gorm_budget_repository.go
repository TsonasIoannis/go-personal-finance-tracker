package repositories

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
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

// GormBudgetRepository handles DB operations for budgets using GORM
type GormBudgetRepository struct {
	db *gorm.DB
}

// NewGormBudgetRepository initializes a new GormBudgetRepository
func NewGormBudgetRepository(db *gorm.DB) *GormBudgetRepository {
	return &GormBudgetRepository{db: db}
}

// CreateBudget inserts a new budget into the database
func (r *GormBudgetRepository) CreateBudget(ctx context.Context, budget *models.Budget) error {
	return r.db.WithContext(ctx).Create(budget).Error
}

// GetBudgetByID retrieves a budget by its ID
func (r *GormBudgetRepository) GetBudgetByID(ctx context.Context, id uint) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.WithContext(ctx).First(&budget, id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// GetBudgetsByUserID fetches all budgets for a specific user
func (r *GormBudgetRepository) GetBudgetsByUserID(ctx context.Context, userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

// UpdateBudget updates an existing budget
func (r *GormBudgetRepository) UpdateBudget(ctx context.Context, budget *models.Budget) error {
	return r.db.WithContext(ctx).Save(budget).Error
}

// DeleteBudget removes a budget from the database
func (r *GormBudgetRepository) DeleteBudget(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Budget{}, id).Error
}
