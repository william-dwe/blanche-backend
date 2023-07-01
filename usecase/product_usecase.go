package usecase

import (
	"fmt"
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type ProductUsecase interface {
	GetProductList(req dto.ProductListReqParamDTO) (*dto.ProductListResDTO, error)
	GetRecommendationProductList(req dto.PaginationRequest) (*dto.ProductListResDTO, error)
	GetMerchantProductList(username string, req dto.ProductListReqParamDTO) (*dto.ProductSellerListResDTO, error)
	GetProductDetailsBySlug(userJwt *dto.AccessTokenPayload, slug string) (*dto.ProductDetailResDTO, error)
	GetAdminProductDetailByProductID(userJwt *dto.AccessTokenPayload, productId uint) (*dto.ProductAdminDetailResDTO, error)
	UploadProductFiles(req dto.UploadImageReqDTO) (*dto.UploadImageResDTO, error)

	CreateProduct(username string, req dto.CreateProductReqDTO) (*dto.CreateProductResDTO, error)
	CheckMerchantProductName(username string, req dto.CheckMerchantProductNameReqDTO) (*dto.CreateProductCheckNameResDTO, error)
	UpdateMerchantProduct(username string, productIdInt uint, req dto.CreateProductReqDTO) (*dto.CreateProductResDTO, error)
	UpdateMerchantProductAvailability(username string, productIdIntList []uint, req dto.UpdateProductAvailabilityReqDTO) ([]dto.CreateProductResDTO, error)
	DeleteMerchantProduct(username string, productIdIntList []uint) ([]dto.CreateProductResDTO, error)
}

type ProductUsecaseConfig struct {
	ProductRepository  repository.ProductRepository
	CategoryRepository repository.CategoryRepository
	UserRepository     repository.UserRepository
	MerchantRepository repository.MerchantRepository
	MediaUsecase       MediaUsecase
}

type productUsecaseImpl struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
	userRepository     repository.UserRepository
	merchantRepository repository.MerchantRepository
	mediaUsecase       MediaUsecase
}

func NewProductUsecase(c ProductUsecaseConfig) ProductUsecase {
	return &productUsecaseImpl{
		productRepository:  c.ProductRepository,
		categoryRepository: c.CategoryRepository,
		userRepository:     c.UserRepository,
		merchantRepository: c.MerchantRepository,
		mediaUsecase:       c.MediaUsecase,
	}
}

func (u *productUsecaseImpl) GetProductList(req dto.ProductListReqParamDTO) (*dto.ProductListResDTO, error) {
	if req.MinPrice > req.MaxPrice {
		return nil, domain.ErrInvalidPriceRange
	}

	category, err := u.categoryRepository.GetCategoryBySlug(req.CategorySlug)
	if err == nil {
		req.CategoryId = category.ID
	}

	req.IsMerchant = false
	products, totalProducts, err := u.productRepository.GetProductList(req)
	if err != nil {
		return nil, err
	}

	productsDTO := make([]dto.ProductResDTO, len(products))
	for i, product := range products {
		productsDTO[i] = dto.ProductResDTO{
			ID:               product.ID,
			Title:            product.Title,
			Slug:             product.Slug,
			MinRealPrice:     product.MinRealPrice,
			MaxRealPrice:     product.MaxRealPrice,
			MinDiscountPrice: product.MinDiscountedPrice,
			MaxDiscountPrice: product.MaxDiscountedPrice,
			NumOfSale:        product.ProductAnalytic.NumOfSale,
			AvgRating:        product.ProductAnalytic.AvgRating,
			ThumbnailImg:     "http://dummyimage.com/173x122.png/ff4444/ffffff",
			SellerCity:       product.Merchant.City.Name,
		}

		if len(product.ProductImages) > 0 {
			productsDTO[i].ThumbnailImg = product.ProductImages[0].ImageUrl
		}
	}

	return &dto.ProductListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalProducts,
			TotalPage:   (totalProducts + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
			CurrentPage: req.Pagination.Page,
		},
		Products: productsDTO,
	}, nil
}

