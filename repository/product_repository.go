package repository

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductList(req dto.ProductListReqParamDTO) ([]entity.Product, int64, error)
	GetRecommendationProductList(req dto.PaginationRequest) ([]entity.Product, int64, error)
	GetProductBySlug(slug string) (*entity.Product, error)
	GetProductVariantDetailByProductID(productId uint) (*entity.Product, error)
	GetProductsByProductIds(productIdList []uint) ([]entity.Product, error)
	GetProductByProductId(productId uint) (*entity.Product, error)
	GetProductDetailBySlug(slug string) (*entity.Product, error)
	GetProductPromotionByProductId(productId uint) (*entity.ProductPromotion, error)
	GetProductMerchantById(productId uint) (*entity.Product, error)
	DecreaseProductStockTx(tx *gorm.DB, productId uint, variantItemId uint, quantity uint) error
	IncreaseProductStockTx(tx *gorm.DB, productId uint, variantItemId uint, quantity uint) error
	ChangeNumOfPendingSaleTx(tx *gorm.DB, productId uint, delta int) error
	IncreaseNumOfSaleTx(tx *gorm.DB, productId uint, delta int) error
	DecreaseProductPromotionTx(tx *gorm.DB, productId uint, quantity uint) error
	IncreaseProductPromotionTx(tx *gorm.DB, productId uint, quantity uint, transactionTime time.Time) error

	CheckMerchantProductName(merchantDomain, name string) (bool, error)
	CreateProduct(product *entity.Product, req dto.CreateProductReqDTO) (*entity.Product, error)
	UpdateMerchantProduct(product *entity.Product, req dto.CreateProductReqDTO) (*entity.Product, error)
	UpdateMerchantProductStatus(productIdIntList []uint, isArchived bool) error
	DeleteMerchantProduct(merchantDomain string, productId uint) error
}

type ProductRepositoryConfig struct {
	DB *gorm.DB
}

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(c ProductRepositoryConfig) ProductRepository {
	return &productRepositoryImpl{
		db: c.DB,
	}
}

func (r *productRepositoryImpl) getBaseProductQuery() *gorm.DB {
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

	return r.db.Unscoped().Table("products as prod, (?) as cte, product_analytics as ProductAnalytic", subQueryDiscountedProduct).
		Preload("ProductImages").
		Preload("ProductAnalytic").
		Preload("Merchant.City")
}

func (r *productRepositoryImpl) GetProductList(req dto.ProductListReqParamDTO) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64
	PageOffset := req.Pagination.Limit * (req.Pagination.Page - 1)
	if req.CategoryId == 0 && req.MerchantDomain == "" && req.Search == "" {
		return products, total, nil
	}

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

	fromTable := "products as prod, (?) as cte, product_analytics as ProductAnalytic"

	query := r.db.Unscoped().Table(fromTable, subQueryDiscountedProduct).
		Preload("ProductImages").
		Preload("ProductAnalytic").
		Preload("Merchant.City").
		Select("prod.*, min_discounted_price, max_discounted_price").
		Limit(req.Pagination.Limit).
		Offset(PageOffset).
		Where("prod.deleted_at is null").
		Where("prod.id = cte.id").
		Where("prod.product_analytic_id = ProductAnalytic.id").
		Where("min_discount_price >= ?", req.MinPrice).
		Where("max_discount_price <= ?", req.MaxPrice).
		Where("avg_rating >= ?", req.MinRating)

	if req.Search != "" {
		query = query.Where("documentq_weights @@ plainto_tsquery(?)", req.Search)
	}
	if !req.IsMerchant {
		query = query.Where("is_archived = ?", false).Order(req.SortBy + " " + req.SortDir)
	}
	if req.IsMerchant {
		query = query.Order("created_at desc")
	}

	if req.CategoryId != 0 {
		fromTable += ", categories"
		query = query.Table(fromTable, subQueryDiscountedProduct).
			Where("categories.id = prod.category_id").
			Where("grandparent_id = ? or parent_id = ? or category_id = ?", req.CategoryId, req.CategoryId, req.CategoryId)
	}
	if req.MerchantDomain != "" {
		query = query.Where("merchant_domain = ?", req.MerchantDomain)
	}
	if req.SellerCityIdList != nil {
		fromTable += ", merchants"
		query = query.Table(fromTable, subQueryDiscountedProduct).Where("merchants.domain = prod.merchant_domain").Where("city_id IN ?", req.SellerCityIdList)
	}

	res := query.Find(&products).Limit(-1).Offset(-1).Count(&total)
	if res.Error != nil {
		return nil, total, domain.ErrGetProducts
	}
	return products, total, nil
}

