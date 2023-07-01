package dto

import (
	"mime/multipart"
	"time"
)

const (
	FilterByProductReviewWithImage   = 1
	FilterByProductReviewWithComment = 2
)

type ProductReviewReqParamDTO struct {
	PaginationRequest
	FilterBy int  `form:"filter_by"`
	Rating   uint `form:"rating"`
}

type ProductReviewResDTO struct {
	PaginationResponse
	Reviews []ProductReviewDTO `json:"reviews"`
}

type ProductReviewDTO struct {
	Username           string     `json:"username"`
	UserProfilePicture *string    `json:"user_profile_picture"`
	ProductId          uint       `json:"product_id"`
	VariantItemId      uint       `json:"variant_item_id"`
	ProductName        string     `json:"product_name,omitempty"`
	ProductVariantName string     `json:"product_variant_name,omitempty"`
	ProductImgUrl      string     `json:"product_img_url,omitempty"`
	ProductPrice       uint       `json:"product_price,omitempty"`
	Rating             uint       `json:"rating"`
	Description        string     `json:"description"`
	ImageUrl           *string    `json:"image_url"`
	ReviewedAt         *time.Time `json:"reviewed_at"`
}

type ReviewProductFormReqDTO struct {
	ProductId     uint                  `form:"product_id" binding:"required"`
	VariantItemId uint                  `form:"variant_item_id" binding:"required"`
	Rating        uint                  `form:"rating" binding:"required,gte=1,lte=5"`
	Description   string                `form:"description" binding:"max=500"`
	Image         *multipart.FileHeader `form:"image,omitempty"`
}
