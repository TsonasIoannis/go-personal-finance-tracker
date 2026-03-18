package repositories

import (
	"context"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openSQLiteTestDB(t)
	err := database.ApplyMigrations(db)
	assert.NoError(t, err)
	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("should create a new user", func(t *testing.T) {
		user := &models.User{Name: "John Doe", Email: "john@example.com", Password: "hashedpassword"}

		err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		var retrievedUser models.User
		err = db.First(&retrievedUser, "email = ?", "john@example.com").Error
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", retrievedUser.Name)
	})

	t.Run("should retrieve a user by email", func(t *testing.T) {
		user := &models.User{Name: "Jane Doe", Email: "jane@example.com", Password: "hashedpassword"}
		err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)

		foundUser, err := repo.GetUserByEmail(ctx, "jane@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "Jane Doe", foundUser.Name)
	})

	t.Run("should delete a user", func(t *testing.T) {
		user := &models.User{Name: "Mark Smith", Email: "mark@example.com", Password: "hashedpassword"}
		err1 := repo.CreateUser(ctx, user)
		assert.NoError(t, err1)

		err := repo.DeleteUser(ctx, user.ID)
		assert.NoError(t, err)

		var retrievedUser models.User
		err = db.First(&retrievedUser, user.ID).Error
		assert.Error(t, err) // Should return an error because user is deleted
	})
	t.Run("should return error if user not found", func(t *testing.T) {
		repo := NewUserRepository(db)

		// Attempt to fetch a non-existent user
		user, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")

		// This should trigger the missing error branch
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

}
