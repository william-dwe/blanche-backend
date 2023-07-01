package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserOrder struct {
	ID         uint        `gorm:"primaryKey"`
	OrderCode  string      `gorm:"uniqueIndex"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderCode;references:OrderCode;"`
	UserId     uint
	User       User

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
