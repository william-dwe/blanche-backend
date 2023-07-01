package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) WalletpayPayReq(c *gin.Context) {
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputWalletPay dto.WalletpayReqDTO
	err = util.ShouldBindJsonWithValidation(c, &inputWalletPay)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.walletpayUsecase.HandleWalletpaySuccessRequest(userJwt.Username, inputWalletPay)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_PAY_WALLET_RESPONSE",
		Message: "Success pay wallet response",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) WalletpayCancelPayReq(c *gin.Context) {
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputWalletPay dto.WalletpayReqDTO
	err = util.ShouldBindJsonWithValidation(c, &inputWalletPay)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.walletpayUsecase.HandleWalletpayCancelRequest(userJwt.Username, inputWalletPay)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CANCEL_WALLET_RESPONSE",
		Message: "Success cancel wallet response",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
