package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserFavoriteProduct struct {
	ID        uint `gorm:"primaryKey"`
	UserId    uint
	User      User
	ProductId uint
	Product   Product

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
