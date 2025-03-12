package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines the required repository methods
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

// GormUserRepository handles DB operations for users
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository initializes a new GormUserRepository
func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *GormUserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByEmail retrieves a user by email
func (r *GormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser removes a user from the database
func (r *GormUserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
