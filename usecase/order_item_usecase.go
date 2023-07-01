package usecase

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type OrderItemUsecase interface {
	GetOrderCheckoutSummary(username string, input dto.PostOrderSummaryReqDTO) (*dto.PostOrderSummaryResDTO, error)
	MakeOrderCheckout(username string, input []dto.MakeOrderCheckoutProductDTO) (*dto.PostOrderSummaryResDTO, error)
}

type OrderItemUsecaseConfig struct {
	OrderItemRepository          repository.OrderItemRepository
	CartItemRepository           repository.CartItemRepository
	UserRepository               repository.UserRepository
	ProductRepository            repository.ProductRepository
	ProductVariantRepository     repository.ProductVariantRepository
	DeliveryRepository           repository.DeliveryRepository
	AddressRepository            repository.AddressRepository
	MarketplaceVoucherRepository repository.MarketplaceVoucherRepository
	MerchantRepository           repository.MerchantRepository
	UserOrderRepository          repository.UserOrderRepository
}

type orderItemUsecaseImpl struct {
	orderItemRepository          repository.OrderItemRepository
	cartItemRepository           repository.CartItemRepository
	userRepository               repository.UserRepository
	productRepository            repository.ProductRepository
	productVariantRepository     repository.ProductVariantRepository
	deliveryRepository           repository.DeliveryRepository
	addressRepository            repository.AddressRepository
	marketplaceVoucherRepository repository.MarketplaceVoucherRepository
	merchantRepository           repository.MerchantRepository
	userOrderRepository          repository.UserOrderRepository
}

func NewOrderItemUsecase(c OrderItemUsecaseConfig) OrderItemUsecase {
	return &orderItemUsecaseImpl{
		orderItemRepository:          c.OrderItemRepository,
		cartItemRepository:           c.CartItemRepository,
		userRepository:               c.UserRepository,
		deliveryRepository:           c.DeliveryRepository,
		productRepository:            c.ProductRepository,
		productVariantRepository:     c.ProductVariantRepository,
		addressRepository:            c.AddressRepository,
		marketplaceVoucherRepository: c.MarketplaceVoucherRepository,
		merchantRepository:           c.MerchantRepository,
		userOrderRepository:          c.UserOrderRepository,
	}
}

func (u *orderItemUsecaseImpl) MakeOrderCheckout(username string, orderItemsInput []dto.MakeOrderCheckoutProductDTO) (*dto.PostOrderSummaryResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	newOrderCode := util.GenerateUUIDWithDate()
	var orderItems []entity.OrderItem

	// check all product that will be ordered
	for _, orderItem := range orderItemsInput {
		orderItemValidated, err := u.checkOrderItemAvailability(user.ID, orderItem)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, entity.OrderItem{
			ProductId:     orderItemValidated.ProductId,
			VariantItemId: *orderItemValidated.VariantItemId,
			Quantity:      uint(orderItemValidated.Quantity),
			Notes:         orderItemValidated.Notes,
		})
	}

	// create order
	order := entity.UserOrder{
		UserId:     user.ID,
		OrderCode:  newOrderCode,
		OrderItems: orderItems,
	}

	createdOrder, err := u.userOrderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return u.GetOrderCheckoutSummary(username, dto.PostOrderSummaryReqDTO{
		OrderCode: createdOrder.OrderCode,
	})
}

