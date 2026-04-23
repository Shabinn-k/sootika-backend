package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	CartID    uuid.UUID `gorm:"type:uuid;not null;index:idx_cart_product,unique" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_cart_product,unique" json:"product_id"`

	Quantity int `gorm:"not null;default:1" json:"quantity"`
	Price    int `gorm:"not null" json:"price"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	ci.ID = uuid.New()
	return nil
}
