package routes

import (
	"golang/middleware"
	"golang/src/controllers"
	"golang/src/repository"
	"golang/utils/jwt"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(
	r *gin.Engine,
	authController *controllers.AuthController,
	jwtManager *jwt.Manager,
	productController *controllers.ProductController,
	wishlistController *controllers.WishlistController,
	cartController *controllers.CartController,
	repo *repository.Repository,
) {

	auth := r.Group("/auth")
	{
		auth.POST("/signup", authController.Signup)
		auth.POST("/check", authController.VerifyOTP)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
		auth.POST("/logout", authController.Logout)
	}

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtManager))
	user.GET("/dashboard", authController.Dashboard)

	products := r.Group("/products")
	{
		products.GET("/", productController.GetAllProduct)
		products.GET("/:id", productController.GetProductByID)
		products.GET("/search", productController.SearchProducts, productController.GetProductsByTitle)
		products.GET("stock", productController.GetInStockProducts)
	}
	wishlist := r.Group("/wishlist")
	wishlist.Use(middleware.AuthMiddleware(jwtManager))
	{
		wishlist.GET("/", wishlistController.GetWishlist)
		wishlist.GET("/count", wishlistController.GetWishlistCount)
		wishlist.POST("/add", wishlistController.AddToWishlist)
		wishlist.DELETE("/remove/:id", wishlistController.RemoveFromWishlist)
		wishlist.GET("/check/:id", wishlistController.IsInWishlist)
		wishlist.DELETE("/clear", wishlistController.ClearWishlist)
	}
	cart := r.Group("/cart")
	cart.Use(middleware.AuthMiddleware(jwtManager))
	{
		cart.GET("/", cartController.GetCart)
		cart.GET("/count", cartController.GetCartCount)
		cart.GET("/total", cartController.GetCartTotal)
		cart.POST("/add", cartController.AddToCart)
		cart.PUT("/update/:id", cartController.UpdateCartItemQuantity)
		cart.DELETE("remove/:id", cartController.RemoveFromCart)
		cart.DELETE("/clear", cartController.ClearCart)
	}
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.AdminMiddleware(repo))
	protected := admin.Group("/products")
	{
		protected.POST("/", productController.CreateProduct)
		protected.PUT("/:id", productController.UpdateProduct)
		protected.DELETE("/:id", productController.DeleteProduct)
		protected.PATCH("/:id/stock", productController.UpdateProductStock)
	}

}
