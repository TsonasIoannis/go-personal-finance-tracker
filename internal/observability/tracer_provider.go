package observability

import (
	"context"
	"strings"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace/noop"
)

// SetupTracing configures global tracing. When no endpoint is configured, it keeps a noop provider.
func SetupTracing(ctx context.Context, cfg config.TracingConfig) (func(context.Context) error, error) {
	ConfigureTracing()

	if strings.TrimSpace(cfg.Endpoint) == "" {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return func(context.Context) error { return nil }, nil
	}

	exporterOptions := []otlptracehttp.Option{}
	endpoint := strings.TrimSpace(cfg.Endpoint)
	switch {
	case strings.HasPrefix(endpoint, "http://"):
		exporterOptions = append(exporterOptions, otlptracehttp.WithEndpointURL(endpoint), otlptracehttp.WithInsecure())
	case strings.HasPrefix(endpoint, "https://"):
		exporterOptions = append(exporterOptions, otlptracehttp.WithEndpointURL(endpoint))
	default:
		exporterOptions = append(exporterOptions, otlptracehttp.WithEndpoint(endpoint))
		if cfg.Insecure {
			exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
		}
	}

	if cfg.Insecure && strings.HasPrefix(endpoint, "https://") {
		exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, exporterOptions...)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithProcess(),
	)
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRatio))),
	)

	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}
