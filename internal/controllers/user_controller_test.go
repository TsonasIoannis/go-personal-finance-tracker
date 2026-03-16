package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) ParseToken(token string) (*auth.Claims, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0).(*auth.Claims), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockUserService implements services.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(name, email, password string) (*models.User, error) {
	args := m.Called(name, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) AuthenticateUser(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		input := map[string]string{
			"name":     "Alice",
			"email":    "alice@example.com",
			"password": "secure123",
		}
		expected := &models.User{
			ID:    1,
			Name:  "Alice",
			Email: "alice@example.com",
		}

		mockService.On("RegisterUser", "Alice", "alice@example.com", "secure123").
			Return(expected, nil).Once()
		mockTokenManager.On("GenerateToken", expected).Return("token-123", nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, err := json.Marshal(input)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered")
		assert.Contains(t, w.Body.String(), `"email":"alice@example.com"`)
		assert.Contains(t, w.Body.String(), `"token":"token-123"`)
		assert.NotContains(t, w.Body.String(), "Password")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_request"`)
		assert.Contains(t, w.Body.String(), "invalid request payload")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		input := map[string]string{
			"name":     "Bob",
			"email":    "bob@example.com",
			"password": "password123",
		}

		mockService.On("RegisterUser", "Bob", "bob@example.com", "password123").
			Return((*models.User)(nil), apperrors.Conflict("email_already_registered", "email already registered")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, err := json.Marshal(input)
		assert.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"email_already_registered"`)
		assert.Contains(t, w.Body.String(), "email already registered")
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		email := "alice@example.com"
		password := "secure123"
		expected := &models.User{ID: 1, Name: "Alice", Email: email}

		mockService.On("AuthenticateUser", email, password).
			Return(expected, nil).Once()
		mockTokenManager.On("GenerateToken", expected).Return("token-abc", nil).Once()

		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		jsonBody, err := json.Marshal(payload)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login successful")
		assert.Contains(t, w.Body.String(), `"email":"alice@example.com"`)
		assert.Contains(t, w.Body.String(), `"token":"token-abc"`)
		assert.NotContains(t, w.Body.String(), "Password")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_request"`)
		assert.Contains(t, w.Body.String(), "invalid request payload")
	})

	t.Run("Authentication Failure", func(t *testing.T) {
		mockService := new(MockUserService)
		mockTokenManager := new(MockTokenManager)
		controller := NewUserController(mockService, mockTokenManager)

		email := "bob@example.com"
		password := "wrongpass"

		mockService.On("AuthenticateUser", email, password).
			Return((*models.User)(nil), apperrors.Unauthorized("invalid_credentials", "invalid credentials")).Once()

		body := map[string]string{"email": email, "password": password}
		jsonBody, err := json.Marshal(body)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), `"code":"invalid_credentials"`)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}
