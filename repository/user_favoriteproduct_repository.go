package repository

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserFavoriteProductRepository interface {
	UpdateFavoriteProduct(user entity.User, product entity.Product, isFavorited bool) error
	GetFavoriteProducts(user entity.User, query dto.UserFavoriteProductReqParamDTO) ([]entity.Product, int64, error)
}

type UserFavoriteProductRepositoryConfig struct {
	DB                        *gorm.DB
	ProductAnalyticRepository ProductAnalyticRepository
}

type userFavoriteProductRepositoryImpl struct {
	db                        *gorm.DB
	productAnalyticRepository ProductAnalyticRepository
}

func NewUserFavoriteProductRepository(c UserFavoriteProductRepositoryConfig) UserFavoriteProductRepository {
	return &userFavoriteProductRepositoryImpl{
		db:                        c.DB,
		productAnalyticRepository: c.ProductAnalyticRepository,
	}
}

func (r *userFavoriteProductRepositoryImpl) UpdateFavoriteProduct(user entity.User, product entity.Product, isFavorited bool) (addFavProdErr error) {
	favoriteProduct := entity.UserFavoriteProduct{
		UserId:    user.ID,
		ProductId: product.ID,
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			addFavProdErr = domain.ErrUpdateFavoriteProduct
		}
	}()

	var delta int
	var err error
	if isFavorited {
		err = tx.Create(&favoriteProduct).Error
		delta = 1
	}

	if !isFavorited {
		res := tx.Unscoped().Where("user_id = ?", user.ID).Where("product_id = ?", product.ID).Delete(&favoriteProduct)
		err = res.Error
		delta = -1 * int(res.RowsAffected)
		if (err == nil) && (res.RowsAffected == 0) {
			tx.Rollback()
			return nil
		}
	}

	if err != nil {
		tx.Rollback()
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "user_favorite_products_product_id_fkey" {
				return domain.ErrUpdateFavoriteProductNotFound
			}
			if pgErr.ConstraintName == "user_favorite_products_user_id_fkey" {
				return domain.ErrUpdateFavoriteUserNotFound
			}
		}
		if pgErr.Code == "23505" {
			if pgErr.ConstraintName == "user_favorite_products_user_id_product_id_key" {
				return nil
			}
		}

		log.Error().Msgf("Error update favorite product: %v", pgErr.Message)
		return domain.ErrUpdateFavoriteProduct
	}

	err = r.productAnalyticRepository.UpdateFavoriteProductTx(tx, product.ProductAnalyticID, delta)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrUpdateFavoriteProduct
	}

	return nil
}

func (r *userFavoriteProductRepositoryImpl) GetFavoriteProducts(user entity.User, query dto.UserFavoriteProductReqParamDTO) ([]entity.Product, int64, error) {
	productIdQuery := ""
	productSearchQuery := ""
	if query.ProductId > 0 {
		productIdQuery = fmt.Sprintf("product_id = %d", query.ProductId)
	}
	if query.Search != "" {
		productSearchQuery = fmt.Sprintf("title ILIKE '%%%s%%'", query.Search)
	}

	queryDb := r.getBaseProductQuery().
		Select("prod.*, min_discounted_price, max_discounted_price").
		Where(productIdQuery).
		Where(productSearchQuery)

	var products []entity.Product
	var totalData int64
	err := queryDb.
		Limit(query.Pagination.Limit).
		Offset((query.Pagination.Page-1)*query.Pagination.Limit).
		Where("UserFavoriteProducts.product_id = prod.id").
		Where("UserFavoriteProducts.user_id = ?", user.ID).
		Where("prod.id = cte.id").
		Order("UserFavoriteProducts.updated_at DESC").
		Find(&products).
		Limit(-1).
		Offset(-1).
		Count(&totalData).
		Error

	if err != nil {
		pgErr := err.(*pgconn.PgError)

		log.Error().Msgf("Error get favorite products: %v", pgErr.Message)
		return nil, 0, domain.ErrGetFavoriteProducts
	}

	return products, totalData, nil
}

func (r *userFavoriteProductRepositoryImpl) getBaseProductQuery() *gorm.DB {
	subQueryValidPromotion := r.db.Raw(`
		select product_id, min_discounted_price, max_discounted_price
		from product_promotions pp
		join promotions p 
		on pp.promotion_id = p.id
		where p.end_at > now()
			and p.start_at <= now() 
			and quota > 0
			and pp.deleted_at is null
	`)

	subQueryDiscountedProduct := r.db.Raw(`
		select 
			p.id,
			coalesce(pp.min_discounted_price, p.min_real_price) as min_discounted_price,
			coalesce(pp.max_discounted_price, p.max_real_price) as max_discounted_price
		from (?) as pp
		right join products p
		on pp.product_id  = p.id
	`, subQueryValidPromotion)

	return r.db.Unscoped().Table("products as prod, (?) as cte, user_favorite_products as UserFavoriteProducts", subQueryDiscountedProduct).
		Preload("ProductImages").
		Preload("ProductAnalytic").
		Preload("Merchant.City")
}
