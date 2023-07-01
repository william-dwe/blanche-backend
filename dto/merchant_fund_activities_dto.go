package dto

import "time"

const MERCHANT_HOLDING_ACC_DEBIT_CODE = "DR"
const MERCHANT_HOLDING_ACC_CREDIT_CODE = "CR"

type MerchantFundActivitiesDTO struct {
	ID       uint      `json:"id"`
	Notes    string    `json:"note"`
	Amount   float64   `json:"amount"`
	Type     string    `json:"type"`
	IssuedAt time.Time `json:"issued_at"`
}

type MerchantFundActivitiesReqParamDTO struct {
	PaginationRequest
	StartDate string `form:"start_date,default="`
	EndDate   string `form:"end_date,default="`
	Type      string `form:"type,default="`
}

type MerchantFundActivitiesResDTO struct {
	PaginationResponse
	FundActivities []MerchantFundActivitiesDTO `json:"fund_activities"`
}

type MerchantFundBalanceDTO struct {
	TotalBalance float64 `json:"total_balance"`
}

type MerchantWithdrawReqDTO struct {
	Amount uint `json:"amount" binding:"required,gte=10000"`
}

type MerchantWithdrawResDTO struct {
	ID     uint    `json:"id"`
	Amount float64 `json:"amount"`
	Notes  string  `json:"notes"`
}
