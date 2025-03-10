package repositories

import (
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupPaymentTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.PaymentMethod{})
	return db
}

// TestCreatePaymentMethod tests different cases for creating a payment method.
func TestCreatePaymentMethod(t *testing.T) {
	db := setupPaymentTestDB()
	repo := NewPaymentMethodRepository(db)

	// Ensure related entity (User) exists
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("Create valid payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{Name: "Credit Card", UserID: user.ID}
		err := repo.CreatePaymentMethod(paymentMethod)
		assert.NoError(t, err)

		var retrievedPaymentMethod models.PaymentMethod
		err = db.First(&retrievedPaymentMethod, paymentMethod.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, paymentMethod.Name, retrievedPaymentMethod.Name)
	})

	t.Run("Create duplicate payment method name for same user", func(t *testing.T) {
		paymentMethod1 := &models.PaymentMethod{Name: "PayPal", UserID: user.ID}
		paymentMethod2 := &models.PaymentMethod{Name: "PayPal", UserID: user.ID} // Duplicate

		err := repo.CreatePaymentMethod(paymentMethod1)
		assert.NoError(t, err)

		err = repo.CreatePaymentMethod(paymentMethod2) // Should fail due to unique constraint
		assert.Error(t, err)
	})
}

// TestGetPaymentMethodByID tests retrieving a payment method by ID.
func TestGetPaymentMethodByID(t *testing.T) {
	db := setupPaymentTestDB()
	repo := NewPaymentMethodRepository(db)

	// Ensure related entity (User) exists
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("Retrieve existing payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{Name: "Debit Card", UserID: user.ID}
		repo.CreatePaymentMethod(paymentMethod)

		foundPaymentMethod, err := repo.GetPaymentMethodByID(paymentMethod.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundPaymentMethod)
		assert.Equal(t, paymentMethod.Name, foundPaymentMethod.Name)
	})

	t.Run("Retrieve non-existent payment method", func(t *testing.T) {
		paymentMethod, err := repo.GetPaymentMethodByID(9999) // Non-existent ID
		assert.Error(t, err)
		assert.Nil(t, paymentMethod)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestGetPaymentMethodsByUserID tests retrieving all payment methods for a specific user.
func TestGetPaymentMethodsByUserID(t *testing.T) {
	db := setupPaymentTestDB()
	repo := NewPaymentMethodRepository(db)

	// Ensure related entity (User) exists
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("User has multiple payment methods", func(t *testing.T) {
		paymentMethod1 := &models.PaymentMethod{Name: "Apple Pay", UserID: user.ID}
		paymentMethod2 := &models.PaymentMethod{Name: "Google Pay", UserID: user.ID}
		repo.CreatePaymentMethod(paymentMethod1)
		repo.CreatePaymentMethod(paymentMethod2)

		paymentMethods, err := repo.GetPaymentMethodsByUserID(user.ID)
		assert.NoError(t, err)
		assert.Len(t, paymentMethods, 2)
	})

	t.Run("User has no payment methods", func(t *testing.T) {
		paymentMethods, err := repo.GetPaymentMethodsByUserID(9999) // Non-existent user
		assert.NoError(t, err)
		assert.Len(t, paymentMethods, 0)
	})
}

// TestUpdatePaymentMethod tests updating an existing payment method.
func TestUpdatePaymentMethod(t *testing.T) {
	db := setupPaymentTestDB()
	repo := NewPaymentMethodRepository(db)

	// Ensure related entity (User) exists
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("Update existing payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{Name: "Bank Transfer", UserID: user.ID}
		repo.CreatePaymentMethod(paymentMethod)

		// Update the payment method
		paymentMethod.Name = "Wire Transfer"
		err := repo.UpdatePaymentMethod(paymentMethod)
		assert.NoError(t, err)

		// Fetch updated payment method
		var updatedPaymentMethod models.PaymentMethod
		err = db.First(&updatedPaymentMethod, paymentMethod.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Wire Transfer", updatedPaymentMethod.Name)
	})
}

// TestDeletePaymentMethod tests deleting a payment method.
func TestDeletePaymentMethod(t *testing.T) {
	db := setupPaymentTestDB()
	repo := NewPaymentMethodRepository(db)

	// Ensure related entity (User) exists
	user := &models.User{Name: "Test User", Email: "test@example.com", Password: "hashedpassword"}
	db.Create(user)

	t.Run("Delete existing payment method", func(t *testing.T) {
		paymentMethod := &models.PaymentMethod{Name: "Venmo", UserID: user.ID}
		repo.CreatePaymentMethod(paymentMethod)

		err := repo.DeletePaymentMethod(paymentMethod.ID)
		assert.NoError(t, err)

		var deletedPaymentMethod models.PaymentMethod
		err = db.First(&deletedPaymentMethod, paymentMethod.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("Delete non-existent payment method", func(t *testing.T) {
		err := repo.DeletePaymentMethod(9999) // Non-existent ID
		assert.NoError(t, err)                // `gorm.Delete` does not return an error if the record doesn't exist
	})
}
