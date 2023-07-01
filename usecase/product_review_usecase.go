package usecase

import (
	"encoding/json"
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
)

type ProductReviewUsecase interface {
	GetProductReviewByInvoiceCode(username string, invoiceCode string) ([]dto.ProductReviewDTO, error)
	AddProductReview(username string, req dto.ReviewProductFormReqDTO, invoiceCode string) (*dto.ProductReviewDTO, error)

	GetProductReviewByProductSlug(productSlug string, reqParam dto.ProductReviewReqParamDTO) (*dto.ProductReviewResDTO, error)
}

type ProductReviewUsecaseConfig struct {
	UserRepository          repository.UserRepository
	ProductReviewRepository repository.ProductReviewRepository
	TransactionRepository   repository.TransactionRepository
	MerchantRepository      repository.MerchantRepository
	GcsUploader             util.GCSUploader
}

type productReviewUsecaseImpl struct {
	userRepository          repository.UserRepository
	productReviewRepository repository.ProductReviewRepository
	transactionRepository   repository.TransactionRepository
	merchantRepository      repository.MerchantRepository
	gcsUploader             util.GCSUploader
}

func NewProductReviewUsecase(c ProductReviewUsecaseConfig) ProductReviewUsecase {
	return &productReviewUsecaseImpl{
		userRepository:          c.UserRepository,
		productReviewRepository: c.ProductReviewRepository,
		transactionRepository:   c.TransactionRepository,
		merchantRepository:      c.MerchantRepository,
		gcsUploader:             c.GcsUploader,
	}
}

func (u *productReviewUsecaseImpl) GetProductReviewByInvoiceCode(username string, invoiceCode string) ([]dto.ProductReviewDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	transaction, err := u.transactionRepository.GetTransactionDetailByInvoiceCode(user.ID, invoiceCode)
	if err != nil {
		return nil, err
	}

	var cartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItem
	}

	productReviewMap := make(map[uint]dto.ProductReviewDTO)
	for _, cartItem := range cartItems {
		productReviewMap[cartItem.ProductId] = dto.ProductReviewDTO{
			ProductId:          cartItem.ProductId,
			VariantItemId:      cartItem.ProductVariantId,
			ProductName:        cartItem.Name,
			ProductVariantName: cartItem.VariantName,
			ProductImgUrl:      cartItem.Image,
			ProductPrice:       uint(cartItem.DiscountPrice),
			Rating:             0,
			Description:        "",
		}
	}

	productReviews, err := u.productReviewRepository.GetProductReviewByTransactionId(transaction.ID)
	if err != nil {
		return nil, err
	}

	for _, productReview := range productReviews {
		prodReviewDTO, ok := productReviewMap[productReview.ProductID]
		if !ok {
			continue
		}

		prodReviewDTO.ReviewedAt = &productReview.CreatedAt
		prodReviewDTO.Description = productReview.Description
		prodReviewDTO.Rating = uint(productReview.Rating)
		prodReviewDTO.ImageUrl = productReview.ImageUrl

		productReviewMap[productReview.ProductID] = prodReviewDTO
	}

	return util.MapValues(productReviewMap), nil
}

func (u *productReviewUsecaseImpl) AddProductReview(username string, input dto.ReviewProductFormReqDTO, invoiceCode string) (*dto.ProductReviewDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	transaction, err := u.transactionRepository.GetTransactionDetailByInvoiceCode(user.ID, invoiceCode)
	if err != nil {
		return nil, err
	}

	if transaction.TransactionStatus.OnCompletedAt == nil {
		return nil, domain.ErrAddProductReviewTransactionNotCompleted
	}

	var cartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItem
	}

	foundProduct := false
	for _, cartItem := range cartItems {
		if cartItem.ProductId == input.ProductId && cartItem.ProductVariantId == input.VariantItemId {
			foundProduct = true
			break
		}
	}

	if !foundProduct {
		return nil, domain.ErrAddProductReviewProductNotRelatedToTransaction
	}

	newProductReview := entity.ProductReview{
		TransactionID: transaction.ID,
		ProductID:     input.ProductId,
		VariantItemID: input.VariantItemId,
		Rating:        int(input.Rating),
		Description:   input.Description,
	}

	if input.Image != nil {
		url, err := u.gcsUploader.UploadFileFromFileHeader(*input.Image,
			fmt.Sprintf("prod_review_%d_%d_%d_%d", user.ID, transaction.ID, input.ProductId, input.VariantItemId))
		if err != nil {
			log.Error().Msgf("Failed to upload file: %v", err)
			return nil, domain.ErrUploadFile
		}
		newProductReview.ImageUrl = &url
	}

	productReview, err := u.productReviewRepository.AddProductReview(newProductReview, transaction.Merchant.ID)
	if err != nil {
		return nil, err
	}

	return &dto.ProductReviewDTO{
		ProductId:     productReview.ProductID,
		VariantItemId: productReview.VariantItemID,
		ImageUrl:      productReview.ImageUrl,
		Description:   productReview.Description,
		Rating:        uint(productReview.Rating),
		ReviewedAt:    &productReview.CreatedAt,
	}, nil
}

func (u *productReviewUsecaseImpl) GetProductReviewByProductSlug(productSlug string, reqParam dto.ProductReviewReqParamDTO) (*dto.ProductReviewResDTO, error) {

	productReviews, totalData, err := u.productReviewRepository.GetProductReviewByProductSlug(productSlug, reqParam)
	if err != nil {
		return nil, err
	}

	productReviewDTOs := make([]dto.ProductReviewDTO, 0)
	for _, productReview := range productReviews {
		tmpProdReview := dto.ProductReviewDTO{
			Username:           productReview.Transaction.User.Username,
			UserProfilePicture: productReview.Transaction.User.UserDetail.ProfilePicture,
			ProductId:          productReview.ProductID,
			ProductName:        productReview.Product.Title,
			VariantItemId:      productReview.VariantItemID,
			ImageUrl:           productReview.ImageUrl,
			Description:        productReview.Description,
			Rating:             uint(productReview.Rating),
			ReviewedAt:         &productReview.CreatedAt,
		}

		if len(productReview.VariantItem.VariantSpecs) > 0 {
			tmpProdReview.ProductVariantName += productReview.VariantItem.VariantSpecs[0].VariationName
		}
		if len(productReview.VariantItem.VariantSpecs) > 1 {
			tmpProdReview.ProductVariantName += "," + productReview.VariantItem.VariantSpecs[1].VariationName
		}

		productReviewDTOs = append(productReviewDTOs, tmpProdReview)
	}

	return &dto.ProductReviewResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalPage:   (totalData + int64(reqParam.Limit) - 1) / int64(reqParam.Limit),
			TotalData:   totalData,
			CurrentPage: reqParam.Page,
		},
		Reviews: productReviewDTOs,
	}, nil
}
