package dto

import "time"

type PromotionResDTO struct {
	ID                    uint                  `json:"id"`
	PromotionType         string                `json:"promotion_type"`
	Title                 string                `json:"title"`
	MaxDiscountedQuantity int                   `json:"max_discounted_quantity"`
	DiscountNominal       float64               `json:"discount_nominal,omitempty"`
	DiscountPercentage    float64               `json:"discount_percentage,omitempty"`
	Quota                 int                   `json:"quota"`
	UsedQuota             int                   `json:"used_quota"`
	StartDate             time.Time             `json:"start_date"`
	EndDate               time.Time             `json:"end_date"`
	Products              []ProductSellerResDTO `json:"products"`
}

type PromotionListResDTO struct {
	PaginationResponse
	Promotions []PromotionResDTO `json:"promotions"`
}

type PromotionListReqParamDTO struct {
	PaginationRequest
	Status uint `form:"status"`
}

type UpsertPromotionReqDTO struct {
	ProductIds            []uint    `json:"product_ids" binding:"required"`
	PromotionTypeId       uint      `json:"promotion_type_id" binding:"required"`
	Title                 string    `json:"title" binding:"required"`
	MaxDiscountedQuantity int       `json:"max_discounted_quantity" binding:"required"`
	Nominal               float64   `json:"nominal" binding:"required"`
	Quota                 int       `json:"quota" binding:"required"`
	StartDate             time.Time `json:"start_date" binding:"required"`
	EndDate               time.Time `json:"end_date" binding:"required"`
	MerchantId            uint
}

type PromotionDetailResDTO struct {
	ID                    uint                  `json:"id"`
	PromotionTypeId       uint                  `json:"promotion_type_id"`
	ProductIds            []uint                `json:"product_ids"`
	Title                 string                `json:"title"`
	MaxDiscountedQuantity int                   `json:"max_discounted_quantity"`
	Nominal               float64               `json:"nominal"`
	Quota                 int                   `json:"quota"`
	UsedQuota             int                   `json:"used_quota"`
	StartDate             time.Time             `json:"start_date"`
	EndDate               time.Time             `json:"end_date"`
	Products              []ProductSellerResDTO `json:"products"`
}
