package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrAddMessageRefundRequest = httperror.InternalServerError("Failed to add message to refund request")
var ErrGetMessageRefundRequest = httperror.InternalServerError("Failed to get messages from refund request")
var ErrRefundRequestClosed = httperror.BadRequestError("Refund request already closed", "REFUND_REQUEST_CLOSED")

var ErrInvalidRefundId = httperror.BadRequestError("Invalid refund id", "INVALID_REFUND_ID")

var ErrUserWalletNotActivated = httperror.BadRequestError("User wallet not activated", "USER_WALLET_NOT_ACTIVATED")
var ErrCreateRefundRequest = httperror.InternalServerError("Failed to create refund request")
var ErrCreateRefundRequestDuplicate = httperror.BadRequestError("Refund request already exists", "REFUND_REQUEST_ALREADY_EXISTS")
var ErrCreateRefundRequestUpdateStatus = httperror.InternalServerError("Failed to update transaction status to request refund")
var ErrAddRefundRequestMessage = httperror.InternalServerError("Failed to add refund request message")

var ErrRefundTransactionNotEligible = httperror.BadRequestError("Transaction is not eligible for refund", "REFUND_TRANSACTION_NOT_ELIGIBLE")
var ErrRefundTransactionAlreadyRequested = httperror.BadRequestError("Transaction already requested for refund", "REFUND_TRANSACTION_ALREADY_REQUESTED")

var ErrGetRefundRequestList = httperror.InternalServerError("Failed to get refund request list")
var ErrGetRefundRequestNotFound = httperror.BadRequestError("Refund request not found", "REFUND_REQUEST_NOT_FOUND")
var ErrGetRefundRequest = httperror.InternalServerError("Failed to get refund request list by merchant id")

var ErrRefundRequestAlreadyCanceledOrProcessed = httperror.BadRequestError("Refund request already canceled or processed", "REFUND_REQUEST_ALREADY_CANCELED_OR_PROCESSED")
var ErrRefundRequestNotYetProcessedBySeller = httperror.BadRequestError("Refund request not yet processed by seller", "REFUND_REQUEST_NOT_YET_PROCESSED_BY_SELLER")
var ErrCancelRefundRequest = httperror.InternalServerError("Failed to cancel refund request")
var ErrCancelRefundRequestUpdateTransactionStatus = httperror.InternalServerError("Failed to update transaction status to cancel refund request")

var ErrMerchantAcceptRefundRequest = httperror.InternalServerError("Failed to accept refund request by merchant")
var ErrMerchantRejectRefundRequest = httperror.InternalServerError("Failed to reject refund request by merchant")
var ErrAdminAcceptRefundRequest = httperror.InternalServerError("Failed to accept refund request by admin")
var ErrAdminRejectRefundRequest = httperror.InternalServerError("Failed to reject refund request by admin")
var ErrUserAcceptRefundRequest = httperror.InternalServerError("Failed to accept refund request by user")
var ErrUserRejectRefundRequest = httperror.InternalServerError("Failed to reject refund request by user")
var ErrUpdateAllRefundRequestStatus = httperror.InternalServerError("Failed to update all refund request status")
var ErrUserAcceptRefundRequestUpdateRefundRequestStatus = httperror.InternalServerError("Failed to update refund request status to accept refund request by user")
var ErrUserAcceptRefundRequestUpdateTransactionStatus = httperror.InternalServerError("Failed to update transaction status to accept refund request by user")
var ErrUserAcceptRefundRequestCommitTransaction = httperror.InternalServerError("Failed to commit transaction to accept refund request by user")
var ErrUserRejectRefundRequestNewRefundRequestStatus = httperror.InternalServerError("Failed to create new refund request status to reject refund request by user")
var ErrUserRejectRefundRequestUpdateRefundRequestStatus = httperror.InternalServerError("Failed to update refund request status to reject refund request by user")
var ErrUserRejectRefundRequestCommit = httperror.InternalServerError("Failed to commit transaction to reject refund request by user")
var ErrRefundRequestUserAlreadyAcceptedOrRejected = httperror.BadRequestError("Refund request already accepted or rejected by user", "REFUND_REQUEST_ALREADY_ACCEPTED_OR_REJECTED_BY_USER")
var ErrRefundRequestNotYetProcessedByAdmin = httperror.BadRequestError("Refund request not yet processed by admin", "REFUND_REQUEST_NOT_YET_PROCESSED_BY_ADMIN")
var ErrRefundRequestUserAlreadyRejectedThreeTimes = httperror.BadRequestError("Refund request already rejected by user three times, refund request has been dismissed", "REFUND_REQUEST_ALREADY_REJECTED_BY_USER_THREE_TIMES")
var ErrAdminRejectRefundRequestUpdateTransactionStatus = httperror.InternalServerError("Failed to update transaction status to reject refund request by admin")
var ErrAdminRejectRefundRequestCommit = httperror.InternalServerError("Failed to commit transaction to reject refund request by admin")

var ErrCronRefundRequestStatusToAcceptedBySeller = httperror.InternalServerError("Failed to update refund request status to accepted by seller")
var ErrCronRefundRequestStatusToAcceptedByBuyer = httperror.InternalServerError("Failed to update refund request status to accepted by buyer")
