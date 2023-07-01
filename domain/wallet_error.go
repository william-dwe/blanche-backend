package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrWalletNotFound = httperror.BadRequestError("cannot find wallet details, maybe user has not activated wallet yet", "WALLET_NOT_FOUND")
var ErrGetWallet = httperror.InternalServerError("cannot get wallet details")
var ErrCreateWalletPinDeficient = httperror.BadRequestError("cannot create wallet pin, ensure pin is provided with lenght of 6 digits", "WALLET_INSUFFICIENT_DATA")

var ErrCreateWalletUserIdNotFound = httperror.BadRequestError("cannot create new wallet pin, user ID not found", "WALLET_USER_ID_NOT_FOUND")
var ErrCreateWalletUserDuplicate = httperror.BadRequestError("cannot create wallet pin, user already has wallet pin", "WALLET_USER_ID_DUPLICATE")
var ErrCreateWalletPinHash = httperror.InternalServerError("cannot create wallet pin hash")
var ErrCreateWallet = httperror.InternalServerError("cannot create wallet pin")

var ErrUpdateWalletPin = httperror.InternalServerError("cannot update wallet pin")
var ErrUpdateWalletUserNotFound = httperror.BadRequestError("cannot update wallet pin, user has not activated wallet yet", "WALLET_USER_ID_NOT_FOUND")
var ErrUpdateWalletPinDeficient = httperror.BadRequestError("cannot update wallet pin, ensure new_pin is provided with lenght of 6 digits", "WALLET_INSUFFICIENT_DATA")
var ErrWalletTransactionNotFound = httperror.BadRequestError("cannot find wallet transaction, maybe user has not activated wallet yet", "WALLET_NOT_FOUND")
var ErrGetWalletTransaction = httperror.InternalServerError("cannot get wallet transaction")

var ErrUpdateWalletBalance = httperror.InternalServerError("cannot update wallet balance")
var ErrAddWalletTransaction = httperror.InternalServerError("cannot add wallet transaction")

var ErrTopUpWallet = httperror.InternalServerError("cannot top up wallet")
var ErrTopUpWalletPaymentRecord = httperror.InternalServerError("cannot update payment record for top up wallet")

var ErrCreateWalletPaymentRecord = httperror.InternalServerError("cannot create payment record for transactions")
var ErrWalletBalanceNotSufficient = httperror.BadRequestError("cannot process transaction, wallet balance is not sufficient", "WALLET_BALANCE_NOT_SUFFICIENT")
var ErrCreateWalletTransactionRecord = httperror.InternalServerError("cannot create wallet transaction record")

var ErrCreateWalletTransactionRecordWithdrawPaymentRec = httperror.InternalServerError("cannot create wallet transaction record for withdraw")
var ErrCreateWalletTransactionRecordWithdraw = httperror.InternalServerError("cannot create wallet transaction record for withdraw")
var ErrCreateWalletTransactionRecordWithdrawUpdateBalance = httperror.InternalServerError("cannot create wallet transaction record for withdraw in updating balance")

var ErrCreateWalletTransactionRecordRefundPaymentRec = httperror.InternalServerError("cannot create wallet transaction record for refund")
var ErrCreateWalletTransactionRecordRefundUpdateBalance = httperror.InternalServerError("cannot create wallet transaction record for refund in updating balance")
var ErrCreateWalletTransactionRecordRefund = httperror.InternalServerError("cannot create wallet transaction record for refund")
var ErrCreateWalletTransactionRecordRefundTransactionPaymentRec = httperror.InternalServerError("cannot create wallet transaction record for refund in updating transaction payment record")
