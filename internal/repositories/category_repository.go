package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// CategoryRepository defines required repository methods
type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
}
