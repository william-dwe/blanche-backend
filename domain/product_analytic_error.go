package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrProductAnalyticUpdateFavoriteProduct = httperror.InternalServerError("failed to update favorite product count")
var ErrProductAnalyticUpdateAvgRatingAndNumReview = httperror.InternalServerError("failed to update avg rating and num review product")
