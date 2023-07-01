package dto

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
)

type MerchantAnalyticsMerchantResponsivenessReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MerchantAnalyticsMerchantResponsivenessResBody struct {
	Type  string `json:"type"`
	Date  string `json:"date"`
	Value string `json:"value"`
}

type MerchantAnalyticsUserConversionReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MerchantAnalyticsUserConversionResBody struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type MerchantAnalyticsSalesReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MerchantAnalyticsSalesResBody struct {
	Date string  `json:"date"`
	Rev  float64 `json:"rev"`
	Trx  int     `json:"trx"`
}

type MerchantAnalyticsCustomerSatisfactionReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MerchantAnalyticsCustomerSatisfactionResBody struct {
	Date   string  `json:"date"`
	Review float64 `json:"review"`
	Count  int     `json:"count"`
}

type MerchantAnalyticsUpdateReqBody struct {
	DatePartition string `json:"date_partition"`
}

type MerchantAnalyticsUpdateResBody struct {
	entity.MerchantAnalytical
}
