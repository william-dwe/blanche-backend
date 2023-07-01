package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetCategories = httperror.InternalServerError("cannot get categories tree record")

var ErrCategoryQueryParam = httperror.BadRequestError("category query param is invalid", "INVALID_PARAMS")
var ErrCreateCategory = httperror.InternalServerError("cannot create category record")
var ErrUpdateCategory = httperror.InternalServerError("cannot update category record")
var ErrDeleteCategory = httperror.InternalServerError("cannot delete category record")
var ErrCategoryInUse = httperror.BadRequestError("category is in use", "CATEGORY_IN_USE")
var ErrCategoryNotAuthorized = httperror.UnauthorizedError()
var ErrCategoryExceedLimit = httperror.BadRequestError("can't create category, category exceed level limit", "CATEGORY_EXCEED_LIMIT")
var ErrCategorySlugAlreadyExist = httperror.BadRequestError("this category name is already exist", "CATEGORY_NAME_ALREADY_EXIST")
