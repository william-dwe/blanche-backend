package dto

const PAYMENT_METHOD_ID_WALLET = 1
const PAYMENT_METHOD_ID_SLP = 2
const PAYMENT_METHOD_ID_MERCHANT_WITHDRAW = 3
const PAYMENT_METHOD_ID_TRANSACTION_REFUND = 4

type PaymentMehodDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
