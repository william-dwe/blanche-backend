package server

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/handler"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/middleware"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/usecase"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AddressUsecase                   usecase.AddressUsecase
	AuthUsecase                      usecase.AuthUsecase
	CategoryUsecase                  usecase.CategoryUsecase
	CartItemUsecase                  usecase.CartItemUsecase
	DeliveryUsecase                  usecase.DeliveryUsecase
	ExampleUsecase                   usecase.ExampleUsecase
	MerchantUsecase                  usecase.MerchantUsecase
	MpVoucherUsecase                 usecase.MarketplaceVoucherUsecase
	OrderItemUsecase                 usecase.OrderItemUsecase
	PaymentMethodUsecase             usecase.PaymentMethodUsecase
	PaymentRecordUsecase             usecase.PaymentRecordUsecase
	ProductUsecase                   usecase.ProductUsecase
	ProductVariantUsecase            usecase.ProductVariantUsecase
	ProductReviewUsecase             usecase.ProductReviewUsecase
	SlpAccountUsecase                usecase.SlpAccountUsecase
	SealabspayUsecase                usecase.SealabspayUsecase
	TransactionUsecase               usecase.TransactionUsecase
	TransactionStatusUsecase         usecase.TransactionStatusUsecase
	TransactionDeliveryStatusUsecase usecase.TransactionDeliveryStatusUsecase
	RefundRequestUsecase             usecase.RefundRequestUsecase
	RefundRequestMessageUsecase      usecase.RefundRequestMessageUsecase
	UserUsecase                      usecase.UserUsecase
	UserFavoriteProductUsecase       usecase.UserFavoriteProductUsecase
	WalletUsecase                    usecase.WalletUsecase
	WalletpayUsecase                 usecase.WalletpayUsecase
	MediaUsecase                     usecase.MediaUsecase
	MarketplaceAnalyticsUsecase      usecase.MarketplaceAnalyticsUsecase
	MerchantAnalyticsUsecase         usecase.MerchantAnalyticsUsecase
	PromotionBannerUsecase           usecase.PromotionBannerUsecase
	PromotionUsecase                 usecase.PromotionUsecase
}

