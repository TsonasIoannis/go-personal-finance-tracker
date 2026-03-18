package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Run("loads required config with defaults", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/personal_finance_db?sslmode=disable")
		t.Setenv("JWT_SECRET", "super-secret")
		t.Setenv("PORT", "")
		t.Setenv("HTTP_READ_TIMEOUT", "")
		t.Setenv("HTTP_READ_HEADER_TIMEOUT", "")
		t.Setenv("HTTP_WRITE_TIMEOUT", "")
		t.Setenv("HTTP_IDLE_TIMEOUT", "")
		t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "")
		t.Setenv("AUTH_TOKEN_TTL", "")
		t.Setenv("OTEL_SERVICE_NAME", "")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "")
		t.Setenv("OTEL_TRACES_SAMPLER_ARG", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("expected config to load, got error %v", err)
		}

		if cfg.DatabaseURL == "" {
			t.Fatal("expected database url to be set")
		}

		if cfg.JWTSecret != "super-secret" {
			t.Fatalf("expected JWT secret to be loaded, got %q", cfg.JWTSecret)
		}

		if cfg.Port != defaultPort {
			t.Fatalf("expected default port %q, got %q", defaultPort, cfg.Port)
		}

		if cfg.HTTP.ReadTimeout != defaultReadTimeout {
			t.Fatalf("expected default read timeout %v, got %v", defaultReadTimeout, cfg.HTTP.ReadTimeout)
		}

		if cfg.Auth.TokenTTL != defaultTokenTTL {
			t.Fatalf("expected default token ttl %v, got %v", defaultTokenTTL, cfg.Auth.TokenTTL)
		}

		if cfg.Tracing.ServiceName != defaultServiceName {
			t.Fatalf("expected default service name %q, got %q", defaultServiceName, cfg.Tracing.ServiceName)
		}

		if cfg.Tracing.Endpoint != "" {
			t.Fatalf("expected tracing endpoint to default to empty, got %q", cfg.Tracing.Endpoint)
		}

		if cfg.Tracing.SampleRatio != defaultTraceSampleRatio {
			t.Fatalf("expected default tracing sample ratio %v, got %v", defaultTraceSampleRatio, cfg.Tracing.SampleRatio)
		}
	})

	t.Run("loads custom durations", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/personal_finance_db?sslmode=disable")
		t.Setenv("JWT_SECRET", "super-secret")
		t.Setenv("PORT", "9090")
		t.Setenv("HTTP_READ_TIMEOUT", "7s")
		t.Setenv("HTTP_READ_HEADER_TIMEOUT", "3s")
		t.Setenv("HTTP_WRITE_TIMEOUT", "12s")
		t.Setenv("HTTP_IDLE_TIMEOUT", "75s")
		t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "15s")
		t.Setenv("AUTH_TOKEN_TTL", "48h")
		t.Setenv("OTEL_SERVICE_NAME", "finance-api")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
		t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
		t.Setenv("OTEL_TRACES_SAMPLER_ARG", "0.25")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("expected config to load, got error %v", err)
		}

		if cfg.Port != "9090" {
			t.Fatalf("expected custom port, got %q", cfg.Port)
		}

		if cfg.HTTP.ReadTimeout != 7*time.Second {
			t.Fatalf("expected custom read timeout, got %v", cfg.HTTP.ReadTimeout)
		}

		if cfg.HTTP.ShutdownTimeout != 15*time.Second {
			t.Fatalf("expected custom shutdown timeout, got %v", cfg.HTTP.ShutdownTimeout)
		}

		if cfg.Auth.TokenTTL != 48*time.Hour {
			t.Fatalf("expected custom token ttl, got %v", cfg.Auth.TokenTTL)
		}

		if cfg.Tracing.ServiceName != "finance-api" {
			t.Fatalf("expected custom service name, got %q", cfg.Tracing.ServiceName)
		}

		if cfg.Tracing.Endpoint != "http://localhost:4318" {
			t.Fatalf("expected custom tracing endpoint, got %q", cfg.Tracing.Endpoint)
		}

		if !cfg.Tracing.Insecure {
			t.Fatal("expected tracing insecure flag to be true")
		}

		if cfg.Tracing.SampleRatio != 0.25 {
			t.Fatalf("expected custom trace sample ratio, got %v", cfg.Tracing.SampleRatio)
		}
	})

	t.Run("fails when required env vars are missing", func(t *testing.T) {
		unsetEnv(t, "DATABASE_URL")
		unsetEnv(t, "JWT_SECRET")

		_, err := Load()
		if err == nil {
			t.Fatal("expected missing env error")
		}
	})

	t.Run("fails when duration env vars are invalid", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/personal_finance_db?sslmode=disable")
		t.Setenv("JWT_SECRET", "super-secret")
		t.Setenv("HTTP_READ_TIMEOUT", "invalid")

		_, err := Load()
		if err == nil {
			t.Fatal("expected invalid duration error")
		}
	})

	t.Run("fails when tracing sample ratio is invalid", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/personal_finance_db?sslmode=disable")
		t.Setenv("JWT_SECRET", "super-secret")
		t.Setenv("OTEL_TRACES_SAMPLER_ARG", "1.5")

		_, err := Load()
		if err == nil {
			t.Fatal("expected invalid tracing sample ratio error")
		}
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	originalValue, hadOriginalValue := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}

	t.Cleanup(func() {
		if !hadOriginalValue {
			_ = os.Unsetenv(key)
			return
		}

		_ = os.Setenv(key, originalValue)
	})
}
