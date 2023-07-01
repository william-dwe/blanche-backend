package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantVoucher struct {
	ID                 uint `gorm:"primary_key"`
	MerchantDomain     string
	DiscountNominal    float64
	Code               string
	StartDate          time.Time
	ExpiredAt          time.Time
	IsInvalid          bool
	Quantity           int
	Quota              int
	MinOrderNominal    float64
	MaxDiscountNominal float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
