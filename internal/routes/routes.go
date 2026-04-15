package routes

import (
	"golang/middleware"
	"golang/src/controllers"
	"golang/utils/jwt"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine, authController *controllers.AuthController, jwtManager *jwt.Manager) {
	auth := r.Group("/auth")
	auth.POST("/signup", authController.Signup)
	auth.POST("/check", authController.VerifyOTP)
	auth.POST("/login", authController.Login)
	auth.POST("/refresh", authController.Refresh)
	auth.POST("/logout", authController.Logout)

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	user.GET("/dashboard", authController.Dashboard)
}
