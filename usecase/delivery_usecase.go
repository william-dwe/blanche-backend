package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type DeliveryUsecase interface {
	GetAllDeliveryOption() (*dto.DeliveryGetAllOptionResDTO, error)
	GetDeliveryOptionByMerchantDomain(domain string) (*dto.DeliveryGetMerchantOptionResDTO, error)
	GetMerchantDeliveryOption(username string) ([]dto.DeliveryOptionUserMerchantResDTO, error)
	UpdateMerchantDeliveryOption(username string, input []dto.DeliveryUpdateMerchantOptionReqDTO) ([]dto.DeliveryUpdateMerchantOptionResDTO, error)
}

type DeliveryUsecaseConfig struct {
	DeliveryRepository repository.DeliveryRepository
	MerchantRepository repository.MerchantRepository
}

type deliveryUsecaseImpl struct {
	deliveryRepository repository.DeliveryRepository
	merchantRepository repository.MerchantRepository
}

func NewDeliveryUsecase(c DeliveryUsecaseConfig) DeliveryUsecase {
	return &deliveryUsecaseImpl{
		deliveryRepository: c.DeliveryRepository,
		merchantRepository: c.MerchantRepository,
	}
}

func (u *deliveryUsecaseImpl) GetAllDeliveryOption() (*dto.DeliveryGetAllOptionResDTO, error) {
	deliveryOptions, err := u.deliveryRepository.GetAllDeliveryOption()
	if err != nil {
		return nil, err
	}

	var options []dto.DeliveryOption
	for _, v := range deliveryOptions {
		options = append(options, dto.DeliveryOption{
			CourierName: v.CourierName,
			CourierCode: v.CourierCode,
			CourierLogo: v.CourierLogo,
		})
	}

	return &dto.DeliveryGetAllOptionResDTO{
		DeliveryOptions: options,
		Total:           len(options),
	}, nil
}

func (u *deliveryUsecaseImpl) GetDeliveryOptionByMerchantDomain(domain string) (*dto.DeliveryGetMerchantOptionResDTO, error) {
	merchant, err := u.merchantRepository.GetByDomain(domain)
	if err != nil {
		return nil, err
	}
	deliveryOptions, err := u.deliveryRepository.GetDeliveryOptionsByMerchantID(merchant.ID)
	if err != nil {
		return nil, err
	}

	options := make([]dto.DeliveryOption, 0)
	for _, v := range deliveryOptions {
		options = append(options, dto.DeliveryOption{
			CourierName: v.CourierName,
			CourierCode: v.CourierCode,
			CourierLogo: v.CourierLogo,
		})
	}

	return &dto.DeliveryGetMerchantOptionResDTO{
		MerchantDomain:  domain,
		MerchantName:    merchant.Name,
		DeliveryOptions: options,
		Total:           len(options),
	}, nil
}

func (u *deliveryUsecaseImpl) GetMerchantDeliveryOption(username string) ([]dto.DeliveryOptionUserMerchantResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	allDeliveryOptions, err := u.deliveryRepository.GetAllDeliveryOption()
	if err != nil {
		return nil, err
	}

	selectedDeliveryOptions, err := u.deliveryRepository.GetDeliveryOptionsByMerchantID(merchant.ID)
	if err != nil {
		return nil, err
	}

	deliveryMap := make(map[string]string)
	for _, v := range selectedDeliveryOptions {
		deliveryMap[v.CourierCode] = v.CourierName
	}

	res := make([]dto.DeliveryOptionUserMerchantResDTO, 0)
	for _, v := range allDeliveryOptions {
		isChecked := false
		if _, ok := deliveryMap[v.CourierCode]; ok {
			isChecked = true
		}

		deliveryOption := dto.DeliveryOptionUserMerchantResDTO{
			CourierName: v.CourierName,
			CourierCode: v.CourierCode,
			CourierLogo: v.CourierLogo,
			IsChecked:   isChecked,
		}
		res = append(res, deliveryOption)
	}

	return res, nil
}

func (u *deliveryUsecaseImpl) UpdateMerchantDeliveryOption(username string, input []dto.DeliveryUpdateMerchantOptionReqDTO) ([]dto.DeliveryUpdateMerchantOptionResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	allDeliveryOptions, err := u.deliveryRepository.GetAllDeliveryOption()
	if err != nil {
		return nil, err
	}
	deliveryCodeToIdMap := make(map[string]uint)
	for _, v := range allDeliveryOptions {
		deliveryCodeToIdMap[v.CourierCode] = v.ID
	}

	curDeliveryOptions, err := u.deliveryRepository.GetMerchantDeliveryOptionsByMerchantID(merchant.ID)
	if err != nil {
		return nil, err
	}

	curDeliveryMap := make(map[string]entity.MerchantDeliveryOption)
	for _, v := range curDeliveryOptions {
		curDeliveryMap[v.DeliveryOption.CourierCode] = v
	}

	var isAnyInsert, isAnyDelete bool = false, false
	insertSlice := []entity.MerchantDeliveryOption{}
	deleteSlice := []uint{}
	resBody := []dto.DeliveryUpdateMerchantOptionResDTO{}
	for _, v := range input {
		if v.IsChecked {
			if _, ok := curDeliveryMap[v.CourierCode]; !ok {
				insertSlice = append(insertSlice,
					entity.MerchantDeliveryOption{
						MerchantId:       merchant.ID,
						DeliveryOptionId: deliveryCodeToIdMap[v.CourierCode],
					},
				)
				isAnyInsert = true
			}
		} else {

			if _, ok := curDeliveryMap[v.CourierCode]; ok {
				deleteSlice = append(deleteSlice, curDeliveryMap[v.CourierCode].ID)
				isAnyDelete = true
			}
		}
		resBody = append(resBody, dto.DeliveryUpdateMerchantOptionResDTO(v))
	}

	if isAnyInsert {
		_, err = u.deliveryRepository.AddMerchantDeliveryOptions(insertSlice)
		if err != nil {
			return nil, err
		}
	}

	if isAnyDelete {

		err = u.deliveryRepository.DeleteMerchantDeliveryOptions(merchant.ID, deleteSlice)
		if err != nil {
			return nil, err
		}
	}

	return resBody, nil
}
