package repositories

import (
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupBudgetTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(
		&models.Budget{},
		&models.User{},
		&models.Transaction{},
	)
	return db
}

// TestCreateBudget tests different cases for creating a budget.
func TestCreateBudget(t *testing.T) {
	db := setupBudgetTestDB()
	repo := NewGormBudgetRepository(db)

	t.Run("Create valid budget", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      1000.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}

		err := repo.CreateBudget(budget)
		assert.NoError(t, err)

		var retrievedBudget models.Budget
		err = db.First(&retrievedBudget, budget.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, budget.Limit, retrievedBudget.Limit)
	})
}

// TestGetBudgetByID tests retrieving a budget by ID.
func TestGetBudgetByID(t *testing.T) {
	db := setupBudgetTestDB()
	repo := NewGormBudgetRepository(db)

	t.Run("Retrieve existing budget", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     1,
			CategoryID: 3,
			Limit:      500.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}
		err := repo.CreateBudget(budget)
		assert.NoError(t, err)

		foundBudget, err := repo.GetBudgetByID(budget.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundBudget)
		assert.Equal(t, budget.Limit, foundBudget.Limit)
	})

	t.Run("Retrieve non-existent budget", func(t *testing.T) {
		nonExistentID := uint(9999)
		budget, err := repo.GetBudgetByID(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, budget)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestGetBudgetsByUserID tests retrieving budgets by user ID.
func TestGetBudgetsByUserID(t *testing.T) {
	db := setupBudgetTestDB()
	repo := NewGormBudgetRepository(db)

	// Ensure a clean state before running the test
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Budget{})

	t.Run("User has multiple budgets", func(t *testing.T) {
		budget1 := &models.Budget{
			UserID:     1,
			CategoryID: 2,
			Limit:      1000.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}
		budget2 := &models.Budget{
			UserID:     1,
			CategoryID: 3,
			Limit:      500.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}
		err1 := repo.CreateBudget(budget1)
		assert.NoError(t, err1)
		err2 := repo.CreateBudget(budget2)
		assert.NoError(t, err2)

		budgets, err := repo.GetBudgetsByUserID(1)
		assert.NoError(t, err)
		assert.Len(t, budgets, 2)
	})

	t.Run("User has no budgets", func(t *testing.T) {
		budgets, err := repo.GetBudgetsByUserID(9999) // Non-existent user
		assert.NoError(t, err)
		assert.Len(t, budgets, 0)
	})
}

// TestDeleteBudget tests deleting a budget.
func TestDeleteBudget(t *testing.T) {
	db := setupBudgetTestDB()
	repo := NewGormBudgetRepository(db)

	t.Run("Delete existing budget", func(t *testing.T) {
		budget := &models.Budget{
			UserID:     2,
			CategoryID: 1,
			Limit:      200.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}
		err1 := repo.CreateBudget(budget)
		assert.NoError(t, err1)

		err := repo.DeleteBudget(budget.ID)
		assert.NoError(t, err)

		var deletedBudget models.Budget
		err = db.First(&deletedBudget, budget.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestUpdateBudget tests updating an existing budget.
func TestUpdateBudget(t *testing.T) {
	db := setupBudgetTestDB()
	repo := NewGormBudgetRepository(db)

	// Ensure related entities exist
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	category := &models.Category{Name: "Utilities", Description: "Bills and utilities"}
	db.Create(user)
	db.Create(category)

	t.Run("Update existing budget", func(t *testing.T) {
		// Create a budget
		budget := &models.Budget{
			UserID:     user.ID,
			CategoryID: category.ID,
			Limit:      500.00,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 1, 0),
		}
		err1 := repo.CreateBudget(budget)
		assert.NoError(t, err1)

		// Update the budget
		budget.Limit = 750.00
		err := repo.UpdateBudget(budget)
		assert.NoError(t, err)

		// Fetch updated budget
		var updatedBudget models.Budget
		err = db.First(&updatedBudget, budget.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, 750.00, updatedBudget.Limit)
	})
}
