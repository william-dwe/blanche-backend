package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateWallet(c *gin.Context) {
	var newWalletDTO dto.CreateWalletPinReq
	err := c.ShouldBindJSON(&newWalletDTO)
	if err != nil {
		_ = c.Error(domain.ErrCreateWalletPinDeficient)
		return
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	newWalletDTO.Username = userJWT.Username

	resCreateWallet, err := h.walletUsecase.CreateWalletPin(newWalletDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_WALLET_PIN",
		Message: "Success create wallet pin",
		Data:    resCreateWallet,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetWalletDetails(c *gin.Context) {
	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resWallet, err := h.walletUsecase.GetWalletDetails(userJWT.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_WALLET_DETAILS",
		Message: "Success get wallet details",
		Data:    resWallet,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetWalletTransactions(c *gin.Context) {
	var walletTransactionReqParamDTO dto.WalletTransactionReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &walletTransactionReqParamDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resWalletTransactions, err := h.walletUsecase.GetWalletTransactions(userJWT.Username, walletTransactionReqParamDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_WALLET_TRANSACTIONS",
		Message: "Success get wallet transactions",
		Data:    resWalletTransactions,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateWalletPin(c *gin.Context) {
	var updateWalletPinDTO dto.WalletUpdatePin
	err := c.ShouldBindJSON(&updateWalletPinDTO)
	if err != nil {
		_ = c.Error(domain.ErrUpdateWalletPinDeficient)
		return
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	updateWalletPinDTO.Username = userJWT.Username
	resWallet, err := h.walletUsecase.UpdateWalletPin(updateWalletPinDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_WALLET_PIN",
		Message: "Success update wallet pin",
		Data:    resWallet,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MakeTopUpWalletSlp(c *gin.Context) {
	var topUpWalletDTO dto.TopUpWalletUsingSlpReqDTO
	err := util.ShouldBindJsonWithValidation(c, &topUpWalletDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resWallet, err := h.walletUsecase.MakeTopUpWalletUsingSlpReq(userJWT.Username, topUpWalletDTO)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MAKE_TOP_UP_WALLET_SLP_REQ",
		Message: "Success make top up wallet slp request",
		Data:    resWallet,
	}

	util.ResponseSuccessJSON(c, response)
}
