package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WishlistID uuid.UUID `gorm:"type:uuid;not null;index" json:"wishlist_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	CreatedAt  time.Time `json:"created_at"`

	Wishlist Wishlist `gorm:"foreignKey:WishlistID" json:"-"`
	Product  Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (wi *WishlistItem) BeforeCreate(tx *gorm.DB) error {
	wi.ID = uuid.New()
	return nil
}
