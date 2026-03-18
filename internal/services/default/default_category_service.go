package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
)

// DefaultCategoryService is the default implementation of CategoryService
type DefaultCategoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService initializes a new DefaultCategoryService
func NewCategoryService(categoryRepo repositories.CategoryRepository) DefaultCategoryService {
	return DefaultCategoryService{categoryRepo: categoryRepo}
}

func (s *DefaultCategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	return s.categoryRepo.CreateCategory(ctx, category)
}

func (s *DefaultCategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepo.GetAllCategories(ctx)
}

func (s *DefaultCategoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	return s.categoryRepo.UpdateCategory(ctx, category)
}

func (s *DefaultCategoryService) DeleteCategory(ctx context.Context, categoryID uint) error {
	return s.categoryRepo.DeleteCategory(ctx, categoryID)
}
