package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"gorm.io/gorm"
)

type CartItemUsecase interface {
	AddCartItem(username string, input dto.AddItemToCartReqDTO) (*dto.AddItemToCartResDTO, error)
	GetCartItems(username string) (*dto.GetCartItemResDTO, error)
	GetHomeCartItems(username string) (*dto.GetHomeCartItemResDTO, error)
	DeleteCartItemByCartId(username string, cartItemId uint) error
	DeleteSelectedCartItem(username string) error
	UpdateCartItem(username string, cartItemId uint, input dto.UpdateCartItemDTO) (*dto.UpdateCartItemDTO, error)
	UpdateAllCartCheckStatus(username string, input []dto.UpdateAllCartCheckStatusDTO) ([]dto.UpdateAllCartCheckStatusDTO, error)
}

type CartItemUsecaseConfig struct {
	CartItemRepository       repository.CartItemRepository
	UserRepository           repository.UserRepository
	MerchantRepository       repository.MerchantRepository
	ProductRepository        repository.ProductRepository
	ProductVariantRepository repository.ProductVariantRepository
}

type cartItemUsecaseImpl struct {
	cartItemRepository       repository.CartItemRepository
	userRepository           repository.UserRepository
	merchantRepository       repository.MerchantRepository
	productRepository        repository.ProductRepository
	productVariantRepository repository.ProductVariantRepository
}

func NewCartItemUsecase(c CartItemUsecaseConfig) CartItemUsecase {
	return &cartItemUsecaseImpl{
		cartItemRepository:       c.CartItemRepository,
		userRepository:           c.UserRepository,
		merchantRepository:       c.MerchantRepository,
		productRepository:        c.ProductRepository,
		productVariantRepository: c.ProductVariantRepository,
	}
}

func (u *cartItemUsecaseImpl) AddCartItem(username string, input dto.AddItemToCartReqDTO) (*dto.AddItemToCartResDTO, error) {
	if input.ProductId == 0 || input.Quantity == 0 {
		return nil, domain.ErrAddCartItemInvalidInput
	}

	product, err := u.productRepository.GetProductMerchantById(input.ProductId)
	if err != nil {
		return nil, err
	}
	if product.IsArchived {
		return nil, domain.ErrAddCartProductNotAvailable
	}

	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if product.Merchant.UserId == user.ID {
		return nil, domain.ErrAddCartItemAddOwnProduct
	}

	newCartItem := entity.CartItem{
		UserId:    user.ID,
		ProductId: input.ProductId,
		Quantity:  input.Quantity,
		IsChecked: false,
	}

	var stock int
	if input.VariantItemId != nil {
		variantItem, err := u.productVariantRepository.GetVariantItemById(*input.VariantItemId)
		if err != nil {
			return nil, err
		}
		if variantItem.ProductId != input.ProductId || variantItem == nil {
			return nil, domain.ErrAddCartVariantItemNotExist
		}
		stock = int(variantItem.Stock)
		newCartItem.VariantItemId = input.VariantItemId
	} else {
		variantItems, err := u.productVariantRepository.GetVariantItemsByProductId(input.ProductId)
		if err != nil {
			return nil, err
		}
		if len(variantItems) == 0 {
			return nil, domain.ErrAddCartInternalError
		}
		if len(variantItems) > 1 {
			return nil, domain.ErrAddCartItemNeedVariantItem
		}
		stock = int(variantItems[0].Stock)
		newCartItem.VariantItemId = &variantItems[0].ID
	}

	if stock < input.Quantity {
		return nil, domain.ErrAddCartItemQuantityExceedStock
	}

	sameExistingCartItem, err := u.cartItemRepository.GetSameExistingCartItem(user.ID, &newCartItem)
	if err != nil {
		return nil, err
	}

	resBody := dto.AddItemToCartResDTO{}
	if sameExistingCartItem != nil {
		sameExistingCartItem.Quantity += newCartItem.Quantity

		if sameExistingCartItem.Quantity > stock {
			return nil, domain.ErrAddCartItemQuantityExceedStock
		}

		err = u.cartItemRepository.UpdateCartItem(sameExistingCartItem)
		if err != nil {
			return nil, err
		}

		resBody = dto.AddItemToCartResDTO{
			ProductId:     input.ProductId,
			VariantItemId: input.VariantItemId,
			Quantity:      sameExistingCartItem.Quantity,
			IsChecked:     sameExistingCartItem.IsChecked,
		}
		if sameExistingCartItem.Notes != nil {
			resBody.Notes = *sameExistingCartItem.Notes
		}
	} else {
		err = u.cartItemRepository.AddCartItem(&newCartItem)
		if err != nil {
			return nil, err
		}

		resBody = dto.AddItemToCartResDTO{
			ProductId:     input.ProductId,
			VariantItemId: input.VariantItemId,
			Quantity:      input.Quantity,
			IsChecked:     false,
			Notes:         "",
		}
	}
	return &resBody, nil

}

