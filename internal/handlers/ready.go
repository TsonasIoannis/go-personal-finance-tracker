package handlers

import (
	"net/http"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/database"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/gin-gonic/gin"
)

// ReadinessCheckHandler checks if the service is ready (i.e., DB is available)
// @Summary Readiness check
// @Description Readiness probe backed by the database connection.
// @Tags system
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 503 {object} httpapi.ErrorResponse
// @Router /ready [get]
func ReadinessCheckHandler(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			httpapi.WriteError(c, apperrors.Unavailable("database_unavailable", "database is not initialized"))
			return
		}

		if err := db.CheckConnection(); err != nil {
			if err.Error() == "database connection is not initialized" {
				httpapi.WriteError(c, apperrors.Unavailable("database_unavailable", "database is not initialized"))
				return
			}

			httpapi.WriteError(c, apperrors.Unavailable("database_unavailable", err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	}
}
