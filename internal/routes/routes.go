// routes.go - FIXED payment routes
package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang/middleware"
	"golang/src/controllers"
	"golang/src/repository"
	"golang/utils/jwt"
	"time"
)

func SetUpRoutes(
	r *gin.Engine,
	authController *controllers.AuthController,
	jwtManager *jwt.Manager,
	productController *controllers.ProductController,
	wishlistController *controllers.WishlistController,
	cartController *controllers.CartController,
	paymentController *controllers.PaymentController,
	addressController *controllers.AddressController,
	orderController *controllers.OrderController,
	adminController *controllers.AdminController,
	repo *repository.Repository,
) {
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:   func(origin string) bool { return true },
		AllowMethods:      []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:      []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:     []string{"Content-Length"},
		AllowCredentials:  true,
		MaxAge:            12 * time.Hour,
	}))

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "backend connected"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/check", authController.VerifyOTP)
		auth.POST("/resend-otp", authController.ResendOTP)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", authController.Logout)
	}

	// User routes
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	user.GET("/dashboard", authController.Dashboard)

	// Product routes (public)
	products := r.Group("/products")
	{
		products.GET("/", productController.GetAllProducts)
		products.GET("/:id", productController.GetProductByID)
		products.GET("/search", productController.SearchProducts)
		products.GET("/in-stock", productController.GetInStockProducts)
	}

	// Wishlist routes
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

	// Cart routes
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

	// ⚠️ FIXED: Payment routes - protected except webhook
	payment := r.Group("/payment") // Changed from /api/payment
	{
		payment.POST("/create-order", middleware.AuthMiddleware(jwtManager), paymentController.CreateOrder)
		payment.POST("/verify", middleware.AuthMiddleware(jwtManager), paymentController.VerifyPayment)
		payment.POST("/webhook", paymentController.Webhook) // NO auth for webhook
	}

	// Address routes
	addresses := r.Group("/addresses")
	addresses.Use(middleware.AuthMiddleware(jwtManager))
	{
		addresses.GET("/", addressController.GetMyAddresses)
		addresses.POST("/", addressController.AddAddress)
		addresses.PUT("/:id", addressController.UpdateAddress)
		addresses.DELETE("/:id", addressController.DeleteAddress)
	}

	// Order routes
	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware(jwtManager))
	{
		orders.POST("/", orderController.CreateOrder)
		orders.GET("/", orderController.GetMyOrders)
		orders.GET("/:id", orderController.GetOrderByID)
		orders.PUT("/:id/cancel", orderController.CancelOrder)
	}

	// Admin routes
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
		admin.GET("/orders", adminController.GetAllOrders)
		admin.GET("/feedbacks", adminController.GetAllFeedbacks)
		admin.PUT("/feedbacks/:id/approve", adminController.ApproveFeedback)
		admin.DELETE("/feedbacks/:id", adminController.DeleteFeedback)
		admin.POST("/products", productController.CreateProduct)
		admin.PUT("/products/:id", productController.UpdateProduct)
		admin.PUT("/products/:id/image/:type", productController.UpdateProductImage)
		admin.DELETE("/products/:id", productController.DeleteProduct)
	}
}