package handler

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/usecase"

type Handler struct {
	exampleUsecase                   usecase.ExampleUsecase
	userUsecase                      usecase.UserUsecase
	userFavoriteProductUsecase       usecase.UserFavoriteProductUsecase
	authUsecase                      usecase.AuthUsecase
	merchantUsecase                  usecase.MerchantUsecase
	walletUsecase                    usecase.WalletUsecase
	walletpayUsecase                 usecase.WalletpayUsecase
	productUsecase                   usecase.ProductUsecase
	productReviewUsecase             usecase.ProductReviewUsecase
	productVariantUsecase            usecase.ProductVariantUsecase
	cartItemUsecase                  usecase.CartItemUsecase
	slpAccountUsecase                usecase.SlpAccountUsecase
	mpVoucherUsecase                 usecase.MarketplaceVoucherUsecase
	transactionUsecase               usecase.TransactionUsecase
	transactionStatusUsecase         usecase.TransactionStatusUsecase
	transactionDeliveryStatusUsecase usecase.TransactionDeliveryStatusUsecase
	refundRequestUsecase             usecase.RefundRequestUsecase
	refundRequestMessageUsecase      usecase.RefundRequestMessageUsecase

	sealabspayUsecase           usecase.SealabspayUsecase
	addressUsecase              usecase.AddressUsecase
	categoryUsecase             usecase.CategoryUsecase
	deliveryUsecase             usecase.DeliveryUsecase
	orderItemUsecase            usecase.OrderItemUsecase
	paymentRecordUsecase        usecase.PaymentRecordUsecase
	paymentMethodUsecase        usecase.PaymentMethodUsecase
	mediaUsecase                usecase.MediaUsecase
	marketplaceAnalyticsUsecase usecase.MarketplaceAnalyticsUsecase
	merchantAnalyticsUsecase    usecase.MerchantAnalyticsUsecase
	promotionBannerUsecase      usecase.PromotionBannerUsecase
	promotionUsecase            usecase.PromotionUsecase
}

type HandlerConfig struct {
	ExampleUsecase                   usecase.ExampleUsecase
	UserUsecase                      usecase.UserUsecase
	UserFavoriteProductUsecase       usecase.UserFavoriteProductUsecase
	AuthUsecase                      usecase.AuthUsecase
	MerchantUsecase                  usecase.MerchantUsecase
	WalletUsecase                    usecase.WalletUsecase
	WalletpayUsecase                 usecase.WalletpayUsecase
	ProductUsecase                   usecase.ProductUsecase
	ProductVariantUsecase            usecase.ProductVariantUsecase
	ProductReviewUsecase             usecase.ProductReviewUsecase
	PaymentRecordUsecase             usecase.PaymentRecordUsecase
	PaymentMethodUsecase             usecase.PaymentMethodUsecase
	CategoryUsecase                  usecase.CategoryUsecase
	CartItemUsecase                  usecase.CartItemUsecase
	AddressUsecase                   usecase.AddressUsecase
	DeliveryUsecase                  usecase.DeliveryUsecase
	SlpAccountUsecase                usecase.SlpAccountUsecase
	SealabspayUsecase                usecase.SealabspayUsecase
	MpVoucherUsecase                 usecase.MarketplaceVoucherUsecase
	TransactionUsecase               usecase.TransactionUsecase
	TransactionStatusUsecase         usecase.TransactionStatusUsecase
	TransactionDeliveryStatusUsecase usecase.TransactionDeliveryStatusUsecase
	RefundRequestUsecase             usecase.RefundRequestUsecase
	RefundRequestMessageUsecase      usecase.RefundRequestMessageUsecase
	OrderItemUsecase                 usecase.OrderItemUsecase
	MediaUsecase                     usecase.MediaUsecase
	MarketplaceAnalyticsUsecase      usecase.MarketplaceAnalyticsUsecase
	MerchantAnalyticsUsecase         usecase.MerchantAnalyticsUsecase
	PromotionBannerUsecase           usecase.PromotionBannerUsecase
	PromotionUsecase                 usecase.PromotionUsecase
}

func New(c HandlerConfig) *Handler {
	return &Handler{
		exampleUsecase:                   c.ExampleUsecase,
		userUsecase:                      c.UserUsecase,
		userFavoriteProductUsecase:       c.UserFavoriteProductUsecase,
		authUsecase:                      c.AuthUsecase,
		merchantUsecase:                  c.MerchantUsecase,
		walletUsecase:                    c.WalletUsecase,
		walletpayUsecase:                 c.WalletpayUsecase,
		productUsecase:                   c.ProductUsecase,
		productVariantUsecase:            c.ProductVariantUsecase,
		productReviewUsecase:             c.ProductReviewUsecase,
		categoryUsecase:                  c.CategoryUsecase,
		cartItemUsecase:                  c.CartItemUsecase,
		addressUsecase:                   c.AddressUsecase,
		slpAccountUsecase:                c.SlpAccountUsecase,
		mpVoucherUsecase:                 c.MpVoucherUsecase,
		transactionUsecase:               c.TransactionUsecase,
		transactionStatusUsecase:         c.TransactionStatusUsecase,
		transactionDeliveryStatusUsecase: c.TransactionDeliveryStatusUsecase,
		refundRequestUsecase:             c.RefundRequestUsecase,
		refundRequestMessageUsecase:      c.RefundRequestMessageUsecase,
		paymentRecordUsecase:             c.PaymentRecordUsecase,
		paymentMethodUsecase:             c.PaymentMethodUsecase,
		deliveryUsecase:                  c.DeliveryUsecase,
		sealabspayUsecase:                c.SealabspayUsecase,
		orderItemUsecase:                 c.OrderItemUsecase,
		mediaUsecase:                     c.MediaUsecase,
		marketplaceAnalyticsUsecase:      c.MarketplaceAnalyticsUsecase,
		merchantAnalyticsUsecase:         c.MerchantAnalyticsUsecase,
		promotionBannerUsecase:           c.PromotionBannerUsecase,
		promotionUsecase:                 c.PromotionUsecase,
	}
}
