package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles DB operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository initializes a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// CreateCategory inserts a new category into the database
func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetCategoryByID retrieves a category by its ID
func (r *CategoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories fetches all categories
func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

// DeleteCategory removes a category from the database
func (r *CategoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
