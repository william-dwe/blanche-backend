package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrAddFavoriteProduct = httperror.InternalServerError("cannot add favorite product")
var ErrAddFavoriteProductAlreadyExist = httperror.BadRequestError("cannot add favorite product, product already favorited", "DATA_DUPLICATE")

var ErrUpdateFavoriteUserNotFound = httperror.BadRequestError("cannot update user favorite product, user not found", "DATA_NOT_VALID")
var ErrUpdateFavoriteProductNotFound = httperror.BadRequestError("cannot update user favorite product, product not found", "DATA_NOT_VALID")
var ErrUpdateFavoriteProductInvalidInput = httperror.BadRequestError("cannot update favorite product, request not valid", "DATA_NOT_VALID")
var ErrUpdateFavoriteProduct = httperror.InternalServerError("cannot update favorite product")

var ErrGetFavoriteProducts = httperror.InternalServerError("cannot get favorite products")
var ErrGetFavoriteProduct = httperror.InternalServerError("cannot get favorite product")
