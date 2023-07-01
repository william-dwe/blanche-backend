package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrCheckUserInternalServer = httperror.InternalServerError("Unable to check user")
var ErrInvalidPin = httperror.BadRequestError("Invalid pin", "INVALID_PIN")
var ErrInvalidPass = httperror.BadRequestError("Invalid password", "INVALID_PASSWORD")
var ErrBlacklistToken = httperror.InternalServerError("Token couldn't be blacklisted")
var ErrInvalidOAuthState = httperror.BadRequestError("OAuth request is expired", "INVALID_OAUTH_STATE")
var ErrForbiddenNotAdmin = httperror.ForbiddenError()
var ErrForbiddenAdmin = httperror.ForbiddenError()

var ErrChangePasswordInternalError = httperror.InternalServerError("Unable to reset password")
var ErrResetPasswordWait = httperror.BadRequestError("Please wait for 1 minute before requesting another verification code", "WAIT_BEFORE_REQUESTING_OTP")
var ErrInvalidVerificationCode = httperror.BadRequestError("Invalid verification code", "INVALID_VERIFICATION_CODE")
var ErrSamePassword = httperror.BadRequestError("New password cannot be the same as the old one", "SAME_PASSWORD")
var ErrResetPasswordCodeExpired = httperror.BadRequestError("Verification code has expired", "VERIFICATION_CODE_EXPIRED")

var ErrGetUserDataFromGoogleInternalError = httperror.InternalServerError("Unable to get user data from Google. Please try again later.")
