package repository

import (
	"encoding/json"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type TransactionStatusRepository interface {
	GetTransactionStatusByTransactionID(transactionID uint) (*entity.TransactionStatus, error)
	UpdateTransactionStatusWaitedAtTx(tx *gorm.DB, trxIds []uint, timeNow time.Time) error
	DeleteTransactionStatusTx(tx *gorm.DB, trxIds []uint) error

	UpdateTransactionStatus(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, int64, error)
	UpdateTransactionStatusTx(tx *gorm.DB, transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, int64, error)

	CronUpdateTransactionWaitingStatusToCanceled(batchNumber int) error
	CronUpdateTransactionProcessedStatusToCanceled(batchNumber int) error
	CronUpdateTransactionDeliveredStatusToCompleted(batchNumber int) error
}

type transactionStatusRepositoryImpl struct {
	db                       *gorm.DB
	transactionRepositoryPtr TransactionRepository
}

type TransactionStatusRepositoryConfig struct {
	DB                       *gorm.DB
	TransactionRepositoryPtr TransactionRepository
}

func NewTransactionStatusRepository(c TransactionStatusRepositoryConfig) TransactionStatusRepository {
	return &transactionStatusRepositoryImpl{
		db:                       c.DB,
		transactionRepositoryPtr: c.TransactionRepositoryPtr,
	}
}

func (r *transactionStatusRepositoryImpl) GetTransactionStatusByTransactionID(transactionID uint) (*entity.TransactionStatus, error) {
	var transactionStatus entity.TransactionStatus
	err := r.db.Model(&transactionStatus).Where("transaction_id = ?", transactionID).First(&transactionStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransactionStatusNotFound
		}

		return nil, domain.ErrGetTransactionStatus
	}

	return &transactionStatus, nil
}

func (r *transactionStatusRepositoryImpl) UpdateTransactionStatusWaitedAtTx(tx *gorm.DB, trxIds []uint, timeNow time.Time) error {
	err := tx.Model(&entity.TransactionStatus{}).Where("transaction_id IN ?", trxIds).Update("on_waited_at", timeNow).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionStatusRepositoryImpl) DeleteTransactionStatusTx(tx *gorm.DB, trxIds []uint) error {
	err := tx.Where("transaction_id IN ?", trxIds).Delete(&entity.TransactionStatus{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionStatusRepositoryImpl) UpdateTransactionStatus(transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, int64, error) {
	res := r.db.Where("id = ?", transactionStatus.ID).Updates(&transactionStatus)
	if res.Error != nil {
		return nil, 0, domain.ErrUpdateTransactionStatus
	}

	return &transactionStatus, res.RowsAffected, nil
}

func (r *transactionStatusRepositoryImpl) UpdateTransactionStatusTx(tx *gorm.DB, transactionStatus entity.TransactionStatus) (*entity.TransactionStatus, int64, error) {
	transactionStatus.Transaction = entity.Transaction{}
	res := tx.Where("id = ?", transactionStatus.ID).Updates(&transactionStatus)
	if res.Error != nil {
		return nil, 0, domain.ErrUpdateTransactionStatus
	}

	return &transactionStatus, res.RowsAffected, nil
}

func (r *transactionStatusRepositoryImpl) updateTransactionToCancelledTx(tx *gorm.DB, transactionStatuses []entity.TransactionStatus) {
	for _, transactionStatus := range transactionStatuses {
		transactionTmp := transactionStatus.Transaction
		transactionTmp.TransactionStatus = &transactionStatus

		nowTime := time.Now()
		transactionTmp.TransactionStatus.OnCanceledAt = &nowTime
		transactionTmp.TransactionStatus.CancellationNotes = "Transaction is cancelled due to merchant is not responding"

		amount, _, err := r.countAmountAndPromotionTrx(transactionStatus.Transaction)
		if err != nil {
			log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Error: %v", err)
			continue
		}
		var trxCartItems []entity.TransactionCartItem
		err = json.Unmarshal([]byte(transactionStatus.Transaction.CartItems.Bytes), &trxCartItems)
		if err != nil {
			log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Error: %v", err)
			continue
		}

		_, err = r.transactionRepositoryPtr.UpdateTransactionStatusCanceledTx(tx, transactionTmp, amount, trxCartItems)
		if err != nil {
			log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Update Status Trx: %v", err)
			continue
		}
	}
}

func (r *transactionStatusRepositoryImpl) CronUpdateTransactionWaitingStatusToCanceled(batchNumber int) (errCron error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in CronUpdateTransactionWaitingStatusToCanceled repo: %v", r)
			errCron = domain.ErrCronUpdateTransactionWaitingStatusToCanceled
		}
	}()

	var transactionStatuses []entity.TransactionStatus
	err := tx.Model(&transactionStatuses).
		Preload("Transaction").
		Preload("Transaction.Merchant").
		Where("on_waited_at IS NOT NULL").
		Where("on_waited_at <= ?", time.Now().Add(-24*time.Hour)).
		Where("on_processed_at IS NULL").
		Where("on_delivered_at IS NULL").
		Where("on_completed_at IS NULL").
		Where("on_canceled_at IS NULL").
		Where("on_refunded_at IS NULL").
		Where("on_request_refund_at IS NULL").
		Where("deleted_at IS NULL").
		Limit(config.Config.CronConfig.TrxBatchSizeWaitingToCanceled).
		Offset(batchNumber * config.Config.CronConfig.TrxBatchSizeWaitingToCanceled).
		Find(&transactionStatuses).
		Error

	if err != nil {
		log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Error: %v", err)
		return domain.ErrCronUpdateTransactionWaitingStatusToCanceled
	}

	r.updateTransactionToCancelledTx(tx, transactionStatuses)

	err = tx.Commit().Error
	if err != nil {
		log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Commit: %v", err)
		return domain.ErrCronUpdateTransactionWaitingStatusToCanceled
	}

	return nil
}

