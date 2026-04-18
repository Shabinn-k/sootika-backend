package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Items     []WishlistItem `gorm:"foreignKey:WishlistID" json:"items,omitempty"`
}

func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	w.ID = uuid.New()
	return nil
}
