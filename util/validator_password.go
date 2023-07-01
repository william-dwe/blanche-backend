package util

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

func ValidatePassword(password, username string) error {
	if len(password) < 8 || len(password) > 32 {
		return domain.ErrInvalidPasswordLength
	}
	if stringContainInsensitive(password, username) {
		return domain.ErrInvalidPasswordContainUsername
	}
	if stringContainLeadingOrTrailingSpaces(password) {
		return domain.ErrInvalidPasswordContainLeadingOrTrailingSpaces
	}

	return nil
}

func stringContainInsensitive(str, subStr string) bool {
	str = strings.ToLower(str)
	subStr = strings.ToLower(subStr)
	return strings.Contains(str, subStr)
}

func stringContainLeadingOrTrailingSpaces(str string) bool {
	return strings.HasPrefix(str, " ") || strings.HasSuffix(str, " ")
}
