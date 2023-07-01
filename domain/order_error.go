package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrNoOrderItem = httperror.BadRequestError("no order item", "NO_ORDER_ITEM")
var ErrOrderProductNotAvailable = httperror.BadRequestError("some product is not available", "PRODUCT_NOT_AVAILABLE")
var ErrOrderQuantityNotValid = httperror.BadRequestError("some product quantity is not valid", "PRODUCT_QUANTITY_NOT_VALID")
var ErrOrderButOwnProduct = httperror.BadRequestError("you can't buy your own product", "ORDER_OWN_PRODUCT")

var ErrGetOrderNotFound = httperror.BadRequestError("order not found", "ORDER_NOT_FOUND")
var ErrGetOrderSummary = httperror.InternalServerError("failed to get order summary")
var ErrOrderAddressNotFound = httperror.BadRequestError("default address not found, please add new address", "ORDER_ADDRESS_NOT_FOUND")
