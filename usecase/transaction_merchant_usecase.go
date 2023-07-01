package usecase

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
)

func (u *transactionUsecaseImpl) UpdateMerchantTransactionStatus(username string, req dto.UpdateMerchantTransactionStatusReqDTO) (*dto.UpdateMerchantTransactionStatusResDTO, error) {
	validatedTransaction, err := u.validateUpdateMerchantTransactionStatus(username, req)
	if err != nil {
		return nil, err
	}

	if req.Status == dto.TransactionStatusCompleted ||
		req.Status == dto.TransactionStatusRequestRefund ||
		req.Status == dto.TransactionStatusRefunded {
		return nil, domain.ErrInvalidTransactionStatusForbidden
	}

	//update transaction delivery status
	updatedTrxDeliveryStatus, err := u.updateTransactionDeliveryStatus(*validatedTransaction.TransactionDeliveryStatus, req)
	if err != nil {
		return nil, err
	}

	//update transaction status
	updatedTrxStatus, err := u.updateMerchantTransactionStatus(*validatedTransaction.TransactionStatus, *validatedTransaction, req)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateMerchantTransactionStatusResDTO{
		InvoiceCode: req.InvoiceCode,
		TransactionStatus: dto.TransactionStatusResDTO{
			OnWaitedAt:        updatedTrxStatus.OnWaitedAt,
			OnProcessedAt:     updatedTrxStatus.OnProcessedAt,
			OnDeliveredAt:     updatedTrxStatus.OnDeliveredAt,
			OnCompletedAt:     updatedTrxStatus.OnCompletedAt,
			OnCanceledAt:      updatedTrxStatus.OnCanceledAt,
			OnRefundedAt:      updatedTrxStatus.OnRefundedAt,
			OnRequestRefundAt: updatedTrxStatus.OnRequestRefundAt,

			CancellationNotes: updatedTrxStatus.CancellationNotes,
		},
		TransactionDeliveryStatus: dto.TransactionDeliveryStatusResDTO{
			OnDeliveryAt:  updatedTrxDeliveryStatus.OnDeliveryAt,
			OnDeliveredAt: updatedTrxDeliveryStatus.OnDeliveredAt,
		},
		UpdatedAt: updatedTrxStatus.UpdatedAt,
	}, nil
}

func (u *transactionUsecaseImpl) updateMerchantTransactionStatus(trxStatus entity.TransactionStatus, transaction entity.Transaction, req dto.UpdateMerchantTransactionStatusReqDTO) (*entity.TransactionStatus, error) {
	if req.Status == dto.TransactionStatusCompleted {
		return nil, domain.ErrInvalidTransactionStatusForbidden
	}

	updatedStatus, err := u.changeMerchantTransactionStatus(trxStatus, req)
	if err != nil {
		return nil, err
	}

	transaction.TransactionStatus = updatedStatus

	if req.Status == dto.TransactionStatusCanceled {
		updatedTrxStatus, err := u.updateTransactionCanceledProcess(transaction)
		if err != nil {
			return nil, err
		}

		return updatedTrxStatus, nil
	}

	updatedTrxStatus, updatedRow, err := u.transactionStatusRepository.UpdateTransactionStatus(*updatedStatus)
	if err != nil {
		return nil, err
	}
	if updatedRow <= 0 {
		return nil, domain.ErrUpdateTransactionStatus
	}

	return updatedTrxStatus, nil
}

func (u *transactionUsecaseImpl) changeMerchantTransactionStatus(transactionStatus entity.TransactionStatus, req dto.UpdateMerchantTransactionStatusReqDTO) (*entity.TransactionStatus, error) {
	timeNow := time.Now()
	switch req.Status {
	case dto.TransactionStatusCanceled:
		transactionStatus.OnCanceledAt = &timeNow
		transactionStatus.CancellationNotes = req.CancellationNotes
	case dto.TransactionStatusDelivered:
		transactionStatus.OnDeliveredAt = &timeNow
	case dto.TransactionStatusProcessed:
		transactionStatus.OnProcessedAt = &timeNow
	}

	return &transactionStatus, nil
}

func (u *transactionUsecaseImpl) changeMerchantTransactionDeliveryStatus(transactionDeliveryStatus entity.TransactionDeliveryStatus, req dto.UpdateMerchantTransactionStatusReqDTO) (*entity.TransactionDeliveryStatus, error) {
	timeNow := time.Now()
	switch req.Status {
	case dto.TransactionStatusOnDelivery:
		transactionDeliveryStatus.OnDeliveryAt = &timeNow
		transactionDeliveryStatus.ReceiptNumber = &req.ReceiptNumber
	case dto.TransactionStatusDelivered:
		transactionDeliveryStatus.OnDeliveredAt = &timeNow
	}

	return &transactionDeliveryStatus, nil
}

func (u *transactionUsecaseImpl) validateUpdateMerchantTransactionStatus(username string, req dto.UpdateMerchantTransactionStatusReqDTO) (*entity.Transaction, error) {
	//if ondelivery, check if receipt number is exist
	if req.Status == dto.TransactionStatusOnDelivery {
		if req.ReceiptNumber == "" {
			return nil, domain.ErrUpdateTransactionStatusReceiptNumberEmpty
		}
	}

	//get userId from username
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	//get merchantDomain from username
	merchant, err := u.merchantRepository.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	//get transaction by invoiceCode
	transaction, err := u.transactionRepository.GetMerchantTransactionDetailByInvoiceCode(merchant.Domain, req.InvoiceCode)
	if err != nil {
		return nil, err
	}

	//check if status is valid
	err = u.checkUpdateStatusIsValid(*transaction, req)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *transactionUsecaseImpl) updateTransactionDeliveryStatus(trxDeliveryStatus entity.TransactionDeliveryStatus, req dto.UpdateMerchantTransactionStatusReqDTO) (*entity.TransactionDeliveryStatus, error) {
	updatedDeliveryStatus, err := u.changeMerchantTransactionDeliveryStatus(trxDeliveryStatus, req)
	if err != nil {
		return nil, err
	}

	if *updatedDeliveryStatus == trxDeliveryStatus {
		return updatedDeliveryStatus, nil
	}

	updatedTrxDeliveryStatus, updatedRow, err := u.transactionDeliveryStatusRepository.UpdateTransactionDeliveryStatus(*updatedDeliveryStatus)
	if err != nil {
		return nil, err
	}
	if updatedRow <= 0 {
		return nil, domain.ErrUpdateTransactionDeliveryStatus
	}

	return updatedTrxDeliveryStatus, nil
}
