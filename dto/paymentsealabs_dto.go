package dto

type PaymentSeaLabsReqDTO struct {
	CardNumber string `json:"card_number" binding:"required"`
	Amount     int    `json:"amount" binding:"required"`
}

type PaymentSeaLabsResDTO struct {
	RedirectUrl string `json:"redirect_url"`
}

type PaymentSeaLabsCallback struct {
}
