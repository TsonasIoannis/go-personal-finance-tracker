package controllers

import (
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/gin-gonic/gin"
)

func currentUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		httpapi.WriteError(c, apperrors.Unauthorized("authentication_required", "authentication required"))
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		httpapi.WriteError(c, apperrors.Unauthorized("invalid_authentication_context", "invalid authentication context"))
		return 0, false
	}

	return userID, true
}
