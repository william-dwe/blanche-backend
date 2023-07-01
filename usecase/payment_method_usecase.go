package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type PaymentMethodUsecase interface {
	GetAll() ([]dto.PaymentMehodDTO, error)
}

type PaymentMethodUsecaseConfig struct {
	PaymentMethodRepository repository.PaymentMethodRepository
}

type paymentMethodUsecaseImpl struct {
	paymentMethodRepository repository.PaymentMethodRepository
}

func NewPaymentMethodUsecase(c PaymentMethodUsecaseConfig) PaymentMethodUsecase {
	return &paymentMethodUsecaseImpl{
		paymentMethodRepository: c.PaymentMethodRepository,
	}
}

func (u *paymentMethodUsecaseImpl) GetAll() ([]dto.PaymentMehodDTO, error) {
	paymentMethods, err := u.paymentMethodRepository.GetAll()
	if err != nil {
		return nil, err
	}

	paymentMethodsDTO := make([]dto.PaymentMehodDTO, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		paymentMethodsDTO[i] = dto.PaymentMehodDTO{
			ID:   int(paymentMethod.ID),
			Name: paymentMethod.Name,
			Code: paymentMethod.Code,
		}
	}

	return paymentMethodsDTO, nil
}
