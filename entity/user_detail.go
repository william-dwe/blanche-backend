package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserDetail struct {
	ID             uint `gorm:"primary_key"`
	UserId         uint
	Fullname       string
	Phone          *string    `gorm:"default:null"`
	Gender         *string    `gorm:"default:null"`
	BirthDate      *time.Time `gorm:"default:null"`
	ProfilePicture *string    `gorm:"default:null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
