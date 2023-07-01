package entity

import (
	"time"

	"gorm.io/gorm"
)

type MerchantDeliveryOption struct {
	ID               uint `gorm:"primaryKey"`
	MerchantId       uint
	DeliveryOptionId uint
	DeliveryOption   DeliveryOption

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
