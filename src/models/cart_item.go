package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Price     int       `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Cart    Cart    `gorm:"foreignKey:CartID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	ci.ID = uuid.New()
	return nil
}

