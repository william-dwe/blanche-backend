package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus struct {
	ID            uint `gorm:"primary_key"`
	TransactionId uint
	Transaction   Transaction

	OnWaitedAt        *time.Time
	OnProcessedAt     *time.Time
	OnDeliveredAt     *time.Time
	OnCompletedAt     *time.Time
	OnCanceledAt      *time.Time
	OnRefundedAt      *time.Time
	OnRequestRefundAt *time.Time

	CancellationNotes string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
