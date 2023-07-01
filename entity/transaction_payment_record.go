package entity

type TransactionPaymentRecord struct {
	ID            uint `gorm:"primary_key"`
	TransactionId uint `gorm:"not null"`
	Transaction   Transaction
	PaymentId     string `gorm:"not null"`
	OrderCode     string
}
