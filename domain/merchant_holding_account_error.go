package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetMerchantFundActivities = httperror.InternalServerError("Failed to get merchant fund activities")
var ErrMerchantHoldingAccountNotFound = httperror.BadRequestError("Merchant holding account not found", "MERCHANT_HOLDING_ACCOUNT_NOT_FOUND")
var ErrMerchantHoldingAccount = httperror.InternalServerError("Failed to get merchant holding account")
var ErrMerchantHoldingAccInsufficientBalance = httperror.BadRequestError("Merchant holding account balance is insufficient", "MERCHANT_HOLDING_ACCOUNT_INSUFFICIENT_BALANCE")
var ErrMerchantHoldingAccWithdraw = httperror.InternalServerError("Failed to withdraw merchant holding account")
var ErrMerchantHoldingAccWithdrawAddHistory = httperror.InternalServerError("Failed to add withdraw history")
var ErrMerchantHoldingAccWithdrawUpdateWallet = httperror.InternalServerError("Failed to update wallet account after withdraw")
