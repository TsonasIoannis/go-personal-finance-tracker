package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLoggerMiddleware writes one structured log entry per completed request.
func StructuredLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		requestID, _ := RequestIDFromGinContext(c)
		path := c.Request.URL.Path
		route := c.FullPath()
		if route == "" {
			route = path
		}

		args := []any{
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
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
			logger.Error("request completed", args...)
		case c.Writer.Status() >= http.StatusBadRequest:
			logger.Warn("request completed", args...)
		default:
			logger.Info("request completed", args...)
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

		logger.Error(
			"panic recovered",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"route", route,
			"panic", recovered,
			"stack_trace", string(debug.Stack()),
		)

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
