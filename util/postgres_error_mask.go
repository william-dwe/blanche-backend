package util

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

func PgConsErrMasker(
	err error,
	consErr entity.ConstraintErrMaskerMap,
	finErr httperror.AppError,
) error {
	if len(consErr) != 0 {
		if errValue, ok := consErr[err.(*pgconn.PgError).ConstraintName]; ok {
			return errValue
		}
		return finErr
	}

	log.Error().Msgf("Unknown error: %s", err.Error())
	return finErr
}