func NewRouter(c RouterConfig) *gin.Engine {
	h := handler.New(handler.HandlerConfig{
		AddressUsecase:                   c.AddressUsecase,
		AuthUsecase:                      c.AuthUsecase,
		CategoryUsecase:                  c.CategoryUsecase,
		CartItemUsecase:                  c.CartItemUsecase,
		DeliveryUsecase:                  c.DeliveryUsecase,
		ExampleUsecase:                   c.ExampleUsecase,
		MerchantUsecase:                  c.MerchantUsecase,
		MpVoucherUsecase:                 c.MpVoucherUsecase,
		OrderItemUsecase:                 c.OrderItemUsecase,
		PaymentRecordUsecase:             c.PaymentRecordUsecase,
		PaymentMethodUsecase:             c.PaymentMethodUsecase,
		ProductUsecase:                   c.ProductUsecase,
		ProductVariantUsecase:            c.ProductVariantUsecase,
		ProductReviewUsecase:             c.ProductReviewUsecase,
		SlpAccountUsecase:                c.SlpAccountUsecase,
		SealabspayUsecase:                c.SealabspayUsecase,
		TransactionUsecase:               c.TransactionUsecase,
		TransactionStatusUsecase:         c.TransactionStatusUsecase,
		TransactionDeliveryStatusUsecase: c.TransactionDeliveryStatusUsecase,
		RefundRequestUsecase:             c.RefundRequestUsecase,
		RefundRequestMessageUsecase:      c.RefundRequestMessageUsecase,
		UserUsecase:                      c.UserUsecase,
		UserFavoriteProductUsecase:       c.UserFavoriteProductUsecase,
		WalletUsecase:                    c.WalletUsecase,
		WalletpayUsecase:                 c.WalletpayUsecase,
		MediaUsecase:                     c.MediaUsecase,
		MarketplaceAnalyticsUsecase:      c.MarketplaceAnalyticsUsecase,
		MerchantAnalyticsUsecase:         c.MerchantAnalyticsUsecase,
		PromotionBannerUsecase:           c.PromotionBannerUsecase,
		PromotionUsecase:                 c.PromotionUsecase,
	})

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler) // global error middleware

	r.NoRoute(func(c *gin.Context) {
		util.ResponseErrorJSON(c, httperror.NotFoundError("endpoint not found"))
	})
	r.GET("/ping", func(c *gin.Context) {
		util.ResponseSuccessJSONData(c, "pong!")
	})

	apiEndpoint := r.Group("/api")
	v1 := apiEndpoint.Group("/v1")

	v1.GET("/provinces", h.GetAllProvinces)
	v1.GET("/cities", h.GetAllCities)
	v1.GET("/cities/:provinceId", h.GetCitiesByProvinceID)
	v1.GET("/districts/:cityId", h.GetDistrictsByCityID)
	v1.GET("/subdistricts/:districtId", h.GetSubDistrictsByDistrictID)

	v1.POST("/example-process", h.ExampleHandler)
	v1.POST("/example-cache", h.CachedExampleHandler)
	v1.POST("/example-process-error", h.ExampleHandlerErrorMiddleware)

	v1.POST("/register/check-email", h.UserRegisterCheckEmailHandler)
	v1.POST("/register/check-username", h.UserRegisterCheckUsernameHandler)
	v1.POST("/register", h.UserRegisterHandler)

	v1.GET("/google/request-login", h.GoogleRequestLogin)
	v1.GET("/google/request-callback", h.GoogleRequestCallback)

	v1.POST("/login", h.UserLoginHandler)
	v1.POST("/login/admin", h.AdminLoginHandler)
	v1.GET("/refresh", h.UserRefreshHandler)
	v1.POST("/logout", h.UserLogoutHandler)

	userEndpoints := v1.Group("/users")
	userEndpoints.Use(middleware.Authenticate)
	userEndpoints.Use(middleware.Authorize(h, dto.ROLE_USER))
	userEndpoints.GET("/profile", h.UserGetProfileHandler)
	userEndpoints.PATCH("/profile", h.UserUpdateProfileHandler)
	userEndpoints.PATCH("/profile-details", h.UserUpdateProfileDetailHandler)
	userEndpoints.GET("/favorite-products", h.GetUserFavoriteProducts)
	userEndpoints.POST("/favorite-products", h.UpdateUserFavoriteProduct)

	changePasswordEndpoints := userEndpoints.Group("/password/change-password")
	changePasswordEndpoints.POST("/send-code", h.ChangePasswordRequestVerificationCode)
	changePasswordEndpoints.POST("/verify-code", h.ChangePasswordVerification)
	forgetPasswordEndpoints := v1.Group("users/password/forget-password")
	forgetPasswordEndpoints.POST("/send-code", h.ForgetPasswordRequestVerificationUrl)
	forgetPasswordEndpoints.POST("/verify-code", h.ForgetPasswordVerification)

	passwordEndpoints := v1.Group("/users/password")
	passwordEndpoints.Use(middleware.Authenticate)
	passwordEndpoints.Use(middleware.Authorize(h, dto.SCOPE_RESET_PASSWORD))
	passwordEndpoints.POST("/reset", middleware.Authorize(h, "reset_password"), h.ResetPassword)

	authEndpoints := v1.Group("/step-up")
	authEndpoints.Use(middleware.Authenticate)
	authEndpoints.Use(middleware.Authorize(h, dto.ROLE_USER))
	authEndpoints.POST("/pin", h.StepUpTokenScopeWithPin)
	authEndpoints.POST("/pass", h.StepUpTokenScopeWithPass)

	userRefundRequest := userEndpoints.Group("/refund-requests")
	userRefundRequest.POST("", h.AddRefundRequest)
	userRefundRequest.GET("", h.GetUserRefundRequestList)
	userRefundRequest.POST("/:refund_id/accept", h.UserAcceptRequestRefund)
	userRefundRequest.POST("/:refund_id/reject", h.UserRejectRequestRefund)
	userRefundRequest.POST("/:refund_id/cancel", h.UserCancelRequestRefund)
	userRefundRequest.GET("/:refund_id/messages", h.UserGetMessageRequestRefund)
	userRefundRequest.POST("/:refund_id/messages", h.BuyerAddMessageRequestRefund)

	transactionEndpoints := userEndpoints.Group("/transactions")
	transactionEndpoints.POST("", h.MakeTransaction)
	transactionEndpoints.GET("", h.GetTransactionList)
	transactionEndpoints.GET("/:invoice_code", h.GetTransactionDetail)
	transactionEndpoints.PUT(":invoice_code/status", h.UpdateUserTransactionStatus)
	transactionEndpoints.GET("/:invoice_code/review", h.GetProductReviewFromTransaction)
	transactionEndpoints.POST("/:invoice_code/review", h.AddProductReviewFromTransaction)

	merchantEndpoints := v1.Group("/merchants")
	merchantEndpoints.GET("/:domain/profile", h.GetMerchantInfo)
	merchantEndpoints.GET("/:domain/categories", h.GetMerchantProductCategories)
	merchantEndpoints.GET("/:domain/vouchers", h.GetMerchantVoucherList)
	merchantEndpoints.GET("/:domain/deliveries", h.GetDeliveryOptionByMerchantDomain)

	merchantEndpoints.Use(middleware.Authenticate)
	merchantEndpoints.Use(middleware.Authorize(h, dto.ROLE_USER))
	merchantEndpoints.POST("/register", h.RegisterMerchant)
	merchantEndpoints.POST("/register/check-domain", h.CheckMerchantDomain)
	merchantEndpoints.POST("/register/check-name", h.CheckMerchantStoreName)

	merchantEndpoints.Use(middleware.Authorize(h, dto.ROLE_MERCHANT))
	merchantEndpoints.GET("/profile", h.GetUserMerchantInfo)
	merchantEndpoints.PATCH("/profile", h.UpdateMerchantProfile)
	merchantEndpoints.POST("/products/images", h.UploadProductImage)
	merchantEndpoints.GET("/products", h.GetMerchantProductList)
	merchantEndpoints.POST("/products", h.CreateProduct)
	merchantEndpoints.GET("/products/:product_id", h.GetMerchantProductDetails)
	merchantEndpoints.GET("/products/:product_id/variants", h.GetMerchantProductVariants)
	merchantEndpoints.PUT("/products/:product_id", h.UpdateMerchantProduct)
	merchantEndpoints.PATCH("/products/:product_id/status", h.UpdateProductAvailability)
	merchantEndpoints.DELETE("/products/:product_id", h.DeleteMerchantProduct)
	merchantEndpoints.POST("/products/check-name", h.CheckMerchantProductName)
	merchantEndpoints.GET("/deliveries", h.GetMerchantUserDeliveryOption)
	merchantEndpoints.PUT("/deliveries", h.ChangeMerchantUserDeliveryOption)
	merchantEndpoints.PUT("/transactions/:invoice_code/status", h.UpdateMerchantTransactionStatus)
	merchantEndpoints.GET("/transactions", h.GetSellerTransactionList)
	merchantEndpoints.GET("/transactions/:invoice_code", h.GetSellerTransactionDetail)
	merchantEndpoints.GET("/vouchers", h.GetMerchantAdminVoucherList)
	merchantEndpoints.POST("/vouchers", h.CreateMerchantVoucher)
	merchantEndpoints.GET("/vouchers/:voucher_code", h.GetMerchantAdminVoucherDetails)
	merchantEndpoints.PUT("/vouchers/:voucher_code", h.UpdateMerchantVoucher)
	merchantEndpoints.DELETE("/vouchers/:voucher_code", h.DeleteMerchantAdminVoucher)
	merchantEndpoints.GET("/funds/activities", h.GetMerchantFundActivities)
	merchantEndpoints.GET("/funds/balance", h.GetMerchantFundBalance)
	merchantEndpoints.POST("/funds/withdraw", h.WithdrawMerchantFundBalance)
	merchantEndpoints.PATCH("/addresses/:address_id", h.UpdateMerchantAddress)
	merchantEndpoints.GET("/promotions", h.GetAllPromotions)
	merchantEndpoints.GET("/promotions/:promotion_id", h.GetPromotionDetails)
	merchantEndpoints.POST("/promotions", h.CreateNewPromotion)
	merchantEndpoints.PUT("/promotions/:promotion_id", h.UpdatePromotion)
	merchantEndpoints.DELETE("/promotions/:promotion_id", h.DeletePromotion)

	merchantDashboardEndpoints := merchantEndpoints.Group("/dashboards")
	merchantDashboardEndpoints.GET("responsiveness", h.GetMerchantDashboardMerchantResponsivenessStatistics)
	merchantDashboardEndpoints.GET("sales", h.GetMerchantDashboardSalesStatistics)
	merchantDashboardEndpoints.GET("customer-satisfactions", h.GetMerchantDashboardCustomerSatisfactionStatistics)

	merchantRefundReqEndpoints := merchantEndpoints.Group("/refund-requests")
	merchantRefundReqEndpoints.GET("", h.GetMerchantRefundRequestList)
	merchantRefundReqEndpoints.POST("/:refund_id/accept", h.MerchantAcceptRequestRefund)
	merchantRefundReqEndpoints.POST("/:refund_id/reject", h.MerchantRejectRequestRefund)
	merchantRefundReqEndpoints.GET("/:refund_id/messages", h.UserGetMessageRequestRefund)
	merchantRefundReqEndpoints.POST("/:refund_id/messages", h.MerchantAddMessageRequestRefund)

	deliveryEndpoints := v1.Group("/deliveries")
	deliveryEndpoints.GET("", h.GetAllDeliveryOption)

	walletEndpoints := userEndpoints.Group("/wallet")
	walletEndpoints.POST("/make-payment", middleware.AuthorizeAndBlacklist(h, dto.SCOPE_PIN), h.WalletpayPayReq)
	walletEndpoints.POST("/cancel-payment", h.WalletpayCancelPayReq)
	walletEndpoints.GET("", h.GetWalletDetails)
	walletEndpoints.GET("/transactions", h.GetWalletTransactions)
	walletEndpoints.POST("/create-pin", h.CreateWallet)
	walletEndpoints.POST("/change-pin", middleware.AuthorizeAndBlacklist(h, dto.SCOPE_PASSWORD), h.UpdateWalletPin)
	walletEndpoints.POST("/topup", h.MakeTopUpWalletSlp)

	slpEndpoints := userEndpoints.Group("/slp-accounts")
	slpEndpoints.GET("", h.GetUserSlpAccountList)
	slpEndpoints.GET(":slp_id", h.GetUserSlpAccount)
	slpEndpoints.POST("", h.RegisterUserSlpAccount)
	slpEndpoints.PATCH(":slp_id/default", h.SetDefaultSlpAccount)
	slpEndpoints.DELETE(":slp_id", h.DeleteUserSlpAccount)

	addressEndpoints := userEndpoints.Group("/addresses")
	addressEndpoints.GET("", h.GetAllUserAddress)
	addressEndpoints.POST("", h.AddUserAddress)
	addressEndpoints.PATCH("/:address_id/default", h.SetDefaultUserAddress)
	addressEndpoints.PUT("/:address_id", h.UpdateUserAddress)
	addressEndpoints.DELETE("/:address_id", h.DeleteUserAddress)

	productEndpoints := v1.Group("/products")
	productEndpoints.GET("", h.GetProductList)
	productEndpoints.GET("/recommendations", h.GetRecommendationProductList)
	productEndpoints.GET("/:domain/:slug/variants", h.GetProductVariants)
	productEndpoints.Use(middleware.AuthenticateWithByPass)
	productEndpoints.GET("/:domain/:slug/details", h.GetProductDetails)
	productEndpoints.GET("/:domain/:slug/reviews", h.GetProductReviewByProductSlug)

	categoryEndpoints := v1.Group("/categories")
	categoryEndpoints.GET("", h.GetCategoryTree)
	categoryEndpoints.GET("/:category_slug/ancestors", h.GetCategoryAncestorsById)
	categoryEndpoints.GET("/:category_slug", h.GetCategoryBySlug)

	cartsEndpoints := v1.Group("/carts")
	cartsEndpoints.Use(middleware.Authenticate)
	cartsEndpoints.Use(middleware.Authorize(h, dto.ROLE_USER))
	cartsEndpoints.POST("", h.AddItemToCart)
	cartsEndpoints.GET("", h.GetCart)
	cartsEndpoints.GET("/home", h.GetHomeCart)
	cartsEndpoints.DELETE("", h.DeleteSelectedCartItem)
	cartsEndpoints.DELETE("/:cart_item_id", h.DeleteCartItem)
	cartsEndpoints.PUT("", h.UpdateAllCart)
	cartsEndpoints.PATCH("/:cart_item_id", h.UpdateCart)

	orderEndpoints := v1.Group("/orders")
	orderEndpoints.Use(middleware.Authenticate)
	orderEndpoints.Use(middleware.Authorize(h, dto.ROLE_USER))
	orderEndpoints.POST("", h.MakeOrderCheckout)
	orderEndpoints.POST("/summary", h.GetOrderSummary)
	orderEndpoints.GET("/waiting-for-payment", h.GetWaitingForPayments)
	orderEndpoints.GET("/waiting-for-payment/:payment_id", h.GetWaitingForPaymentDetails)

	marketplaceEndpoints := v1.Group("/marketplace")
	marketplacePromotionBannerEndpoints := marketplaceEndpoints.Group("/promotion-banners")
	marketplacePromotionBannerEndpoints.GET("", h.GetPromotionBannerList)
	marketplacePromotionBannerEndpoints.Use(middleware.AuthenticateAdmin)
	marketplacePromotionBannerEndpoints.Use(middleware.Authorize(h, dto.ROLE_ADMIN))
	marketplacePromotionBannerEndpoints.GET("/:banner_id", h.GetPromotionBannerByID)
	marketplacePromotionBannerEndpoints.POST("", h.CreatePromotionBanner)
	marketplacePromotionBannerEndpoints.PUT("/:banner_id", h.UpdatePromotionBanner)
	marketplacePromotionBannerEndpoints.DELETE("/:banner_id", h.DeletePromotionBanner)

	marketplaceEndpoints.GET("/vouchers", h.GetMarketplaceVoucherList)
	marketplaceEndpoints.Use(middleware.AuthenticateAdmin)
	marketplaceEndpoints.Use(middleware.Authorize(h, dto.ROLE_ADMIN))
	marketplaceEndpoints.POST("/vouchers", h.CreateMarketplaceVoucher)
	marketplaceEndpoints.GET("/vouchers/admin", h.GetMarketplaceAdminVoucherList)
	marketplaceEndpoints.GET("/vouchers/:voucher_code", h.GetMarketplaceVoucherDetails)
	marketplaceEndpoints.PUT("/vouchers/:voucher_code", h.UpdateMarketplaceVoucher)
	marketplaceEndpoints.DELETE("/vouchers/:voucher_code", h.DeleteMarketplaceVoucher)

	marketplaceCategoryEndpoints := marketplaceEndpoints.Group("/categories")
	marketplaceCategoryEndpoints.POST("", h.CreateCategory)
	marketplaceCategoryEndpoints.GET("", h.GetCategoryList)
	marketplaceCategoryEndpoints.GET("/:category_id", h.GetCategoryDetailByID)
	marketplaceCategoryEndpoints.PUT("/:category_id", h.UpdateCategory)
	marketplaceCategoryEndpoints.DELETE("/:category_id", h.DeleteCategory)

	marketplaceDashboardEndpoints := marketplaceEndpoints.Group("/dashboards")
	marketplaceDashboardEndpoints.GET("active-users", h.GetMarketplaceDashboardActiveUserStatistics)
	marketplaceDashboardEndpoints.GET("user-conversions", h.GetMarketplaceDashboardUserConversionStatistics)
	marketplaceDashboardEndpoints.GET("sales", h.GetMarketplaceDashboardSalesStatistics)
	marketplaceDashboardEndpoints.GET("customer-satisfactions", h.GetMarketplaceDashboardCustomerSatisfactionStatistics)
	marketplaceDashboardEndpoints.PATCH("", h.UpdateMarketplaceDashboard)
	marketplaceDashboardEndpoints.PATCH("/merchants", h.UpdateMerchantDashboard)

	marketplaceRefundReqEndpoints := marketplaceEndpoints.Group("/refund-requests")
	marketplaceRefundReqEndpoints.GET("", h.GetAdminRefundRequestList)
	marketplaceRefundReqEndpoints.POST("/:refund_id/accept", h.AdminAcceptRequestRefund)
	marketplaceRefundReqEndpoints.POST("/:refund_id/reject", h.AdminRejectRequestRefund)
	marketplaceRefundReqEndpoints.GET("/:refund_id/messages", h.AdminGetMessageRequestRefund)
	marketplaceRefundReqEndpoints.POST("/:refund_id/messages", h.AdminAddMessageRequestRefund)

	paymentEndpoints := v1.Group("/payments")
	paymentEndpoints.POST("/sealabspay/response", h.SealabspayResponseHandler)
	paymentEndpoints.GET("", h.GetAllPaymentMethod)

	return r
}
