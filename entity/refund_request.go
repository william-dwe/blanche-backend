package entity

import (
	"time"

	"gorm.io/gorm"
)

type RefundRequest struct {
	ID            uint `gorm:"primarykey"`
	TransactionID uint
	Transaction   Transaction

	RefundRequestStatuses []RefundRequestStatus `gorm:"foreignKey:RefundRequestId"`

	Reason   string
	ImageUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type RefundRequestStatus struct {
	ID              uint `gorm:"primarykey"`
	RefundRequestId uint
	RefundRequest   RefundRequest

	CanceledByBuyerAt *time.Time
	AcceptedByBuyerAt *time.Time
	RejectedByBuyerAt *time.Time

	AcceptedBySellerAt *time.Time
	RejectedBySellerAt *time.Time

	AcceptedByAdminAt *time.Time
	RejectedByAdminAt *time.Time

	ClosedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type RefundReqMessageRole struct {
	ID       uint `gorm:"primarykey"`
	RoleName string
}

type RefundReqMessage struct {
	ID                     uint `gorm:"primarykey"`
	RefundRequestId        uint
	RefundRequest          RefundRequest
	RefundReqMessageRoleId uint
	RefundReqMessageRole   RefundReqMessageRole

	Message  string
	ImageUrl *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
