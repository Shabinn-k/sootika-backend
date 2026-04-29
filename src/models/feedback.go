package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `json:"name"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	Feed      string    `json:"feed"`
	CreatedAt time.Time
}

func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	f.ID = uuid.New()
	return nil
}
