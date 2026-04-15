package middleware

import (
	"golang/utils/constant"
	"golang/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "invalid authorization"})
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
