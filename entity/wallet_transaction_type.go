package entity

type WalletTransactionType struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Code string
}
