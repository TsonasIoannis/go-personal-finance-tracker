package handlers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/gin-gonic/gin"
)

// ReadinessCheckHandler checks if the service is ready (i.e., DB is available)
func ReadinessCheckHandler(c *gin.Context) {
	if err := database.CheckDBConnection(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
