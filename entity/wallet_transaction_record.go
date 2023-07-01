package entity

import (
	"time"

	"gorm.io/gorm"
)

type WalletTransactionRecord struct {
	ID                      uint `gorm:"primary_key"`
	WalletId                uint
	WalletTransactionTypeId uint
	WalletTransactionType   WalletTransactionType
	PaymentId               string
	PaymentRecord           PaymentRecord `gorm:"foreignKey:PaymentId;references:PaymentId;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