func (r *productRepositoryImpl) GetRecommendationProductList(req dto.PaginationRequest) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64
	PageOffset := req.Limit * (req.Page - 1)
	RecommendationLimit := int64(150)

	baseQuery := r.getBaseProductQuery().
		Select("prod.*, min_discounted_price, max_discounted_price, ProductAnalytic.avg_rating").
		Limit(req.Limit).
		Offset(PageOffset).
		Where("prod.deleted_at is null").
		Where("prod.id = cte.id").
		Where("prod.product_analytic_id = ProductAnalytic.id").
		Where("is_archived = ?", false).
		Order("avg_rating desc").
		Preload("ProductAnalytic").
		Preload("ProductImages").
		Preload("Merchant.City")

	res := baseQuery.Find(&products).Limit(-1).Offset(-1).Count(&total)
	if res.Error != nil {
		return nil, total, res.Error
	}
	if total > RecommendationLimit {
		total = RecommendationLimit
	}

	return products, total, nil
}

func (r *productRepositoryImpl) GetProductBySlug(slug string) (*entity.Product, error) {
	var product entity.Product
	res := r.db.Where("slug = ?", slug).First(&product)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductSlugNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductsByProductIds(productIdList []uint) ([]entity.Product, error) {
	var products []entity.Product
	res := r.db.
		Preload("Merchant").
		Preload("ProductAnalytic").
		Preload("ProductImages").
		Preload("Category").
		Preload("ProductPromotion").
		Preload("ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now()).Where("quota > 0")
		}).
		Where("id IN ?", productIdList).Find(&products)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductIdNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return products, nil
}

func (r *productRepositoryImpl) GetProductByProductId(productId uint) (*entity.Product, error) {
	var product entity.Product
	res := r.db.
		Preload("Merchant").
		Preload("ProductAnalytic").
		Preload("ProductImages").
		Preload("Category").
		Preload("ProductPromotion").
		Preload("ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now()).Where("quota > 0")
		}).
		Where("id = ?", productId).First(&product)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductIdNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductMerchantById(productId uint) (*entity.Product, error) {
	var product entity.Product
	res := r.db.Preload("Merchant").
		Where("id = ?", productId).First(&product)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductSlugNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductDetailBySlug(slug string) (*entity.Product, error) {
	var product entity.Product
	res := r.db.Debug().
		Preload("Merchant").
		Preload("ProductAnalytic").
		Preload("ProductImages").
		Preload("Category").
		Preload("ProductPromotion").
		Preload("ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now()).Where("quota > 0")
		}).
		Where("slug = ?", slug).First(&product)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductDetailSlugNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return &product, nil
}

func (r *productRepositoryImpl) GetProductPromotionByProductId(productId uint) (*entity.ProductPromotion, error) {
	var productPromotion entity.ProductPromotion
	res := r.db.
		Joins("Promotion").
		Where("product_id = ?", productId).
		Where("end_at > now()").
		Where("start_at <= now()").
		Where("quota > 0").
		First(&productPromotion)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductPromotionNotFound
		}
		return nil, domain.ErrGetProductPromotion
	}
	return &productPromotion, nil
}

