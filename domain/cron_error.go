package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrCronJobAddJobFailed = httperror.InternalServerError("cannot add cron job")
var ErrCronJobIsDisabled = httperror.InternalServerError("cron job is disabled")
