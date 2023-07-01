package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type TransactionDeliveryStatusUsecase interface {
	GetTransactionDeliveryStatusByTransactionID(transactionID uint) (*dto.TransactionDeliveryStatusResDTO, *string, error)
}

type transactionDeliveryStatusUsecaseImpl struct {
	transactionDeliveryStatusRepo repository.TransactionDeliveryStatusRepository
}

type TransactionDeliveryStatusUsecaseConfig struct {
	TransactionDeliveryStatusRepo repository.TransactionDeliveryStatusRepository
}

func NewTransactionDeliveryStatusUsecase(c TransactionDeliveryStatusUsecaseConfig) TransactionDeliveryStatusUsecase {
	return &transactionDeliveryStatusUsecaseImpl{
		transactionDeliveryStatusRepo: c.TransactionDeliveryStatusRepo,
	}
}

func (u *transactionDeliveryStatusUsecaseImpl) GetTransactionDeliveryStatusByTransactionID(transactionID uint) (*dto.TransactionDeliveryStatusResDTO, *string, error) {
	transactionDeliveryStatus, err := u.transactionDeliveryStatusRepo.GetTransactionDeliveryStatusByTransactionID(transactionID)
	if err != nil {
		return nil, nil, err
	}

	return &dto.TransactionDeliveryStatusResDTO{
		OnDeliveryAt:  transactionDeliveryStatus.OnDeliveryAt,
		OnDeliveredAt: transactionDeliveryStatus.OnDeliveredAt,
	}, transactionDeliveryStatus.ReceiptNumber, nil
}
