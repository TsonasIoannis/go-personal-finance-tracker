package observability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/TsonasIoannis/go-personal-finance-tracker/http"

// ConfigureTracing installs trace context and baggage propagation hooks.
func ConfigureTracing() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}

// TracingMiddleware starts a server span for each incoming HTTP request.
func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer(tracerName)

	return func(c *gin.Context) {
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		route := c.FullPath()
		path := c.Request.URL.Path
		if route == "" {
			route = path
		}

		attributes := []attribute.KeyValue{
			attribute.String("http.request.method", c.Request.Method),
			attribute.String("url.path", path),
			attribute.String("http.route", route),
		}

		if requestID, ok := requestIDFromGinContext(c); ok {
			attributes = append(attributes, attribute.String("request.id", requestID))
		}

		ctx, span := tracer.Start(
			ctx,
			spanName(c.Request.Method, route),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attributes...),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		route = c.FullPath()
		if route == "" {
			route = path
		}

		statusCode := c.Writer.Status()
		span.SetName(spanName(c.Request.Method, route))
		span.SetAttributes(
			attribute.String("http.route", route),
			attribute.Int("http.response.status_code", statusCode),
		)

		if statusCode >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		}
	}
}

// RecordError records an error on the current span and annotates it with response details.
func RecordError(ctx context.Context, err error, statusCode int) {
	if err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.RecordError(err)
	span.SetAttributes(attribute.Int("http.response.status_code", statusCode))
	if statusCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, err.Error())
	}
}

// RecordPanic records a recovered panic on the current span.
func RecordPanic(ctx context.Context, recovered any) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.RecordError(fmt.Errorf("panic recovered: %v", recovered))
	span.SetStatus(codes.Error, "panic recovered")
	span.AddEvent("panic.recovered")
}

// SetAuthenticatedUser annotates the current span with the authenticated user ID.
func SetAuthenticatedUser(ctx context.Context, userID uint) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.SetAttributes(attribute.Int64("user.id", int64(userID)))
	span.AddEvent("auth.user_authenticated")
}

// TraceIDsFromContext returns the current trace and span IDs when tracing is active.
func TraceIDsFromContext(ctx context.Context) (string, string, bool) {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return "", "", false
	}

	return spanContext.TraceID().String(), spanContext.SpanID().String(), true
}

func requestIDFromGinContext(c *gin.Context) (string, bool) {
	requestID, ok := c.Get("requestID")
	if !ok {
		return "", false
	}

	value, ok := requestID.(string)
	return value, ok && value != ""
}

func spanName(method, route string) string {
	return method + " " + route
}
