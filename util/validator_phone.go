package util

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

func ValidatePhone(phone string) error {
	if len(phone) < 11 || len(phone) > 15 {
		return domain.ErrInvalidPhoneLength
	}
	if !strings.HasPrefix(phone, "62") {
		return domain.ErrInvalidPhoneFalsePrefix
	}

	return nil
}
