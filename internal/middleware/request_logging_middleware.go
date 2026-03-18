package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/observability"
	"github.com/gin-gonic/gin"
)

// StructuredLoggerMiddleware writes one structured log entry per completed request.
func StructuredLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return func(c *gin.Context) {
		requestID, _ := RequestIDFromGinContext(c)
		path := c.Request.URL.Path
		requestLogger := logger.With(
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
		)
		observability.SetLoggerOnGinContext(c, requestLogger)

		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = path
		}

		requestLogger = observability.LoggerFromGinContext(c)
		args := []any{
			"route", route,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		if len(c.Errors) > 0 {
			args = append(args, "errors", c.Errors.String())
		}

		switch {
		case c.Writer.Status() >= http.StatusInternalServerError:
			requestLogger.Error("request completed", args...)
		case c.Writer.Status() >= http.StatusBadRequest:
			requestLogger.Warn("request completed", args...)
		default:
			requestLogger.Info("request completed", args...)
		}
	}
}

// RecoveryMiddleware logs panics using the shared structured logger.
func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		requestID, _ := RequestIDFromGinContext(c)
		path := c.Request.URL.Path
		route := c.FullPath()
		if route == "" {
			route = path
		}

		requestLogger := logger.With(
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
		)
		if userID, exists := c.Get("userID"); exists {
			requestLogger = requestLogger.With("user_id", userID)
		}

		requestLogger.Error(
			"panic recovered",
			"route", route,
			"panic", recovered,
			"stack_trace", string(debug.Stack()),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
