package util

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

func ValidateFullname(fullname string) error {
	if len(fullname) < 2 || len(fullname) > 32 {
		return domain.ErrInvalidFullnameLength
	}

	return nil
}
