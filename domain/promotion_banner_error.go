package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrPromotionBannersNotFound = httperror.BadRequestError("Promotion banners not found", "PROMOTION_BANNERS_NOT_FOUND")
var ErrGetPromotionBanners = httperror.InternalServerError("Failed to get promotion banners")
var ErrPromotionBannerNotFound = httperror.BadRequestError("Promotion banner not found", "PROMOTION_BANNER_NOT_FOUND")
var ErrGetPromotionBanner = httperror.InternalServerError("Failed to get promotion banner")
var ErrCreatePromotionBanner = httperror.InternalServerError("Failed to create promotion banner")
var ErrDuplicatePromotionBanner = httperror.BadRequestError("The promotion banner already exists", "DUPLICATE_PROMOTION_BANNER")
var ErrUpdatePromotionBanner = httperror.InternalServerError("Failed to update promotion banner")
var ErrDeletePromotionBanner = httperror.InternalServerError("Failed to delete promotion banner")
var ErrPromotionBannerIdNotValid = httperror.BadRequestError("Promotion banner id is not valid", "PROMOTION_BANNER_ID_NOT_VALID")