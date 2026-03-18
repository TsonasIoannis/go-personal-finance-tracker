package services

import (
	"context"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
)

// UserService defines the interface for user operations
type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string) (*models.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
}
