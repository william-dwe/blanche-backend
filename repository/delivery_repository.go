package repository

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type DeliveryRepository interface {
	GetAllDeliveryOption() ([]entity.DeliveryOption, error)
	GetDeliveryOptionByMerchantID(merchantId uint, courierCode string) (*entity.DeliveryOption, error)
	GetDeliveryOptionsByMerchantID(merchantID uint) ([]entity.DeliveryOption, error)
	GetMerchantDeliveryOptionsByMerchantID(merchantId uint) ([]entity.MerchantDeliveryOption, error)
	GetDeliveryInfo(input dto.RajaOngkirDeliveryInfoReqDTO) (*dto.RajaOngkirDeliveryInfoResDTO, error)
	AddMerchantDeliveryOptions(input []entity.MerchantDeliveryOption) ([]entity.MerchantDeliveryOption, error)
	AddMerchantDeliveryOptionsTx(tx *gorm.DB, input []entity.MerchantDeliveryOption) ([]entity.MerchantDeliveryOption, error)
	DeleteMerchantDeliveryOptions(merchantId uint, input []uint) error
}

type DeliveryRepositoryConfig struct {
	DB *gorm.DB
}

type deliveryRepositoryImpl struct {
	db *gorm.DB
}

func NewDeliveryRepository(c DeliveryRepositoryConfig) DeliveryRepository {
	return &deliveryRepositoryImpl{
		db: c.DB,
	}
}

func (r *deliveryRepositoryImpl) GetAllDeliveryOption() ([]entity.DeliveryOption, error) {
	var deliveryOptions []entity.DeliveryOption
	err := r.db.Find(&deliveryOptions).Error
	if err != nil {
		return nil, domain.ErrGetAllDeliveryInternalError
	}
	return deliveryOptions, nil
}

func (r *deliveryRepositoryImpl) GetDeliveryOptionByMerchantID(merchantId uint, courierCode string) (*entity.DeliveryOption, error) {
	var deliveryOptions entity.DeliveryOption
	sq := r.db.Select("delivery_option_id").
		Model(&entity.MerchantDeliveryOption{}).
		Where("merchant_id = ? AND courier_code = ?", merchantId, courierCode)
	err := r.db.Where("id in (?)", sq).
		First(&deliveryOptions).
		Error

	if err != nil {
		return nil, domain.ErrGetAllDeliveryInternalError
	}
	return &deliveryOptions, nil
}

func (r *deliveryRepositoryImpl) GetDeliveryOptionsByMerchantID(merchantId uint) ([]entity.DeliveryOption, error) {
	var res []entity.DeliveryOption
	sq := r.db.Select("delivery_option_id").
		Model(&entity.MerchantDeliveryOption{}).
		Where("merchant_id = ?", merchantId)
	err := r.db.Where("id in (?)", sq).
		Find(&res).
		Error

	if err != nil {
		return nil, domain.ErrGetAllDeliveryInternalError
	}
	return res, nil
}

func (r *deliveryRepositoryImpl) GetMerchantDeliveryOptionsByMerchantID(merchantId uint) ([]entity.MerchantDeliveryOption, error) {
	var res []entity.MerchantDeliveryOption
	err := r.db.Preload("DeliveryOption").
		Where("merchant_id = ?", merchantId).
		Find(&res).
		Error

	if err != nil {
		return nil, domain.ErrUpdateDeliveryOptionInternalError
	}

	return res, nil
}

func (r *deliveryRepositoryImpl) GetDeliveryInfo(input dto.RajaOngkirDeliveryInfoReqDTO) (*dto.RajaOngkirDeliveryInfoResDTO, error) {
	c := config.Config.RajaOngkirConfig
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, domain.ErrGetDeliveryFeeInternalError
	}
	req, err := http.NewRequest("POST", c.Url+"/cost", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, domain.ErrGetDeliveryFeeInternalError
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("key", c.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, domain.ErrGetDeliveryFeeInternalError
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.ErrGetDeliveryFeeInternalError
	}

	var output dto.RajaOngkirDeliveryInfoResDTO
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

func (r *deliveryRepositoryImpl) AddMerchantDeliveryOptions(input []entity.MerchantDeliveryOption) ([]entity.MerchantDeliveryOption, error) {
	err := r.db.Create(&input).Error
	if err != nil {
		return nil, domain.ErrUpdateDeliveryOptionInternalError
	}
	return input, err
}

func (r *deliveryRepositoryImpl) AddMerchantDeliveryOptionsTx(tx *gorm.DB, input []entity.MerchantDeliveryOption) ([]entity.MerchantDeliveryOption, error) {
	err := tx.Create(&input).Error
	if err != nil {
		return nil, domain.ErrUpdateDeliveryOptionInternalError
	}
	return input, err
}

func (r *deliveryRepositoryImpl) DeleteMerchantDeliveryOptions(merchantId uint, input []uint) error {

	err := r.db.Where("merchant_id = ?", merchantId).
		Where("id in (?)", input).
		Delete(&entity.MerchantDeliveryOption{}).
		Error
	if err != nil {
		return domain.ErrUpdateDeliveryOptionInternalError
	}
	return err
}
