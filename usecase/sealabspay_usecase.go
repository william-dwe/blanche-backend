package usecase

import (
	"fmt"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type SealabspayUsecase interface {
	HandlePaymentResponse(input dto.SealabspayReqDTO) (*dto.SealabspayReqDTO, error)
}

type SealabspayUsecaseConfig struct {
	PaymentRecordUsecase PaymentRecordUsecase
}

type sealabspayUsecaseImpl struct {
	paymentRecordUsecase PaymentRecordUsecase
}

func NewSealabspayUsecase(c SealabspayUsecaseConfig) SealabspayUsecase {
	return &sealabspayUsecaseImpl{
		paymentRecordUsecase: c.PaymentRecordUsecase,
	}
}

func (u *sealabspayUsecaseImpl) HandlePaymentResponse(input dto.SealabspayReqDTO) (*dto.SealabspayReqDTO, error) {
	err := util.ValidateResSlpSignature(input)
	if err != nil {
		return nil, err
	}

	paymentId, err := strconv.ParseUint(input.TxnId, 10, 64)
	if err != nil {
		return nil, domain.ErrSealabspayTxnIdNotValid
	}

	paymentAmount, err := strconv.ParseUint(input.Amount, 10, 64)
	if err != nil {
		return nil, domain.ErrSealabspayAmountNotValid
	}

	if input.Status == dto.SLP_FAILED_CODE {
		err = u.paymentRecordUsecase.UpdatePaymentRecordStatus(fmt.Sprintf("SLP%d", paymentId), uint(paymentAmount), false)
		if err != nil {
			return nil, err
		}
	}

	if input.Status == dto.SLP_SUCCESS_CODE {
		err = u.paymentRecordUsecase.UpdatePaymentRecordStatus(fmt.Sprintf("SLP%d", paymentId), uint(paymentAmount), true)
		if err != nil {
			return nil, err
		}
	}

	return &input, nil
}
