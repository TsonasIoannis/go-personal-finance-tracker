package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	// Create a new Gin router
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Create a request to test the endpoint
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert that we got a 200 status code and expected response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}
