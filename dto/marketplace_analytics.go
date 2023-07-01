package dto

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"

type DashboardReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type MarketplaceAnalyticsActiveUserReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MarketplaceAnalyticsActiveUserResBody struct {
	Type  string `json:"type"`
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type MarketplaceAnalyticsUserConversionReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MarketplaceAnalyticsUserConversionResBody struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type MarketplaceAnalyticsSalesReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MarketplaceAnalyticsSalesResBody struct {
	Date string  `json:"date"`
	Rev  float64 `json:"rev"`
	Trx  int     `json:"trx"`
}

type MarketplaceAnalyticsCustomerSatisfactionReqBody struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}
type MarketplaceAnalyticsCustomerSatisfactionResBody struct {
	Date   string  `json:"date"`
	Review float64 `json:"review"`
	Count  int     `json:"count"`
}

type MarketplaceAnalyticsUpdateReqBody struct {
	DatePartition string `json:"date_partition"`
}

type MarketplaceAnalyticsUpdateResBody struct {
	entity.MerchantAnalytical
}
