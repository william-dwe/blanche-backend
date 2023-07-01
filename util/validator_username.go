package util

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

func ValidateUsername(username string) error {
	if len(username) < 8 || len(username) > 16 {
		return domain.ErrInvalidUsernameLength
	}

	return nil
}