func (u *orderItemUsecaseImpl) GetOrderCheckoutSummary(username string, input dto.PostOrderSummaryReqDTO) (*dto.PostOrderSummaryResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	address, err := u.getUserAddress(uint(input.AddressId), *user)
	if err != nil {
		return nil, domain.ErrOrderAddressNotFound
	}

	userOrder, err := u.userOrderRepository.FindOrderByOrderCode(input.OrderCode, user.ID)
	if err != nil {
		return nil, err
	}

	orderItems := userOrder.OrderItems
	if len(orderItems) == 0 {
		return nil, domain.ErrNoOrderItem
	}

	var isOrderValid = true
	var orderMerchantMap = make(map[uint][]dto.OrderItemDTO)
	var cartMerchantKeys []uint
	var merchantTotalMap = make(map[uint]float64)
	var merchantWeightMap = make(map[uint]int)
	var trxTotal float64
	for _, orderItem := range orderItems {
		newOrderPH := u.fillOrderItemDTO(orderItem)
		isOrderValid = isOrderValid && newOrderPH.IsValid

		if orderItem.Product.Merchant.UserId == user.ID {
			return nil, domain.ErrOrderButOwnProduct
		}

		if orderMerchantMap[orderItem.Product.MerchantId] == nil {
			cartMerchantKeys = append(cartMerchantKeys, orderItem.Product.MerchantId)
		}
		orderMerchantMap[orderItem.Product.MerchantId] = append(orderMerchantMap[orderItem.Product.MerchantId], newOrderPH)

		total := newOrderPH.DiscountPrice * float64(newOrderPH.Quantity)
		weight := orderItem.Product.Weight * int(orderItem.Quantity)
		trxTotal += total
		merchantTotalMap[orderItem.Product.MerchantId] += total
		merchantWeightMap[orderItem.Product.MerchantId] += weight
	}

	var mapMerchantVoucher = make(map[uint]string)
	var mapMerchantDeliveryOption = make(map[uint]string)
	if len(input.Merchants) > 0 {
		for _, v := range input.Merchants {
			mapMerchantVoucher[uint(v.MerchantId)] = v.VoucherMerchant
			mapMerchantDeliveryOption[uint(v.MerchantId)] = v.DeliveryOption
		}
	}

	var isVoucherInvalid = false
	var trxDelivery float64
	var trxSellerDiscount float64
	var order []dto.OrderItemPerMerchantDTO
	for _, key := range cartMerchantKeys {
		userCity, err := u.addressRepository.GetCityById(address.CityId)
		if err != nil {
			return nil, err
		}
		merchantCity, err := u.addressRepository.GetCityById(orderMerchantMap[key][0].MerchantCityId)
		if err != nil {
			return nil, err
		}

		var sellerDiscount float64 = 0
		var merchantVoucherId *uint = nil
		if len(mapMerchantVoucher) > 0 {
			if val, ok := mapMerchantVoucher[key]; ok && val != "" {
				sellerVoucher, _ := u.merchantRepository.GetMerchantVoucherByCode(orderMerchantMap[key][0].MerchantDomain, val)
				if sellerVoucher != nil {
					merchantVoucherId = &sellerVoucher.ID
					if sellerVoucher.MinOrderNominal <= merchantTotalMap[key] {
						sellerDiscount = sellerVoucher.DiscountNominal
						if sellerVoucher.DiscountNominal >= merchantTotalMap[key] {
							sellerDiscount = merchantTotalMap[key]
						}
					} else {
						isVoucherInvalid = true
					}
				}
				if sellerVoucher == nil {
					isVoucherInvalid = true
				}
			}
		}

		var orderItem = dto.OrderItemPerMerchantDTO{
			Merchant: dto.OrderMerchantDTO{
				MerchantId:     key,
				MerchantName:   orderMerchantMap[key][0].MerchantName,
				MerchantImage:  orderMerchantMap[key][0].MerchantImage,
				MerchantDomain: orderMerchantMap[key][0].MerchantDomain,
			},
			Items:             orderMerchantMap[key],
			SubTotal:          merchantTotalMap[key],
			Discount:          sellerDiscount,
			IsVoucherInvalid:  isVoucherInvalid,
			MerchantVoucherId: merchantVoucherId,
		}
		trxSellerDiscount += sellerDiscount

		if mapMerchantDeliveryOption[key] != "" {
			deliveryInput := dto.RajaOngkirDeliveryInfoReqDTO{
				Origin:      int(merchantCity.RoId),
				Destination: int(userCity.RoId),
				Weight:      merchantWeightMap[key],
				Courier:     mapMerchantDeliveryOption[key],
			}

			deliveryOption, err := u.deliveryRepository.GetDeliveryOptionByMerchantID(key, mapMerchantDeliveryOption[key])
			if err != nil {
				return nil, err
			}

			deliveryInfo, err := u.deliveryRepository.GetDeliveryInfo(deliveryInput)
			if err != nil {
				return nil, err
			}

			var isSelectedDelivery = false
			for _, v := range deliveryInfo.Rajaongkir.Results[0].Costs {
				if util.IsSliceContainString(strings.Split(deliveryOption.ServiceCode, ","), v.Service) {
					isSelectedDelivery = true
					deliveryCost := v.Cost[0].Value

					orderItem.DeliveryService = dto.DeliveryServiceDTO{
						DeliveryOption: mapMerchantDeliveryOption[key],
						Name:           deliveryInfo.Rajaongkir.Results[0].Name,
						Service:        v.Service,
						Description:    v.Description,
						MerchantCity:   merchantCity.Name,
						UserCity:       userCity.Name,
						Etd:            v.Cost[0].Etd,
						Note:           v.Cost[0].Note,
					}
					orderItem.DeliveryCost = deliveryCost
					orderItem.Total = merchantTotalMap[key] + deliveryCost - sellerDiscount
					trxDelivery += deliveryCost
				}
			}
			if !isSelectedDelivery {
				isOrderValid = false
			}
		} else {
			isOrderValid = false
		}
		order = append(order, orderItem)
	}

	var isMpVoucherInvalid bool
	var marketplaceDiscount float64 = 0
	var marketplaceVoucherId *uint = nil
	if input.VoucherMarketplace != "" {
		mpVoucher, err := u.marketplaceVoucherRepository.GetMarketplaceVoucherByCode(input.VoucherMarketplace)
		isMpVoucherInvalid = true
		if mpVoucher != nil && err == nil {
			if mpVoucher.MinOrderNominal <= trxTotal {
				marketplaceVoucherId = &mpVoucher.ID
				marketplaceDiscount = float64(mpVoucher.DiscountPercentage) / 100 * (trxTotal - trxSellerDiscount)
				if marketplaceDiscount > float64(mpVoucher.MaxDiscountNominal) {
					marketplaceDiscount = float64(mpVoucher.MaxDiscountNominal)
				}
				isMpVoucherInvalid = false
			}
		}
	}

	resBody := dto.PostOrderSummaryResDTO{
		OrderCode:           userOrder.OrderCode,
		Orders:              order,
		SubTotal:            trxTotal,
		DeliveryCost:        trxDelivery,
		DiscountMerchant:    trxSellerDiscount,
		DiscountMarketplace: marketplaceDiscount,
		Total:               trxTotal + trxDelivery - trxSellerDiscount - marketplaceDiscount,
		IsVouchervalid:      !isMpVoucherInvalid,
		IsOrderEligible:     !userOrder.DeletedAt.Valid,
		IsOrderValid:        isOrderValid,

		MarketplaceVoucherId: marketplaceVoucherId,
		Address:              *address,
	}
	return &resBody, nil
}

