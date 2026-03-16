package httpapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
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
