package usecase

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type PromotionUsecase interface {
	GetAllPromotions(username string, req dto.PromotionListReqParamDTO) (*dto.PromotionListResDTO, error)
	CreateNewPromotion(username string, req dto.UpsertPromotionReqDTO) (*dto.PromotionResDTO, error)
	GetPromotionByID(id uint) (*dto.PromotionDetailResDTO, error)
	UpdatePromotion(username string, promotionId int, req dto.UpsertPromotionReqDTO) (*dto.PromotionResDTO, error)
	DeletePromotion(promotionId int) (*dto.PromotionResDTO, error)
}

type PromotionUsecaseConfig struct {
	PromotionRepository repository.PromotionRepository
	MerchantRepository  repository.MerchantRepository
	ProductRepository   repository.ProductRepository
}

type promotionUsecaseImpl struct {
	promotionRepository repository.PromotionRepository
	merchantRepository  repository.MerchantRepository
	productRepository   repository.ProductRepository
}

func NewPromotionUsecase(c PromotionUsecaseConfig) PromotionUsecase {
	return &promotionUsecaseImpl{
		promotionRepository: c.PromotionRepository,
		merchantRepository:  c.MerchantRepository,
		productRepository:   c.ProductRepository,
	}
}

func (u *promotionUsecaseImpl) GetAllPromotions(username string, req dto.PromotionListReqParamDTO) (*dto.PromotionListResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	promotions, total, err := u.promotionRepository.GetAllPromotions(req, merchant.ID)
	if err != nil {
		return nil, err
	}

	promotionsDTO := make([]dto.PromotionResDTO, len(promotions))
	for i, promotion := range promotions {
		promotionsDTO[i] = dto.PromotionResDTO{
			ID:                    promotion.ID,
			PromotionType:         promotion.PromotionType.Name,
			Title:                 promotion.Title,
			MaxDiscountedQuantity: promotion.MaxDiscountedQty,
			Quota:                 promotion.Quantity,
			UsedQuota:             promotion.Quantity - promotion.Quota,
			StartDate:             promotion.StartAt,
			EndDate:               promotion.EndAt,
			Products:              make([]dto.ProductSellerResDTO, len(promotion.ProductPromotions)),
		}
		if promotion.PromotionType.ID == dto.NOMINAL_PROMOTION_ID {
			promotionsDTO[i].DiscountNominal = promotion.Nominal
		}
		if promotion.PromotionType.ID == dto.PERCENTAGE_PROMOTION_ID {
			promotionsDTO[i].DiscountPercentage = promotion.Nominal
		}

		for j, productPromotion := range promotion.ProductPromotions {
			promotionsDTO[i].Products[j] = dto.ProductSellerResDTO{
				ID:           productPromotion.Product.ID,
				Title:        productPromotion.Product.Title,
				Slug:         productPromotion.Product.Slug,
				NumOfSale:    productPromotion.Product.ProductAnalytic.NumOfSale,
				TotalStock:   productPromotion.Product.ProductAnalytic.TotalStock,
				MinRealPrice: productPromotion.Product.MinRealPrice,
				MaxRealPrice: productPromotion.Product.MaxRealPrice,
			}
			if productPromotion.Product.ProductImages != nil && len(productPromotion.Product.ProductImages) > 0 {
				promotionsDTO[i].Products[j].ThumbnailImg = productPromotion.Product.ProductImages[0].ImageUrl
			}
		}
	}

	return &dto.PromotionListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   total,
			TotalPage:   (total + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
		Promotions: promotionsDTO,
	}, nil
}

