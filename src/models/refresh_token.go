package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID uuid.UUID `gorm :"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Token string	`gorm:"type:text;not null"`
	ExpiresAt time.Time	
	CreatedAt time.Time  
}

func (r *RefreshToken)BeforeCreate(tx *gorm.DB)error{
	r.ID=uuid.New()
	return nil
}