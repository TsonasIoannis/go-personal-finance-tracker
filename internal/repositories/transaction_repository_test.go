package repositories

import (
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTransactionTestDB initializes an in-memory SQLite database for testing.
func setupTransactionTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{}, &models.Transaction{})
	return db
}

func TestTransactionRepository(t *testing.T) {
	db := setupTransactionTestDB()
	repo := NewTransactionRepository(db)

	// Create a test user for transactions
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("CreateTransaction", func(t *testing.T) {
		transaction := &models.Transaction{
			UserID:     user.ID,
			Type:       "expense",
			Amount:     100.50,
			CategoryID: 1,
			Date:       time.Now(),
			Note:       "Test transaction",
		}

		err := repo.CreateTransaction(transaction)
		assert.NoError(t, err)

		// Verify transaction was inserted
		var retrievedTransaction models.Transaction
		err = db.First(&retrievedTransaction, transaction.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, transaction.Amount, retrievedTransaction.Amount)
	})

	t.Run("GetTransactionByID", func(t *testing.T) {
		transaction := &models.Transaction{
			UserID:     user.ID,
			Type:       "income",
			Amount:     200.00,
			CategoryID: 2,
			Date:       time.Now(),
			Note:       "Salary payment",
		}
		err := repo.CreateTransaction(transaction)
		assert.NoError(t, err)

		foundTransaction, err := repo.GetTransactionByID(transaction.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundTransaction)
		assert.Equal(t, transaction.Amount, foundTransaction.Amount)
	})

	t.Run("GetTransactionByID_NotFound", func(t *testing.T) {
		transaction, err := repo.GetTransactionByID(9999) // Non-existent ID
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("GetTransactionsByUserID", func(t *testing.T) {
		// Reset transactions before running test
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Transaction{})

		err1 := repo.CreateTransaction(&models.Transaction{UserID: user.ID, Type: "expense", Amount: 50, CategoryID: 1, Date: time.Now()})
		assert.NoError(t, err1)
		err2 := repo.CreateTransaction(&models.Transaction{UserID: user.ID, Type: "income", Amount: 150, CategoryID: 2, Date: time.Now()})
		assert.NoError(t, err2)

		transactions, err := repo.GetTransactionsByUserID(user.ID)
		assert.NoError(t, err)
		assert.Len(t, transactions, 2) // Ensure we only have 2 transactions
	})

	t.Run("UpdateTransaction", func(t *testing.T) {
		transaction := &models.Transaction{
			UserID:     user.ID,
			Type:       "expense",
			Amount:     500.00,
			CategoryID: 1,
			Date:       time.Now(),
			Note:       "Old transaction",
		}
		err1 := repo.CreateTransaction(transaction)
		assert.NoError(t, err1)

		transaction.Amount = 600.00 // Update amount
		transaction.Note = "Updated transaction"
		err := repo.UpdateTransaction(transaction)
		assert.NoError(t, err)

		// Verify update
		updatedTransaction, _ := repo.GetTransactionByID(transaction.ID)
		assert.Equal(t, 600.00, updatedTransaction.Amount)
		assert.Equal(t, "Updated transaction", updatedTransaction.Note)
	})

	t.Run("DeleteTransaction", func(t *testing.T) {
		transaction := &models.Transaction{
			UserID:     user.ID,
			Type:       "expense",
			Amount:     20.00,
			CategoryID: 3,
			Date:       time.Now(),
			Note:       "To be deleted",
		}
		err1 := repo.CreateTransaction(transaction)
		assert.NoError(t, err1)

		err := repo.DeleteTransaction(transaction.ID)
		assert.NoError(t, err)

		// Verify deletion
		var deletedTransaction models.Transaction
		err = db.First(&deletedTransaction, transaction.ID).Error
		assert.Error(t, err) // Should return an error because it's deleted
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
