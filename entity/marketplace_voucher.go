package entity

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceVoucher struct {
	ID                 uint `gorm:"primary_key"`
	DiscountPercentage uint
	StartDate          time.Time
	ExpiredAt          time.Time
	IsInvalid          bool
	Quantity           int
	Quota              int
	Code               string
	MaxDiscountNominal float64
	MinOrderNominal    float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