func (u *productUsecaseImpl) GetMerchantProductList(username string, req dto.ProductListReqParamDTO) (*dto.ProductSellerListResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	req.MerchantDomain = merchant.Domain
	req.IsMerchant = true

	products, totalProducts, err := u.productRepository.GetProductList(req)
	if err != nil {
		return nil, err
	}

	productsDTO := make([]dto.ProductSellerResDTO, len(products))
	for i, product := range products {
		productsDTO[i] = dto.ProductSellerResDTO{
			ID:           product.ID,
			Title:        product.Title,
			Slug:         product.Slug,
			MinRealPrice: product.MinRealPrice,
			MaxRealPrice: product.MaxRealPrice,
			NumOfSale:    product.ProductAnalytic.NumOfSale,
			AvgRating:    product.ProductAnalytic.AvgRating,
			ThumbnailImg: "http://dummyimage.com/173x122.png/ff4444/ffffff",
			TotalStock:   product.ProductAnalytic.TotalStock,
			IsArchived:   product.IsArchived,
		}

		if len(product.ProductImages) > 0 {
			productsDTO[i].ThumbnailImg = product.ProductImages[0].ImageUrl
		}
	}

	return &dto.ProductSellerListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalProducts,
			TotalPage:   (totalProducts + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
			CurrentPage: req.Pagination.Page,
		},
		Products: productsDTO,
	}, nil
}

func (u *productUsecaseImpl) GetRecommendationProductList(req dto.PaginationRequest) (*dto.ProductListResDTO, error) {
	products, totalProducts, err := u.productRepository.GetRecommendationProductList(req)
	if err != nil {
		return nil, err
	}
	totalPages := (totalProducts + int64(req.Limit) - 1) / int64(req.Limit)
	if int64(req.Page) > totalPages {
		return &dto.ProductListResDTO{
			Products: []dto.ProductResDTO{},
		}, nil
	}

	productsDTO := make([]dto.ProductResDTO, len(products))
	for i, product := range products {
		productsDTO[i] = dto.ProductResDTO{
			ID:               product.ID,
			Title:            product.Title,
			Slug:             product.Slug,
			MinRealPrice:     product.MinRealPrice,
			MaxRealPrice:     product.MaxRealPrice,
			MinDiscountPrice: product.MinDiscountedPrice,
			MaxDiscountPrice: product.MaxDiscountedPrice,
			NumOfSale:        product.ProductAnalytic.NumOfSale,
			AvgRating:        product.ProductAnalytic.AvgRating,
			ThumbnailImg:     "http://dummyimage.com/173x122.png/ff4444/ffffff",
			SellerCity:       product.Merchant.City.Name,
		}

		if len(product.ProductImages) > 0 {
			productsDTO[i].ThumbnailImg = product.ProductImages[0].ImageUrl
		}
	}

	return &dto.ProductListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalProducts,
			TotalPage:   totalPages,
			CurrentPage: req.Page,
		},
		Products: productsDTO,
	}, nil
}

func (u *productUsecaseImpl) GetAdminProductDetailByProductID(userJwt *dto.AccessTokenPayload, productId uint) (*dto.ProductAdminDetailResDTO, error) {
	product, err := u.productRepository.GetProductByProductId(productId)
	if err != nil {
		return nil, err
	}

	productImages := make([]string, 0)
	for _, image := range product.ProductImages {
		productImages = append(productImages, image.ImageUrl)
	}

	productDTO := dto.ProductAdminDetailResDTO{
		ID:          product.ID,
		Title:       product.Title,
		Price:       nil,
		Images:      productImages,
		Description: product.Description,
		IsUsed:      product.IsUsed,
		TotalStock:  product.ProductAnalytic.TotalStock,
		IsArchived:  product.IsArchived,
		Weight:      product.Weight,
		Dimension: dto.ProductDetailDimension{
			Height: product.Height,
			Width:  product.Width,
			Length: product.Length,
		},
	}

	categories := make([]uint, 3)
	categories[0] = product.Category.GrandparentId
	categories[1] = product.Category.ParentId
	categories[2] = product.Category.ID
	productDTO.Categories = categories

	if product.MinRealPrice == product.MaxRealPrice {
		productDTO.Price = &product.MinRealPrice
	}

	return &productDTO, nil
}

