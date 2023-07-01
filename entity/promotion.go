package entity

import (
	"time"

	"gorm.io/gorm"
)

type Promotion struct {
	ID              uint `gorm:"primary_key" json:"id"`
	MerchantId      uint `json:"merchant_id"`
	PromotionTypeId uint `json:"promotion_type_id"`
	PromotionType   PromotionType

	Title            string    `json:"title"`
	Nominal          float64   `json:"nominal"`
	MaxDiscountedQty int       `json:"max_discounted_price"`
	Quota            int       `json:"quota"`
	Quantity         int       `json:"quantity"`
	StartAt          time.Time `json:"start_at"`
	EndAt            time.Time `json:"end_at"`

	ProductPromotions []ProductPromotion

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
