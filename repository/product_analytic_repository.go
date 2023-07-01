package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type ProductAnalyticRepository interface {
	UpdateFavoriteProductTx(txGorm *gorm.DB, analyticId uint, delta int) error
	UpdateAvgRatingAndNumReviewTx(txGorm *gorm.DB, analyticId uint, newRating float64) error
	UpdateAvgRatingAndNumReviewProductIdTx(txGorm *gorm.DB, productId uint, newRating float64) error
}

type ProductAnalyticRepositoryConfig struct {
	DB *gorm.DB
}

type ProductAnalyticRepositoryImpl struct {
	db *gorm.DB
}

func NewProductAnalyticRepository(c ProductAnalyticRepositoryConfig) ProductAnalyticRepository {
	return &ProductAnalyticRepositoryImpl{
		db: c.DB,
	}
}

func (p *ProductAnalyticRepositoryImpl) UpdateFavoriteProductTx(txGorm *gorm.DB, analyticId uint, delta int) error {
	err := txGorm.Model(&entity.ProductAnalytic{}).Where("id = ?", analyticId).Update("num_of_favorite", gorm.Expr("num_of_favorite + ?", delta)).Error

	if err != nil {
		return domain.ErrProductAnalyticUpdateFavoriteProduct
	}

	return nil
}

func (p *ProductAnalyticRepositoryImpl) UpdateAvgRatingAndNumReviewTx(txGorm *gorm.DB, analyticId uint, newRating float64) error {
	err := txGorm.
		Model(&entity.ProductAnalytic{}).
		Where("id = ?", analyticId).
		Updates(map[string]interface{}{"avg_rating": gorm.Expr("((avg_rating * num_of_review) + ?)/(num_of_review + 1)", newRating), "num_of_review": gorm.Expr("num_of_review + ?", 1)}).
		Error

	if err != nil {
		return domain.ErrProductAnalyticUpdateAvgRatingAndNumReview
	}

	return nil
}

func (p *ProductAnalyticRepositoryImpl) UpdateAvgRatingAndNumReviewProductIdTx(txGorm *gorm.DB, productId uint, newRating float64) error {
	prodQuery := txGorm.Select("product_analytic_id").Model(&entity.Product{}).Where("id = ?", productId)

	err := txGorm.
		Model(&entity.ProductAnalytic{}).
		Where("id = (?)", prodQuery).
		Updates(map[string]interface{}{"avg_rating": gorm.Expr("((avg_rating * num_of_review) + ?)/(num_of_review + 1)", newRating), "num_of_review": gorm.Expr("num_of_review + ?", 1)}).
		Error

	if err != nil {
		return domain.ErrProductAnalyticUpdateAvgRatingAndNumReview
	}

	return nil
}
