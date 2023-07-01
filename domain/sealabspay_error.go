package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrSlpCannotGenerateSignature = httperror.InternalServerError("cannot generate payment signature")
var ErrSlpRedirect = httperror.NotFoundError("cannot redirect to payment page")
var ErrSlpParsePaymentId = httperror.InternalServerError("cannot parse payment id")
var ErrSlpRequest = httperror.InternalServerError("cannot contact to slp services")
var ErrSlpInvalidSignature = httperror.BadRequestError("invalid signature", "SLP_INVALID_SIGNATURE")

var ErrSlpUserInsufficientFund = httperror.BadRequestError("user's card does not have sufficient balance", "SLP_INSUFFICIENT_BALANCE")
var ErrSlpCardNotFound = httperror.BadRequestError("card number is not found", "SLP_CARD_NOT_FOUND")

var ErrSealabspayTxnIdNotValid = httperror.BadRequestError("sealabspay transaction id is not valid", "SEALABSPAY_TXN_ID_NOT_VALID")
var ErrSealabspayAmountNotValid = httperror.BadRequestError("sealabspay amount is not valid", "SEALABSPAY_AMOUNT_NOT_VALID")
