package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/persistence"
)

type stubDatabase struct{}

func (stubDatabase) Connect() error         { return nil }
func (stubDatabase) Migrate() error         { return nil }
func (stubDatabase) Close() error           { return nil }
func (stubDatabase) CheckConnection() error { return nil }

func TestNewRouterExposesMetrics(t *testing.T) {
	cfg := config.Config{
		JWTSecret: "test-secret",
		Port:      "8080",
		HTTP: config.HTTPConfig{
			ReadTimeout:       time.Second,
			ReadHeaderTimeout: time.Second,
			WriteTimeout:      time.Second,
			IdleTimeout:       time.Second,
			ShutdownTimeout:   time.Second,
		},
		Auth: config.AuthConfig{
			TokenTTL: time.Hour,
		},
	}

	router := newRouter(cfg, stubDatabase{}, persistence.Repositories{})

	healthRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthResponse := httptest.NewRecorder()
	router.ServeHTTP(healthResponse, healthRequest)

	if healthResponse.Code != http.StatusOK {
		t.Fatalf("expected /health status %d, got %d", http.StatusOK, healthResponse.Code)
	}

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsResponse := httptest.NewRecorder()
	router.ServeHTTP(metricsResponse, metricsRequest)

	if metricsResponse.Code != http.StatusOK {
		t.Fatalf("expected /metrics status %d, got %d", http.StatusOK, metricsResponse.Code)
	}

	body := metricsResponse.Body.String()
	if body == "" {
		t.Fatal("expected /metrics response body to be populated")
	}

	if !strings.Contains(body, "personal_finance_tracker_http_requests_total") {
		t.Fatal("expected requests total metric to be exposed")
	}

	if !strings.Contains(body, "personal_finance_tracker_http_request_duration_seconds") {
		t.Fatal("expected request duration metric to be exposed")
	}

	if !strings.Contains(body, `route="/health"`) {
		t.Fatal("expected /health route labels to be present in metrics output")
	}

	if !strings.Contains(body, `status_code="200"`) {
		t.Fatal("expected successful status code labels to be present in metrics output")
	}

	swaggerRequest := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	swaggerResponse := httptest.NewRecorder()
	router.ServeHTTP(swaggerResponse, swaggerRequest)

	if swaggerResponse.Code != http.StatusOK {
		t.Fatalf("expected /openapi.json status %d, got %d", http.StatusOK, swaggerResponse.Code)
	}

	if !strings.Contains(swaggerResponse.Body.String(), "\"swagger\"") {
		t.Fatal("expected generated swagger document to be served")
	}
}
