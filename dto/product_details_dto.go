package dto

type ProductDetailCategory struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ProductDetailRating struct {
	AvgRating float64 `json:"avg_rating"`
	Count     int     `json:"count"`
}

type ProductDetailDimension struct {
	Width  int `json:"width"`
	Length int `json:"length"`
	Height int `json:"height"`
}

type ProductDetailMerchant struct {
	Name       string  `json:"name"`
	Domain     string  `json:"domain"`
	Image      string  `json:"image"`
	AvgRating  float64 `json:"avg_rating"`
	SellerCity string  `json:"seller_city"`
}

type ProductDetailResDTO struct {
	ID               uint                   `json:"id"`
	Title            string                 `json:"title"`
	MinRealPrice     float64                `json:"min_real_price"`
	MaxRealPrice     float64                `json:"max_real_price"`
	MinDiscountPrice float64                `json:"min_discount_price"`
	MaxDiscountPrice float64                `json:"max_discount_price"`
	Category         ProductDetailCategory  `json:"category"`
	Images           []string               `json:"images"`
	Description      string                 `json:"description"`
	IsUsed           bool                   `json:"is_used"`
	SKU              string                 `json:"SKU"`
	FavouriteCount   int                    `json:"favourite_count"`
	UnitSold         int                    `json:"unit_sold"`
	TotalStock       int                    `json:"total_stock"`
	IsArchived       bool                   `json:"is_archived"`
	Rating           ProductDetailRating    `json:"rating"`
	Weight           int                    `json:"weight"`
	Dimension        ProductDetailDimension `json:"dimension"`

	IsMyProduct bool `json:"is_my_product"`
}

type ProductAdminDetailResDTO struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Price       *float64               `json:"price"`
	Categories  []uint                 `json:"categories"`
	Images      []string               `json:"images"`
	Description string                 `json:"description"`
	IsUsed      bool                   `json:"is_used"`
	TotalStock  int                    `json:"total_stock"`
	IsArchived  bool                   `json:"is_archived"`
	Weight      int                    `json:"weight"`
	Dimension   ProductDetailDimension `json:"dimension"`
}
