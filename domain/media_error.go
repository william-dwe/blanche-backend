package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrUploadFile = httperror.InternalServerError("failed to upload file")
var ErrDeleteFile = httperror.InternalServerError("failed to delete file")
