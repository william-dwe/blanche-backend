package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetMerchantAnalytics = httperror.InternalServerError("failed to get merchant analytics")
var ErrUpdateMerchantAnalytics = httperror.InternalServerError("failed to update merchant analytics")
