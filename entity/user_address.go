package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserAddress struct {
	ID            uint `gorm:"primary_key"`
	UserID        uint `gorm:"not null"`
	ProvinceId    uint
	Province      Province
	CityId        uint
	City          City
	DistrictId    uint
	District      District
	SubdistrictId uint
	Subdistrict   Subdistrict
	Label         string
	Details       string
	Name          string
	PhoneNumber   string
	IsDefault     bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
