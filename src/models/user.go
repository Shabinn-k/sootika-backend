package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	Name     string `json:"name" binding:"required,min=2,max=50,name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `binding:"required,min=6,password" gorm:"not null" json:"-"`
	Phone    string `gorm:"type:varchar(15);uniqueIndex;not null"`

	Role       string `gorm:"type:varchar(20);default:user" json:"role"`
	IsBlocked  bool   `gorm:"default:false" json:"is_blocked"`
	IsVerified bool   `gorm:"default:false" json:"is_verified"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
