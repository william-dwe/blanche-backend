package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrCreatePaymentRecord = httperror.InternalServerError("failed to create payment record")
var ErrUpdatePaymentRecord = httperror.InternalServerError("failed to update payment record")
var ErrPaymentIdNotFound = httperror.BadRequestError("payment id not found", "PAYMENT_DATA_NOT_FOUND")

var ErrGetPaymentRecord = httperror.InternalServerError("failed to get payment record")
var ErrPaymentIdExpired = httperror.BadRequestError("payment is expired", "PAYMENT_ID_EXPIRED")
var ErrPaymentAmountNotMatch = httperror.BadRequestError("payment amount not match", "PAYMENT_AMOUNT_NOT_MATCH")
