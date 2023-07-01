package util

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
)

func GenerateReqSlpSignature(cardNumber string, amount uint, merchantCode string) (string, error) {
	apiKey := config.Config.SeaLabsPayConfig.ApiKey
	signature, err := HashSHA256(fmt.Sprintf("%s:%d:%s", cardNumber, amount, merchantCode), apiKey)
	if err != nil {
		return "", err
	}

	return signature, nil
}

func ValidateResSlpSignature(input dto.SealabspayReqDTO) error {
	apiKey := config.Config.SeaLabsPayConfig.ApiKey
	signature, err := HashSHA256(fmt.Sprintf("%s:%s:%s:%s:%s", input.TxnId, input.Amount, input.MerchantCode, input.Status, input.Message), apiKey)
	if err != nil {
		return domain.ErrSlpCannotGenerateSignature
	}

	if signature != input.Signature {
		return domain.ErrSlpInvalidSignature
	}

	return nil
}
