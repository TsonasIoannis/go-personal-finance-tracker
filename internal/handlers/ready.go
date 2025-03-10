package handlers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/gin-gonic/gin"
)

// ReadinessCheckHandler checks if the service is ready (i.e., DB is available)
func ReadinessCheckHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil || db.GetDB() == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "database is not initialized"})
			return
		}

		if err := db.CheckConnection(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}
