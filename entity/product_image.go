package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID        uint `gorm:"primarykey"`
	ProductID uint
	ImageUrl  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
