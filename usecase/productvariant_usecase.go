package usecase

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type ProductVariantUsecase interface {
	GetByProductSlug(string) (*dto.ProductDetailVariant, error)
	GetByProductId(uint) (*dto.ProductDetailVariant, error)
}

type ProductVariantUsecaseConfig struct {
	ProductRepository        repository.ProductRepository
	ProductVariantRepository repository.ProductVariantRepository
}

type productVariantUsecaseImpl struct {
	productRepository        repository.ProductRepository
	productVariantRepository repository.ProductVariantRepository
}

func NewProductVariantUsecase(c ProductVariantUsecaseConfig) ProductVariantUsecase {
	return &productVariantUsecaseImpl{
		productRepository:        c.ProductRepository,
		productVariantRepository: c.ProductVariantRepository,
	}
}

func (u *productVariantUsecaseImpl) GetByProductSlug(slug string) (*dto.ProductDetailVariant, error) {
	product, err := u.productRepository.GetProductBySlug(slug)
	if err != nil {
		return nil, err
	}

	productPromotion := entity.ProductPromotion{}
	productPromotionRes, err := u.productRepository.GetProductPromotionByProductId(product.ID)
	if productPromotionRes != nil && err == nil {
		productPromotion = *productPromotionRes
	}

	productVariantItems, err := u.productVariantRepository.GetVariantItemsByProductId(product.ID)
	if err != nil {
		return nil, err
	}

	return u.buildDTOVariant(productVariantItems, productPromotion)
}

func (u *productVariantUsecaseImpl) GetByProductId(productId uint) (*dto.ProductDetailVariant, error) {
	product, err := u.productRepository.GetProductByProductId(productId)
	if err != nil {
		return nil, err
	}

	productPromotion := entity.ProductPromotion{}
	productPromotionRes, err := u.productRepository.GetProductPromotionByProductId(product.ID)
	if productPromotionRes != nil && err == nil {
		productPromotion = *productPromotionRes
	}

	productVariantItems, err := u.productVariantRepository.GetVariantItemsByProductId(product.ID)
	if err != nil {
		return nil, err
	}

	return u.buildDTOVariant(productVariantItems, productPromotion)
}

func (u *productVariantUsecaseImpl) calculateDiscountPrice(price float64, promotion entity.ProductPromotion) float64 {
	if promotion.Promotion.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
		dcPrice := price - promotion.Promotion.Nominal
		if dcPrice < 100 {
			dcPrice = 100
		}
		return dcPrice
	}

	if promotion.Promotion.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
		dcPrice := price - (price * promotion.Promotion.Nominal / 100)
		if dcPrice < 100 {
			dcPrice = 100
		}
		return dcPrice
	}

	if price < 100 {
		price = 100
	}
	return price
}

func (u *productVariantUsecaseImpl) buildDTONestedVariant(variantItems []entity.VariantItem, promotion entity.ProductPromotion) (*dto.ProductDetailVariant, error) {
	productVariantOptions := make([]dto.ProductDetailVariantOption, 0)
	productVariantItems := make([]dto.ProductDetailVariantItem, 0)

	variantMap := make(map[string]map[string]dto.ProductDetailVariantItem)
	for _, variantItem := range variantItems {
		if len(variantItem.VariantSpecs) != 2 {
			return nil, domain.ErrGetProductVariantInconsistent
		}

		variantSpecs := variantItem.VariantSpecs
		if _, ok := variantMap[variantSpecs[0].VariationName]; !ok {
			variantMap[variantSpecs[0].VariationName] = make(map[string]dto.ProductDetailVariantItem)
		}
		variantMap[variantSpecs[0].VariationName][variantSpecs[1].VariationName] = dto.ProductDetailVariantItem{
			ID:            variantItem.ID,
			Key:           fmt.Sprintf("%s,%s", variantSpecs[0].VariationName, variantSpecs[1].VariationName),
			Image:         variantItem.ImageUrl,
			Price:         variantItem.Price,
			DiscountPrice: u.calculateDiscountPrice(variantItem.Price, promotion),
			Stock:         variantItem.Stock,
		}

		for i := 0; i < 2; i++ {
			if len(productVariantOptions) < (i + 1) {
				productVariantOptions = append(productVariantOptions, dto.ProductDetailVariantOption{
					Name: variantSpecs[i].VariationGroup.Name,
				})
			}
			if idx := util.FindIdxArrString(variantSpecs[i].VariationName, productVariantOptions[i].Type); idx == -1 {
				productVariantOptions[i].Type = append(productVariantOptions[i].Type, variantSpecs[i].VariationName)
			}
		}
	}

	for _, variantType := range productVariantOptions[0].Type {
		for _, variantType2 := range productVariantOptions[1].Type {
			productVariantItems = append(productVariantItems, variantMap[variantType][variantType2])
		}
	}

	return &dto.ProductDetailVariant{
		VariantOptions: productVariantOptions,
		VariantItems:   productVariantItems,
	}, nil
}

func (u *productVariantUsecaseImpl) buildDTOSingleVariant(variantItems []entity.VariantItem, promotion entity.ProductPromotion) (*dto.ProductDetailVariant, error) {
	productVariantOptions := make([]dto.ProductDetailVariantOption, 0)
	productVariantItems := make([]dto.ProductDetailVariantItem, 0)

	for _, variantItem := range variantItems {
		if len(variantItem.VariantSpecs) != 1 {
			return nil, domain.ErrGetProductVariantInconsistent
		}

		variantSpecs := variantItem.VariantSpecs
		productVariantItems = append(productVariantItems, dto.ProductDetailVariantItem{
			ID:            variantItem.ID,
			Key:           variantSpecs[0].VariationName,
			Image:         variantItem.ImageUrl,
			Price:         variantItem.Price,
			DiscountPrice: u.calculateDiscountPrice(variantItem.Price, promotion),
			Stock:         variantItem.Stock,
		})

		if len(productVariantOptions) <= 0 {
			productVariantOptions = append(productVariantOptions, dto.ProductDetailVariantOption{
				Name: variantSpecs[0].VariationGroup.Name,
			})
		}
		productVariantOptions[0].Type = append(productVariantOptions[0].Type, variantSpecs[0].VariationName)
	}

	return &dto.ProductDetailVariant{
		VariantOptions: productVariantOptions,
		VariantItems:   productVariantItems,
	}, nil
}

func (u *productVariantUsecaseImpl) buildDTOVariant(variantItems []entity.VariantItem, promotion entity.ProductPromotion) (*dto.ProductDetailVariant, error) {
	if len(variantItems) <= 1 {
		return &dto.ProductDetailVariant{
			VariantOptions: make([]dto.ProductDetailVariantOption, 0),
			VariantItems:   make([]dto.ProductDetailVariantItem, 0),
		}, nil
	}

	if len(variantItems[0].VariantSpecs) == 1 {
		return u.buildDTOSingleVariant(variantItems, promotion)
	} else if len(variantItems[0].VariantSpecs) == 2 {
		return u.buildDTONestedVariant(variantItems, promotion)
	}

	return nil, domain.ErrGetProductVariantInconsistent
}
