package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductPromotion struct {
	ID                 uint `gorm:"primaryKey"`
	PromotionId        uint
	Promotion          Promotion
	ProductId          uint
	Product            Product
	MinDiscountedPrice float64
	MaxDiscountedPrice float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
