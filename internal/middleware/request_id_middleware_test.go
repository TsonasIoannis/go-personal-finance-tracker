package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("generates request id when header is missing", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestIDMiddleware())
		router.GET("/test", func(c *gin.Context) {
			requestID, ok := RequestIDFromGinContext(c)
			assert.True(t, ok)
			assert.NotEmpty(t, requestID)

			ctxRequestID, ok := RequestIDFromContext(c.Request.Context())
			assert.True(t, ok)
			assert.Equal(t, requestID, ctxRequestID)

			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Header().Get(RequestIDHeader))
		assert.Len(t, rec.Header().Get(RequestIDHeader), 32)
	})

	t.Run("reuses incoming request id header", func(t *testing.T) {
		router := gin.New()
		router.Use(RequestIDMiddleware())
		router.GET("/test", func(c *gin.Context) {
			requestID, ok := RequestIDFromGinContext(c)
			assert.True(t, ok)
			assert.Equal(t, "client-request-id", requestID)

			ctxRequestID, ok := RequestIDFromContext(c.Request.Context())
			assert.True(t, ok)
			assert.Equal(t, "client-request-id", ctxRequestID)

			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(RequestIDHeader, "client-request-id")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "client-request-id", rec.Header().Get(RequestIDHeader))
	})
}
