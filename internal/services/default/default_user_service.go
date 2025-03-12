package services

import (
	"errors"

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
func (s *DefaultUserService) RegisterUser(name, email, password string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{Name: name, Email: email, Password: string(hashedPassword)}
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AuthenticateUser checks email & password for login
func (s *DefaultUserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
