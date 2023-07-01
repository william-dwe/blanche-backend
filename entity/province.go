package entity

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID   uint `gorm:"primaryKey"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
