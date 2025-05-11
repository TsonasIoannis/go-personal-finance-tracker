package services

import (
	"errors"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("Register valid user", func(t *testing.T) {
		mockRepo.On("CreateUser", mock.MatchedBy(func(u *models.User) bool {
			// Check that the password is hashed
			err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("mypassword"))
			return err == nil
		})).Return(nil)

		createdUser, err := service.RegisterUser("John Doe", "john@example.com", "mypassword")
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, "John Doe", createdUser.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail when repository returns an error", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		mockRepo.On("CreateUser", mock.Anything).Return(errors.New("database error"))

		createdUser, err := service.RegisterUser("Jane Doe", "jane@example.com", "securepassword")

		// Ensure the error is returned
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.Equal(t, "database error", err.Error())

		// Ensure repository was called
		mockRepo.AssertExpectations(t)
	})
	t.Run("Fail when password is too long for bcrypt", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil // Reset expectations

		// Create an excessively long password
		longPassword := string(make([]byte, 1000)) // 1000 bytes long

		createdUser, err := service.RegisterUser("John Doe", "john@example.com", longPassword)

		// Ensure bcrypt failure occurs
		assert.Error(t, err)
		assert.Nil(t, createdUser)

		// Ensure repository was NOT called since hashing failed
		mockRepo.AssertNotCalled(t, "CreateUser")
	})

}

func TestAuthenticateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	t.Run("Authenticate valid user", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
		user := &models.User{Name: "John Doe", Email: "john@example.com", Password: string(hashedPassword)}

		mockRepo.On("GetUserByEmail", "john@example.com").Return(user, nil)

		authenticatedUser, err := service.AuthenticateUser("john@example.com", "mypassword")
		assert.NoError(t, err)
		assert.NotNil(t, authenticatedUser)
		assert.Equal(t, "John Doe", authenticatedUser.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail when user does not exist", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "nonexistent@example.com").Return(nil, errors.New("not found"))

		authenticatedUser, err := service.AuthenticateUser("nonexistent@example.com", "password")
		assert.Error(t, err)
		assert.Nil(t, authenticatedUser)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
	t.Run("Fail when password is incorrect", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := &models.User{Name: "John Doe", Email: "john@example.com", Password: string(hashedPassword)}

		mockRepo.On("GetUserByEmail", "john@example.com").Return(user, nil)

		authenticatedUser, err := service.AuthenticateUser("john@example.com", "wrongpassword")
		assert.Error(t, err)
		assert.Nil(t, authenticatedUser)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
