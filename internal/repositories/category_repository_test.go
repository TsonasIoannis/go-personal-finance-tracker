package repositories

import (
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupCategoryTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.Category{})
	return db
}

// TestCreateCategory tests different cases for creating a category.
func TestCreateCategory(t *testing.T) {
	db := setupCategoryTestDB()
	repo := NewCategoryRepository(db)

	t.Run("Create valid category", func(t *testing.T) {
		category := &models.Category{Name: "Groceries", Description: "Food and drinks"}
		err := repo.CreateCategory(category)
		assert.NoError(t, err)

		var retrievedCategory models.Category
		err = db.First(&retrievedCategory, category.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, category.Name, retrievedCategory.Name)
	})

	t.Run("Create duplicate category name", func(t *testing.T) {
		category1 := &models.Category{Name: "Health", Description: "Medical expenses"}
		category2 := &models.Category{Name: "Health", Description: "Duplicate category"}

		err := repo.CreateCategory(category1)
		assert.NoError(t, err)

		err = repo.CreateCategory(category2) // Should fail due to unique constraint
		assert.Error(t, err)
	})
}

// TestGetCategoryByID tests retrieving a category by ID.
func TestGetCategoryByID(t *testing.T) {
	db := setupCategoryTestDB()
	repo := NewCategoryRepository(db)

	t.Run("Retrieve existing category", func(t *testing.T) {
		category := &models.Category{Name: "Utilities", Description: "Electricity, water, gas"}
		err := repo.CreateCategory(category)
		assert.NoError(t, err)

		foundCategory, err := repo.GetCategoryByID(category.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundCategory)
		assert.Equal(t, category.Name, foundCategory.Name)
	})

	t.Run("Retrieve non-existent category", func(t *testing.T) {
		category, err := repo.GetCategoryByID(9999) // Non-existent ID
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestGetAllCategories tests retrieving all categories.
func TestGetAllCategories(t *testing.T) {
	db := setupCategoryTestDB()
	repo := NewCategoryRepository(db)

	t.Run("Retrieve multiple categories", func(t *testing.T) {
		category1 := &models.Category{Name: "Education", Description: "School and learning"}
		category2 := &models.Category{Name: "Entertainment", Description: "Movies, concerts"}
		err1 := repo.CreateCategory(category1)
		assert.NoError(t, err1)
		err2 := repo.CreateCategory(category2)
		assert.NoError(t, err2)

		categories, err := repo.GetAllCategories()
		assert.NoError(t, err)
		assert.Len(t, categories, 2)
	})

	t.Run("Retrieve categories when none exist", func(t *testing.T) {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Category{}) // Ensure no data
		categories, err := repo.GetAllCategories()

		assert.NoError(t, err)
		assert.Len(t, categories, 0)
	})
}

// TestDeleteCategory tests deleting a category.
func TestDeleteCategory(t *testing.T) {
	db := setupCategoryTestDB()
	repo := NewCategoryRepository(db)

	t.Run("Delete existing category", func(t *testing.T) {
		category := &models.Category{Name: "Gaming", Description: "Video games"}
		err1 := repo.CreateCategory(category)
		assert.NoError(t, err1)

		err := repo.DeleteCategory(category.ID)
		assert.NoError(t, err)

		var deletedCategory models.Category
		err = db.First(&deletedCategory, category.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete non-existent category", func(t *testing.T) {
		err := repo.DeleteCategory(9999) // Non-existent ID
		assert.NoError(t, err)           // `gorm.Delete` doesn't return an error if the record doesn't exist
	})
}

// TestUpdateCategory tests updating an existing category.
func TestUpdateCategory(t *testing.T) {
	db := setupCategoryTestDB()
	repo := NewCategoryRepository(db)

	t.Run("Update existing category", func(t *testing.T) {
		// Create a category
		category := &models.Category{Name: "Transport", Description: "Car, public transport"}
		err1 := repo.CreateCategory(category)
		assert.NoError(t, err1)

		// Update the category
		category.Name = "Travel"
		category.Description = "Flights, hotels, and transport"
		err := repo.UpdateCategory(category)
		assert.NoError(t, err)

		// Fetch updated category
		var updatedCategory models.Category
		err = db.First(&updatedCategory, category.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Travel", updatedCategory.Name)
		assert.Equal(t, "Flights, hotels, and transport", updatedCategory.Description)
	})
}
