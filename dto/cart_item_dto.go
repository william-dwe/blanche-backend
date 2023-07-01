package dto

type AddItemToCartReqDTO struct {
	ProductId     uint  `json:"product_id"`
	VariantItemId *uint `json:"variant_item_id,omitempty"`
	Quantity      int   `json:"quantity"`
}

type AddItemToCartResDTO struct {
	ProductId     uint   `json:"product_id"`
	VariantItemId *uint  `json:"variant_item_id"`
	Quantity      int    `json:"quantity"`
	IsChecked     bool   `json:"is_checked"`
	Notes         string `json:"notes"`
}

type CartItemDTO struct {
	CartItemId            uint    `json:"cart_item_id"`
	ProductId             uint    `json:"product_id"`
	ProductSlug           string  `json:"product_slug"`
	VariantItemId         *uint   `json:"variant_item_id"`
	VariantName           string  `json:"variant_name"`
	MerchantId            uint    `json:"merchant_id"`
	MerchantDomain        string  `json:"merchant_domain"`
	MerchantName          string  `json:"merchant_name"`
	MerchantImage         string  `json:"merchant_image"`
	Image                 string  `json:"image"`
	Name                  string  `json:"name"`
	RealPrice             float64 `json:"real_price"`
	DiscountPrice         float64 `json:"discount_price"`
	Quantity              int     `json:"quantity"`
	Stock                 int     `json:"stock"`
	Notes                 *string `json:"notes"`
	IsChecked             bool    `json:"is_checked"`
	IsValid               bool    `json:"is_valid"`
	IsPromotionPriceValid bool    `json:"is_promotion_price_valid"`
}

type CartItemPerStoreDTO struct {
	MerchantId     uint          `json:"merchant_id"`
	MerchantName   string        `json:"merchant_name"`
	MerchantImage  string        `json:"merchant_image"`
	MerchantDomain string        `json:"merchant_domain"`
	Items          []CartItemDTO `json:"items"`
}

type GetCartItemResDTO struct {
	Carts    []CartItemPerStoreDTO `json:"carts"`
	Quantity int                   `json:"quantity"`
	Total    float64               `json:"total"`
}

type GetHomeCartItemResDTO struct {
	Carts    []CartItemDTO `json:"carts"`
	Quantity int           `json:"quantity"`
}

type UpdateAllCartCheckStatusDTO struct {
	CartItemId uint `json:"cart_item_id"`
	IsChecked  bool `json:"is_checked"`
}

type UpdateCartItemDTO struct {
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}
