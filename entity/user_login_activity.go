package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserLoginActivity struct {
	ID           uint `gorm:"primaryKey"`
	UserId       uint
	RefreshToken string
	IsValid      bool
	ExpiredAt    time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
