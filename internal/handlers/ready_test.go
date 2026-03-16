package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDatabase is a mock implementation of the database.Database interface.
type MockDatabase struct {
	mock.Mock
	MockDB *gorm.DB
}

func (m *MockDatabase) Connect() error {
	args := m.Called()
	if m.MockDB == nil {
		m.MockDB = &gorm.DB{}
	}

	return args.Error(0)
}

func (m *MockDatabase) CheckConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) GetDB() *gorm.DB {
	args := m.Called()
	if db, ok := args.Get(0).(*gorm.DB); ok {
		return db
	}

	return nil
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestReadinessCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 OK when DB is available", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockGormDB := &gorm.DB{}
		mockDB.MockDB = mockGormDB

		mockDB.On("CheckConnection").Return(nil)
		mockDB.On("GetDB").Return(mockGormDB)

		router := gin.New()
		router.GET("/readiness", ReadinessCheckHandler(mockDB))

		req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status": "ready"}`, w.Body.String())

		mockDB.AssertExpectations(t)
	})

	t.Run("should return 503 Service Unavailable when DB is down", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockGormDB := &gorm.DB{}
		mockDB.MockDB = mockGormDB

		mockDB.On("CheckConnection").Return(errors.New("database not reachable"))
		mockDB.On("GetDB").Return(mockGormDB)

		router := gin.New()
		router.GET("/readiness", ReadinessCheckHandler(mockDB))

		req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.JSONEq(t, `{"status": "unavailable", "error": "database not reachable"}`, w.Body.String())

		mockDB.AssertExpectations(t)
	})

	t.Run("should return 503 when DB is not initialized", func(t *testing.T) {
		mockDB := new(MockDatabase)

		mockDB.On("GetDB").Return(nil)

		router := gin.New()
		router.GET("/readiness", ReadinessCheckHandler(mockDB))

		req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.JSONEq(t, `{"status": "unavailable", "error": "database is not initialized"}`, w.Body.String())

		mockDB.AssertExpectations(t)
	})
}
