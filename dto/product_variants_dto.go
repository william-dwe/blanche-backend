package dto

type ProductDetailVariantOption struct {
	Name string   `json:"name"`
	Type []string `json:"type"`
}

type ProductDetailVariantItem struct {
	ID            uint    `json:"id"`
	Key           string  `json:"key"`
	Image         string  `json:"image"`
	Price         float64 `json:"price"`
	DiscountPrice float64 `json:"discount_price"`
	Stock         uint    `json:"stock"`
}

type ProductDetailVariant struct {
	VariantOptions []ProductDetailVariantOption `json:"variant_options"`
	VariantItems   []ProductDetailVariantItem   `json:"variant_items"`
}
