package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository defines required repository methods
type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint) error
}

// CategoryRepository handles DB operations for categories
type GormCategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository initializes a new GormCategoryRepository
func NewCategoryRepository(db *gorm.DB) *GormCategoryRepository {
	return &GormCategoryRepository{db: db}
}

// CreateCategory inserts a new category into the database
func (r *GormCategoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetCategoryByID retrieves a category by its ID
func (r *GormCategoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories fetches all categories
func (r *GormCategoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

// UpdateCategory updates an existing category
func (r *GormCategoryRepository) UpdateCategory(category *models.Category) error {
	return r.db.Save(category).Error
}

// DeleteCategory removes a category from the database
func (r *GormCategoryRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
