package entity

import (
	"time"

	"gorm.io/gorm"
)

type Subdistrict struct {
	ID         uint `gorm:"primaryKey"`
	DistrictId uint
	District   District
	Name       string
	ZipCode    string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
