package dto

import (
	"mime/multipart"
	"time"
)

const REFUND_REQ_MSG_ROLE_ADMIN_ID = 1
const REFUND_REQ_MSG_ROLE_MERCHANT_ID = 2
const REFUND_REQ_MSG_ROLE_BUYER_ID = 3

var RefundRequestStatusMap = map[uint]string{
	1: "created_at",
	2: "rejected_by_seller_at",
	3: "accepted_by_seller_at",
	4: "rejected_by_admin_at",
	5: "accepted_by_admin_at",
	6: "cancelled_by_buyer_at",
	7: "completed_at",
}

type RefundRequestFilter uint

const (
	RefundRequestFilterAll                    RefundRequestFilter = 0
	RefundRequestFilterWaitingMerchantAproval RefundRequestFilter = 1
	RefundRequestFilterWaitingAdminAproval    RefundRequestFilter = 2
	RefundRequestFilterClosed                 RefundRequestFilter = 3
	RefundRequestFilterCanceled               RefundRequestFilter = 4
	RefundRequestFilterRejected               RefundRequestFilter = 5
	RefundRequestFilterRefunded               RefundRequestFilter = 6
	RefundRequestFilterWaitingBuyerAproval    RefundRequestFilter = 7
)

type RefundRequestListReqParamDTO struct {
	PaginationRequest
	Status uint `form:"status"`
}

type RefundRequestStatusDTO struct {
	CanceledByBuyerAt *time.Time `json:"canceled_by_buyer_at"`
	AcceptedByBuyerAt *time.Time `json:"accepted_by_buyer_at"`
	RejectedByBuyerAt *time.Time `json:"rejected_by_buyer_at"`

	AcceptedBySellerAt *time.Time `json:"accepted_by_seller_at"`
	RejectedBySellerAt *time.Time `json:"rejected_by_seller_at"`
	AcceptedByAdminAt  *time.Time `json:"accepted_by_admin_at"`
	RejectedByAdminAt  *time.Time `json:"rejected_by_admin_at"`

	ClosedAt *time.Time `json:"closed_at"`
}

type RefundRequestDTO struct {
	ID             uint   `json:"id"`
	TransactionId  uint   `json:"transaction_id"`
	InvoiceCode    string `json:"invoice_code"`
	Username       string `json:"username"`
	MerchantDomain string `json:"merchant_domain"`
	Reason         string `json:"reason"`
	ImageUrl       string `json:"image_url"`

	CreatedAt time.Time `json:"created_at"`

	RefundRequestStatusesDTO []RefundRequestStatusDTO `json:"refund_request_statuses"`
}

type RefundRequestListResDTO struct {
	PaginationResponse
	RefundRequests []RefundRequestDTO `json:"refund_requests"`
}

type RefundRequestFormReqDTO struct {
	InvoiceCode string                `form:"invoice_code" binding:"required"`
	Reason      string                `form:"reason" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
}

type RefundRequestFormResDTO struct {
	ID            uint      `json:"id"`
	TransactionId uint      `json:"transaction_id"`
	Reason        string    `json:"reason"`
	ImageUrl      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
}

type RefundRequestMsgFormReqDTO struct {
	Message string                `form:"message" binding:"required"`
	Image   *multipart.FileHeader `form:"image,omitempty"`
}

type RefundRequestMsgDetailsDTO struct {
	RefundId       uint       `json:"refund_id"`
	TransactionId  uint       `json:"transaction_id"`
	InvoiceCode    string     `json:"invoice_code"`
	BuyerUsername  string     `json:"buyer_username"`
	MerchantDomain string     `json:"merchant_domain"`
	Reason         string     `json:"reason"`
	ImageUrl       string     `json:"image_url"`
	ClosedAt       *time.Time `json:"closed_at"`
}

type RefundRequestMsgListResDTO struct {
	Details             RefundRequestMsgDetailsDTO `json:"details"`
	RefundRequestStatus []RefundRequestStatusDTO   `json:"refund_request_status"`
	Messages            []RefundRequestMsgResDTO   `json:"messages"`
}

type RefundRequestMsgRoleDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type RefundRequestMsgResDTO struct {
	ID         uint                    `json:"id"`
	Role       RefundRequestMsgRoleDTO `json:"role"`
	SenderName string                  `json:"sender_name"`
	Message    string                  `json:"message"`
	Image      *string                 `json:"image"`
	CreatedAt  time.Time               `json:"created_at"`
}
