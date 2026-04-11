package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID         uint
    Name       string
    Email      string
    Password   string
    Role       string
    IsBlocked  bool
    IsVerified bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt
}