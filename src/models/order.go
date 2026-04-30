package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
    ID              uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
    OrderNumber     string       `gorm:"uniqueIndex;not null" json:"order_number"`
    UserID          uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
    User            User         `json:"user,omitempty"`
    Items           []OrderItem  `gorm:"foreignKey:OrderID" json:"items"`
    Total           int64        `json:"total"`
    Subtotal        int64        `json:"subtotal"`
    ShippingCost    int64        `json:"shipping_cost"`
    Tax             int64        `json:"tax"`
    Discount        int64        `json:"discount"`
    PaymentMethod   string       `json:"payment_method"` // cod, razorpay, card, upi
    PaymentStatus   string       `json:"payment_status"` // pending, paid, failed
    OrderStatus     string       `json:"order_status"`   // pending, confirmed, shipped, delivered, cancelled
    Track           string       `json:"track"`
    ShippingAddress Address      `gorm:"embedded" json:"shipping_address"`
    RazorpayOrderID string       `json:"razorpay_order_id,omitempty"`
    RazorpayPaymentID string     `json:"razorpay_payment_id,omitempty"`
    CreatedAt       time.Time    `json:"created_at"`
    UpdatedAt       time.Time    `json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
    o.ID = uuid.New()
    return nil
}