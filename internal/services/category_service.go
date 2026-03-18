package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// CategoryService defines the interface for category operations
type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context) ([]models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, categoryID uint) error
}
