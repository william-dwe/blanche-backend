package dto

import "time"

const WALLET_TRANSACTION_TYPE_ID_TOP_UP_SLP = 1
const WALLET_TRANSACTION_TYPE_ID_MERCHANT_WITHDRAWAL = 2
const WALLET_TRANSACTION_TYPE_ID_TRANSACTION = 3
const WALLET_TRANSACTION_TYPE_ID_REFUND = 4

type WalletTransactionTypeDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type WalletTransactionRecordDTO struct {
	WalletTransactionTypeDTO WalletTransactionTypeDTO `json:"wallet_transaction_type"`
	PaymentId                *string                  `json:"payment_id"`
	Amount                   float64                  `json:"amount"`
	Notes                    string                   `json:"notes"`
	Title                    string                   `json:"title"`

	IssuedAt time.Time `json:"issued_at"`
}

type WalletTransactionReqParamDTO struct {
	PaginationRequest
	StartDate string `form:"start_date,default="`
	EndDate   string `form:"end_date,default="`
}

type WalletTransactionRecordDataDTO struct {
	PaginationResponse
	Transactions []WalletTransactionRecordDTO `json:"transactions"`
}
