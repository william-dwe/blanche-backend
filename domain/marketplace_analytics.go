package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrInvalidStartDateFormat = httperror.BadRequestError("invalid start date format", "INVALID_START_DATE_FORMAT")
var ErrInvalidEndDateFormat = httperror.BadRequestError("invalid end date format", "INVALID_END_DATE_FORMAT")

var ErrGetMarketplaceAnalytics = httperror.InternalServerError("failed to get marketplace analytics")
var ErrUpdateMarketplaceAnalytics = httperror.InternalServerError("failed to update marketplace analytics")
