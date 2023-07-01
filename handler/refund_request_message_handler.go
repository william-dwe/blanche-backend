package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AdminAddMessageRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	var reqBody dto.RefundRequestMsgFormReqDTO
	err = util.ShouldBindWithValidation(c, &reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestMessageUsecase.AdminAddMessage(uint(refundId), reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADMIN_ADD_MESSAGE_REFUND_REQUEST",
		Message: "Success admin add message refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MerchantAddMessageRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.RefundRequestMsgFormReqDTO
	err = util.ShouldBindWithValidation(c, &reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestMessageUsecase.MerchantAddMessage(userJwt.Username, uint(refundId), reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MERCHANT_ADD_MESSAGE_REFUND_REQUEST",
		Message: "Success merchant add message refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) BuyerAddMessageRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.RefundRequestMsgFormReqDTO
	err = util.ShouldBindWithValidation(c, &reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestMessageUsecase.BuyerAddMessage(userJwt.Username, uint(refundId), reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_BUYER_ADD_MESSAGE_REFUND_REQUEST",
		Message: "Success buyer add message refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AdminGetMessageRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	resBody, err := h.refundRequestMessageUsecase.GetAdminListMessageByRefundRequestId(uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADMIN_GET_MESSAGE_REFUND_REQUEST",
		Message: "Success admin get message refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserGetMessageRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestMessageUsecase.GetListMessageByRefundRequestId(userJwt.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_USER_GET_MESSAGE_REFUND_REQUEST",
		Message: "Success user get message refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
