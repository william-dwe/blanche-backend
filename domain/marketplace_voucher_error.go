package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetMarketplaceVoucherList = httperror.InternalServerError("failed to get marketplace voucher list")
var ErrGetMarketplaceVoucher = httperror.InternalServerError("failed to get marketplace voucher record")
var ErrInvalidVoucherID = httperror.BadRequestError("invalid voucher ID", "INVALID_VOUCHER_ID")
var ErrMarketplaceVoucherNotFound = httperror.BadRequestError("marketplace voucher not found", "MARKETPLACE_VOUCHER_NOT_FOUND")

var ErrDecreaseMarketplaceVoucherQuota = httperror.InternalServerError("failed to decrease marketplace voucher quota")
var ErrIncreaseMarketplaceVoucherQuota = httperror.InternalServerError("failed to increase marketplace voucher quota")
var ErrCreateMarketplaceVoucher = httperror.InternalServerError("failed to create marketplace voucher")
var ErrDeleteMarketplaceVoucher = httperror.InternalServerError("failed to delete marketplace voucher")
var ErrUpdateMarketplaceVoucher = httperror.InternalServerError("failed to update marketplace voucher")
var ErrMarketplaceVoucherCodeAlreadyExist = httperror.BadRequestError("voucher code already exist", "VOUCHER_CODE_ALREADY_EXIST")
