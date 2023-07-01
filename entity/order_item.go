package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrderItem struct {
	ID            uint   `gorm:"primaryKey"`
	OrderCode     string `gorm:"uniqueIndex"`
	ProductId     uint
	Product       Product
	VariantItemId uint `gorm:"foreignKey:VariantItemId;references:ID"`
	VariantItem   VariantItem
	Quantity      uint
	Notes         string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