func (r *productRepositoryImpl) DecreaseProductStockTx(tx *gorm.DB, productId uint, variantItemId uint, quantity uint) error {
	var product entity.Product
	err := tx.
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	err = tx.Model(&entity.ProductAnalytic{}).
		Where("id = ?", product.ProductAnalyticID).
		Update("total_stock", gorm.Expr("total_stock - ?", quantity)).
		Error
	if err != nil {
		return err
	}

	err = tx.
		Model(&entity.VariantItem{}).
		Where("id = ?", variantItemId).
		Where("product_id = ?", productId).
		Update("stock", gorm.Expr("stock - ?", quantity)).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) ChangeNumOfPendingSaleTx(tx *gorm.DB, productId uint, delta int) error {
	var product entity.Product
	err := tx.
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	err = tx.Model(&entity.ProductAnalytic{}).
		Where("id = ?", product.ProductAnalyticID).
		Update("num_of_pending_sale", gorm.Expr("num_of_pending_sale + ?", delta)).
		Error
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) IncreaseProductStockTx(tx *gorm.DB, productId uint, variantItemId uint, quantity uint) error {
	var product entity.Product
	err := tx.
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	err = tx.Model(&entity.ProductAnalytic{}).
		Where("id = ?", product.ProductAnalyticID).
		Update("total_stock", gorm.Expr("total_stock + ?", quantity)).
		Error
	if err != nil {
		return err
	}

	err = tx.
		Model(&entity.VariantItem{}).
		Where("id = ?", variantItemId).
		Where("product_id = ?", productId).
		Update("stock", gorm.Expr("stock + ?", quantity)).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) DecreaseProductPromotionTx(tx *gorm.DB, productId uint, quantity uint) error {
	var product entity.Product
	err := tx.
		Preload("ProductPromotion").
		Preload("ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", time.Now()).Where("end_at >= ?", time.Now()).Where("quota > 0")
		}).
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	if product.ProductPromotion != nil {
		if int(quantity) > product.ProductPromotion.Promotion.Quota {
			quantity = uint(product.ProductPromotion.Promotion.Quota)
		}

		err = tx.Model(&entity.Promotion{}).
			Where("id = ?", product.ProductPromotion.PromotionId).
			Update("quota", gorm.Expr("quota - ?", quantity)).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepositoryImpl) IncreaseProductPromotionTx(tx *gorm.DB, productId uint, quantity uint, transactionTime time.Time) error {
	var product entity.Product
	err := tx.
		Preload("ProductPromotion").
		Preload("ProductPromotion.Promotion", func(db *gorm.DB) *gorm.DB {
			return db.Where("start_at <= ?", transactionTime).Where("end_at >= ?", transactionTime)
		}).
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	if product.ProductPromotion != nil {
		err = tx.Model(&entity.Promotion{}).
			Where("id = ?", product.ProductPromotion.PromotionId).
			Update("quota", gorm.Expr("quota + ?", quantity)).
			Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepositoryImpl) GetProductVariantDetailByProductID(productId uint) (*entity.Product, error) {
	var product entity.Product
	res := r.db.
		Preload("ProductAnalytic").
		Preload("ProductImages").
		Preload("VariantItems.VariantSpecs.VariationGroup").
		Where("id = ?", productId).
		First(&product)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProductDetailSlugNotFound
		}
		return nil, domain.ErrGetProduct
	}
	return &product, nil
}

