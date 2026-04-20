package middleware

import (
	"strings"
	"golang/utils/constant"
	"golang/utils/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header only
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate access token
		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user info from claims
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)

		// Set user info in context
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}