func (u *productUsecaseImpl) GetProductDetailsBySlug(userJwt *dto.AccessTokenPayload, slug string) (*dto.ProductDetailResDTO, error) {
	product, err := u.productRepository.GetProductDetailBySlug(slug)
	if err != nil {
		return nil, err
	}

	var userId uint
	if userJwt != nil {
		user, err := u.userRepository.GetUserByUsername(userJwt.Username)
		if err == nil {
			userId = user.ID
		}
	}

	productImages := make([]string, 0)
	for _, image := range product.ProductImages {
		productImages = append(productImages, image.ImageUrl)
	}

	minDiscountPrice, maxDiscountPrice := product.MinRealPrice, product.MaxRealPrice
	if product.ProductPromotion != nil {
		minDiscountPrice, maxDiscountPrice = product.ProductPromotion.MinDiscountedPrice, product.ProductPromotion.MaxDiscountedPrice
	}

	productDTO := dto.ProductDetailResDTO{
		ID:               product.ID,
		Title:            product.Title,
		MinRealPrice:     product.MinRealPrice,
		MaxRealPrice:     product.MaxRealPrice,
		MinDiscountPrice: minDiscountPrice,
		MaxDiscountPrice: maxDiscountPrice,
		Category: dto.ProductDetailCategory{
			Name: product.Category.Name,
			URL:  product.Category.Slug,
		},
		Images:         productImages,
		Description:    product.Description,
		IsUsed:         product.IsUsed,
		SKU:            product.SKU,
		FavouriteCount: product.ProductAnalytic.NumOfFavorite,
		UnitSold:       product.ProductAnalytic.NumOfSale,
		TotalStock:     product.ProductAnalytic.TotalStock,
		IsArchived:     product.IsArchived,

		Rating: dto.ProductDetailRating{
			AvgRating: product.ProductAnalytic.AvgRating,
			Count:     product.ProductAnalytic.NumOfReview,
		},

		Weight: product.Weight,
		Dimension: dto.ProductDetailDimension{
			Height: product.Height,
			Width:  product.Width,
			Length: product.Length,
		},

		IsMyProduct: product.Merchant.UserId == userId,
	}

	return &productDTO, nil
}

func (u *productUsecaseImpl) UploadProductFiles(req dto.UploadImageReqDTO) (*dto.UploadImageResDTO, error) {
	url, err := u.mediaUsecase.UploadFileForBinding(req.File, req.File.Filename)
	if err != nil {
		return nil, err
	}

	return &dto.UploadImageResDTO{
		ImageURLs: url,
	}, nil
}

func (u *productUsecaseImpl) CreateProduct(username string, req dto.CreateProductReqDTO) (*dto.CreateProductResDTO, error) {
	product := &entity.Product{
		CategoryId:  req.CategoryId,
		Title:       req.Title,
		IsArchived:  req.IsArchived,
		Description: req.Description,
		Weight:      req.Weight,
		Height:      req.Dimension.Height,
		Width:       req.Dimension.Width,
		Length:      req.Dimension.Length,
		IsUsed:      req.IsUsed,
	}

	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if merchant != nil {
		product.MerchantId = merchant.ID
		product.MerchantDomain = merchant.Domain
	}

	slugTitle := strings.Join(strings.Split(strings.ToLower(product.Title), " "), "-")
	slug := fmt.Sprintf("%s/%s", merchant.Domain, slugTitle)
	product.Slug = slug

	for _, url := range req.Images {
		image := entity.ProductImage{ImageUrl: url}
		product.ProductImages = append(product.ProductImages, image)
	}

	productRes, err := u.productRepository.CreateProduct(product, req)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProductResDTO{
		ID:             productRes.ID,
		Slug:           productRes.Slug,
		MerchantDomain: productRes.MerchantDomain,
		Title:          productRes.Title,
		MinRealPrice:   productRes.MinRealPrice,
		MaxRealPrice:   productRes.MaxRealPrice,
		Description:    productRes.Description,
		IsArchived:     productRes.IsArchived,
		IsUsed:         productRes.IsUsed,
		TotalStock:     productRes.ProductAnalytic.TotalStock,
		Weight:         productRes.Weight,
		Images:         req.Images,
		Dimension: dto.CreateProductDimensionReqDTO{
			Width:  productRes.Width,
			Length: productRes.Length,
			Height: productRes.Height,
		},
	}, nil
}

func (u *productUsecaseImpl) CheckMerchantProductName(username string, req dto.CheckMerchantProductNameReqDTO) (*dto.CreateProductCheckNameResDTO, error) {
	title := strings.ToLower(req.ProductName)
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return &dto.CreateProductCheckNameResDTO{
			ProductName: title,
			IsAvailable: false,
		}, err
	}

	isAvailable, err := u.productRepository.CheckMerchantProductName(merchant.Domain, title)
	if err != nil {
		return &dto.CreateProductCheckNameResDTO{
			ProductName: title,
			IsAvailable: false,
		}, err
	}

	return &dto.CreateProductCheckNameResDTO{
		ProductName: title,
		IsAvailable: isAvailable,
	}, nil
}

