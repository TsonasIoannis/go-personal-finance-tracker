package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService implements services.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(name, email, password string) (*models.User, error) {
	args := m.Called(name, email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) AuthenticateUser(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		input := models.User{
			Name:     "Alice",
			Email:    "alice@example.com",
			Password: "secure123",
		}
		expected := &models.User{
			ID:    1,
			Name:  "Alice",
			Email: "alice@example.com",
		}

		mockService.On("RegisterUser", input.Name, input.Email, input.Password).
			Return(expected, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(input)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered")
		assert.Contains(t, w.Body.String(), `"Email":"alice@example.com"`)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		input := models.User{
			Name:     "Bob",
			Email:    "bob@example.com",
			Password: "pass123",
		}

		mockService.On("RegisterUser", input.Name, input.Email, input.Password).
			Return((*models.User)(nil), errors.New("email already registered")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonBody, _ := json.Marshal(input)
		c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Register(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "email already registered")
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		email := "alice@example.com"
		password := "secure123"
		expected := &models.User{ID: 1, Name: "Alice", Email: email}

		mockService.On("AuthenticateUser", email, password).
			Return(expected, nil).Once()

		payload := map[string]string{
			"email":    email,
			"password": password,
		}
		jsonBody, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Login successful")
		assert.Contains(t, w.Body.String(), `"Email":"alice@example.com"`)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":`))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request payload")
	})

	t.Run("Authentication Failure", func(t *testing.T) {
		mockService := new(MockUserService)
		controller := NewUserController(mockService)

		email := "bob@example.com"
		password := "wrongpass"

		mockService.On("AuthenticateUser", email, password).
			Return((*models.User)(nil), errors.New("invalid credentials")).Once()

		body := map[string]string{"email": email, "password": password}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})
}
