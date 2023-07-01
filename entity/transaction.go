package entity

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                   uint `gorm:"primary_key"`
	InvoiceCode          string
	MerchantDomain       string
	Merchant             Merchant `gorm:"foreignKey:MerchantDomain;references:Domain"`
	UserId               uint
	User                 User
	MerchantVoucherId    *uint
	MerchantVoucher      MerchantVoucher
	MarketplaceVoucherId *uint
	MarketplaceVoucher   MarketplaceVoucher
	CartItems            pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	PaymentMethod        pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	PaymentDetails       pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	DeliveryOption       pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`
	Address              pgtype.JSONB `gorm:"type:jsonb;default:'[]'"`

	TransactionStatus         *TransactionStatus
	TransactionDeliveryStatus *TransactionDeliveryStatus

	PaymentRecords []PaymentRecord `gorm:"many2many:transaction_payment_records;foreignKey:ID;joinForeignKey:TransactionId;references:PaymentId;joinReferences:PaymentId"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type TransactionCartItem struct {
	ProductId        uint    `json:"product_id"`
	Name             string  `json:"name"`
	Image            string  `json:"image"`
	Notes            *string `json:"notes"`
	RealPrice        float64 `json:"real_price"`
	DiscountPrice    float64 `json:"discount_price"`
	ProductSlug      string  `json:"product_slug"`
	ProductVariantId uint    `json:"product_variant_id"`
	VariantName      string  `json:"variant_name"`
	Quantity         int     `json:"quantity"`
}

type TransactionPaymentDetails struct {
	Subtotal                  float64 `json:"subtotal"`
	DeliveryFee               float64 `json:"delivery_fee"`
	MarketplaceVoucherNominal float64 `json:"marketplace_voucher_nominal"`
	MerchantVoucherNominal    float64 `json:"merchant_voucher_nominal"`
	Total                     float64 `json:"total"`
}

type TransactionPaymentMethod struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	Code                 string `json:"code"`
	AccountRelatedNumber string `json:"account_related_number"`
}

type TransactionDeliveryOption struct {
	CourierName string `json:"courier_name"`
}

type TransactionAddress struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Label           string `json:"label"`
	Details         string `json:"details"`
	ZipCode         string `json:"zip_code"`
	CityName        string `json:"city_name"`
	DistrictName    string `json:"district_name"`
	ProvinceName    string `json:"province_name"`
	SubdistrictName string `json:"subdistrict_name"`
}
