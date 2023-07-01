package dto

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
)

const (
	TransactionStatusWaited        int = 1
	TransactionStatusProcessed     int = 2
	TransactionStatusCanceled      int = 3
	TransactionStatusOnDelivery    int = 4
	TransactionStatusDelivered     int = 5
	TransactionStatusRequestRefund int = 6
	TransactionStatusCompleted     int = 7
	TransactionStatusRefunded      int = 8
)

type TransactionDetailResDTO struct {
	InvoiceCode       string                          `json:"invoice_code"`
	TransactionStatus TransactionStatusResDTO         `json:"transaction_status"`
	ProductDetails    TransactionDetailProductResDTO  `json:"product_details"`
	ShippingDetails   TransactionDetailShippingResDTO `json:"shipping_details"`
	PaymentDetails    TransactionDetailPaymentResDTO  `json:"payment_details"`
}

type TransactionSellerDetailResDTO struct {
	InvoiceCode       string                               `json:"invoice_code"`
	TransactionStatus TransactionStatusResDTO              `json:"transaction_status"`
	ProductDetails    TransactionSellerDetailProductResDTO `json:"product_details"`
	ShippingDetails   TransactionDetailShippingResDTO      `json:"shipping_details"`
	PaymentDetails    TransactionDetailPaymentResDTO       `json:"payment_details"`
}

type TransactionSellerDetailProductResDTO struct {
	User      TransactionUserResDTO        `json:"user"`
	CartItems []entity.TransactionCartItem `json:"products"`
}

type TransactionDetailProductResDTO struct {
	Merchant  TransactionMerchantResDTO    `json:"merchant"`
	CartItems []entity.TransactionCartItem `json:"products"`
}

type TransactionDeliveryOptionResDTO struct {
	CourierName   string `json:"courier_name"`
	ReceiptNumber string `json:"receipt_number"`
}

type TransactionDetailShippingResDTO struct {
	Address                   entity.TransactionAddress       `json:"address"`
	DeliveryOption            TransactionDeliveryOptionResDTO `json:"delivery_option"`
	TransactionDeliveryStatus TransactionDeliveryStatusResDTO `json:"transaction_delivery_status"`
}

type TransactionDetailPaymentResDTO struct {
	PaymentMethod  entity.TransactionPaymentMethod  `json:"method"`
	PaymentDetails entity.TransactionPaymentDetails `json:"summary"`
}

type TransactionResDTO struct {
	InvoiceCode       string                          `json:"invoice_code"`
	Total             float64                         `json:"total_price"`
	Merchant          TransactionMerchantResDTO       `json:"merchant"`
	ProductOverview   TransactionListCartResDTO       `json:"product_overview"`
	TransactionStatus TransactionStatusResDTO         `json:"transaction_status"`
	ShippingDetails   TransactionDetailShippingResDTO `json:"shipping_details"`
}

type TransactionSellerResDTO struct {
	InvoiceCode       string                          `json:"invoice_code"`
	Total             float64                         `json:"total_price"`
	Username          string                          `json:"buyer_username"`
	TransactionDate   time.Time                       `json:"transaction_date"`
	ProductOverview   TransactionListCartResDTO       `json:"product_overview"`
	TransactionStatus TransactionStatusResDTO         `json:"transaction_status"`
	ShippingDetails   TransactionDetailShippingResDTO `json:"shipping_details"`
}

type TransactionListCartResDTO struct {
	Product      entity.TransactionCartItem `json:"product"`
	TotalProduct int                        `json:"total_product"`
}

type TransactionReqParamDTO struct {
	StartDate  time.Time `form:"start_date,default=2023-01-01T00:00:00Z"`
	EndDate    time.Time `form:"end_date,default=9999-01-01T00:00:00Z"`
	Search     string    `form:"q,default="`
	Sort       uint      `form:"sort,default=1"`
	Status     uint      `form:"status,default=0"`
	Pagination PaginationRequest
}

type TransactionListResDTO struct {
	PaginationResponse
	Transactions []TransactionResDTO `json:"transactions"`
}

type TransactionSellerListResDTO struct {
	PaginationResponse
	Transactions []TransactionSellerResDTO `json:"transactions"`
}

type TransactionMerchantResDTO struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type TransactionUserResDTO struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type MakeTransactionReqDTO struct {
	OrderCode string `json:"order_code" binding:"required"`

	AddressId          int                            `json:"address_id" binding:"required"`
	Merchants          []PostOrderSummaryMerchantsDTO `json:"merchants"`
	VoucherMarketplace string                         `json:"voucher_marketplace"`

	PaymentTotal         float64 `json:"payment_total" binding:"required"`
	PaymentMethodCode    string  `json:"payment_method_code" binding:"required"`
	PaymentAccountNumber string  `json:"payment_account_number" binding:"required"`
}

type MakeTransactionResDTO struct {
	PaymentId          string  `json:"payment_id"`
	OrderCode          string  `json:"order_code"`
	Amount             float64 `json:"amount"`
	PaymentRedirectUrl string  `json:"payment_redirect_url"`
}

type UpdateMerchantTransactionStatusReqDTO struct {
	InvoiceCode       string `json:"-"`
	Status            int    `json:"status" binding:"required"`
	ReceiptNumber     string `json:"receipt_number"`
	CancellationNotes string `json:"cancellation_notes"`
}

type UpdateMerchantTransactionStatusResDTO struct {
	InvoiceCode               string                          `json:"invoice_code"`
	TransactionStatus         TransactionStatusResDTO         `json:"transaction_status"`
	TransactionDeliveryStatus TransactionDeliveryStatusResDTO `json:"transaction_delivery_status"`
	UpdatedAt                 time.Time                       `json:"updated_at"`
}

type UpdateUserTransactionStatusReqDTO struct {
	InvoiceCode string `json:"-"`
	Status      int    `json:"status" binding:"required"`
}

type UpdateUserTransactionStatusResDTO struct {
	InvoiceCode       string                  `json:"invoice_code"`
	TransactionStatus TransactionStatusResDTO `json:"transaction_status"`
	UpdatedAt         time.Time               `json:"updated_at"`
}