func (r *productRepositoryImpl) CreateProduct(product *entity.Product, req dto.CreateProductReqDTO) (*entity.Product, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var variationGroups []entity.VariationGroup
		for _, option := range req.Variant.VariantOptions {
			variationGroup := entity.VariationGroup{
				Name: option.Name,
			}
			err := tx.Create(&variationGroup).Error
			if err != nil {
				return domain.ErrCreateVariantGroup
			}

			variationGroups = append(variationGroups, variationGroup)
		}

		var variantItemsIndex int
		var productTotalStock uint
		minPrice := float64(0)
		maxPrice := float64(0)
		if len(req.Variant.VariantItems) > 1 {
			for _, value := range req.Variant.VariantOptions[0].Type {
				if len(req.Variant.VariantOptions) > 1 {
					for _, type2 := range req.Variant.VariantOptions[1].Type {
						variant := entity.VariantItem{
							Price:    req.Variant.VariantItems[variantItemsIndex].Price,
							ImageUrl: req.Variant.VariantItems[variantItemsIndex].Image,
							Stock:    uint(req.Variant.VariantItems[variantItemsIndex].Stock),
							VariantSpecs: []entity.VariantSpec{
								{
									VariationGroup: variationGroups[0],
									VariationName:  value,
								},
								{
									VariationGroup: variationGroups[1],
									VariationName:  type2,
								},
							},
						}
						if variant.Price < minPrice || minPrice == 0 {
							minPrice = variant.Price
						}
						if variant.Price > maxPrice {
							maxPrice = variant.Price
						}
						product.VariantItems = append(product.VariantItems, variant)
						productTotalStock += variant.Stock
						variantItemsIndex++
					}
				}

				if len(req.Variant.VariantOptions) == 1 {
					variant := entity.VariantItem{
						Price:    req.Variant.VariantItems[variantItemsIndex].Price,
						ImageUrl: req.Variant.VariantItems[variantItemsIndex].Image,
						Stock:    uint(req.Variant.VariantItems[variantItemsIndex].Stock),
						VariantSpecs: []entity.VariantSpec{
							{
								VariationGroup: variationGroups[0],
								VariationName:  value,
							},
						},
					}
					if variant.Price < minPrice || minPrice == 0 {
						minPrice = variant.Price
					}
					if variant.Price > maxPrice {
						maxPrice = variant.Price
					}
					product.VariantItems = append(product.VariantItems, variant)
					productTotalStock += variant.Stock
					variantItemsIndex++
				}
			}
		}

		if len(req.Variant.VariantItems) == 1 {
			variant := entity.VariantItem{
				Price:    req.Variant.VariantItems[0].Price,
				ImageUrl: req.Variant.VariantItems[0].Image,
				Stock:    uint(req.Variant.VariantItems[0].Stock),
			}
			if variant.Price < minPrice || minPrice == 0 {
				minPrice = variant.Price
			}
			if variant.Price > maxPrice {
				maxPrice = variant.Price
			}
			product.VariantItems = append(product.VariantItems, variant)
			productTotalStock += variant.Stock
		}

		product.ProductAnalytic = entity.ProductAnalytic{
			TotalStock: int(productTotalStock),
		}
		product.MinRealPrice = minPrice
		product.MaxRealPrice = maxPrice

		err := tx.Create(&product).Error
		if err != nil {
			maskedErr := util.PgConsErrMasker(
				err,
				entity.ConstraintErrMaskerMap{
					"products_pkey":     domain.ErrProductAlreadyExist,
					"products_slug_key": domain.ErrProductAlreadyExist,
				},
				domain.ErrCreateProduct,
			)
			return maskedErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepositoryImpl) CheckMerchantProductName(merchantDomain, name string) (bool, error) {
	var product entity.Product
	err := r.db.Where("LOWER(title) = ? AND merchant_domain = ?", name, merchantDomain).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, domain.ErrCheckMerchantProductName
	}

	return false, nil
}

