package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WishlistItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	WishlistID uuid.UUID `gorm:"type:uuid;not null;index:idx_wishlist_product,unique"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null;index:idx_wishlist_product,unique"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (wi *WishlistItem) BeforeCreate(tx *gorm.DB) error {
	wi.ID = uuid.New()
	return nil
}
