package dto

import (
	"mime/multipart"
	"time"
)

type MerchantAddressDTO struct {
	Province string `json:"province"`
	City     string `json:"city"`
}

type MerchantInfoResDTO struct {
	ID           uint               `json:"id"`
	Domain       string             `json:"domain"`
	Name         string             `json:"name"`
	Address      MerchantAddressDTO `json:"address"`
	AvgRating    float64            `json:"avg_rating"`
	JoinDate     string             `json:"join_date"`
	NumOfProduct uint               `json:"num_of_product"`
	NumOfSale    uint               `json:"num_of_sale"`
	NumOfReview  uint               `json:"num_of_review"`
	Image        string             `json:"image"`
}

type MerchantProductCategory struct {
	CategoryId uint `json:"category_id"`
	Quantity   int  `json:"quantity"`
}

type RegisterMerchantReqDTO struct {
	Name      string `json:"name" binding:"required"`
	Domain    string `json:"domain" binding:"required"`
	AddressId uint   `json:"address_id" binding:"required"`
}

type RegisterMerchantResDTO struct {
	Name     string             `json:"name"`
	Domain   string             `json:"domain"`
	Address  MerchantAddressDTO `json:"address"`
	JoinDate time.Time          `json:"join_date"`
}

type CheckMerchantDomainReqDTO struct {
	Domain string `json:"domain" binding:"required"`
}

type CheckMerchantDomainResDTO struct {
	Domain      string `json:"domain"`
	IsAvailable bool   `json:"is_available"`
}

type UpdateMerchantProfileFormReqDTO struct {
	Name        *string               `form:"name"`
	Description *string               `form:"description"`
	Image       *multipart.FileHeader `form:"image,omitempty"`
}

type UpdateMerchantProfileResDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type CheckMerchantStoreNameReqDTO struct {
	Name string `json:"name" binding:"required"`
}

type CheckMerchantStoreNameResDTO struct {
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}

type MerchantVoucherResDTO struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`
	DiscountNominal float64   `json:"discount_nominal"`
	ExpiredAt       time.Time `json:"expired_at"`
	Quota           int       `json:"quota"`
	MinOrderNominal float64   `json:"min_order_nominal"`
}

type MerchantAdminVoucherResDTO struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`
	MerchantDomain  string    `json:"merchant_domain,omitempty"`
	CodeSuffix      string    `json:"code_suffix,omitempty"`
	DiscountNominal float64   `json:"discount_nominal"`
	MinOrderNominal float64   `json:"min_order_nominal"`
	StartDate       time.Time `json:"start_date"`
	ExpiredAt       time.Time `json:"expired_at"`
	Quota           int       `json:"quota"`
	UsedQuota       int       `json:"used_quota"`
}

type MerchantAdminVoucherListResDTO struct {
	PaginationResponse
	Vouchers []MerchantAdminVoucherResDTO `json:"vouchers"`
}

type UpsertMerchantVoucherReqDTO struct {
	Code            string    `json:"code"`
	DiscountNominal float64   `json:"discount_nominal" binding:"required"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
	Quota           int       `json:"quota" binding:"required"`
	MinOrderNominal float64   `json:"min_order_nominal" binding:"required"`
}

type UpsertMerchantVoucherResDTO struct {
	ID              uint      `json:"id"`
	MerchantDomain  string    `json:"merchant_domain"`
	Code            string    `json:"code"`
	DiscountNominal float64   `json:"discount_nominal"`
	MinOrderNominal float64   `json:"min_order_nominal"`
	StartDate       time.Time `json:"start_date"`
	ExpiredAt       time.Time `json:"expired_at"`
	Quantity        int       `json:"quantity"`
	Quota           int       `json:"quota"`
}

type MerchantVoucherListParamReqDTO struct {
	PaginationRequest
	Status uint `form:"status"`
}

const (
	VoucherStatusOngoing  uint = 1
	VoucherStatusIncoming uint = 2
	VoucherStatusExpired  uint = 3
)
