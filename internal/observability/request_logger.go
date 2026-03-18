package observability

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

const loggerGinKey = "requestLogger"

type contextKey string

const loggerContextKey contextKey = "request_logger"

// WithLogger attaches a structured logger to a standard context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// LoggerFromContext returns a logger from a standard context or the default logger.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	logger, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}

	return logger
}

// SetLoggerOnGinContext stores the logger on both Gin and request contexts.
func SetLoggerOnGinContext(c *gin.Context, logger *slog.Logger) {
	if c == nil || logger == nil {
		return
	}

	c.Set(loggerGinKey, logger)
	c.Request = c.Request.WithContext(WithLogger(c.Request.Context(), logger))
}

// LoggerFromGinContext returns the request-scoped logger or the default logger.
func LoggerFromGinContext(c *gin.Context) *slog.Logger {
	if c == nil {
		return slog.Default()
	}

	loggerValue, ok := c.Get(loggerGinKey)
	if ok {
		logger, ok := loggerValue.(*slog.Logger)
		if ok && logger != nil {
			return logger
		}
	}

	return LoggerFromContext(c.Request.Context())
}
