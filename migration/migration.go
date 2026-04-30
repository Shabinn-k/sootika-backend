package migration

import (
	"golang/src/models"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
    err := db.AutoMigrate(
        &models.User{},
        &models.RefreshToken{},
        &models.Product{},
        &models.Cart{},
        &models.CartItem{},
        &models.Wishlist{},
        &models.WishlistItem{},
        &models.Address{},      // Add this
        &models.Order{},        // Add this
        &models.OrderItem{},    // Add this
        &models.Payment{},      // Add this
    )
    if err != nil {
        log.Fatal("Migration failed:", err)
    }
    log.Println("Migration success")
}