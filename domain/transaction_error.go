package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrTransactionNotFound = httperror.NotFoundError("transaction not found")
var ErrTransactionsNotFound = httperror.NotFoundError("transactions not found")
var ErrGetTransactions = httperror.InternalServerError("failed to get transactions")
var ErrGetTransaction = httperror.InternalServerError("failed to get transaction")

var ErrTransactionStatusNotFound = httperror.NotFoundError("transaction status not found")
var ErrGetTransactionStatus = httperror.InternalServerError("failed to get transaction status")

var ErrTransactionDeliveryStatusNotFound = httperror.NotFoundError("transaction delivery status not found")
var ErrGetTransactionDeliveryStatus = httperror.InternalServerError("failed to get transaction delivery status")
var ErrUnmarshalJSONCartItem = httperror.InternalServerError("failed to unmarshal json cart item")
var ErrUnmarshalJSONDeliveryOption = httperror.InternalServerError("failed to unmarshal json delivery option")
var ErrUnmarshalJSONPaymentMethod = httperror.InternalServerError("failed to unmarshal json payment method")
var ErrUnmarshalJSONTransactionAddress = httperror.InternalServerError("failed to unmarshal json transaction address")

var ErrCreateTransaction = httperror.InternalServerError("failed to create transaction invoice")
var ErrMakeTransactionNotValid = httperror.BadRequestError("transaction not valid", "TRANSACTION_NOT_VALID")
var ErrPaymentTotalNotMatch = httperror.BadRequestError("payment total not match", "PAYMENT_TOTAL_NOT_MATCH")
var ErrCreateTransactionPayment = httperror.InternalServerError("failed to create transaction payment")
var ErrCreateTransactionWalletNotSupported = httperror.BadRequestError("wallet payment method not supported", "WALLET_PAYMENT_METHOD_NOT_SUPPORTED")
var ErrMakeTransaction = httperror.InternalServerError("failed to make transaction, maybe some data has been changed. please try again")
var ErrInvalidTransactionId = httperror.BadRequestError("invalid transaction id", "INVALID_TRANSACTION_ID")

var ErrUpdateTransactionPayment = httperror.InternalServerError("failed to update transaction payment")
var ErrUpdateTransactionStatusPayment = httperror.InternalServerError("failed to update transaction status after payment")

var ErrUpdateMerchantTransactionStatus = httperror.InternalServerError("failed to update merchant transaction status")
var ErrUpdateMerchantTransactionForbidden = httperror.ForbiddenErrorMsg("failed to update merchant transaction status, merchant not allowed to update this transaction")
var ErrUpdateTransactionStatus = httperror.InternalServerError("failed to update transaction status")
var ErrUpdateTransactionDeliveryStatus = httperror.InternalServerError("failed to update transaction delivery status")
var ErrUpdateTransactionStatusCannotReverse = httperror.BadRequestError("failed to update transaction status, cannot reverse transaction status", "CANNOT_REVERSE_TRANSACTION_STATUS")
var ErrUpdateTransactionStatusCannotSkip = httperror.BadRequestError("failed to update transaction status, cannot skip transaction status", "CANNOT_SKIP_TRANSACTION_STATUS")
var ErrUpdateTransactionStatusReceiptNumberEmpty = httperror.BadRequestError("failed to update transaction status to on delivery, receipt number cannot be empty", "RECEIPT_NUMBER_EMPTY")

var ErrUpdateTransactionStatusToCancel = httperror.InternalServerError("failed to update transaction status to canceled")
var ErrUpdateTransactionStatusToCompleted = httperror.InternalServerError("failed to update transaction status to completed")
var ErrUpdateTransactionStatusToCompletedFundActivities = httperror.InternalServerError("failed to update transaction status to completed, failed to update fund activities")
var ErrUnmarshalJSONPaymentDetails = httperror.InternalServerError("failed to to get transaction payment details data structure")
var ErrUnmarshalJSONCartItems = httperror.InternalServerError("failed to to get transaction cart items data structure")
var ErrInvalidTransactionStatusForbidden = httperror.ForbiddenErrorMsg("invalid transaction status, not eligible to update transaction status.")

var ErrUpdateTransactionStatusToRefunded = httperror.InternalServerError("failed to update transaction status to refunded")
var ErrUpdateTransactionStatusToRefundedDeductMpBalance = httperror.InternalServerError("failed to update transaction status to refunded, failed to deduct balance")
var ErrUpdateTransactionStatusToRefundedAddUserBalance = httperror.InternalServerError("failed to update transaction status to refunded, failed to add balance")
var ErrUpdateTransactionStatusToRefundedMarketplaceVoucher = httperror.InternalServerError("failed to update transaction status to refunded, failed to update marketplace voucher")
var ErrUpdateTransactionStatusToRefundedPendingTransaction = httperror.InternalServerError("failed to update transaction status to refunded, failed to update pending transaction")
var ErrAdminAcceptRefundRequestUpdateTransactionStatus = httperror.InternalServerError("failed to update transaction status to refunded, failed to update transaction status")
var ErrAdminAcceptRefundRequestCommit = httperror.InternalServerError("failed to update transaction status to refunded, failed to commit transaction")

var ErrUpdateTransactionStatusToRefundedAddWalletHistory = httperror.InternalServerError("failed to update transaction status to refunded, failed to add wallet history")

var ErrCronUpdateTransactionWaitingStatusToCanceled = httperror.InternalServerError("failed to update transaction status to canceled by cron")
var ErrCronUpdateTransactionStatusToCanceled = httperror.InternalServerError("failed to update transaction status to canceled by cron")
var ErrCronUpdateTransactionStatusToCompleted = httperror.InternalServerError("failed to update transaction status to completed by cron")
