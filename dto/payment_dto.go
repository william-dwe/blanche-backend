package dto

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
)

type WaitingForPaymentDTO struct {
	PaymentId   string    `json:"payment_id"`
	OrderCode   string    `json:"order_code"`
	Amount      uint      `json:"amount"`
	RedirectUrl string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
	PayBefore   time.Time `json:"pay_before"`

	PaymentMethod         string `json:"payment_method"`
	PaymentRelatedAccount string `json:"payment_related_account"`
}

type WaitingForPaymentTransactions struct {
	TransactionDetailProductResDTO
	PaymentDetails entity.TransactionPaymentDetails `json:"payment_details"`
}

type WaitingForPaymentDetailDTO struct {
	PaymentId             string    `json:"payment_id"`
	OrderCode             string    `json:"order_code"`
	Amount                uint      `json:"amount"`
	RedirectUrl           string    `json:"redirect_url"`
	CreatedAt             time.Time `json:"created_at"`
	PayBefore             time.Time `json:"pay_before"`
	PaymentMethod         string    `json:"payment_method"`
	PaymentRelatedAccount string    `json:"payment_related_account"`

	Transactions []WaitingForPaymentTransactions `json:"transactions"`
}
