package dto

type WalletpayReqDTO struct {
	Amount    uint   `json:"amount" binding:"required,gte=1"`
	PaymentId string `json:"payment_id" binding:"required"`
}

type WalletpayResDTO struct {
	Amount    uint   `json:"amount"`
	PaymentId string `json:"payment_id"`
	Status    string `json:"status"`
}
