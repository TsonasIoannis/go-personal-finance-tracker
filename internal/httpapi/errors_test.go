package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/gin-gonic/gin"
)

func TestBuildErrorResponse(t *testing.T) {
	t.Run("maps typed validation errors", func(t *testing.T) {
		status, response := buildErrorResponse(apperrors.Validation("invalid_request", "invalid request payload"))

		if status != http.StatusBadRequest {
			t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
		}

		if response.Error.Code != "invalid_request" {
			t.Fatalf("expected code invalid_request, got %s", response.Error.Code)
		}

		if response.Error.Message != "invalid request payload" {
			t.Fatalf("expected validation message, got %s", response.Error.Message)
		}
	})

	t.Run("maps unknown errors to generic internal server errors", func(t *testing.T) {
		status, response := buildErrorResponse(errors.New("database exploded"))

		if status != http.StatusInternalServerError {
			t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, status)
		}

		if response.Error.Code != "internal_error" {
			t.Fatalf("expected internal_error code, got %s", response.Error.Code)
		}

		if response.Error.Message != "internal server error" {
			t.Fatalf("expected generic internal error message, got %s", response.Error.Message)
		}
	})
}

func TestWriteErrorLogsRequestScopedFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logBuffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuffer, nil)).With(
		"request_id", "req-123",
		"method", http.MethodGet,
		"path", "/transactions",
		"user_id", uint(42),
	)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/transactions", nil).WithContext(
		observability.WithLogger(context.Background(), logger),
	)
	observability.SetLoggerOnGinContext(c, logger)

	WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))

	var entry map[string]any
	err := json.Unmarshal([]byte(strings.TrimSpace(logBuffer.String())), &entry)
	if err != nil {
		t.Fatalf("expected valid log entry, got error %v", err)
	}

	if entry["level"] != "WARN" {
		t.Fatalf("expected WARN level, got %v", entry["level"])
	}

	if entry["msg"] != "request failed" {
		t.Fatalf("expected request failed log message, got %v", entry["msg"])
	}

	if entry["request_id"] != "req-123" {
		t.Fatalf("expected request id req-123, got %v", entry["request_id"])
	}

	if entry["user_id"] != float64(42) {
		t.Fatalf("expected user id 42, got %v", entry["user_id"])
	}

	if entry["error_code"] != "invalid_request" {
		t.Fatalf("expected invalid_request code, got %v", entry["error_code"])
	}
}
