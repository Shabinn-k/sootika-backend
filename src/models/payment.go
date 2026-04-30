package models

import (
    "time"
    "github.com/google/uuid"
)

type Payment struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    OrderID       string    `json:"order_id"`        // Your internal order ID
    RazorpayOrderID string  `json:"razorpay_order_id"`
    RazorpayPaymentID string `json:"razorpay_payment_id,omitempty"`
    RazorpaySignature string `json:"razorpay_signature,omitempty"`
    Amount        int64     `json:"amount"`          // Amount in paise
    Currency      string    `json:"currency"`        // INR
    Status        string    `json:"status"`          // created, paid, failed
    UserID        uuid.UUID `json:"user_id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}