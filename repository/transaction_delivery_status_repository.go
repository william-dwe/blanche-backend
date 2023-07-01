package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type TransactionDeliveryStatusRepository interface {
	GetTransactionDeliveryStatusByTransactionID(transactionID uint) (*entity.TransactionDeliveryStatus, error)
	DeleteTransactionDeliveryStatusTx(tx *gorm.DB, trxIds []uint) error

	UpdateTransactionDeliveryStatus(transactionDeliveryStatus entity.TransactionDeliveryStatus) (*entity.TransactionDeliveryStatus, int64, error)
}

type TransactionDeliveryStatusRepositoryConfig struct {
	DB *gorm.DB
}

type transactionDeliveryStatusRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionDeliveryStatusRepository(c TransactionDeliveryStatusRepositoryConfig) TransactionDeliveryStatusRepository {
	return &transactionDeliveryStatusRepositoryImpl{
		db: c.DB,
	}
}

func (r *transactionDeliveryStatusRepositoryImpl) GetTransactionDeliveryStatusByTransactionID(transactionID uint) (*entity.TransactionDeliveryStatus, error) {
	var transactionDeliveryStatus entity.TransactionDeliveryStatus
	err := r.db.Where("transaction_id = ?", transactionID).First(&transactionDeliveryStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransactionDeliveryStatusNotFound
		}

		return nil, domain.ErrGetTransactionDeliveryStatus
	}

	return &transactionDeliveryStatus, nil
}

func (r *transactionDeliveryStatusRepositoryImpl) DeleteTransactionDeliveryStatusTx(tx *gorm.DB, trxIds []uint) error {
	err := tx.Where("transaction_id IN ?", trxIds).Delete(&entity.TransactionDeliveryStatus{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionDeliveryStatusRepositoryImpl) UpdateTransactionDeliveryStatus(transactionDeliveryStatus entity.TransactionDeliveryStatus) (*entity.TransactionDeliveryStatus, int64, error) {
	res := r.db.Where("id = ?", transactionDeliveryStatus.ID).Updates(&transactionDeliveryStatus)
	if res.Error != nil {
		return nil, 0, domain.ErrUpdateTransactionDeliveryStatus
	}

	return &transactionDeliveryStatus, res.RowsAffected, nil
}
