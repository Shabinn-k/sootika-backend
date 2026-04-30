package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    OrderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
    ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
    Product   Product   `json:"product,omitempty"`
    Title     string    `json:"title"`
    Name      string    `json:"name"`
    Image     string    `json:"image"`
    Price     int64     `json:"price"`
    Quantity  int       `json:"quantity"`
    Total     int64     `json:"total"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
    oi.ID = uuid.New()
    return nil
}