package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantAnalytical struct {
	ID           uint `gorm:"primary_key"`
	MerchantId   uint
	AvgRating    float64
	NumOfSale    uint
	NumOfProduct uint
	NumOfReview  uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
