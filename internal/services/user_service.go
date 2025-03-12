package services

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// UserService defines the interface for user operations
type UserService interface {
	RegisterUser(name, email, password string) (*models.User, error)
	AuthenticateUser(email, password string) (*models.User, error)
}
