package dto

import "time"

type TransactionDeliveryStatusResDTO struct {
	OnDeliveryAt  *time.Time `json:"on_delivery_at"`
	OnDeliveredAt *time.Time `json:"on_delivered_at"`
}
