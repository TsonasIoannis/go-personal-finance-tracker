package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// CategoryService defines the interface for category operations
type CategoryService interface {
	CreateCategory(category *models.Category) error
	GetCategories() ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(categoryID uint) error
}
