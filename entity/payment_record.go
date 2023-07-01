package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentRecord struct {
	ID              uint   `gorm:"primary_key"`
	PaymentId       string `gorm:"not null"`
	PaymentMethodId uint   `gorm:"not null"`
	PaymentMethod   PaymentMethod
	Amount          float64 `gorm:"not null"`
	PaymentUrl      string

	Transactions             []Transaction            `gorm:"many2many:transaction_payment_records;references:ID;joinReferences:TransactionId;foreignKey:PaymentId;joinForeignKey:PaymentId;"`
	TransactionPaymentRecord TransactionPaymentRecord `gorm:"foreignKey:PaymentId;references:PaymentId;"`

	PaidAt     *time.Time
	CanceledAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
