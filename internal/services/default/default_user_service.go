package services

import (
	"context"
	"strings"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type DefaultUserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *DefaultUserService {
	return &DefaultUserService{userRepo: userRepo}
}

// RegisterUser creates a new user with a hashed password
func (s *DefaultUserService) RegisterUser(ctx context.Context, name, email, password string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Internal("user_registration_failed", "failed to register user", err)
	}

	// Create user
	user := &models.User{Name: name, Email: email, Password: string(hashedPassword)}
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, apperrors.Conflict("email_already_registered", "email already registered")
		}

		return nil, apperrors.Internal("user_registration_failed", "failed to register user", err)
	}
	return user, nil
}

// AuthenticateUser checks email & password for login
func (s *DefaultUserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid_credentials", "invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.Unauthorized("invalid_credentials", "invalid credentials")
	}

	return user, nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	errMessage := strings.ToLower(err.Error())
	return strings.Contains(errMessage, "duplicate key") ||
		strings.Contains(errMessage, "unique constraint") ||
		strings.Contains(errMessage, "unique failed")
}
