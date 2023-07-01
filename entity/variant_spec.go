package entity

import (
	"time"

	"gorm.io/gorm"
)

type VariantSpec struct {
	ID               uint `gorm:"primarykey"`
	VariantItemID    uint
	VariationGroupID uint
	VariationGroup   VariationGroup
	VariationName    string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
