package server

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cache"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/db"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/usecase"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func initRouter() *gin.Engine {
	authUtil := util.NewAuthUtil()
	exampleRepo := repository.NewExampleRepository(repository.ExampleRepositoryConfig{
		DB:  db.Get(),
		RDB: cache.GetClientRDB(),
	})
	paymentRecordRepo := repository.NewPaymentRecordRepository(repository.PaymentRecordRepositoryConfig{
		DB: db.Get(),
	})
	userRepo := repository.NewUserRepository(repository.UserRepositoryConfig{
		DB: db.Get(),
	})
	authRepo := repository.NewAuthRepository(repository.AuthRepositoryConfig{
		DB:  db.Get(),
		RDB: cache.GetClientRDB(),
	})
	deliveryRepo := repository.NewDeliveryRepository(repository.DeliveryRepositoryConfig{
		DB: db.Get(),
	})
	merchantRepo := repository.NewMerchantRepository(repository.MerchantRepositoryConfig{
		DB:                 db.Get(),
		UserRepository:     userRepo,
		DeliveryRepository: deliveryRepo,
	})
	merchantHoldingAccountHistoryRepo := repository.NewMerchantHoldingAccountHistoryRepository(repository.MerchantHoldingAccountHistoryRepositoryConfig{
		DB: db.Get(),
	})
	walletRepo := repository.NewWalletRepository(repository.WalletRepositoryConfig{
		DB:                      db.Get(),
		PaymentRecordRepository: paymentRecordRepo,
	})
	merchantHoldingAccountRepo := repository.NewMerchantHoldingAccountRepository(repository.MerchantHoldingAccountRepositoryConfig{
		DB:               db.Get(),
		WalletRepository: walletRepo,
	})
	productRepo := repository.NewProductRepository(repository.ProductRepositoryConfig{
		DB: db.Get(),
	})
	productVariantRepo := repository.NewProductVariantRepository(repository.ProductVariantRepositoryConfig{
		DB: db.Get(),
	})
	categoryRepo := repository.NewCategoryRepository(repository.CategoryRepositoryConfig{
		DB: db.Get(),
	})
	cartItemRepo := repository.NewCartItemRepository(repository.CartItemRepositoryConfig{
		DB: db.Get(),
	})
	addressRepo := repository.NewAddressRepository(repository.AddressRepositoryConfig{
		DB: db.Get(),
	})
	orderItemRepo := repository.NewOrderItemRepository(repository.OrderItemRepositoryConfig{
		DB: db.Get(),
	})
	productAnalyticRepo := repository.NewProductAnalyticRepository(repository.ProductAnalyticRepositoryConfig{
		DB: db.Get(),
	})
	productReviewRepo := repository.NewProductReviewRepository(repository.ProductReviewRepositoryConfig{
		DB:                        db.Get(),
		ProductAnalyticRepository: productAnalyticRepo,
		MerchantRepository:        merchantRepo,
	})
	userFavoriteProductRepo := repository.NewUserFavoriteProductRepository(repository.UserFavoriteProductRepositoryConfig{
		DB:                        db.Get(),
		ProductAnalyticRepository: productAnalyticRepo,
	})
	userOrderRepo := repository.NewUserOrderRepository(repository.UserOrderRepositoryConfig{
		DB: db.Get(),
	})
	slpAccountRepo := repository.NewSlpAccountsRepository(repository.SlpAccountsRepositoryConfig{
		DB: db.Get(),
	})
	mpVoucherRepo := repository.NewMarketplaceVoucherRepository(repository.MarketplaceVoucherRepositoryConfig{
		DB: db.Get(),
	})
	transactionStatusRepo := repository.NewTransactionStatusRepository(repository.TransactionStatusRepositoryConfig{
		DB: db.Get(),
	})
	transactionDeliveryStatusRepo := repository.NewTransactionDeliveryStatusRepository(repository.TransactionDeliveryStatusRepositoryConfig{
		DB: db.Get(),
	})
	transactionPaymentRecordRepo := repository.NewTransactionPaymentRecordRepository(repository.TransactionPaymentRecordRepositoryConfig{
		DB: db.Get(),
	})
	paymentMethodRepo := repository.NewPaymentMethodRepository(repository.PaymentMethodRepositoryConfig{
		DB: db.Get(),
	})
	transactionRepo := repository.NewTransactionRepository(repository.TransactionRepositoryConfig{
		DB:                                      db.Get(),
		MarketplaceVoucherRepository:            mpVoucherRepo,
		MerchantRepository:                      merchantRepo,
		ProductRepository:                       productRepo,
		PaymentRecordRepository:                 paymentRecordRepo,
		TransactionDeliveryStatusRepository:     transactionDeliveryStatusRepo,
		TransactionStatusRepository:             transactionStatusRepo,
		WalletRepository:                        walletRepo,
		TransactionPaymentRecordRepository:      transactionPaymentRecordRepo,
		MerchantHoldingAccountRepository:        merchantHoldingAccountRepo,
		MerchantHoldingAccountHistoryRepository: merchantHoldingAccountHistoryRepo,
	})
	transactionStatusRepo = repository.NewTransactionStatusRepository(repository.TransactionStatusRepositoryConfig{
		DB:                       db.Get(),
		TransactionRepositoryPtr: transactionRepo,
	})
	refundRequestMessageRepo := repository.NewRefundRequestMessageRepository(repository.RefundRequestMessageRepositoryConfig{
		DB: db.Get(),
	})
	refundRequestRepo := repository.NewRefundRequestRepository(repository.RefundRequestRepositoryConfig{
		DB:                    db.Get(),
		TransactionRepository: transactionRepo,
	})
	promotionBannerRepo := repository.NewPromotionBannerRepository(repository.PromotionBannerRepositoryConfig{
		DB: db.Get(),
	})
	promotionRepo := repository.NewPromotionRepository(repository.PromotionRepositoryConfig{
		DB: db.Get(),
	})
	sealabspayRepo := repository.NewSealabspayRepository(repository.SealabspayRepositoryConfig{})
	marketplaceAnalyticsRepo := repository.NewMarketplaceAnalyticsRepository(repository.MarketplaceAnalyticsRepositoryConfig{
		DB: db.Get(),
	})
	merchantAnalyticsRepo := repository.NewMerchantAnalyticsRepository(repository.MerchantAnalyticsRepositoryConfig{
		DB: db.Get(),
	})
	gscUploader := util.NewGCSUploader(util.GCSUploaderConfig{
		ClientUploader: util.NewClientUploader(),
	})

	exampleUsecase := usecase.NewExampleUsecase(usecase.ExampleUsecaseConfig{
		ExampleRepository: exampleRepo,
	})
	userFavoriteProductUsecase := usecase.NewUserFavoriteProductUsecase(usecase.UserFavoriteProductUsecaseConfig{
		UserRepository:                userRepo,
		ProductAnalyticRepository:     productAnalyticRepo,
		UserFavoriteProductRepository: userFavoriteProductRepo,
		ProductRepository:             productRepo,
	})
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseConfig{
		UserRepository:   userRepo,
		AuthRepository:   authRepo,
		WalletRepository: walletRepo,
		AuthUtil:         authUtil,
	})
	merchantUsecase := usecase.NewMerchantUsecase(usecase.MerchantUsecaseConfig{
		MerchantRepository:                      merchantRepo,
		CategoryRepository:                      categoryRepo,
		UserRepository:                          userRepo,
		MerchantHoldingAccountHistoryRepository: merchantHoldingAccountHistoryRepo,
		MerchantHoldingAccountRepository:        merchantHoldingAccountRepo,
		GcsUploader:                             gscUploader,
	})
	mediaUsecase := usecase.NewMediaUsecase(usecase.MediaUsecaseConfig{
		GCSUploader: gscUploader,
	})
	userUsecase := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
		UserRepository:     userRepo,
		AddressRepository:  addressRepo,
		MerchantRepository: merchantRepo,
		MediaUsecase:       mediaUsecase,
	})
	productUsecase := usecase.NewProductUsecase(usecase.ProductUsecaseConfig{
		ProductRepository:  productRepo,
		CategoryRepository: categoryRepo,
		UserRepository:     userRepo,
		MerchantRepository: merchantRepo,
		MediaUsecase:       mediaUsecase,
	})
	productVariantUsecase := usecase.NewProductVariantUsecase(usecase.ProductVariantUsecaseConfig{
		ProductVariantRepository: productVariantRepo,
		ProductRepository:        productRepo,
	})
	productReviewUsecase := usecase.NewProductReviewUsecase(usecase.ProductReviewUsecaseConfig{
		ProductReviewRepository: productReviewRepo,
		TransactionRepository:   transactionRepo,
		MerchantRepository:      merchantRepo,
		UserRepository:          userRepo,
		GcsUploader:             gscUploader,
	})
	categoryUsecase := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
		CategoryRepository: categoryRepo,
		MediaUsecase:       mediaUsecase,
	})
	cartItemUsecase := usecase.NewCartItemUsecase(usecase.CartItemUsecaseConfig{
		CartItemRepository:       cartItemRepo,
		UserRepository:           userRepo,
		MerchantRepository:       merchantRepo,
		ProductRepository:        productRepo,
		ProductVariantRepository: productVariantRepo,
	})
	addressUsecase := usecase.NewAddressUsecase(usecase.AddressUsecaseConfig{
		AddressRepository: addressRepo,
	})
	orderItemUsecase := usecase.NewOrderItemUsecase(usecase.OrderItemUsecaseConfig{
		OrderItemRepository:          orderItemRepo,
		CartItemRepository:           cartItemRepo,
		UserRepository:               userRepo,
		UserOrderRepository:          userOrderRepo,
		ProductRepository:            productRepo,
		ProductVariantRepository:     productVariantRepo,
		DeliveryRepository:           deliveryRepo,
		AddressRepository:            addressRepo,
		MarketplaceVoucherRepository: mpVoucherRepo,
		MerchantRepository:           merchantRepo,
	})
	slpAccountUsecase := usecase.NewSlpAccountUsecase(usecase.SlpAccountUsecaseConfig{
		UserRepository:        userRepo,
		SlpAccountsRepository: slpAccountRepo,
	})
	mpVoucherUsecase := usecase.NewMarketplaceVoucherUsecase(usecase.MarketplaceVoucherUsecaseConfig{
		MarketplaceVoucherRepo: mpVoucherRepo,
	})
	transactionStatusUsecase := usecase.NewTransactionStatusUsecase(usecase.TransactionStatusUsecaseConfig{
		TransactionStatusRepo: transactionStatusRepo,
		Cron:                  cronjob.GetCron(),
	})
	transactionDeliveryStatusUsecase := usecase.NewTransactionDeliveryStatusUsecase(usecase.TransactionDeliveryStatusUsecaseConfig{
		TransactionDeliveryStatusRepo: transactionDeliveryStatusRepo,
	})
	transactionUsecase := usecase.NewTransactionUsecase(usecase.TransactionUsecaseConfig{
		CartItemRepository:                  cartItemRepo,
		TransactionRepository:               transactionRepo,
		UserRepository:                      userRepo,
		MerchantRepository:                  merchantRepo,
		TransactionStatusUsecase:            transactionStatusUsecase,
		TransactionDeliveryStatusUsecase:    transactionDeliveryStatusUsecase,
		TransactionDeliveryStatusRepository: transactionDeliveryStatusRepo,
		TransactionStatusRepository:         transactionStatusRepo,
		OrderItemUsecase:                    orderItemUsecase,
		SealabspayRepository:                sealabspayRepo,
		PaymentMethodRepository:             paymentMethodRepo,
		WalletRepository:                    walletRepo,
	})
	paymentMethodUsecase := usecase.NewPaymentMethodUsecase(usecase.PaymentMethodUsecaseConfig{
		PaymentMethodRepository: paymentMethodRepo,
	})
	refundRequestUsecase := usecase.NewRefundRequestUsecase(usecase.RefundRequestUsecaseConfig{
		RefundRequestRepository: refundRequestRepo,
		TransactionRepository:   transactionRepo,
		UserRepository:          userRepo,
		GCSUploader:             gscUploader,
		MerchantRepository:      merchantRepo,
		WalletRepository:        walletRepo,
		Cron:                    cronjob.GetCron(),
	})
	refundRequestMessageUsecase := usecase.NewRefundRequestMessageUsecase(usecase.RefundRequestMessageUsecaseConfig{
		RefundRequestRepository:        refundRequestRepo,
		RefundRequestMessageRepository: refundRequestMessageRepo,
		GcsUploader:                    gscUploader,
		UserRepository:                 userRepo,
	})
	deliveryUsecase := usecase.NewDeliveryUsecase(usecase.DeliveryUsecaseConfig{
		DeliveryRepository: deliveryRepo,
		MerchantRepository: merchantRepo,
	})

	walletUsecase := usecase.NewWalletUsecase(usecase.WalletUsecaseConfig{
		WalletRepository:     walletRepo,
		UserRepository:       userRepo,
		SealabspayRepository: sealabspayRepo,
	})

	promotionBannerUsecase := usecase.NewPromotionBannerUsecase(usecase.PromotionBannerUsecaseConfig{
		PromotionBannerRepo: promotionBannerRepo,
		MediaUsecase:        mediaUsecase,
	})

	promotionUsecase := usecase.NewPromotionUsecase(usecase.PromotionUsecaseConfig{
		PromotionRepository: promotionRepo,
		MerchantRepository:  merchantRepo,
		ProductRepository:   productRepo,
	})

	paymentRecordUsecase := usecase.NewPaymentRecordUsecase(usecase.PaymentRecordUsecaseConfig{
		WalletRepository:                   walletRepo,
		WalletUsecase:                      walletUsecase,
		TransactionPaymentRecordRepository: transactionPaymentRecordRepo,
		UserRepository:                     userRepo,
		TransactionUsecase:                 transactionUsecase,
		PaymentRecordRepository:            paymentRecordRepo,
	})

	sealabspayUsecase := usecase.NewSealabspayUsecase(usecase.SealabspayUsecaseConfig{
		PaymentRecordUsecase: paymentRecordUsecase,
	})

	walletpayUsecase := usecase.NewWalletpayUsecase(usecase.WalletpayUsecaseConfig{
		WalletRepository:     walletRepo,
		UserRepository:       userRepo,
		PaymentRecordUsecase: paymentRecordUsecase,
	})

	marketplaceAnalyticsUsecase := usecase.NewMarketplaceAnalyticsUsecase(usecase.MarketplaceAnalyticsUsecaseConfig{
		MarketplaceAnalyticsRepository: marketplaceAnalyticsRepo,
	})

	merchantAnalyticsUsecase := usecase.NewMerchantAnalyticsUsecase(usecase.MerchantAnalyticsUsecaseConfig{
		MerchantAnalyticsRepository: merchantAnalyticsRepo,
		MerchantRepository:          merchantRepo,
	})

	r := NewRouter(RouterConfig{
		ExampleUsecase:                   exampleUsecase,
		UserUsecase:                      userUsecase,
		UserFavoriteProductUsecase:       userFavoriteProductUsecase,
		AuthUsecase:                      authUsecase,
		MerchantUsecase:                  merchantUsecase,
		WalletUsecase:                    walletUsecase,
		ProductUsecase:                   productUsecase,
		ProductVariantUsecase:            productVariantUsecase,
		ProductReviewUsecase:             productReviewUsecase,
		PaymentMethodUsecase:             paymentMethodUsecase,
		PaymentRecordUsecase:             paymentRecordUsecase,
		CategoryUsecase:                  categoryUsecase,
		CartItemUsecase:                  cartItemUsecase,
		AddressUsecase:                   addressUsecase,
		OrderItemUsecase:                 orderItemUsecase,
		SlpAccountUsecase:                slpAccountUsecase,
		SealabspayUsecase:                sealabspayUsecase,
		MpVoucherUsecase:                 mpVoucherUsecase,
		DeliveryUsecase:                  deliveryUsecase,
		TransactionStatusUsecase:         transactionStatusUsecase,
		TransactionDeliveryStatusUsecase: transactionDeliveryStatusUsecase,
		TransactionUsecase:               transactionUsecase,
		RefundRequestUsecase:             refundRequestUsecase,
		RefundRequestMessageUsecase:      refundRequestMessageUsecase,
		MediaUsecase:                     mediaUsecase,
		WalletpayUsecase:                 walletpayUsecase,
		MarketplaceAnalyticsUsecase:      marketplaceAnalyticsUsecase,
		MerchantAnalyticsUsecase:         merchantAnalyticsUsecase,
		PromotionBannerUsecase:           promotionBannerUsecase,
		PromotionUsecase:                 promotionUsecase,
	})
	return r
}

func Init() {
	r := initRouter()
	err := r.Run()
	// err := r.RunTLS(":8080", "localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal().Msgf("error while running server %v", err)
		return
	}
}
