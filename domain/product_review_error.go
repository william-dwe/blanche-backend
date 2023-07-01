package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrAddProductReviewRatingNotValid = httperror.BadRequestError("Product review rating must be between 1 and 5", "ERR_PRODUCT_REVIEW_RATING_NOT_VALID")
var ErrAddProductReview = httperror.InternalServerError("Failed to add product review")
var ErrAddProductReviewTransactionIdNotValid = httperror.BadRequestError("Transaction id not valid", "ERR_ADD_PRODUCT_REVIEW_TRANSACTION_ID_NOT_VALID")
var ErrAddProductReviewProductIdNotValid = httperror.BadRequestError("Product id not valid", "ERR_ADD_PRODUCT_REVIEW_PRODUCT_ID_NOT_VALID")
var ErrAddProductReviewVariantItemIdNotValid = httperror.BadRequestError("Variant item id not valid", "ERR_ADD_PRODUCT_REVIEW_VARIANT_ITEM_ID_NOT_VALID")
var ErrAddProductReviewDuplicate = httperror.BadRequestError("Product review already exists", "ERR_ADD_PRODUCT_REVIEW_DUPLICATE")
var ErrAddProductReviewProductNotRelatedToTransaction = httperror.BadRequestError("cannot review product which not related to transaction", "ERR_ADD_PRODUCT_REVIEW_PRODUCT_NOT_RELATED_TO_TRANSACTION")
var ErrAddProductReviewTransactionNotCompleted = httperror.BadRequestError("cannot review product which transaction not completed", "ERR_ADD_PRODUCT_REVIEW_TRANSACTION_NOT_COMPLETED")

var ErrProductReviewInvalidParamReq = httperror.BadRequestError("Invalid request parameter", "ERR_PRODUCT_REVIEW_INVALID_PARAM_REQ")
var ErrGetProductReviewByProductSlug = httperror.InternalServerError("Failed to get product review by product slug")
