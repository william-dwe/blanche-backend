package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetAllPaymentMethod = httperror.InternalServerError("Error when getting all payment methods")
var ErrPaymentMethodNotFound = httperror.BadRequestError("Payment method not found", "PAYMENT_METHOD_NOT_FOUND")
var ErrGetPaymentMethodByCode = httperror.InternalServerError("Error when getting payment method by code")
