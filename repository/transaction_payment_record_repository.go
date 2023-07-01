package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type TransactionPaymentRecordRepository interface {
	CreateTx(tx *gorm.DB, input entity.TransactionPaymentRecord) (*entity.TransactionPaymentRecord, error)
	GetByPaymentId(paymentId string) ([]entity.TransactionPaymentRecord, error)
	UpdateRecordOrderCodebyPaymentIdTx(tx *gorm.DB, orderCode string, paymentId string) error
}

type TransactionPaymentRecordRepositoryConfig struct {
	DB *gorm.DB
}

type transactionPaymentRecordRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionPaymentRecordRepository(c TransactionPaymentRecordRepositoryConfig) TransactionPaymentRecordRepository {
	return &transactionPaymentRecordRepositoryImpl{
		db: c.DB,
	}
}

func (r *transactionPaymentRecordRepositoryImpl) CreateTx(tx *gorm.DB, input entity.TransactionPaymentRecord) (*entity.TransactionPaymentRecord, error) {
	err := tx.Create(&input).Error
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func (r *transactionPaymentRecordRepositoryImpl) GetByPaymentId(paymentId string) ([]entity.TransactionPaymentRecord, error) {
	var result []entity.TransactionPaymentRecord
	err := r.db.
		Preload("Transaction").
		Where("payment_id = ?", paymentId).
		Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *transactionPaymentRecordRepositoryImpl) UpdateRecordOrderCodebyPaymentIdTx(tx *gorm.DB, orderCode string, paymentId string) error {
	err := tx.Model(&entity.TransactionPaymentRecord{}).
		Where("payment_id = ?", paymentId).
		Update("order_code", orderCode).Error

	if err != nil {
		return err
	}

	return nil
}
