package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Price       int64     `gorm:"not null" json:"price"`

	MainImage   string `gorm:"not null" json:"main_image"`
	SecondImage string `json:"second_image,omitempty"`
	ThirdImage  string `json:"third_image,omitempty"`

	MainImagePublicID   string `json:"-"`
	SecondImagePublicID string `json:"-"`
	ThirdImagePublicID  string `json:"-"`

	InStock   bool           `gorm:"default:true" json:"in_stock"`
	Stock     int            `gorm:"default:0" json:"stock"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}
