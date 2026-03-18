package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler checks if the service is running
// @Summary Health check
// @Description Liveness probe for the API process.
// @Tags system
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /health [get]
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
