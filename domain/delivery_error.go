package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetAllDeliveryInternalError = httperror.InternalServerError("cannot get all delivery option")
var ErrGetDeliveryFeeInternalError = httperror.InternalServerError("cannot get delivery fee")
var ErrUpdateDeliveryOptionInternalError = httperror.InternalServerError("cannot update delivery fee")
