package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantHoldingAccount struct {
	ID         uint `gorm:"primary_key"`
	MerchantID uint
	Balance    float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
