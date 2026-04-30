package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Address struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    City      string    `json:"city"`
    State     string    `json:"state"`
    Pincode   string    `json:"pincode"`
    Phone     string    `json:"phone"`
    IsDefault bool      `gorm:"default:false" json:"is_default"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
    a.ID = uuid.New()
    return nil
}