package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserSlpAccountList(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slpAccounts, err := h.slpAccountUsecase.GetSlpAccountListByUsername(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_SLP_ACCOUNTS",
		Message: "Success retrieve SLP accounts",
		Data:    slpAccounts,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetUserSlpAccount(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slpAccountID, err := strconv.Atoi(c.Param("slp_id"))
	if err != nil {
		_ = c.Error(domain.ErrInvalidSlpAccountID)
		return
	}

	slpAccount, err := h.slpAccountUsecase.GetUserSlpAccount(user.Username, slpAccountID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_SLP_ACCOUNT",
		Message: "Success retrieve SLP account",
		Data:    slpAccount,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) RegisterUserSlpAccount(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.SlpAccountReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	slpAccount, err := h.slpAccountUsecase.RegisterSlpAccount(user.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_REGISTER_SLP_ACCOUNT",
		Message: "Success register SLP account",
		Data:    slpAccount,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteUserSlpAccount(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slpAccountID := c.Param("slp_id")
	slpAccountIDInt, err := strconv.Atoi(slpAccountID)
	if err != nil {
		_ = c.Error(domain.ErrInvalidSlpAccountID)
		return
	}

	resBody, err := h.slpAccountUsecase.DeleteSlpAccount(user.Username, slpAccountIDInt)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_SLP_ACCOUNT",
		Message: "Success delete SLP account",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) SetDefaultSlpAccount(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	slpAccountID := c.Param("slp_id")
	slpAccountIDInt, err := strconv.Atoi(slpAccountID)
	if err != nil {
		_ = c.Error(domain.ErrInvalidSlpAccountID)
		return
	}

	resBody, err := h.slpAccountUsecase.SetDefaultSlpAccount(user.Username, slpAccountIDInt)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_SET_DEFAULT_SLP_ACCOUNT",
		Message: "Success set default SLP account",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
