package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetProducts = httperror.InternalServerError("failed to get products record")
var ErrGetProduct = httperror.InternalServerError("failed to get product record")
var ErrGetProductNotFound = httperror.BadRequestError("failed to retrieve product details, product not found", "PRODUCT_NOT_FOUND")
var ErrInvalidPriceRange = httperror.BadRequestError("the minimum price must be lower than the maximum price", "INVALID_PRICE_RANGE")

var ErrProductQueryParam = httperror.BadRequestError("invalid query param", "INVALID_PARAMS")

var ErrGetProductSlugNotFound = httperror.BadRequestError("failed to retrieve product, product slug not found", "PRODUCT_NOT_FOUND")
var ErrGetProductIdNotFound = httperror.BadRequestError("failed to retrieve product, product id not found", "PRODUCT_NOT_FOUND")
var ErrGetProductDetailSlugNotFound = httperror.BadRequestError("failed to retrieve product details, product slug not found", "PRODUCT_NOT_FOUND")

var ErrGetProductVariant = httperror.InternalServerError("failed to get product variant record")
var ErrGetProductVariantInconsistent = httperror.InternalServerError("failed to get product variant record, encounter inconsistent data ")

var ErrIncreaseNumOfSale = httperror.InternalServerError("failed to increase number of sale")

var ErrGetProductPromotionNotFound = httperror.BadRequestError("failed to retrieve product promotion, product promotion not found", "PRODUCT_PROMOTION_NOT_FOUND")
var ErrGetProductPromotion = httperror.InternalServerError("failed to get product promotion record")
var ErrInvalidPage = httperror.BadRequestError("invalid page", "INVALID_PAGE")

var ErrPriceVariantEmpty = httperror.BadRequestError("price variant failed to be empty", "PRICE_VARIANT_EMPTY")
var ErrPriceVariantInvalid = httperror.BadRequestError("price and variant items failed to be filled at the same time", "PRICE_VARIANT_INVALID")
var ErrPriceVariantInvalidPrice = httperror.BadRequestError("price must be greater than 0", "PRICE_VARIANT_INVALID_PRICE")

var ErrCreateVariantGroup = httperror.InternalServerError("failed to create variant group")
var ErrCreateProduct = httperror.InternalServerError("failed to create product")
var ErrProductAlreadyExist = httperror.BadRequestError("this product is already exist", "PRODUCT_ALREADY_EXIST")
var ErrCheckMerchantProductName = httperror.InternalServerError("failed to check merchant product name")

var ErrDeleteVariantSpec = httperror.InternalServerError("failed to delete variant spec")
var ErrDeleteVariantGroup = httperror.InternalServerError("failed to delete variant group")
var ErrDeleteVariantItem = httperror.InternalServerError("failed to delete variant item")
var ErrDeleteProduct = httperror.InternalServerError("failed to delete product")
var ErrDeleteProductImages = httperror.InternalServerError("failed to delete product images")
var ErrDeleteProductAnalytics = httperror.InternalServerError("failed to delete product analytics")

var ErrUpdateProduct = httperror.InternalServerError("failed to update product")
var ErrUpdateProductUnauthorized = httperror.UnauthorizedError()
var ErrUpdateMerchantProductStatus = httperror.InternalServerError("failed to update merchant product status")

var ErrProductIdNotValid = httperror.BadRequestError("product id is not valid", "PRODUCT_ID_NOT_VALID")
