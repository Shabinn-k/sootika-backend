package main

import (
    "github.com/gin-gonic/gin"
    "golang/config"
    "golang/internal/cache"
    "golang/internal/routes"
    "golang/migration"
    "golang/src/controllers"
    "golang/src/database"
    "golang/src/repository"
    "golang/src/services"
    "golang/utils/email"
    "golang/utils/jwt"
    "golang/utils/logger"
    "golang/utils/validation"
    "log"
)

func main() {

    cfg := config.LoadConfig()

    logger.InitLogger()

    validation.InitValidation()

    db := database.SetupDatabase(cfg)

    migration.Migrate(db)

    repo := repository.SetUpRepo(db)

    redis := cache.NewRedis()

    jwtManager := jwt.NewJWTManager(cfg)

    emailService := email.NewEmailService(cfg)

    authService := services.NewAuthService(repo, jwtManager, emailService, redis, cfg)
    productService := services.NewProductService(repo)
    wishlistService := services.NewWishlistService(repo)
    cartService := services.NewCartService(repo)
    paymentService := services.NewPaymentService(repo, cfg)
    orderService := services.NewOrderService(repo)
    addressService := services.NewAddressService(repo) // ⚠️ CREATE THIS

    authController := controllers.NewAuthController(authService)
    productController := controllers.NewProductController(productService)
    wishlistController := controllers.NewWishlistController(wishlistService)
    cartController := controllers.NewCartController(cartService)
    paymentController := controllers.NewPaymentController(paymentService)
    orderController := controllers.NewOrderController(orderService)
    addressController := controllers.NewAddressController(addressService) // ⚠️ PASS SERVICE
    adminController := controllers.NewAdminController(productService, repo)

    r := gin.Default()

    routes.SetUpRoutes(
        r,
        authController,
        jwtManager,
        productController,
        wishlistController,
        cartController,
        paymentController,
        addressController,
        orderController,
        adminController,
        repo,
    )

    logger.Log.Info("Server running on port", cfg.Server.Port)
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}