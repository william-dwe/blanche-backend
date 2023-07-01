package entity

import (
	"time"

	"gorm.io/gorm"
)

type DeliveryOption struct {
	ID          uint `gorm:"primaryKey"`
	CourierName string
	CourierCode string
	CourierLogo string
	ServiceCode string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
