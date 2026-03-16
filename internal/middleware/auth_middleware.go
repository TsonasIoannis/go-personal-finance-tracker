package middleware

import (
	"strings"

	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/apperrors"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/auth"
	"github.com/TsonasIoannis/go-personal-finance-tracker/internal/httpapi"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenManager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpapi.AbortWithError(c, apperrors.Unauthorized("missing_authorization_header", "missing authorization header"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpapi.AbortWithError(c, apperrors.Unauthorized("invalid_authorization_header", "invalid authorization header"))
			return
		}

		claims, err := tokenManager.ParseToken(parts[1])
		if err != nil {
			httpapi.AbortWithError(c, apperrors.Unauthorized("invalid_token", "invalid or expired token"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Next()
	}
}
