package entity

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	ID            uint `gorm:"primary_key"`
	UserId        uint `gorm:"foreign_key"`
	ProductId     uint `gorm:"foreign_key"`
	Product       Product
	VariantItemId *uint `gorm:"foreign_key"`
	VariantItem   VariantItem
	Quantity      int
	IsChecked     bool
	Notes         *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
