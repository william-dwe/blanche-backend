package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type UserFavoriteProductUsecase interface {
	UpdateFavoriteProduct(dto.UserFavoriteProductReqDTO) (*dto.UserFavoriteProductResDTO, error)
	GetFavoriteProducts(username string, queryParam dto.UserFavoriteProductReqParamDTO) (*dto.ProductListResDTO, error)
}

type UserFavoriteProductUsecaseConfig struct {
	UserRepository                repository.UserRepository
	UserFavoriteProductRepository repository.UserFavoriteProductRepository
	ProductAnalyticRepository     repository.ProductAnalyticRepository
	ProductRepository             repository.ProductRepository
}

type userFavoriteProductUsecaseImpl struct {
	userRepository                repository.UserRepository
	userFavoriteProductRepository repository.UserFavoriteProductRepository
	productAnalyticRepository     repository.ProductAnalyticRepository
	productRepository             repository.ProductRepository
}

func NewUserFavoriteProductUsecase(c UserFavoriteProductUsecaseConfig) UserFavoriteProductUsecase {
	return &userFavoriteProductUsecaseImpl{
		userRepository:                c.UserRepository,
		userFavoriteProductRepository: c.UserFavoriteProductRepository,
		productAnalyticRepository:     c.ProductAnalyticRepository,
		productRepository:             c.ProductRepository,
	}
}

func (u *userFavoriteProductUsecaseImpl) UpdateFavoriteProduct(input dto.UserFavoriteProductReqDTO) (*dto.UserFavoriteProductResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(input.Username)
	if err != nil {
		return nil, err
	}

	product, err := u.productRepository.GetProductByProductId(input.ProductId)
	if err != nil {
		return nil, err
	}

	err = u.userFavoriteProductRepository.UpdateFavoriteProduct(*user, *product, *input.IsFavorited)
	if err != nil {
		return nil, err
	}

	return &dto.UserFavoriteProductResDTO{
		ProductId:   input.ProductId,
		IsFavorited: *input.IsFavorited,
	}, nil

}

func (u *userFavoriteProductUsecaseImpl) GetFavoriteProducts(username string, query dto.UserFavoriteProductReqParamDTO) (*dto.ProductListResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	favoriteProducts, totalRows, err := u.userFavoriteProductRepository.GetFavoriteProducts(*user, query)
	if err != nil {
		return nil, err
	}

	productsDTO := make([]dto.ProductResDTO, len(favoriteProducts))
	for i, product := range favoriteProducts {
		image := ""
		if len(product.ProductImages) > 0 {
			image = product.ProductImages[0].ImageUrl
		}
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
			ThumbnailImg:     image,
		}
	}

	pageLimit := int64(query.Pagination.Limit)
	productFavDTO := dto.ProductListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalRows,
			TotalPage:   (totalRows + pageLimit - 1) / pageLimit,
			CurrentPage: query.Pagination.Page,
		},
		Products: productsDTO,
	}

	return &productFavDTO, nil
}
