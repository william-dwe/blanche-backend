package dto

type UserFavoriteProductReqDTO struct {
	Username    string `json:"-"`
	ProductId   uint   `json:"product_id" binding:"required"`
	IsFavorited *bool  `json:"is_favorited" binding:"required"`
}

type UserFavoriteProductResDTO struct {
	ProductId   uint `json:"product_id"`
	IsFavorited bool `json:"is_favorited"`
}

type UserFavoriteProductReqParamDTO struct {
	Pagination PaginationRequest
	Search     string `form:"q,default=" binding:"max=255"`
	ProductId  uint   `form:"productId,default=0" binding:"number"`
}
