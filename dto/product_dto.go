package dto

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

type ProductListReqParamDTO struct {
	CategoryId       uint
	MerchantDomain   string  `form:"merchant,default=" binding:"max=255"`
	CategorySlug     string  `form:"cat,default=" binding:"max=255"`
	Search           string  `form:"q,default=" binding:"max=255"`
	SortBy           string  `form:"sort_by,default=avg_rating" binding:"oneof=avg_rating num_of_sale created_at min_discount_price"`
	SortDir          string  `form:"sort_dir,default=desc" binding:"oneof=asc desc"`
	MinPrice         float64 `form:"min_price,default=0" binding:"number,min=0,max=999999999999"`
	MaxPrice         float64 `form:"max_price,default=999999999999" binding:"number,min=0,max=999999999999"`
	MinRating        int     `form:"min_rating,default=0" binding:"number,min=0,max=5"`
	SellerCityId     string  `form:"seller_city_id,default="`
	IsMerchant       bool
	SellerCityIdList []uint
	Pagination       PaginationRequest
}

type ProductResDTO struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Slug             string  `json:"slug"`
	MinRealPrice     float64 `json:"min_real_price"`
	MaxRealPrice     float64 `json:"max_real_price"`
	MinDiscountPrice float64 `json:"min_discount_price"`
	MaxDiscountPrice float64 `json:"max_discount_price"`
	NumOfSale        int     `json:"num_of_sale"`
	AvgRating        float64 `json:"avg_rating"`
	ThumbnailImg     string  `json:"thumbnail_img"`
	SellerCity       string  `json:"seller_city"`
}

type ProductSellerResDTO struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	MinRealPrice float64 `json:"min_real_price"`
	MaxRealPrice float64 `json:"max_real_price"`
	NumOfSale    int     `json:"num_of_sale"`
	AvgRating    float64 `json:"avg_rating,omitempty"`
	TotalStock   int     `json:"total_stock"`
	ThumbnailImg string  `json:"thumbnail_img"`
	IsArchived   bool    `json:"is_archived,omitempty"`
}

type ProductSellerListResDTO struct {
	PaginationResponse
	Products []ProductSellerResDTO `json:"products"`
}

type ProductListResDTO struct {
	PaginationResponse
	Products []ProductResDTO `json:"products"`
}

type CreateProductReqDTO struct {
	Slug           string
	MerchantId     uint
	MerchantDomain string
	Title          string                       `json:"title" binding:"required,max=255"`
	Price          *float64                     `json:"price"`
	CategoryId     uint                         `json:"category_id" binding:"required"`
	Description    string                       `json:"description"`
	IsArchived     bool                         `json:"is_archived"`
	IsUsed         bool                         `json:"is_used"`
	TotalStock     int                          `json:"total_stock" binding:"min=0"`
	Weight         int                          `json:"weight" binding:"required,min=1"`
	Images         []string                     `json:"images"`
	Dimension      CreateProductDimensionReqDTO `json:"dimension" binding:"required"`
	Variant        CreateProductVariantReqDTO   `json:"variant"`
}

func (p *CreateProductReqDTO) Validate() error {
	if p.Price == nil && len(p.Variant.VariantItems) == 0 {
		return domain.ErrPriceVariantEmpty
	}

	if p.Price != nil && len(p.Variant.VariantItems) > 0 {
		return domain.ErrPriceVariantInvalid
	}

	if len(p.Variant.VariantItems) == 0 && p.Price != nil {
		if *p.Price <= 0 {
			return domain.ErrPriceVariantInvalidPrice
		}
		p.Variant.VariantItems = []CreateProductVariantItemsReqDTO{
			{
				Price: *p.Price,
				Stock: p.TotalStock,
			},
		}
	}

	return nil
}

type CreateProductDimensionReqDTO struct {
	Width  int `json:"width"`
	Length int `json:"length"`
	Height int `json:"height"`
}

type CreateProductVariantReqDTO struct {
	VariantOptions []CreateProductVariantOptionsReqDTO `json:"variant_options"`
	VariantItems   []CreateProductVariantItemsReqDTO   `json:"variant_items"`
}

type CreateProductVariantOptionsReqDTO struct {
	Name string   `json:"name" binding:"min=3,max=16"`
	Type []string `json:"type"`
}

type CreateProductVariantItemsReqDTO struct {
	Key   string  `json:"key"`
	Image string  `json:"image"`
	Price float64 `json:"price" binding:"min=100"`
	Stock int     `json:"stock"`
}

type CreateProductResDTO struct {
	ID             uint                         `json:"id"`
	Slug           string                       `json:"slug"`
	MerchantDomain string                       `json:"merchant_domain"`
	Title          string                       `json:"title"`
	MinRealPrice   float64                      `json:"min_real_price"`
	MaxRealPrice   float64                      `json:"max_real_price"`
	Description    string                       `json:"description"`
	IsArchived     bool                         `json:"is_archived"`
	IsUsed         bool                         `json:"is_used"`
	TotalStock     int                          `json:"total_stock"`
	Weight         int                          `json:"weight"`
	Images         []string                     `json:"images,omitempty"`
	Dimension      CreateProductDimensionReqDTO `json:"dimension"`
}

type CreateProductCheckNameResDTO struct {
	ProductName string `json:"product_name"`
	IsAvailable bool   `json:"is_available"`
}

type CheckMerchantProductNameReqDTO struct {
	ProductName string `json:"product_name" binding:"required"`
}

type UpdateProductAvailabilityReqDTO struct {
	IsArchived bool `json:"is_archived"`
}
