package observability

import (
	"context"
	"testing"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestSetupTracingWithoutExporterUsesNoopProvider(t *testing.T) {
	previousProvider := otel.GetTracerProvider()
	previousPropagator := otel.GetTextMapPropagator()
	defer func() {
		otel.SetTracerProvider(previousProvider)
		otel.SetTextMapPropagator(previousPropagator)
	}()

	shutdown, err := SetupTracing(context.Background(), config.TracingConfig{
		ServiceName: "test-service",
		Endpoint:    "",
		SampleRatio: 1,
	})
	if err != nil {
		t.Fatalf("expected setup tracing without exporter to succeed, got %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected shutdown function to be returned")
	}

	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("expected noop shutdown to succeed, got %v", err)
	}

	spanContext := trace.SpanFromContext(context.Background()).SpanContext()
	if spanContext.IsValid() {
		t.Fatal("expected background context span to remain invalid")
	}

	if _, ok := otel.GetTextMapPropagator().(propagation.TextMapPropagator); !ok {
		t.Fatal("expected text map propagator to be configured")
	}
}
