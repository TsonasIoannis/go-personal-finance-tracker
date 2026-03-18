package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDHeader = "X-Request-ID"
	requestIDKey    = "requestID"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

// RequestIDMiddleware ensures each request carries a stable request ID.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.GetHeader(RequestIDHeader))
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := context.WithValue(c.Request.Context(), requestIDContextKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(requestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// RequestIDFromContext returns the request ID stored on a standard context.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	return requestID, ok && requestID != ""
}

// RequestIDFromGinContext returns the request ID stored on a Gin context.
func RequestIDFromGinContext(c *gin.Context) (string, bool) {
	requestID, ok := c.Get(requestIDKey)
	if !ok {
		return "", false
	}

	value, ok := requestID.(string)
	return value, ok && value != ""
}

func newRequestID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic("request ID generation failed: " + err.Error())
	}

	return hex.EncodeToString(bytes[:])
}
