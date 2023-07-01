package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type SealabspayRepository interface {
	MakePayment(cardNumber string, amount uint) (redirectUrl string, paymentId string, slpError error)
	MakePaymentCustomRedirect(cardNumber string, amount uint, redirectPaymentUrl string) (redirectUrl string, paymentId string, slpError error)
}

type SealabspayRepositoryConfig struct {
}

type SealabspayRepositoryImpl struct {
}

func NewSealabspayRepository(c SealabspayRepositoryConfig) SealabspayRepository {
	return &SealabspayRepositoryImpl{}
}

func (r *SealabspayRepositoryImpl) MakePaymentCustomRedirect(cardNumber string, amount uint, redirectPaymentUrl string) (redirectUrl string, paymentId string, slpError error) {
	apiUrl := config.Config.SeaLabsPayConfig.Url + "/transaction/pay"
	merchantCode := config.Config.SeaLabsPayConfig.MerchantCode
	signature, err := util.GenerateReqSlpSignature(cardNumber, amount, merchantCode)
	if err != nil {
		return "", "", domain.ErrSlpCannotGenerateSignature
	}

	data := url.Values{}
	data.Set("card_number", cardNumber)
	data.Set("amount", strconv.FormatInt(int64(amount), 10))
	data.Set("merchant_code", merchantCode)
	data.Set("redirect_url", redirectPaymentUrl)
	data.Set("callback_url", config.Config.SeaLabsPayConfig.CallbackUrl)
	data.Set("signature", signature)

	urlStr := apiUrl

	client := &http.Client{}
	rHttp, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	rHttp.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(rHttp)
	slpErrDTO := dto.SealabspayResDTO{}
	if err == nil {
		json.NewDecoder(response.Body).Decode(&slpErrDTO)
		if slpErrDTO.Code == "user:insufficient-fund" {
			return "", "", domain.ErrSlpUserInsufficientFund
		}
		if slpErrDTO.Code == "user:not-found" {
			return "", "", domain.ErrSlpCardNotFound
		}

		return "", "", domain.ErrSlpRedirect
	}

	redirectUrl = strings.Split(err.Error(), `"`)[1]

	queryUrl, err := url.Parse(redirectUrl)
	if err != nil {
		return "", "", domain.ErrSlpParsePaymentId
	}

	q, err := url.ParseQuery(queryUrl.RawQuery)
	if err != nil {
		return "", "", domain.ErrSlpParsePaymentId
	}

	txnId, err := strconv.ParseUint(q.Get("txn_id"), 10, 64)
	if err != nil {
		return "", "", domain.ErrSlpParsePaymentId
	}

	return redirectUrl, fmt.Sprintf("SLP%d", txnId), nil
}

func (r *SealabspayRepositoryImpl) MakePayment(cardNumber string, amount uint) (redirectUrl string, paymentId string, slpError error) {
	return r.MakePaymentCustomRedirect(cardNumber, amount, config.Config.SeaLabsPayConfig.RedirectUrl)
}