func (r *transactionStatusRepositoryImpl) CronUpdateTransactionProcessedStatusToCanceled(batchNumber int) (errCron error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in CronUpdateTransactionProcessedStatusToCanceled repo: %v", r)
			errCron = domain.ErrCronUpdateTransactionStatusToCanceled
		}
	}()

	var transactionStatuses []entity.TransactionStatus
	err := tx.Model(&transactionStatuses).
		Preload("Transaction").
		Preload("Transaction.Merchant").
		Where("on_processed_at IS NOT NULL").
		Where("on_processed_at <= ?", time.Now().Add(-48*time.Hour)).
		Where("on_delivered_at IS NULL").
		Where("on_completed_at IS NULL").
		Where("on_canceled_at IS NULL").
		Where("on_refunded_at IS NULL").
		Where("on_request_refund_at IS NULL").
		Limit(config.Config.CronConfig.TrxBatchSizeProcessedToCanceled).
		Offset(batchNumber * config.Config.CronConfig.TrxBatchSizeProcessedToCanceled).
		Find(&transactionStatuses).
		Error

	if err != nil {
		log.Error().Msgf("CronUpdateTransactionProcessedStatusToCanceled Error: %v", err)
		return domain.ErrCronUpdateTransactionStatusToCanceled
	}

	r.updateTransactionToCancelledTx(tx, transactionStatuses)

	err = tx.Commit().Error
	if err != nil {
		log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Commit: %v", err)
		return domain.ErrCronUpdateTransactionStatusToCanceled
	}

	return nil
}

func (r *transactionStatusRepositoryImpl) CronUpdateTransactionDeliveredStatusToCompleted(batchNumber int) (errCron error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in CronUpdateTransactionDeliveredStatusToCompleted repo: %v", r)
			errCron = domain.ErrCronUpdateTransactionStatusToCompleted
		}
	}()

	var transactionStatuses []entity.TransactionStatus
	err := tx.Model(&transactionStatuses).
		Preload("Transaction").
		Preload("Transaction.Merchant").
		Where("on_delivered_at IS NOT NULL").
		Where("on_delivered_at <= ?", time.Now().Add(-72*time.Hour)).
		Where("on_completed_at IS NULL").
		Where("on_canceled_at IS NULL").
		Where("on_refunded_at IS NULL").
		Where("on_request_refund_at IS NULL").
		Limit(config.Config.CronConfig.TrxBatchSizeDeliveredToCompleted).
		Offset(batchNumber * config.Config.CronConfig.TrxBatchSizeDeliveredToCompleted).
		Find(&transactionStatuses).
		Error

	if err != nil {
		log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Error: %v", err)
		return domain.ErrCronUpdateTransactionStatusToCompleted
	}

	for _, transactionStatus := range transactionStatuses {
		transactionTmp := transactionStatus.Transaction
		transactionTmp.TransactionStatus = &transactionStatus
		nowTime := time.Now()
		transactionTmp.TransactionStatus.OnCompletedAt = &nowTime

		amount, mpAmount, err := r.countAmountAndPromotionTrx(transactionStatus.Transaction)
		if err != nil {
			log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Error: %v", err)
			continue
		}
		_, err = r.transactionRepositoryPtr.UpdateTransactionStatusCompletedTx(tx, transactionTmp, amount, mpAmount)
		if err != nil {
			log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Update Status Trx: %v", err)
			continue
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Commit: %v", err)
		return domain.ErrCronUpdateTransactionStatusToCompleted
	}

	return nil
}

func (r *transactionStatusRepositoryImpl) countAmountAndPromotionTrx(transaction entity.Transaction) (float64, float64, error) {
	var amount float64
	var promotion float64

	var trxPaymentDetails entity.TransactionPaymentDetails
	err := json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &trxPaymentDetails)
	if err != nil {
		return 0, 0, domain.ErrUnmarshalJSONPaymentDetails
	}

	amount = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MerchantVoucherNominal
	promotion = trxPaymentDetails.MarketplaceVoucherNominal

	return amount, promotion, nil
}