func (u *promotionUsecaseImpl) CreateNewPromotion(username string, req dto.UpsertPromotionReqDTO) (*dto.PromotionResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	newPromotion := &entity.Promotion{
		MerchantId:       merchant.ID,
		PromotionTypeId:  req.PromotionTypeId,
		Title:            req.Title,
		Nominal:          req.Nominal,
		MaxDiscountedQty: req.MaxDiscountedQuantity,
		Quota:            req.Quota,
		Quantity:         req.Quota,
		StartAt:          req.StartDate,
		EndAt:            req.EndDate,
	}

	ok, err := u.promotionRepository.CheckProductPromotionOngoing(req.ProductIds, req.StartDate)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrCheckProductPromotionOngoing
	}

	if req.StartDate.After(req.EndDate) {
		return nil, domain.ErrInvalidPromotionDateRange
	}

	if req.PromotionTypeId != dto.NOMINAL_PROMOTION_ID && req.PromotionTypeId != dto.PERCENTAGE_PROMOTION_ID {
		return nil, domain.ErrInvalidPromotionType
	}

	if req.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
		if req.Nominal <= 100 {
			return nil, domain.ErrInvalidNominal
		}
	}
	if req.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
		if req.Nominal <= 1 || req.Nominal > 100 {
			return nil, domain.ErrInvalidPercentage
		}
	}

	for _, productId := range req.ProductIds {
		product, err := u.productRepository.GetProductByProductId(productId)
		if err != nil {
			return nil, err
		}
		if product.MerchantDomain != merchant.Domain {
			return nil, domain.ErrInvalidProduct
		}

		minDiscountPrice, maxDiscountPrice := u.calculateMinMaxDiscountPrice(product, newPromotion)
		newPromotion.ProductPromotions = append(newPromotion.ProductPromotions, entity.ProductPromotion{
			ProductId:          product.ID,
			MinDiscountedPrice: minDiscountPrice,
			MaxDiscountedPrice: maxDiscountPrice,
		})
	}

	promotion, err := u.promotionRepository.CreateNewPromotion(*newPromotion)
	if err != nil {
		return nil, err
	}

	promotionDTO := dto.PromotionResDTO{
		ID:                    promotion.ID,
		PromotionType:         promotion.PromotionType.Name,
		Title:                 promotion.Title,
		MaxDiscountedQuantity: promotion.MaxDiscountedQty,
		Quota:                 promotion.Quantity,
		UsedQuota:             promotion.Quantity - promotion.Quota,
		StartDate:             promotion.StartAt,
		EndDate:               promotion.EndAt,
	}
	if promotion.PromotionType.ID == dto.NOMINAL_PROMOTION_ID {
		promotionDTO.DiscountNominal = promotion.Nominal
	}
	if promotion.PromotionType.ID == dto.PERCENTAGE_PROMOTION_ID {
		promotionDTO.DiscountPercentage = promotion.Nominal
	}

	return &promotionDTO, nil
}

func (u *promotionUsecaseImpl) calculateMinMaxDiscountPrice(product *entity.Product, promotion *entity.Promotion) (float64, float64) {
	var minDiscountPrice, maxDiscountPrice float64
	if promotion.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
		minDiscountPrice = product.MinRealPrice - promotion.Nominal
		maxDiscountPrice = product.MaxRealPrice - promotion.Nominal
	}
	if promotion.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
		minDiscountPrice = product.MinRealPrice - (product.MinRealPrice * promotion.Nominal / 100)
		maxDiscountPrice = product.MaxRealPrice - (product.MaxRealPrice * promotion.Nominal / 100)
	}

	if minDiscountPrice <= 100 {
		minDiscountPrice = 100
	}
	if maxDiscountPrice <= 100 {
		maxDiscountPrice = 100
	}
	return minDiscountPrice, maxDiscountPrice
}

func (u *promotionUsecaseImpl) GetPromotionByID(id uint) (*dto.PromotionDetailResDTO, error) {
	promotion, err := u.promotionRepository.GetPromotionByID((int64(id)))
	if err != nil {
		return nil, err
	}

	promotionDTO := dto.PromotionDetailResDTO{
		ID:                    promotion.ID,
		PromotionTypeId:       promotion.PromotionTypeId,
		Nominal:               promotion.Nominal,
		Title:                 promotion.Title,
		MaxDiscountedQuantity: promotion.MaxDiscountedQty,
		Quota:                 promotion.Quantity,
		UsedQuota:             promotion.Quantity - promotion.Quota,
		StartDate:             promotion.StartAt,
		EndDate:               promotion.EndAt,
	}
	for _, productPromotion := range promotion.ProductPromotions {
		promotionDTO.ProductIds = append(promotionDTO.ProductIds, productPromotion.Product.ID)
		promotionDTO.Products = append(promotionDTO.Products, dto.ProductSellerResDTO{
			ID:           productPromotion.Product.ID,
			Title:        productPromotion.Product.Title,
			Slug:         productPromotion.Product.Slug,
			NumOfSale:    productPromotion.Product.ProductAnalytic.NumOfSale,
			TotalStock:   productPromotion.Product.ProductAnalytic.TotalStock,
			MinRealPrice: productPromotion.Product.MinRealPrice,
			MaxRealPrice: productPromotion.Product.MaxRealPrice,
		})
		if productPromotion.Product.ProductImages != nil && len(productPromotion.Product.ProductImages) > 0 {
			promotionDTO.Products[len(promotionDTO.Products)-1].ThumbnailImg = productPromotion.Product.ProductImages[0].ImageUrl
		}
	}

	return &promotionDTO, nil
}

