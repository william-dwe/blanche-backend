package util

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
)

func ValidateVoucherCode(voucherCode, merchantDomain string) error {
	res := strings.HasPrefix(strings.ToLower(voucherCode), strings.ToLower(merchantDomain))
	if !res {
		return domain.ErrInvalidVoucherCode
	}

	if strings.ToUpper(voucherCode) != voucherCode {
		return domain.ErrInvalidVoucherCodeCapitalize
	}

	merchantDomainLen := len(merchantDomain)
	voucherCodeLen := len(voucherCode)
	suffixLen := voucherCodeLen - merchantDomainLen
	if suffixLen < 3 || suffixLen > 5 {
		return domain.ErrInvalidVoucherCode
	}

	return nil
}
