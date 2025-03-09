package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the database.Database interface
type MockDatabase struct {
	mock.Mock
	MockDB *sql.DB // Allow storing a mock SQL database
}

func (m *MockDatabase) Connect(openDB func(driverName, dataSourceName string) (*sql.DB, error)) error {
	args := m.Called(openDB)

	// Simulate a real database connection
	if openDB != nil {
		mockDB, err := openDB("mock", "mock_dsn")
		if err != nil {
			return err
		}
		m.MockDB = mockDB
	}

	return args.Error(0)
}

func (m *MockDatabase) CheckConnection() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) GetDB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestReadinessCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 OK when DB is available", func(t *testing.T) {
		mockDB := new(MockDatabase)
		mockDB.On("CheckConnection").Return(nil)

		router := gin.New()
		router.GET("/readiness", ReadinessCheckHandler(mockDB))

		// Explicitly check for error
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
		mockDB.On("CheckConnection").Return(errors.New("database not reachable"))

		router := gin.New()
		router.GET("/readiness", ReadinessCheckHandler(mockDB))

		// Explicitly check for error
		req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.JSONEq(t, `{"status": "unavailable", "error": "database not reachable"}`, w.Body.String())

		mockDB.AssertExpectations(t)
	})
}
