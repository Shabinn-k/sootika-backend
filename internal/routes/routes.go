package routes

import (
	"github.com/gin-gonic/gin"
	"golang/middleware"
	"golang/src/controllers"
	"golang/src/repository"
	"golang/utils/jwt"
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
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/verify", authController.VerifyOTP)
		auth.POST("/resend-otp", authController.ResendOTP)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", authController.Logout)
	}

	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	{
		user.GET("/dashboard", authController.Dashboard)
	}

	products := api.Group("/products")
	{
		products.GET("", productController.GetAllProducts)
		products.GET("/:id", productController.GetProductByID)
		products.GET("/search", productController.SearchProducts)
		products.GET("/in-stock", productController.GetInStockProducts)
	}

	wishlist := api.Group("/wishlist")
	wishlist.Use(middleware.AuthMiddleware(jwtManager))
	{
		wishlist.GET("", wishlistController.GetWishlist)
		wishlist.GET("/count", wishlistController.GetWishlistCount)
		wishlist.POST("", wishlistController.AddToWishlist)
		wishlist.GET("/:product_id", wishlistController.IsInWishlist)
		wishlist.DELETE("/:product_id", wishlistController.RemoveFromWishlist)
		wishlist.DELETE("", wishlistController.ClearWishlist)
	}

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware(jwtManager))
	{
		cart.GET("", cartController.GetCart)
		cart.GET("/count", cartController.GetCartCount)
		cart.GET("/total", cartController.GetCartTotal)
		cart.POST("", cartController.AddToCart)
		cart.PUT("/:item_id", cartController.UpdateCartItemQuantity)
		cart.DELETE("/:item_id", cartController.RemoveFromCart)
		cart.DELETE("", cartController.ClearCart)
	}

	payment := api.Group("/payment")
	{
		payment.POST("/create-order", middleware.AuthMiddleware(jwtManager), paymentController.CreateOrder)
		payment.POST("/verify", middleware.AuthMiddleware(jwtManager), paymentController.VerifyPayment)
		payment.POST("/webhook", paymentController.Webhook)
	}

	addresses := api.Group("/addresses")
	addresses.Use(middleware.AuthMiddleware(jwtManager))
	{
		addresses.GET("", addressController.GetMyAddresses)
		addresses.POST("", addressController.AddAddress)
		addresses.PUT("/:id", addressController.UpdateAddress)
		addresses.DELETE("/:id", addressController.DeleteAddress)
	}

	orders := api.Group("/orders")
	orders.Use(middleware.AuthMiddleware(jwtManager))
	{
		orders.POST("", orderController.CreateOrder)
		orders.GET("", orderController.GetMyOrders)
		orders.GET("/:id", orderController.GetOrderByID)
		orders.PUT("/:id/cancel", orderController.CancelOrder)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.AdminMiddleware(repo))
	{
		admin.GET("/dashboard", adminController.Dashboard)

		admin.GET("/users", adminController.GetAllUsers)
		admin.GET("/users/:id", adminController.GetUserByID)
		admin.PUT("/users/:id/role", adminController.UpdateUserRole)
		admin.PUT("/users/:id/block", adminController.ToggleBlockUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)

		admin.GET("/products/count", adminController.GetTotalProducts)
		admin.GET("/orders", adminController.GetAllOrders)
		admin.PUT("/orders/:id/status", adminController.UpdateOrderStatus)
 

		admin.POST("/products", productController.CreateProduct)
		admin.PUT("/products/:id", productController.UpdateProduct)
		admin.PUT("/products/:id/image/:type", productController.UpdateProductImage)
		admin.DELETE("/products/:id", productController.DeleteProduct)
	}
}
