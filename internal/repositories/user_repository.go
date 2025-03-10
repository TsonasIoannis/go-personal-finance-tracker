package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles DB operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository initializes a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user from the database
func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
