package middleware

import (
	"golang/src/models"
	"golang/src/repository"
	"golang/utils/constant"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware(repo repository.PgSQLRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists || userID == "" {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}
		var user models.User
		if err := repo.FindOneWhere(&user, "id = ?", userID); err != nil {
			c.JSON(constant.UNAUTHORIZED, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if user.IsBlocked {
			c.JSON(constant.FORBIDDEN, gin.H{"error": "Your account has been blocked. Please contact support."})
			c.Abort()
			return
		}
		if user.Role != "admin" {
			c.JSON(constant.FORBIDDEN, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
