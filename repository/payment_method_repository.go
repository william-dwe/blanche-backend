package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	GetAll() ([]entity.PaymentMethod, error)
	GetByCode(paymentCode string) (*entity.PaymentMethod, error)
}

type PaymentMethodRepositoryConfig struct {
	DB *gorm.DB
}

type paymentMethodRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(c PaymentMethodRepositoryConfig) PaymentMethodRepository {
	return &paymentMethodRepositoryImpl{
		db: c.DB,
	}
}

func (r *paymentMethodRepositoryImpl) GetAll() ([]entity.PaymentMethod, error) {
	var paymentMethods []entity.PaymentMethod
	err := r.db.Find(&paymentMethods).Error
	if err != nil {
		return nil, domain.ErrGetAllPaymentMethod
	}
	return paymentMethods, nil
}

func (r *paymentMethodRepositoryImpl) GetByCode(paymentCode string) (*entity.PaymentMethod, error) {
	var paymentMethod entity.PaymentMethod
	err := r.db.Where("code = ?", paymentCode).First(&paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPaymentMethodNotFound
		}
		return nil, domain.ErrGetPaymentMethodByCode
	}
	return &paymentMethod, nil
}
