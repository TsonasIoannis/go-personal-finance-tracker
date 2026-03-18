package observability_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/middleware"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracingMiddlewareCreatesAnnotatedSpan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))

	previousProvider := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	otel.SetTracerProvider(provider)
	observability.ConfigureTracing()
	defer func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
		_ = provider.Shutdown(context.Background())
	}()

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware(), observability.TracingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Set("userID", uint(42))
		observability.SetAuthenticatedUser(c.Request.Context(), 42)
		httpapi.WriteError(c, apperrors.Validation("invalid_request", "invalid request payload"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(middleware.RequestIDHeader, "req-123")
	req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]
	if span.Name() != "GET /test" {
		t.Fatalf("expected span name GET /test, got %s", span.Name())
	}

	if !span.Parent().IsValid() {
		t.Fatal("expected parent span context to be propagated")
	}

	if got := span.Parent().TraceID().String(); got != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("expected parent trace id to match traceparent header, got %s", got)
	}

	attributes := attributesByKey(span.Attributes())
	assertStringAttribute(t, attributes, "request.id", "req-123")
	assertStringAttribute(t, attributes, "http.request.method", http.MethodGet)
	assertStringAttribute(t, attributes, "http.route", "/test")
	assertInt64Attribute(t, attributes, "http.response.status_code", int64(http.StatusBadRequest))
	assertInt64Attribute(t, attributes, "user.id", 42)

	if !hasEvent(span.Events(), "exception") {
		t.Fatal("expected tracing hooks to record an exception event")
	}
}

func attributesByKey(attributes []attribute.KeyValue) map[string]attribute.Value {
	values := make(map[string]attribute.Value, len(attributes))
	for _, kv := range attributes {
		values[string(kv.Key)] = kv.Value
	}

	return values
}

func assertStringAttribute(t *testing.T, attributes map[string]attribute.Value, key, want string) {
	t.Helper()

	value, ok := attributes[key]
	if !ok {
		t.Fatalf("expected attribute %s to be present", key)
	}

	if got := value.AsString(); got != want {
		t.Fatalf("expected attribute %s to be %q, got %q", key, want, got)
	}
}

func assertInt64Attribute(t *testing.T, attributes map[string]attribute.Value, key string, want int64) {
	t.Helper()

	value, ok := attributes[key]
	if !ok {
		t.Fatalf("expected attribute %s to be present", key)
	}

	if got := value.AsInt64(); got != want {
		t.Fatalf("expected attribute %s to be %d, got %d", key, want, got)
	}
}

func hasEvent(events []sdktrace.Event, want string) bool {
	for _, event := range events {
		if event.Name == want {
			return true
		}
	}

	return false
}
