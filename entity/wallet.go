package entity

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID                       uint `gorm:"primary_key"`
	UserId                   uint
	User                     User
	Balance                  float64
	Pin                      string `gorm:"not null"`
	WalletTransactionRecords []WalletTransactionRecord

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
