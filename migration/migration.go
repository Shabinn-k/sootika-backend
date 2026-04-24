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
		&models.Cart{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.WishlistItem{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration success")
}
