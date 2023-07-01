package entity

import (
	"time"

	"gorm.io/gorm"
)

type PromotionBanner struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
