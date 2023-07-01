package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrExampleAuth = httperror.UnauthorizedError()

var ErrExampleFormatReq = httperror.BadRequestError("email and password are required and must be correct", "DATA_NOT_VALID")
var ErrExampleIdNotFound = httperror.BadRequestError("example ID not found", "DATA_NOT_FOUND")

var ErrCreateExample = httperror.InternalServerError("cannot create example record")
var ErrGetExample = httperror.InternalServerError("cannot get example record")
var ErrExampleHash = httperror.InternalServerError("cannot generate encrypted data for user")
var ErrExampleUnexpected = httperror.InternalServerError("unexpected error occured in example process")
