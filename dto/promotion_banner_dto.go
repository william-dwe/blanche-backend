package dto

import "mime/multipart"

type PromotionBannerResDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type PromotionBannerListResDTO struct {
	PaginationResponse
	PromotionBanners []PromotionBannerResDTO `json:"promotion_banners"`
}

type UpsertPromotionBannerReqDTO struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description" binding:"required"`
	Image       *multipart.FileHeader `form:"image,omitempty"`
}
