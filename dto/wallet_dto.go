package dto

const WALLET_TRANSACTION_DEBIT_CODE = "DR"
const WALLET_TRANSACTION_CREDIT_CODE = "CR"

const WALLET_ID_ADMIN = 90000
const WALLET_ID_ADMIN_PROMOTION = 90001

type CreateWalletPinReq struct {
	Username string `json:"-"`
	Pin      string `json:"pin" binding:"required,min=6,max=6"`
}

type CreateWalletPinRes struct {
	WalletId uint `json:"wallet_id"`
}

type WalletDetails struct {
	ID      uint    `json:"id"`
	Balance float64 `json:"balance"`
}

type WalletUpdatePin struct {
	Username string `json:"-"`
	NewPin   string `json:"new_pin" binding:"required,min=6,max=6"`
}

type TopUpWalletUsingSlpReqDTO struct {
	Amount        float64 `json:"amount" binding:"required,number,min=10000,max=2000000"`
	SlpCardNumber string  `json:"slp_card_number" binding:"required,len=16,numeric"`
}

type TopUpWalletUsingSlpResDTO struct {
	PaymentId      string  `json:"payment_id"`
	Amount         float64 `json:"amount"`
	WalletId       uint    `json:"wallet_id"`
	SlpCardNumber  string  `json:"slp_card_number"`
	SlpRedirectUrl string  `json:"slp_redirect_url"`
}
