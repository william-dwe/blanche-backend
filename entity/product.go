package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                 uint   `gorm:"primaryKey"`
	Slug               string `gorm:"uniqueIndex"`
	MerchantId         uint
	MerchantDomain     string `gorm:"foreignKey:MerchantDomain;references:Domain;"`
	Merchant           Merchant
	ProductAnalyticID  uint
	CategoryId         uint
	Category           Category
	Title              string
	MinRealPrice       float64
	MaxRealPrice       float64
	MinDiscountedPrice float64 `gorm:"<-:false;"`
	MaxDiscountedPrice float64 `gorm:"<-:false;"`
	IsArchived         bool
	Description        string
	Weight             int
	Height             int
	Width              int
	Length             int
	IsUsed             bool
	SKU                string

	ProductImages    []ProductImage
	ProductAnalytic  ProductAnalytic
	ProductPromotion *ProductPromotion
	VariantItems     []VariantItem

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
