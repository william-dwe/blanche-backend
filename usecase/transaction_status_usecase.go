package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"github.com/rs/zerolog/log"
)

type TransactionStatusUsecase interface {
	GetTransactionStatusByTransactionID(transactionID uint) (*dto.TransactionStatusResDTO, error)

	CronUpdateTransactionWaitingStatusToCanceled()
	CronUpdateTransactionProcessedStatusToCanceled()
	CronUpdateTransactionDeliveredStatusToCompleted()
}

type transactionStatusUsecaseImpl struct {
	transactionStatusRepo repository.TransactionStatusRepository
	cron                  *cronjob.CronJob
}

type TransactionStatusUsecaseConfig struct {
	TransactionStatusRepo repository.TransactionStatusRepository
	Cron                  *cronjob.CronJob
}

func NewTransactionStatusUsecase(c TransactionStatusUsecaseConfig) TransactionStatusUsecase {
	transactionStatusUsecaseImpl := &transactionStatusUsecaseImpl{
		transactionStatusRepo: c.TransactionStatusRepo,
		cron:                  c.Cron,
	}

	c.Cron.AddJob("* * * * *", transactionStatusUsecaseImpl.CronUpdateTransactionWaitingStatusToCanceled)
	c.Cron.AddJob("* * * * *", transactionStatusUsecaseImpl.CronUpdateTransactionProcessedStatusToCanceled)
	c.Cron.AddJob("* * * * *", transactionStatusUsecaseImpl.CronUpdateTransactionDeliveredStatusToCompleted)

	return transactionStatusUsecaseImpl
}

func (u *transactionStatusUsecaseImpl) GetTransactionStatusByTransactionID(transactionID uint) (*dto.TransactionStatusResDTO, error) {
	transactionStatus, err := u.transactionStatusRepo.GetTransactionStatusByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}

	return &dto.TransactionStatusResDTO{
		OnWaitedAt:        transactionStatus.OnWaitedAt,
		OnProcessedAt:     transactionStatus.OnProcessedAt,
		OnDeliveredAt:     transactionStatus.OnDeliveredAt,
		OnCompletedAt:     transactionStatus.OnCompletedAt,
		OnCanceledAt:      transactionStatus.OnCanceledAt,
		OnRefundedAt:      transactionStatus.OnRefundedAt,
		OnRequestRefundAt: transactionStatus.OnRequestRefundAt,
		CancellationNotes: transactionStatus.CancellationNotes,
	}, nil
}

func (u *transactionStatusUsecaseImpl) CronUpdateTransactionWaitingStatusToCanceled() {
	cronConfig := config.Config.CronConfig
	concurrentSize := cronConfig.TrxQueueSizeWaitingToCanceled / cronConfig.TrxBatchSizeWaitingToCanceled
	for i := 0; i < concurrentSize; i++ {
		go func(batchNum int) {
			err := u.transactionStatusRepo.CronUpdateTransactionWaitingStatusToCanceled(batchNum)
			if err != nil {
				log.Error().Msgf("CronUpdateTransactionWaitingStatusToCanceled Error: %v", err)
			}
		}(i)
	}
}

func (u *transactionStatusUsecaseImpl) CronUpdateTransactionProcessedStatusToCanceled() {
	cronConfig := config.Config.CronConfig
	concurrentSize := cronConfig.TrxQueueSizeWaitingToCanceled / cronConfig.TrxBatchSizeWaitingToCanceled
	for i := 0; i < concurrentSize; i++ {
		go func(batchNum int) {
			err := u.transactionStatusRepo.CronUpdateTransactionProcessedStatusToCanceled(batchNum)
			if err != nil {
				log.Error().Msgf("CronUpdateTransactionProcessedStatusToCanceled Error: %v", err)
			}
		}(i)
	}
}

func (u *transactionStatusUsecaseImpl) CronUpdateTransactionDeliveredStatusToCompleted() {
	cronConfig := config.Config.CronConfig
	concurrentSize := cronConfig.TrxQueueSizeWaitingToCanceled / cronConfig.TrxBatchSizeWaitingToCanceled
	for i := 0; i < concurrentSize; i++ {
		go func(batchNum int) {
			err := u.transactionStatusRepo.CronUpdateTransactionDeliveredStatusToCompleted(batchNum)
			if err != nil {
				log.Error().Msgf("CronUpdateTransactionDeliveredStatusToCompleted Error: %v", err)
			}
		}(i)
	}
}
