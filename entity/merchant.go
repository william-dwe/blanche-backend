package entity

import (
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	ID                 uint `gorm:"primary_key"`
	UserId             uint
	UserAddressId      uint
	UserAddress        UserAddress
	CityId             uint
	City               City
	MerchantAnalytical MerchantAnalytical `gorm:"foreignKey:MerchantId"`
	Domain             string             `gorm:"unique"`
	Name               string             `gorm:"unique"`
	Description        string
	JoinDate           time.Time
	ImageUrl           string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
