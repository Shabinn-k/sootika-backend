package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAddress struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	AddressLine string    `gorm:"not null" json:"address_line"`
	City        string    `gorm:"not null" json:"city"`
	State       string    `gorm:"not null" json:"state"`
	Pincode     string    `gorm:"not null" json:"pincode"`
	Country     string    `gorm:"default:India" json:"country"`
	Phone       string    `json:"phone"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (a *UserAddress) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
