package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	// ⚠️ CRITICAL FIX: Correct unique composite index
	CartID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_cart_product" json:"product_id"`

	Quantity int `gorm:"not null;default:1" json:"quantity"`
	Price    int `gorm:"not null" json:"price"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	ci.ID = uuid.New()
	return nil
}