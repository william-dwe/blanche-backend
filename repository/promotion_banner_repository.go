package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type PromotionBannerRepository interface {
	GetPromotionBannerList(req dto.PaginationRequest) ([]entity.PromotionBanner, int64, error)
	GetPromotionBannerByID(id uint) (*entity.PromotionBanner, error)
	CreatePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error)
	UpdatePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error)
	DeletePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error)
}

type PromotionBannerRepositoryConfig struct {
	DB *gorm.DB
}

type promotionBannerRepositoryImpl struct {
	db *gorm.DB
}

func NewPromotionBannerRepository(c PromotionBannerRepositoryConfig) PromotionBannerRepository {
	return &promotionBannerRepositoryImpl{
		db: c.DB,
	}
}

func (r *promotionBannerRepositoryImpl) GetPromotionBannerList(req dto.PaginationRequest) ([]entity.PromotionBanner, int64, error) {
	var promotionBanners []entity.PromotionBanner
	var total int64
	pageOffset := req.Limit * (req.Page - 1)
	err := r.db.Model(&promotionBanners).
		Order("created_at desc").
		Limit(req.Limit).
		Offset(pageOffset).
		Find(&promotionBanners).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, total, domain.ErrPromotionBannersNotFound
		}

		return nil, total, domain.ErrGetPromotionBanners
	}

	return promotionBanners, total, nil
}

func (r *promotionBannerRepositoryImpl) GetPromotionBannerByID(id uint) (*entity.PromotionBanner, error) {
	var promotionBanner entity.PromotionBanner
	err := r.db.Model(&promotionBanner).Where("id = ?", id).First(&promotionBanner).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPromotionBannerNotFound
		}

		return nil, domain.ErrGetPromotionBanner
	}

	return &promotionBanner, nil
}

func (r *promotionBannerRepositoryImpl) CreatePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error) {
	err := r.db.Model(&promotionBanner).Create(&promotionBanner).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"promotion_banners_pkey": domain.ErrDuplicatePromotionBanner,
			},
			domain.ErrCreatePromotionBanner,
		)
		return nil, maskedErr
	}

	return &promotionBanner, nil
}

func (r *promotionBannerRepositoryImpl) UpdatePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error) {
	err := r.db.Model(&promotionBanner).Updates(&promotionBanner).Error
	if err != nil {
		return nil, domain.ErrUpdatePromotionBanner
	}

	return &promotionBanner, nil
}

func (r *promotionBannerRepositoryImpl) DeletePromotionBanner(promotionBanner entity.PromotionBanner) (*entity.PromotionBanner, error) {
	err := r.db.Model(&promotionBanner).Delete(&promotionBanner).Error
	if err != nil {
		return nil, domain.ErrDeletePromotionBanner
	}

	return &promotionBanner, nil
}
