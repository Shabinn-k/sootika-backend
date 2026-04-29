package routes

import (
	"golang/middleware"
	"golang/src/controllers"
	"golang/src/repository"
	"golang/utils/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetUpRoutes(
    r *gin.Engine,
    authController *controllers.AuthController,
    jwtManager *jwt.Manager,
    productController *controllers.ProductController,
    wishlistController *controllers.WishlistController,
    cartController *controllers.CartController,
    adminController *controllers.AdminController,
    repo *repository.Repository,
)  {
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins for development
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/api/test", func(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "backend connected",
	})
})
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/check", authController.VerifyOTP)
		auth.POST("/resend-otp", authController.ResendOTP) 
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", authController.Logout)
	}

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	user.GET("/dashboard", authController.Dashboard)

	products := r.Group("/products")
	{
		products.GET("/", productController.GetAllProducts)
		products.GET("/:id", productController.GetProductByID)
		products.GET("/search", productController.SearchProducts)
		products.GET("/in-stock", productController.GetInStockProducts)
	}

	wishlist := r.Group("/wishlist")
	wishlist.Use(middleware.AuthMiddleware(jwtManager))
	{
		wishlist.GET("/", wishlistController.GetWishlist)
		wishlist.GET("/count", wishlistController.GetWishlistCount)
		wishlist.POST("/add", wishlistController.AddToWishlist)
		wishlist.GET("/check/:product_id", wishlistController.IsInWishlist)
		wishlist.DELETE("/remove/:product_id", wishlistController.RemoveFromWishlist)
		wishlist.DELETE("/clear", wishlistController.ClearWishlist)
	}

	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware(jwtManager))
	{
		cart.GET("/", cartController.GetCart)
		cart.GET("/count", cartController.GetCartCount)
		cart.GET("/total", cartController.GetCartTotal)
		cart.POST("/add", cartController.AddToCart)
		cart.PUT("/update/:item_id", cartController.UpdateCartItemQuantity)
		cart.DELETE("/remove/:item_id", cartController.RemoveFromCart)
		cart.DELETE("/clear", cartController.ClearCart)
	}

	admin := r.Group("/admin")
admin.Use(middleware.AuthMiddleware(jwtManager))
admin.Use(middleware.AdminMiddleware(repo))
{
    admin.GET("/dashboard", adminController.Dashboard)
    admin.GET("/users", adminController.GetAllUsers)
    admin.GET("/users/:id", adminController.GetUserByID)
    admin.PUT("/users/:id/role", adminController.UpdateUserRole)
    admin.PUT("/users/:id/toggle-block", adminController.ToggleBlockUser)
    admin.DELETE("/users/:id", adminController.DeleteUser)
    admin.GET("/stats/products", adminController.GetTotalProducts)

    admin.GET("/feedbacks", adminController.GetAllFeedbacks)
    admin.PUT("/feedbacks/:id/approve", adminController.ApproveFeedback)
    admin.DELETE("/feedbacks/:id", adminController.DeleteFeedback)


    admin.POST("/products", productController.CreateProduct)
    admin.PUT("/products/:id", productController.UpdateProduct)
    admin.PUT("/products/:id/image/:type", productController.UpdateProductImage)
    admin.DELETE("/products/:id", productController.DeleteProduct)
}
admin.GET("/test", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Admin endpoint working"})
})
}
