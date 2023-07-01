package dto

const SLP_SUCCESS_CODE = "TXN_PAID"
const SLP_FAILED_CODE = "TXN_FAILED"

type SealabspayResDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type SealabspayReqDTO struct {
	Amount       string `json:"amount"`
	MerchantCode string `json:"merchant_code"`
	Message      string `json:"message"`
	Signature    string `json:"signature"`
	Status       string `json:"status"`
	TxnId        string `json:"txn_id"`
}
