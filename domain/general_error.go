package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrRequestBodyJSONInvalid = httperror.BadRequestError("request body is not valid", "DATA_NOT_VALID")
var ErrConvertQueryParams = httperror.BadRequestError("please check again your inputs", "DATA_NOT_VALID")
var ErrRequestBodyInvalid = httperror.BadRequestError("request body is not valid", "DATA_NOT_VALID")
var ErrUnmarshalJSON = httperror.InternalServerError("failed to unmarshal json")
var ErrInvalidVoucherCode = httperror.BadRequestError("invalid voucher code", "INVALID_VOUCHER_CODE")
var ErrInvalidVoucherCodeLength = httperror.BadRequestError("code length maximum is 5 characters", "INVALID_VOUCHER_CODE_LENGTH")
var ErrInvalidVoucherCodeCapitalize = httperror.BadRequestError("voucher code must be in uppercase", "INVALID_VOUCHER_CODE_CAPITALIZE")
