package entity

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID         uint `gorm:"primaryKey"`
	ProvinceID uint
	Province   Province
	Name       string
	RoId       uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
