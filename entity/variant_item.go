package entity

import (
	"time"

	"gorm.io/gorm"
)

type VariantItem struct {
	ID        uint `gorm:"primarykey"`
	ProductId uint
	Price     float64
	ImageUrl  string
	Stock     uint

	VariantSpecs []VariantSpec

	VariantSpec VariantSpec

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
