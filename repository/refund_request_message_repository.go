package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type RefundRequestMessageRepository interface {
	AddMessage(message entity.RefundReqMessage) (*entity.RefundReqMessage, error)
	GetAllMessageByRefundRequestId(refundRequestId uint) ([]entity.RefundReqMessage, error)
}

type RefundRequestMessageRepositoryConfig struct {
	DB *gorm.DB
}

type refundRequestMessageRepositoryImpl struct {
	db *gorm.DB
}

func NewRefundRequestMessageRepository(c RefundRequestMessageRepositoryConfig) RefundRequestMessageRepository {
	return &refundRequestMessageRepositoryImpl{
		db: c.DB,
	}
}

func (r *refundRequestMessageRepositoryImpl) AddMessage(message entity.RefundReqMessage) (*entity.RefundReqMessage, error) {
	err := r.db.Create(&message).Error
	if err != nil {
		return nil, domain.ErrAddMessageRefundRequest
	}

	return &message, nil
}

func (r *refundRequestMessageRepositoryImpl) GetAllMessageByRefundRequestId(refundRequestId uint) ([]entity.RefundReqMessage, error) {
	var messages []entity.RefundReqMessage

	err := r.db.
		Where("refund_request_id = ?", refundRequestId).
		Preload("RefundReqMessageRole").
		Order("created_at DESC").
		Find(&messages).Error
	if err != nil {
		return nil, domain.ErrGetMessageRefundRequest
	}

	return messages, nil
}
