package entity

import (
	"time"

	"gorm.io/gorm"
)

type District struct {
	ID     uint `gorm:"primaryKey"`
	CityId uint
	City   City
	Name   string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
