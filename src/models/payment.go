package models

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID           string    `json:"order_id"`
	RazorpayOrderID   string    `json:"razorpay_order_id"`
	RazorpayPaymentID string    `json:"razorpay_payment_id,omitempty"`
	RazorpaySignature string    `json:"razorpay_signature,omitempty"`
	Amount            int64     `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status" gorm:"type:varchar(20);default:created"` 
	UserID            uuid.UUID `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
