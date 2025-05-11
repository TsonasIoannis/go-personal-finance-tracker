package repositories

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// UserRepository defines the required repository methods
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}
