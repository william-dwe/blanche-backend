package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrAddCartItemInvalidInput = httperror.BadRequestError("invalid input, should provide product_id, variant_item_id, and quantity", "INVALID_INPUT")
var ErrAddCartItemNeedVariantItem = httperror.BadRequestError("please specify the variant_item_id", "NEED_VARIANT_ITEM")
var ErrAddCartProductNotExist = httperror.BadRequestError("product not exist", "PRODUCT_NOT_EXIST")
var ErrAddCartUserNotExist = httperror.BadRequestError("user not exist", "USER_NOT_EXIST")
var ErrAddCartVariantItemNotExist = httperror.BadRequestError("variant item not exist", "VARIANT_ITEM_NOT_EXIST")
var ErrAddCartVariantItemNotEnoughStock = httperror.BadRequestError("variant item not enough stock", "VARIANT_ITEM_NOT_ENOUGH_STOCK")
var ErrAddCartProductNotAvailable = httperror.BadRequestError("product not available", "PRODUCT_NOT_AVAILABLE")

var ErrAddCartItemAddOwnProduct = httperror.BadRequestError("cannot add your own product to cart", "ADD_OWN_PRODUCT")
var ErrAddCartItemQuantityExceedStock = httperror.BadRequestError("quantity exceed stock", "QUANTITY_EXCEED_STOCK")
var ErrAddCartInternalError = httperror.InternalServerError("failed to add cart item")
var ErrGetCartInternalError = httperror.InternalServerError("failed to get cart item")
var ErrDeleteCartInternalError = httperror.InternalServerError("failed to delete cart item")

var ErrInvalidCartItemID = httperror.BadRequestError("invalid cart item id", "INVALID_CART_ITEM_ID")
var ErrInvalidSelectedCartItemId = httperror.BadRequestError("invalid selected cart item id", "INVALID_SELECTED_CART_ITEM_ID")

var ErrUpdateCartItemInvalidInput = httperror.BadRequestError("invalid input, should provide quantity", "INVALID_INPUT")
var ErrUpdateCartVariantItemNotExist = httperror.InternalServerError("variant item not exist")
var ErrUpdateCartItemQuantityExceedStock = httperror.BadRequestError("quantity exceed stock", "QUANTITY_EXCEED_STOCK")