func (u *orderItemUsecaseImpl) checkOrderItemAvailability(userId uint, orderItem dto.MakeOrderCheckoutProductDTO) (*dto.MakeOrderCheckoutProductDTO, error) {
	if orderItem.Quantity <= 0 {
		return nil, domain.ErrOrderQuantityNotValid
	}

	product, err := u.productRepository.GetProductByProductId(orderItem.ProductId)
	if err != nil {
		return nil, domain.ErrOrderProductNotAvailable
	}

	if product.Merchant.UserId == userId {
		return nil, domain.ErrOrderButOwnProduct
	}

	if product.IsArchived {
		return nil, domain.ErrOrderProductNotAvailable
	}

	stock := uint(0)
	if orderItem.VariantItemId == nil {
		variantItems, err := u.productVariantRepository.GetVariantItemsByProductId(orderItem.ProductId)
		if err != nil {
			return nil, err
		}
		if len(variantItems) != 1 {
			return nil, domain.ErrOrderProductNotAvailable
		}
		orderItem.VariantItemId = &variantItems[0].ID
		stock = variantItems[0].Stock
	} else {
		variantItem, err := u.productVariantRepository.GetVariantItemById(*orderItem.VariantItemId)
		if err != nil {
			return nil, domain.ErrOrderProductNotAvailable
		}
		stock = variantItem.Stock
	}

	if stock < uint(orderItem.Quantity) {
		return nil, domain.ErrOrderProductNotAvailable
	}

	return &orderItem, nil
}

func (u *orderItemUsecaseImpl) getUserAddress(addressId uint, user entity.User) (*entity.UserAddress, error) {
	if addressId == 0 {
		address, err := u.userRepository.GetDefaultUserAddress(user)
		if err != nil {
			return nil, err
		}

		return address, nil
	}

	address, err := u.userRepository.GetUserAddressById(user.ID, addressId)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (u *orderItemUsecaseImpl) fillOrderItemDTO(orderItem entity.OrderItem) dto.OrderItemDTO {
	discountPrice := orderItem.VariantItem.Price
	if orderItem.Product.ProductPromotion != nil && orderItem.Product.ProductPromotion.Promotion.MaxDiscountedQty >= int(orderItem.Quantity) {
		if orderItem.Product.ProductPromotion.Promotion.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
			discountPrice -= orderItem.Product.ProductPromotion.Promotion.Nominal
		} else if orderItem.Product.ProductPromotion.Promotion.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
			discountPrice -= (orderItem.VariantItem.Price * orderItem.Product.ProductPromotion.Promotion.Nominal / 100)
		}
	}
	if discountPrice < 100 {
		discountPrice = 100
	}

	newOrderPH := dto.OrderItemDTO{
		CartItemId:     orderItem.ID,
		ProductId:      orderItem.ProductId,
		ProductSlug:    orderItem.Product.Slug,
		VariantItemId:  &orderItem.VariantItemId,
		VariantName:    "",
		MerchantId:     orderItem.Product.MerchantId,
		MerchantName:   orderItem.Product.Merchant.Name,
		MerchantImage:  orderItem.Product.Merchant.ImageUrl,
		MerchantDomain: orderItem.Product.Merchant.Domain,
		MerchantCityId: orderItem.Product.Merchant.CityId,
		Name:           orderItem.Product.Title,
		Weight:         orderItem.Product.Weight,
		RealPrice:      orderItem.VariantItem.Price,
		DiscountPrice:  discountPrice,
		Quantity:       int(orderItem.Quantity),
		Stock:          int(orderItem.VariantItem.Stock),
		Notes:          &orderItem.Notes,
		IsValid: orderItem.Product.DeletedAt == gorm.DeletedAt{} &&
			orderItem.VariantItem.DeletedAt == gorm.DeletedAt{} &&
			orderItem.Quantity <= orderItem.VariantItem.Stock,
	}

	if len(orderItem.Product.ProductImages) != 0 {
		newOrderPH.Image = orderItem.Product.ProductImages[0].ImageUrl
	}

	productVariant, err := u.productVariantRepository.GetVariantItemById(orderItem.VariantItemId)
	if err == nil && productVariant != nil {
		for idx, specs := range productVariant.VariantSpecs {
			if idx > 0 {
				newOrderPH.VariantName += ","
			}
			newOrderPH.VariantName += specs.VariationName
		}
	}

	return newOrderPH
}
