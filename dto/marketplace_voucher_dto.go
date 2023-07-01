package dto

import "time"

type MarketplaceVoucherResDTO struct {
	ID                 uint      `json:"id"`
	Code               string    `json:"code"`
	DiscountPercentage uint      `json:"discount_percentage"`
	MaxDiscountNominal float64   `json:"max_discount_nominal"`
	MinOrderNominal    float64   `json:"min_order_nominal"`
	ExpiredAt          time.Time `json:"expired_at"`
	Quota              int       `json:"quota"`
}

type MarketplaceAdminVoucherResDTO struct {
	ID                 uint      `json:"id"`
	Code               string    `json:"code"`
	MpDomain           string    `json:"mp_domain,omitempty"`
	CodeSuffix         string    `json:"code_suffix,omitempty"`
	DiscountPercentage uint      `json:"discount_percentage"`
	MaxDiscountNominal float64   `json:"max_discount_nominal"`
	MinOrderNominal    float64   `json:"min_order_nominal"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	Quota              int       `json:"quota"`
	UsedQuota          int       `json:"used_quota"`
}

type MarketplaceAdminVoucherListResDTO struct {
	PaginationResponse
	Vouchers []MarketplaceAdminVoucherResDTO `json:"vouchers"`
}

type UpsertMarketplaceVoucherReqDTO struct {
	Code               string    `json:"code"`
	DiscountPercentage uint      `json:"discount_percentage" binding:"required,min=1,max=100"`
	MaxDiscountNominal float64   `json:"max_discount_nominal" binding:"required"`
	StartDate          time.Time `json:"start_date" binding:"required"`
	EndDate            time.Time `json:"end_date" binding:"required"`
	Quota              int       `json:"quota" binding:"required"`
	MinOrderNominal    float64   `json:"min_order_nominal" binding:"required"`
}

type UpsertMarketplaceVoucherResDTO struct {
	Code               string    `json:"code"`
	DiscountPercentage uint      `json:"discount_percentage"`
	MaxDiscountNominal float64   `json:"max_discount_nominal"`
	MinOrderNominal    float64   `json:"min_order_nominal"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	Quota              int       `json:"quota"`
}
