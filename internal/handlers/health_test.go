package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 OK with status 'ok'", func(t *testing.T) {
		router := gin.New()
		router.GET("/health", HealthCheckHandler)

		// Explicitly check for error
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())
	})
}
