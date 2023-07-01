package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserEmailBlacklist struct {
	ID        uint `gorm:"primary_key"`
	UserId    uint
	Email     string `gorm:"unique"`
	IsBlocked bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
