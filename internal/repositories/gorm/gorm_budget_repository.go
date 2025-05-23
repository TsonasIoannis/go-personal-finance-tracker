package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
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

// GormBudgetRepository handles DB operations for budgets using GORM
type GormBudgetRepository struct {
	db *gorm.DB
}

// NewGormBudgetRepository initializes a new GormBudgetRepository
func NewGormBudgetRepository(db *gorm.DB) *GormBudgetRepository {
	return &GormBudgetRepository{db: db}
}

// CreateBudget inserts a new budget into the database
func (r *GormBudgetRepository) CreateBudget(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

// GetBudgetByID retrieves a budget by its ID
func (r *GormBudgetRepository) GetBudgetByID(id uint) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.First(&budget, id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// GetBudgetsByUserID fetches all budgets for a specific user
func (r *GormBudgetRepository) GetBudgetsByUserID(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

// UpdateBudget updates an existing budget
func (r *GormBudgetRepository) UpdateBudget(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

// DeleteBudget removes a budget from the database
func (r *GormBudgetRepository) DeleteBudget(id uint) error {
	return r.db.Delete(&models.Budget{}, id).Error
}
