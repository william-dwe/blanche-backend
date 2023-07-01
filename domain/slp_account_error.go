package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrSlpAccountsNotFound = httperror.BadRequestError("SLP Accounts not found", "SLP_ACCOUNTS_NOT_FOUND")
var ErrGetSlpAccounts = httperror.InternalServerError("Failed to retrieve SLP Accounts")
var ErrRegisterSlpAccount = httperror.InternalServerError("Failed to register SLP Account")
var ErrInvalidCardNumber = httperror.BadRequestError("card number must be in correct format", "INVALID_CARD_NUMBER")
var ErrInvalidActiveDate = httperror.BadRequestError("this card is no longer active", "INVALID_ACTIVE_DATE")
var ErrSlpAccountNotFound = httperror.BadRequestError("SLP Account not found", "SLP_ACCOUNT_NOT_FOUND")
var ErrGetSlpAccount = httperror.InternalServerError("Failed to retrieve SLP Account")
var ErrDeleteSlpAccount = httperror.InternalServerError("Failed to delete SLP Account")
var ErrUnauthorizedSlpAccount = httperror.BadRequestError("This card is not belong to this user", "INVALID_SLP_ACCOUNT")
var ErrInvalidSlpAccountID = httperror.BadRequestError("Invalid SLP Account ID", "INVALID_SLP_ACCOUNT_ID")
var ErrDeleteDefaultSlpAccount = httperror.BadRequestError("Cannot delete default SLP Account", "DELETE_DEFAULT_SLP_ACCOUNT")
var ErrDuplicateSlpAccount = httperror.BadRequestError("You already registered with this card", "DUPLICATE_SLP_ACCOUNT")
