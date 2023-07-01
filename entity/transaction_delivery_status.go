package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionDeliveryStatus struct {
	ID            uint `gorm:"primary_key"`
	TransactionId uint
	OnDeliveryAt  *time.Time
	OnDeliveredAt *time.Time

	ReceiptNumber *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
