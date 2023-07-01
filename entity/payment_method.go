package entity

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID   uint
	Name string
	Code string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
