package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetPromotionList = httperror.InternalServerError("failed to get promotion list")
var ErrCreateNewPromotion = httperror.InternalServerError("failed to create new promotion")
var ErrInvalidNominal = httperror.BadRequestError("minimum discount price is 100", "INVALID_NOMINAL")
var ErrInvalidPercentage = httperror.BadRequestError("please check your again you input, discount percentage range is between 1-100", "INVALID_PERCENTAGE")
var ErrInvalidPromotionType = httperror.BadRequestError("please check your again you input, promotion type is either percentage or nominal", "INVALID_PROMOTION_TYPE")
var ErrGetPromotionByID = httperror.InternalServerError("failed to get promotion by id")
var ErrPromotionNotFound = httperror.NotFoundError("promotion not found")
var ErrUpdatePromotion = httperror.InternalServerError("failed to update promotion")
var ErrForbiddenMerchant = httperror.ForbiddenError()
var ErrPromotionAlreadyEnded = httperror.BadRequestError("the promotion is already ended", "PROMOTION_ALREADY_ENDED")
var ErrDeletePromotion = httperror.InternalServerError("failed to delete promotion")
var ErrPromotionAlreadyStarted = httperror.BadRequestError("the promotion is already started", "PROMOTION_ALREADY_STARTED")
var ErrCheckProductPromotionOngoing = httperror.BadRequestError("one of the product is already in ongoing promotion", "PRODUCT_PROMOTION_ONGOING")
var ErrCheckProductPromotion = httperror.InternalServerError("failed to check product promotion")
var ErrInvalidPromotionDateRange = httperror.BadRequestError("promotion start date must before than promotion end date", "INVALID_PROMOTION_DATE_RANGE")
var ErrInvalidProduct = httperror.BadRequestError("one of the products is doesn't belong to this merchant", "INVALID_PRODUCT")
var ErrPromotionIDNotValid = httperror.BadRequestError("promotion id is not valid", "PROMOTION_ID_NOT_VALID")