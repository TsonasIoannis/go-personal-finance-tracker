package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func currentUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication context"})
		return 0, false
	}

	return userID, true
}
