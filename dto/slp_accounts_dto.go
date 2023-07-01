package dto

import "time"

type SlpAccountResDTO struct {
	ID         uint      `json:"id"`
	CardNumber string    `json:"card_number"`
	NameOnCard string    `json:"name_on_card"`
	ActiveDate time.Time `json:"active_date"`
	IsDefault  bool      `json:"is_default"`
}

type SlpAccountReqDTO struct {
	UserID     uint
	CardNumber string    `json:"card_number" binding:"required"`
	NameOnCard string    `json:"name_on_card" binding:"required"`
	ActiveDate time.Time `json:"active_date" binding:"required"`
	IsDefault  bool
}
