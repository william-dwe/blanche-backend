package dto

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"

type PostOrderSummaryMerchantsDTO struct {
	MerchantId      int    `json:"merchant_id"`
	VoucherMerchant string `json:"voucher_merchant"`
	DeliveryOption  string `json:"delivery_option"`
}

type PostOrderSummaryReqDTO struct {
	OrderCode string `json:"order_code"`

	AddressId          int                            `json:"address_id"`
	Merchants          []PostOrderSummaryMerchantsDTO `json:"merchants"`
	VoucherMarketplace string                         `json:"voucher_marketplace"`
}

type OrderItemDTO struct {
	CartItemId     uint    `json:"cart_item_id"`
	ProductId      uint    `json:"product_id"`
	ProductSlug    string  `json:"product_slug"`
	VariantItemId  *uint   `json:"variant_item_id"`
	VariantName    string  `json:"variant_name"`
	MerchantId     uint    `json:"merchant_id"`
	MerchantDomain string  `json:"merchant_domain"`
	MerchantName   string  `json:"merchant_name"`
	MerchantImage  string  `json:"merchant_image"`
	MerchantCityId uint    `json:"merchant_city_id"`
	Image          string  `json:"image"`
	Name           string  `json:"name"`
	Weight         int     `json:"weight"`
	RealPrice      float64 `json:"real_price"`
	DiscountPrice  float64 `json:"discount_price"`
	Quantity       int     `json:"quantity"`
	Stock          int     `json:"stock"`
	Notes          *string `json:"notes"`
	IsValid        bool    `json:"is_valid"`
}

type OrderMerchantDTO struct {
	MerchantId     uint   `json:"merchant_id"`
	MerchantName   string `json:"merchant_name"`
	MerchantImage  string `json:"merchant_image"`
	MerchantDomain string `json:"merchant_domain"`
}

type DeliveryServiceDTO struct {
	DeliveryOption string `json:"delivery_option"`
	Name           string `json:"name"`
	Service        string `json:"service"`
	Description    string `json:"description"`
	MerchantCity   string `json:"merchant_city"`
	UserCity       string `json:"user_city"`
	Etd            string `json:"etd"`
	Note           string `json:"note"`
}

type OrderItemPerMerchantDTO struct {
	Merchant          OrderMerchantDTO   `json:"merchant"`
	Items             []OrderItemDTO     `json:"items"`
	DeliveryService   DeliveryServiceDTO `json:"delivery_service"`
	SubTotal          float64            `json:"sub_total"`
	DeliveryCost      float64            `json:"delivery_cost"`
	Discount          float64            `json:"discount"`
	Total             float64            `json:"total"`
	IsVoucherInvalid  bool               `json:"is_voucher_invalid"`
	MerchantVoucherId *uint              `json:"-"`
}

type PostOrderSummaryResDTO struct {
	OrderCode string `json:"order_code"`

	Orders              []OrderItemPerMerchantDTO `json:"orders"`
	SubTotal            float64                   `json:"sub_total"`
	DeliveryCost        float64                   `json:"delivery_cost"`
	DiscountMerchant    float64                   `json:"discount_merchant"`
	DiscountMarketplace float64                   `json:"discount_marketplace"`
	Total               float64                   `json:"total"`
	IsVouchervalid      bool                      `json:"is_voucher_valid"`
	IsOrderEligible     bool                      `json:"is_order_eligible"`
	IsOrderValid        bool                      `json:"is_order_valid"`

	MarketplaceVoucherId *uint              `json:"-"`
	Address              entity.UserAddress `json:"-"`
}

type MakeOrderCheckoutProductDTO struct {
	ProductId     uint   `json:"product_id"`
	VariantItemId *uint  `json:"variant_item_id"`
	Quantity      int    `json:"quantity"`
	Notes         string `json:"notes"`
}
