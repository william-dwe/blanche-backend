package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantHoldingAccountHistory struct {
	ID                       uint `gorm:"primary_key"`
	MerchantHoldingAccountID uint
	Amount                   float64
	Type                     string
	Notes                    string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
