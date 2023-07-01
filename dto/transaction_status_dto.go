package dto

import "time"

type TransactionStatusResDTO struct {
	OnWaitedAt        *time.Time `json:"on_waited_at"`
	OnProcessedAt     *time.Time `json:"on_processed_at"`
	OnDeliveredAt     *time.Time `json:"on_delivered_at"`
	OnCompletedAt     *time.Time `json:"on_completed_at"`
	OnCanceledAt      *time.Time `json:"on_canceled_at"`
	OnRefundedAt      *time.Time `json:"on_refunded_at"`
	OnRequestRefundAt *time.Time `json:"on_request_refund_at"`

	CancellationNotes string `json:"cancellation_notes"`
}
