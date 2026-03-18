package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStructuredLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	router := gin.New()
	router.Use(RequestIDMiddleware(), StructuredLoggerMiddleware(logger))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, "req-123")
	req.Header.Set("User-Agent", "integration-test")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var entry map[string]any
	err := json.Unmarshal([]byte(strings.TrimSpace(logBuffer.String())), &entry)
	assert.NoError(t, err)
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "request completed", entry["msg"])
	assert.Equal(t, "req-123", entry["request_id"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/test", entry["path"])
	assert.Equal(t, "/test", entry["route"])
	assert.Equal(t, "integration-test", entry["user_agent"])
	assert.Equal(t, float64(http.StatusCreated), entry["status"])
}

func TestStructuredLoggerMiddlewareIncludesEnrichedRequestFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	router := gin.New()
	router.Use(RequestIDMiddleware(), StructuredLoggerMiddleware(logger))
	router.GET("/test", func(c *gin.Context) {
		observability.SetLoggerOnGinContext(c, observability.LoggerFromGinContext(c).With("user_id", uint(42)))
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, "req-user")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var entry map[string]any
	err := json.Unmarshal([]byte(strings.TrimSpace(logBuffer.String())), &entry)
	assert.NoError(t, err)
	assert.Equal(t, "req-user", entry["request_id"])
	assert.Equal(t, float64(42), entry["user_id"])
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil))

	router := gin.New()
	router.Use(RequestIDMiddleware(), RecoveryMiddleware(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req.Header.Set(RequestIDHeader, "panic-123")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var entry map[string]any
	err := json.Unmarshal([]byte(strings.TrimSpace(logBuffer.String())), &entry)
	assert.NoError(t, err)
	assert.Equal(t, "ERROR", entry["level"])
	assert.Equal(t, "panic recovered", entry["msg"])
	assert.Equal(t, "panic-123", entry["request_id"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/panic", entry["path"])
	assert.Equal(t, "/panic", entry["route"])
	assert.Equal(t, "boom", entry["panic"])
}