func (u *promotionUsecaseImpl) UpdatePromotion(username string, promotionId int, req dto.UpsertPromotionReqDTO) (*dto.PromotionResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	promotion, err := u.promotionRepository.GetPromotionByID(int64(promotionId))
	if err != nil {
		return nil, err
	}

	var productIds []uint
	for _, productPromotion := range promotion.ProductPromotions {
		productIds = append(productIds, productPromotion.ProductId)
	}

	productIdsDiff := util.CompareArrayDifference(req.ProductIds, productIds)
	if len(productIdsDiff) > 0 {
		ok, err := u.promotionRepository.CheckProductPromotionOngoing(productIdsDiff, req.StartDate)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrCheckProductPromotionOngoing
		}
	}

	if promotion.MerchantId != merchant.ID {
		return nil, domain.ErrForbiddenMerchant
	}

	if promotion.EndAt.Before(time.Now()) {
		return nil, domain.ErrPromotionAlreadyEnded
	}

	if req.PromotionTypeId != dto.NOMINAL_PROMOTION_ID && req.PromotionTypeId != dto.PERCENTAGE_PROMOTION_ID {
		return nil, domain.ErrInvalidPromotionType
	}

	if req.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
		if req.Nominal <= 100 {
			return nil, domain.ErrInvalidNominal
		}
	}
	if req.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
		if req.Nominal <= 1 || req.Nominal > 100 {
			return nil, domain.ErrInvalidPercentage
		}
	}

	promotion.PromotionTypeId = req.PromotionTypeId
	promotion.Title = req.Title
	promotion.Nominal = req.Nominal
	promotion.MaxDiscountedQty = req.MaxDiscountedQuantity
	promotion.Quota = req.Quota
	promotion.Quantity = req.Quota
	promotion.StartAt = req.StartDate
	promotion.EndAt = req.EndDate

	var productPromotions []entity.ProductPromotion
	for _, productId := range req.ProductIds {
		product, err := u.productRepository.GetProductByProductId(productId)
		if err != nil {
			return nil, err
		}
		if product.MerchantDomain != merchant.Domain {
			return nil, domain.ErrInvalidProduct
		}

		minDiscountPrice, maxDiscountPrice := u.calculateMinMaxDiscountPrice(product, promotion)
		productPromotions = append(productPromotions, entity.ProductPromotion{
			PromotionId:        promotion.ID,
			ProductId:          product.ID,
			MinDiscountedPrice: minDiscountPrice,
			MaxDiscountedPrice: maxDiscountPrice,
		})
	}

	updatedPromotion, err := u.promotionRepository.UpdatePromotion(*promotion, productPromotions)
	if err != nil {
		return nil, err
	}

	promotionDTO := dto.PromotionResDTO{
		ID:                    updatedPromotion.ID,
		PromotionType:         updatedPromotion.PromotionType.Name,
		Title:                 updatedPromotion.Title,
		MaxDiscountedQuantity: updatedPromotion.MaxDiscountedQty,
		Quota:                 updatedPromotion.Quantity,
		UsedQuota:             updatedPromotion.Quantity - updatedPromotion.Quota,
		StartDate:             updatedPromotion.StartAt,
		EndDate:               updatedPromotion.EndAt,
	}
	if updatedPromotion.PromotionTypeId == dto.NOMINAL_PROMOTION_ID {
		promotionDTO.DiscountNominal = updatedPromotion.Nominal
	}
	if updatedPromotion.PromotionTypeId == dto.PERCENTAGE_PROMOTION_ID {
		promotionDTO.DiscountPercentage = updatedPromotion.Nominal
	}

	return &promotionDTO, nil
}

func (u *promotionUsecaseImpl) DeletePromotion(promotionId int) (*dto.PromotionResDTO, error) {
	promotion, err := u.promotionRepository.GetPromotionByID(int64(promotionId))
	if err != nil {
		return nil, err
	}

	if promotion.StartAt.Before(time.Now()) {
		return nil, domain.ErrPromotionAlreadyStarted
	}

	promotion, err = u.promotionRepository.DeletePromotion(promotion)
	if err != nil {
		return nil, err
	}

	promotionDTO := dto.PromotionResDTO{
		ID:                    promotion.ID,
		PromotionType:         promotion.PromotionType.Name,
		Title:                 promotion.Title,
		MaxDiscountedQuantity: promotion.MaxDiscountedQty,
		Quota:                 promotion.Quantity,
		UsedQuota:             promotion.Quantity - promotion.Quota,
		StartDate:             promotion.StartAt,
		EndDate:               promotion.EndAt,
	}
	if promotion.PromotionType.ID == dto.NOMINAL_PROMOTION_ID {
		promotionDTO.DiscountNominal = promotion.Nominal
	}
	if promotion.PromotionType.ID == dto.PERCENTAGE_PROMOTION_ID {
		promotionDTO.DiscountPercentage = promotion.Nominal
	}

	return &promotionDTO, nil
}