func (u *cartItemUsecaseImpl) GetCartItems(username string) (*dto.GetCartItemResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	cartItems, err := u.cartItemRepository.GetCartItems(user.ID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, nil
	}

	var cartMerchantMap = make(map[uint][]dto.CartItemDTO)
	var cartMerchantKeys []uint
	var total float64
	var quantity int
	for _, cartItem := range cartItems {
		newCartPH := u.fillInCartDTO(cartItem)

		if cartMerchantMap[cartItem.Product.MerchantId] == nil {
			cartMerchantKeys = append(cartMerchantKeys, cartItem.Product.MerchantId)
		}
		cartMerchantMap[cartItem.Product.MerchantId] = append(cartMerchantMap[cartItem.Product.MerchantId], newCartPH)

		if newCartPH.IsChecked {
			if newCartPH.IsPromotionPriceValid {
				total += newCartPH.DiscountPrice * float64(newCartPH.Quantity)
			} else {
				total += newCartPH.RealPrice * float64(newCartPH.Quantity)
			}
			quantity += newCartPH.Quantity
		}
	}

	var cart []dto.CartItemPerStoreDTO
	for _, key := range cartMerchantKeys {
		cart = append(cart, dto.CartItemPerStoreDTO{
			MerchantId:     key,
			MerchantName:   cartMerchantMap[key][0].MerchantName,
			MerchantImage:  cartMerchantMap[key][0].MerchantImage,
			MerchantDomain: cartMerchantMap[key][0].MerchantDomain,
			Items:          cartMerchantMap[key],
		})
	}

	resBody := dto.GetCartItemResDTO{
		Carts:    cart,
		Total:    total,
		Quantity: quantity,
	}
	return &resBody, nil
}

func (u *cartItemUsecaseImpl) GetHomeCartItems(username string) (*dto.GetHomeCartItemResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	cartItems, err := u.cartItemRepository.GetCartItems(user.ID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, nil
	}

	var cartList []dto.CartItemDTO
	var qty int
	for _, cartItem := range cartItems {
		newCartPH := u.fillInCartDTO(cartItem)

		qty += newCartPH.Quantity
		cartList = append(cartList, newCartPH)
	}

	resBody := dto.GetHomeCartItemResDTO{
		Carts:    cartList,
		Quantity: qty,
	}
	return &resBody, nil
}

