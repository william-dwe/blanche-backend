package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductReview struct {
	ID            uint `gorm:"primarykey"`
	ProductID     uint
	Product       Product
	VariantItemID uint
	VariantItem   VariantItem
	TransactionID uint
	Transaction   Transaction
	Rating        int
	Description   string
	ImageUrl      *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
