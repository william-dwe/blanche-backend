package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrCreateUserOrder = httperror.InternalServerError("Failed to create user order")
