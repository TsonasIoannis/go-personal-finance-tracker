package services

import (
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

func (s *DefaultCategoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.CreateCategory(category)
}

func (s *DefaultCategoryService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *DefaultCategoryService) UpdateCategory(category *models.Category) error {
	return s.categoryRepo.UpdateCategory(category)
}

func (s *DefaultCategoryService) DeleteCategory(categoryID uint) error {
	return s.categoryRepo.DeleteCategory(categoryID)
}