func (r *productRepositoryImpl) DeleteMerchantProduct(merchantDomain string, productId uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {

		var product entity.Product
		err := tx.Preload("ProductAnalytic").Preload("VariantItems.VariantSpecs.VariationGroup").Where("id = ? AND merchant_domain = ?", productId, merchantDomain).First(&product).Error
		if err != nil {
			return domain.ErrGetProduct
		}

		for _, variant := range product.VariantItems {
			for _, variantSpec := range variant.VariantSpecs {
				err = tx.Delete(&variantSpec).Error
				if err != nil {
					return domain.ErrDeleteVariantSpec
				}
				err = tx.Delete(&variantSpec.VariationGroup).Error
				if err != nil {
					return domain.ErrDeleteVariantGroup
				}
			}
			err = tx.Delete(&variant).Error
			if err != nil {
				return domain.ErrDeleteVariantItem
			}
		}

		err = tx.Where("product_id = ?", product.ID).Delete(&entity.ProductImage{}).Error
		if err != nil {
			return domain.ErrDeleteProductImages
		}

		err = tx.Delete(&product).Error
		if err != nil {
			return domain.ErrDeleteProduct
		}

		err = tx.Delete(&product.ProductAnalytic).Error
		if err != nil {
			return domain.ErrDeleteProductAnalytics
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) UpdateMerchantProduct(product *entity.Product, req dto.CreateProductReqDTO) (*entity.Product, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, variant := range product.VariantItems {
			for _, variantSpec := range variant.VariantSpecs {
				err := tx.Delete(&variantSpec).Error
				if err != nil {
					return domain.ErrDeleteVariantSpec
				}
				err = tx.Delete(&variantSpec.VariationGroup).Error
				if err != nil {
					return domain.ErrDeleteVariantGroup
				}
			}
			err := tx.Delete(&variant).Error
			if err != nil {
				return domain.ErrDeleteVariantItem
			}
		}

		err := tx.Where("product_id = ?", product.ID).Delete(&entity.ProductImage{}).Error
		if err != nil {
			return domain.ErrDeleteProductImages
		}

		for _, url := range req.Images {
			image := entity.ProductImage{ImageUrl: url}
			product.ProductImages = append(product.ProductImages, image)
		}

		var variationGroups []entity.VariationGroup
		for _, option := range req.Variant.VariantOptions {
			variationGroup := entity.VariationGroup{
				Name: option.Name,
			}
			err := tx.Create(&variationGroup).Error
			if err != nil {
				return err
			}

			variationGroups = append(variationGroups, variationGroup)
		}

		var variantItemsIndex int
		var productTotalStock uint
		minPrice := float64(0)
		maxPrice := float64(0)
		if len(req.Variant.VariantItems) > 1 {
			for _, value := range req.Variant.VariantOptions[0].Type {
				if len(req.Variant.VariantOptions) > 1 {
					for _, type2 := range req.Variant.VariantOptions[1].Type {
						variant := entity.VariantItem{
							Price:    req.Variant.VariantItems[variantItemsIndex].Price,
							ImageUrl: req.Variant.VariantItems[variantItemsIndex].Image,
							Stock:    uint(req.Variant.VariantItems[variantItemsIndex].Stock),
							VariantSpecs: []entity.VariantSpec{
								{
									VariationGroup: variationGroups[0],
									VariationName:  value,
								},
								{
									VariationGroup: variationGroups[1],
									VariationName:  type2,
								},
							},
						}
						if variant.Price < minPrice || minPrice == 0 {
							minPrice = variant.Price
						}
						if variant.Price > maxPrice {
							maxPrice = variant.Price
						}
						product.VariantItems = append(product.VariantItems, variant)
						productTotalStock += variant.Stock
						variantItemsIndex++
					}
				}

				if len(req.Variant.VariantOptions) == 1 {
					variant := entity.VariantItem{
						Price:    req.Variant.VariantItems[variantItemsIndex].Price,
						ImageUrl: req.Variant.VariantItems[variantItemsIndex].Image,
						Stock:    uint(req.Variant.VariantItems[variantItemsIndex].Stock),
						VariantSpecs: []entity.VariantSpec{
							{
								VariationGroup: variationGroups[0],
								VariationName:  value,
							},
						},
					}
					if variant.Price < minPrice || minPrice == 0 {
						minPrice = variant.Price
					}
					if variant.Price > maxPrice {
						maxPrice = variant.Price
					}
					product.VariantItems = append(product.VariantItems, variant)
					productTotalStock += variant.Stock
					variantItemsIndex++
				}
			}
		}

		if len(req.Variant.VariantItems) == 1 {
			variant := entity.VariantItem{
				Price:    req.Variant.VariantItems[0].Price,
				ImageUrl: req.Variant.VariantItems[0].Image,
				Stock:    uint(req.Variant.VariantItems[0].Stock),
			}
			if variant.Price < minPrice || minPrice == 0 {
				minPrice = variant.Price
			}
			if variant.Price > maxPrice {
				maxPrice = variant.Price
			}
			product.VariantItems = append(product.VariantItems, variant)
			productTotalStock += variant.Stock
		}

		product.MinRealPrice = minPrice
		product.MaxRealPrice = maxPrice
		err = tx.Model(&product.ProductAnalytic).Update("total_stock", productTotalStock).Error
		if err != nil {
			return err
		}

		err = tx.Save(&product).Error
		if err != nil {
			maskedErr := util.PgConsErrMasker(
				err,
				entity.ConstraintErrMaskerMap{
					"products_pkey":     domain.ErrProductAlreadyExist,
					"products_slug_key": domain.ErrProductAlreadyExist,
				},
				domain.ErrCreateProduct,
			)
			return maskedErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepositoryImpl) UpdateMerchantProductStatus(productIdIntList []uint, isArchived bool) error {
	err := r.db.Model(&entity.Product{}).Where("id IN (?)", productIdIntList).Update("is_archived", isArchived).Error
	if err != nil {
		return domain.ErrUpdateMerchantProductStatus
	}

	return nil
}

func (r *productRepositoryImpl) IncreaseNumOfSaleTx(tx *gorm.DB, productId uint, delta int) error {
	var product entity.Product
	err := tx.
		Where("id = ?", productId).
		First(&product).Error
	if err != nil {
		return err
	}

	err = tx.Model(&entity.ProductAnalytic{}).
		Where("id = ?", product.ProductAnalyticID).
		Update("num_of_sale", gorm.Expr("num_of_sale + ?", delta)).
		Error
	if err != nil {
		return err
	}

	return nil
}
