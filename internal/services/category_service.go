package services

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

type CategoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(categoryRepo CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory adds a new category
func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.CreateCategory(category)
}

// GetCategories retrieves all categories
func (s *CategoryService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAllCategories()
}

// UpdateCategory modifies an existing category
func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.categoryRepo.UpdateCategory(category)
}

// DeleteCategory removes a category
func (s *CategoryService) DeleteCategory(categoryID uint) error {
	return s.categoryRepo.DeleteCategory(categoryID)
}
