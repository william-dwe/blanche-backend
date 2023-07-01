package repository

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type PaymentRecordRepository interface {
	CreateTx(tx *gorm.DB, input entity.PaymentRecord) (*entity.PaymentRecord, error)
	UpdateTx(tx *gorm.DB, input entity.PaymentRecord) (*entity.PaymentRecord, error)
	GetByPaymentId(paymentId string) (*entity.PaymentRecord, error)
	GetDetailByPaymentId(paymentId string) (*entity.PaymentRecord, error)

	GetWaitingForPayment(userId uint) ([]entity.PaymentRecord, error)
}

type PaymentRecordRepositoryConfig struct {
	DB *gorm.DB
}

type paymentRecordRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRecordRepository(c PaymentRecordRepositoryConfig) PaymentRecordRepository {
	return &paymentRecordRepositoryImpl{
		db: c.DB,
	}
}

func (r *paymentRecordRepositoryImpl) GetWaitingForPayment(userId uint) ([]entity.PaymentRecord, error) {
	var result []entity.PaymentRecord
	err := r.db.Select("distinct payment_records.*").
		Joins("left join transaction_payment_records on transaction_payment_records.payment_id = payment_records.payment_id").
		Joins("left join transactions on transaction_payment_records.transaction_id = transactions.id").
		Preload("TransactionPaymentRecord", func(db *gorm.DB) *gorm.DB {
			return db.Order("transaction_payment_records.created_at DESC")
		}).
		Preload("Transactions").
		Where("transactions.user_id = ?", userId).
		Where("paid_at IS NULL AND canceled_at IS NULL").
		Where(`(payment_records.payment_method_id = ? and payment_records.created_at >= ?)
		OR (payment_records.payment_method_id = ? and payment_records.created_at >= ?)`, dto.PAYMENT_METHOD_ID_WALLET, time.Now().Add(-24*time.Hour),
			dto.PAYMENT_METHOD_ID_SLP, time.Now().Add(-10*time.Minute)).
		Where("payment_url is not null and payment_url != ''").
		Order("payment_records.created_at DESC").
		Find(&result).Error

	if err != nil {
		return nil, domain.ErrGetPaymentRecord
	}

	return result, nil
}

func (r *paymentRecordRepositoryImpl) CreateTx(tx *gorm.DB, input entity.PaymentRecord) (*entity.PaymentRecord, error) {
	err := tx.Create(&input).Error
	if err != nil {
		return nil, domain.ErrCreatePaymentRecord
	}

	return &input, nil
}

func (r *paymentRecordRepositoryImpl) UpdateTx(tx *gorm.DB, input entity.PaymentRecord) (*entity.PaymentRecord, error) {
	err := tx.Save(&input).Error
	if err != nil {
		return nil, domain.ErrUpdatePaymentRecord
	}

	return &input, nil
}

func (r *paymentRecordRepositoryImpl) GetByPaymentId(paymentId string) (*entity.PaymentRecord, error) {
	var result entity.PaymentRecord
	err := r.db.
		Where("payment_id = ?", paymentId).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentIdNotFound
		}
		return nil, domain.ErrGetPaymentRecord
	}

	return &result, nil
}

func (r *paymentRecordRepositoryImpl) GetDetailByPaymentId(paymentId string) (*entity.PaymentRecord, error) {
	var result entity.PaymentRecord
	err := r.db.
		Preload("TransactionPaymentRecord", func(db *gorm.DB) *gorm.DB {
			return db.Order("transaction_payment_records.created_at DESC")
		}).
		Preload("Transactions").
		Preload("Transactions.Merchant").
		Where("paid_at IS NULL AND canceled_at IS NULL").
		Where("payment_id = ?", paymentId).
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentIdNotFound
		}
		return nil, domain.ErrGetPaymentRecord
	}

	return &result, nil
}
