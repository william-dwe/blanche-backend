package repository

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProductReviewRepository interface {
	GetProductReviewByTransactionId(transactionId uint) ([]entity.ProductReview, error)
	AddProductReview(productReview entity.ProductReview, merchantId uint) (*entity.ProductReview, error)

	GetProductReviewByProductSlug(productSlug string, reqParam dto.ProductReviewReqParamDTO) ([]entity.ProductReview, int64, error)
}

type ProductReviewRepositoryConfig struct {
	DB                        *gorm.DB
	ProductAnalyticRepository ProductAnalyticRepository
	MerchantRepository        MerchantRepository
}

type productReviewRepositoryImpl struct {
	db                        *gorm.DB
	productAnalyticRepository ProductAnalyticRepository
	merchantRepository        MerchantRepository
}

func NewProductReviewRepository(c ProductReviewRepositoryConfig) ProductReviewRepository {
	return &productReviewRepositoryImpl{
		db:                        c.DB,
		productAnalyticRepository: c.ProductAnalyticRepository,
		merchantRepository:        c.MerchantRepository,
	}
}

func (r *productReviewRepositoryImpl) GetProductReviewByTransactionId(transactionId uint) ([]entity.ProductReview, error) {
	var productReviews []entity.ProductReview

	err := r.db.
		Where("transaction_id = ?", transactionId).
		Find(&productReviews).Error
	if err != nil {
		return nil, err
	}

	return productReviews, err
}

func (r *productReviewRepositoryImpl) AddProductReview(productReview entity.ProductReview, merchantId uint) (resReview *entity.ProductReview, addReviewErr error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in AddProductReview repo: %v", r)
			addReviewErr = domain.ErrAddProductReview
		}
	}()

	err := tx.Create(&productReview).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"product_reviews_rating_check":                     domain.ErrAddProductReviewRatingNotValid,
				"product_reviews_transaction_id_fkey":              domain.ErrAddProductReviewTransactionIdNotValid,
				"product_reviews_product_id_fkey":                  domain.ErrAddProductReviewProductIdNotValid,
				"product_reviews_variant_item_id_fkey":             domain.ErrAddProductReviewVariantItemIdNotValid,
				"unique_product_reviews_trx_id_prod_id_variant_id": domain.ErrAddProductReviewDuplicate,
			},
			domain.ErrAddProductReview,
		)
		return nil, maskedErr
	}

	err = r.productAnalyticRepository.UpdateAvgRatingAndNumReviewProductIdTx(tx, productReview.ProductID, float64(productReview.Rating))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//update merchant rating review
	err = r.merchantRepository.UpdateMerchantRatingAndNumOfReviewTx(tx, merchantId, float64(productReview.Rating))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, domain.ErrAddProductReview
	}

	return &productReview, nil
}

func (r *productReviewRepositoryImpl) GetProductReviewByProductSlug(productSlug string, reqParam dto.ProductReviewReqParamDTO) ([]entity.ProductReview, int64, error) {
	qRating := ""
	if reqParam.Rating != 0 {
		qRating = fmt.Sprintf("rating = %d", reqParam.Rating)
	}

	qWithImage := ""
	if reqParam.FilterBy == dto.FilterByProductReviewWithImage {
		qWithImage = "image_url != '' or image_url is not null"
	}

	qWithComment := ""
	if reqParam.FilterBy == dto.FilterByProductReviewWithComment {
		qWithComment = "description != '' or description is not null"
	}

	var productReviews []entity.ProductReview

	PageOffset := reqParam.Limit * (reqParam.Page - 1)
	var countData int64

	err := r.db.Unscoped().
		Joins("LEFT JOIN products ON products.id = product_reviews.product_id").
		Where("products.slug = ?", productSlug).
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Transaction.User.UserDetail").
		Preload("Product").
		Preload("VariantItem").
		Preload("VariantItem.VariantSpecs").
		Where(qRating).
		Where(qWithImage).
		Where(qWithComment).
		Order("created_at desc").
		Offset(PageOffset).
		Limit(reqParam.Limit).
		Find(&productReviews).
		Limit(-1).
		Offset(-1).
		Count(&countData).
		Error
	if err != nil {
		return nil, 0, domain.ErrGetProductReviewByProductSlug
	}

	return productReviews, countData, nil
}