func (u *cartItemUsecaseImpl) fillInCartDTO(cartItem entity.CartItem) dto.CartItemDTO {
	discountPrice := cartItem.VariantItem.Price

	if cartItem.Product.ProductPromotion != nil && cartItem.Product.ProductPromotion.Promotion.Quota >= cartItem.Quantity {
		if cartItem.Product.ProductPromotion.Promotion.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
			discountPrice -= cartItem.Product.ProductPromotion.Promotion.Nominal
		} else if cartItem.Product.ProductPromotion.Promotion.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
			discountPrice -= (cartItem.VariantItem.Price * cartItem.Product.ProductPromotion.Promotion.Nominal / 100)
		}
	}
	if discountPrice < 100 {
		discountPrice = 100
	}

	newCartPH := dto.CartItemDTO{
		CartItemId:     cartItem.ID,
		ProductId:      cartItem.ProductId,
		ProductSlug:    cartItem.Product.Slug,
		VariantItemId:  cartItem.VariantItemId,
		VariantName:    "",
		MerchantId:     cartItem.Product.MerchantId,
		MerchantName:   cartItem.Product.Merchant.Name,
		MerchantImage:  cartItem.Product.Merchant.ImageUrl,
		MerchantDomain: cartItem.Product.Merchant.Domain,
		Name:           cartItem.Product.Title,
		RealPrice:      cartItem.VariantItem.Price,
		DiscountPrice:  discountPrice,
		Quantity:       cartItem.Quantity,
		Stock:          int(cartItem.VariantItem.Stock),
		Notes:          cartItem.Notes,
		IsChecked:      cartItem.IsChecked,
		IsValid: cartItem.Product.DeletedAt == gorm.DeletedAt{} &&
			cartItem.VariantItem.DeletedAt == gorm.DeletedAt{} &&
			cartItem.Quantity <= int(cartItem.VariantItem.Stock),
		IsPromotionPriceValid: cartItem.Product.ProductPromotion != nil &&
			cartItem.Product.ProductPromotion.Promotion.Quota >= cartItem.Quantity &&
			cartItem.Product.ProductPromotion.Promotion.MaxDiscountedQty >= cartItem.Quantity,
	}

	if !newCartPH.IsPromotionPriceValid {
		newCartPH.DiscountPrice = cartItem.VariantItem.Price
	}

	if cartItem.VariantItemId != nil {
		newCartPH.VariantName = cartItem.VariantItem.VariantSpec.VariationName
	}

	if len(cartItem.Product.ProductImages) != 0 {
		newCartPH.Image = cartItem.Product.ProductImages[0].ImageUrl
	}

	return newCartPH
}

func (u *cartItemUsecaseImpl) DeleteSelectedCartItem(username string) error {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return err
	}
	err = u.cartItemRepository.DeleteSelectedCartItem(user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartItemUsecaseImpl) DeleteCartItemByCartId(username string, cartItemId uint) error {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return err
	}
	err = u.cartItemRepository.DeleteCartItemByCartId(user.ID, cartItemId)
	if err != nil {
		return err
	}
	return nil
}

func (u *cartItemUsecaseImpl) UpdateCartItem(username string, cartItemId uint, input dto.UpdateCartItemDTO) (*dto.UpdateCartItemDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	cartItem, err := u.cartItemRepository.GetCartItem(user.ID, cartItemId)
	if err != nil {
		return nil, err
	}

	variantItem, err := u.productVariantRepository.GetVariantItemById(*cartItem.VariantItemId)
	if err != nil {
		return nil, err
	}
	if variantItem == nil {
		return nil, domain.ErrUpdateCartVariantItemNotExist
	}
	if int(variantItem.Stock) < input.Quantity {
		return nil, domain.ErrUpdateCartItemQuantityExceedStock
	}

	updatedCartItem := entity.CartItem{
		ID:       cartItemId,
		UserId:   user.ID,
		Quantity: input.Quantity,
		Notes:    &input.Notes,
	}

	err = u.cartItemRepository.UpdateCartItem(&updatedCartItem)
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (u *cartItemUsecaseImpl) UpdateAllCartCheckStatus(username string, input []dto.UpdateAllCartCheckStatusDTO) ([]dto.UpdateAllCartCheckStatusDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	cartItemIds := make(map[bool][]uint)
	for _, v := range input {
		cartItemIds[v.IsChecked] = append(cartItemIds[v.IsChecked], v.CartItemId)
	}

	for k, v := range cartItemIds {
		if len(v) > 0 {
			err = u.cartItemRepository.UpdateCartItems(user.ID, v, &map[string]interface{}{"is_checked": k})
			if err != nil {
				return nil, err
			}
		}
	}

	return input, nil
}
