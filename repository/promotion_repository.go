package repository

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type PromotionRepository interface {
	GetAllPromotions(req dto.PromotionListReqParamDTO, merchantID uint) ([]entity.Promotion, int64, error)
	GetPromotionByID(id int64) (*entity.Promotion, error)
	CreateNewPromotion(promotion entity.Promotion) (*entity.Promotion, error)
	CheckProductPromotionOngoing(productIds []uint, startDate time.Time) (bool, error)
	UpdatePromotion(promotion entity.Promotion, productPromotions []entity.ProductPromotion) (*entity.Promotion, error)
	DeletePromotion(promotion *entity.Promotion) (*entity.Promotion, error)
}

type promotionRepository struct {
	db *gorm.DB
}

type PromotionRepositoryConfig struct {
	DB *gorm.DB
}

func NewPromotionRepository(config PromotionRepositoryConfig) PromotionRepository {
	return &promotionRepository{
		db: config.DB,
	}
}

func (r *promotionRepository) GetAllPromotions(req dto.PromotionListReqParamDTO, merchantID uint) ([]entity.Promotion, int64, error) {
	var promotions []entity.Promotion
	var total int64
	pageOffset := req.Limit * (req.Page - 1)

	baseQuery := r.db.Preload("PromotionType").
		Preload("ProductPromotions.Product.ProductAnalytic").
		Preload("ProductPromotions.Product.ProductImages").
		Where("merchant_id = ?", merchantID).
		Limit(req.Limit).
		Offset(pageOffset)

	if req.Status == dto.VoucherStatusOngoing {
		baseQuery = baseQuery.Where("start_at <= now() AND end_at >= now()")
	}
	if req.Status == dto.VoucherStatusExpired {
		baseQuery = baseQuery.Where("end_at < now()")
	}
	if req.Status == dto.VoucherStatusIncoming {
		baseQuery = baseQuery.Where("start_at > now()")
	}

	err := baseQuery.Find(&promotions).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		return promotions, total, domain.ErrGetPromotionList
	}
	return promotions, total, nil
}

func (r *promotionRepository) GetPromotionByID(id int64) (*entity.Promotion, error) {
	var promotion entity.Promotion
	err := r.db.Preload("PromotionType").
		Preload("ProductPromotions.Product.ProductAnalytic").
		Preload("ProductPromotions.Product.ProductImages").
		Where("id = ?", id).
		First(&promotion).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPromotionNotFound
		}

		return nil, domain.ErrGetPromotionByID
	}
	return &promotion, nil
}

func (r *promotionRepository) CreateNewPromotion(promotion entity.Promotion) (*entity.Promotion, error) {
	err := r.db.Create(&promotion).Error
	if err != nil {
		return nil, domain.ErrCreateNewPromotion
	}
	return &promotion, nil
}

func (r *promotionRepository) CheckProductPromotionOngoing(productIds []uint, startDate time.Time) (bool, error) {
	var productPromotions []entity.ProductPromotion
	err := r.db.Where("product_id IN (?)", productIds).
		Joins("JOIN promotions ON product_promotions.promotion_id = promotions.id").
		Where("start_at <= ? OR end_at >= now()", startDate).
		Find(&productPromotions).Error
	if err != nil {
		return false, domain.ErrCheckProductPromotion
	}
	if len(productPromotions) > 0 {
		return false, nil
	}
	return true, nil
}

func (r *promotionRepository) UpdatePromotion(promotion entity.Promotion, productPromotions []entity.ProductPromotion) (*entity.Promotion, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("promotion_id = ?", promotion.ID).Delete(&entity.ProductPromotion{}).Error
		if err != nil {
			return domain.ErrUpdatePromotion
		}

		for _, productPromotion := range productPromotions {
			err := tx.Create(&productPromotion).Error
			if err != nil {
				return domain.ErrUpdatePromotion
			}
		}

		err = tx.Save(&promotion).Error
		if err != nil {
			return domain.ErrUpdatePromotion
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &promotion, nil
}
func (r *promotionRepository) DeletePromotion(promotion *entity.Promotion) (*entity.Promotion, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("promotion_id = ?", promotion.ID).Delete(&entity.ProductPromotion{}).Error
		if err != nil {
			return domain.ErrDeletePromotion
		}

		err = tx.Delete(&promotion).Error
		if err != nil {
			return domain.ErrDeletePromotion
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return promotion, nil
}
