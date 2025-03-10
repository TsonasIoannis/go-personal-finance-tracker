package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// BudgetRepository handles DB operations for budgets
type BudgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository initializes a new BudgetRepository
func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// CreateBudget inserts a new budget into the database
func (r *BudgetRepository) CreateBudget(budget *models.Budget) error {
	return r.db.Create(budget).Error
}

// GetBudgetByID retrieves a budget by its ID
func (r *BudgetRepository) GetBudgetByID(id uint) (*models.Budget, error) {
	var budget models.Budget
	err := r.db.First(&budget, id).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// GetBudgetsByUserID fetches all budgets for a specific user
func (r *BudgetRepository) GetBudgetsByUserID(userID uint) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.Where("user_id = ?", userID).Find(&budgets).Error
	return budgets, err
}

// UpdateBudget updates an existing budget
func (r *BudgetRepository) UpdateBudget(budget *models.Budget) error {
	return r.db.Save(budget).Error
}

// DeleteBudget removes a budget from the database
func (r *BudgetRepository) DeleteBudget(id uint) error {
	return r.db.Delete(&models.Budget{}, id).Error
}