func (u *productUsecaseImpl) DeleteMerchantProduct(username string, productIdIntList []uint) ([]dto.CreateProductResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]dto.CreateProductResDTO, 0)
	for _, productIdInt := range productIdIntList {
		productRes, err := u.productRepository.GetProductVariantDetailByProductID(productIdInt)
		if err != nil {
			return nil, err
		}

		if merchant.Domain != productRes.MerchantDomain {
			return nil, domain.ErrUpdateProductUnauthorized
		}

		err = u.productRepository.DeleteMerchantProduct(merchant.Domain, productIdInt)
		if err != nil {
			return nil, err
		}

		productDTOs = append(productDTOs, dto.CreateProductResDTO{
			ID:             productRes.ID,
			Slug:           productRes.Slug,
			MerchantDomain: productRes.MerchantDomain,
			Title:          productRes.Title,
			MinRealPrice:   productRes.MinRealPrice,
			MaxRealPrice:   productRes.MaxRealPrice,
			Description:    productRes.Description,
			IsArchived:     productRes.IsArchived,
			IsUsed:         productRes.IsUsed,
			TotalStock:     productRes.ProductAnalytic.TotalStock,
			Weight:         productRes.Weight,
			Dimension: dto.CreateProductDimensionReqDTO{
				Width:  productRes.Width,
				Length: productRes.Length,
				Height: productRes.Height,
			},
		})
	}

	return productDTOs, nil
}

func (u *productUsecaseImpl) UpdateMerchantProduct(username string, productIdInt uint, req dto.CreateProductReqDTO) (*dto.CreateProductResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	productRes, err := u.productRepository.GetProductVariantDetailByProductID(productIdInt)
	if err != nil {
		return nil, err
	}

	if merchant.Domain != productRes.MerchantDomain {
		return nil, domain.ErrUpdateProductUnauthorized
	}

	productNew := productRes
	productNew.CategoryId = req.CategoryId
	if productRes.ProductAnalytic.NumOfSale == 0 {
		productNew.Title = req.Title
	}
	productNew.Title = productRes.Title
	productNew.IsArchived = req.IsArchived
	productNew.Description = req.Description
	productNew.Weight = req.Weight
	productNew.Height = req.Dimension.Height
	productNew.Width = req.Dimension.Width
	productNew.Length = req.Dimension.Length
	productNew.IsUsed = req.IsUsed

	if merchant != nil {
		productNew.MerchantId = merchant.ID
		productNew.MerchantDomain = merchant.Domain
	}

	slugTitle := strings.Join(strings.Split(strings.ToLower(productNew.Title), " "), "-")
	slug := fmt.Sprintf("%s/%s", merchant.Domain, slugTitle)
	productNew.Slug = slug

	productNew, err = u.productRepository.UpdateMerchantProduct(productNew, req)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProductResDTO{
		ID:             productNew.ID,
		Slug:           productNew.Slug,
		MerchantDomain: productNew.MerchantDomain,
		Title:          productNew.Title,
		MinRealPrice:   productNew.MinRealPrice,
		MaxRealPrice:   productNew.MaxRealPrice,
		Description:    productNew.Description,
		IsArchived:     productNew.IsArchived,
		IsUsed:         productNew.IsUsed,
		TotalStock:     productNew.ProductAnalytic.TotalStock,
		Weight:         productNew.Weight,
		Images:         req.Images,
		Dimension: dto.CreateProductDimensionReqDTO{
			Width:  productNew.Width,
			Length: productNew.Length,
			Height: productNew.Height,
		},
	}, nil
}

func (u *productUsecaseImpl) UpdateMerchantProductAvailability(username string, productIdIntList []uint, req dto.UpdateProductAvailabilityReqDTO) ([]dto.CreateProductResDTO, error) {
	err := u.productRepository.UpdateMerchantProductStatus(productIdIntList, req.IsArchived)
	if err != nil {
		return nil, err
	}

	updatedProducts, err := u.productRepository.GetProductsByProductIds(productIdIntList)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]dto.CreateProductResDTO, 0)
	for _, product := range updatedProducts {
		productDTOs = append(productDTOs, dto.CreateProductResDTO{
			ID:             product.ID,
			Slug:           product.Slug,
			MerchantDomain: product.MerchantDomain,
			Title:          product.Title,
			MinRealPrice:   product.MinRealPrice,
			MaxRealPrice:   product.MaxRealPrice,
			Description:    product.Description,
			IsArchived:     product.IsArchived,
			IsUsed:         product.IsUsed,
			TotalStock:     product.ProductAnalytic.TotalStock,
			Weight:         product.Weight,
			Dimension: dto.CreateProductDimensionReqDTO{
				Width:  product.Width,
				Length: product.Length,
				Height: product.Height,
			},
		})
	}

	return productDTOs, nil

}
