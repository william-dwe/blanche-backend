package entity

import (
	"time"

	"gorm.io/gorm"
)

type VariationGroup struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
