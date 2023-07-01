package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AddRefundRequest(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.RefundRequestFormReqDTO
	if err := util.ShouldBindWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.RequestRefundProcess(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADD_REFUND_REQUEST",
		Message: "Success add refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetAdminRefundRequestList(c *gin.Context) {
	var reqParam dto.RefundRequestListReqParamDTO
	if err := util.ShouldBindQueryWithValidation(c, &reqParam); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.GetAdminRefundRequestList(reqParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ADMIN_REFUND_REQUEST_LIST",
		Message: "Success get admin refund request list",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetUserRefundRequestList(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqParam dto.RefundRequestListReqParamDTO
	if err := util.ShouldBindQueryWithValidation(c, &reqParam); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.GetUserRefundRequestList(user.Username, reqParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_USER_REFUND_REQUEST_LIST",
		Message: "Success get user refund request list",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantRefundRequestList(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqParam dto.RefundRequestListReqParamDTO
	if err := util.ShouldBindQueryWithValidation(c, &reqParam); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.GetMerchantRefundRequestList(user.Username, reqParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_REFUND_REQUEST_LIST",
		Message: "Success get merchant refund request list",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserAcceptRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.UserAcceptRefundProcess(user.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_USER_ACCEPT_REFUND_REQUEST",
		Message: "Success user accept refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserRejectRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.UserRejectRefundProcess(user.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_USER_REJECT_REFUND_REQUEST",
		Message: "Success user reject refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserCancelRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.UserCancelRefundProcess(user.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_USER_CANCEL_REFUND_REQUEST",
		Message: "Success user cancel refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MerchantAcceptRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.MerchantAcceptRefundProsess(user.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MERCHANT_ACCEPT_REFUND_REQUEST",
		Message: "Success merchant accept refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MerchantRejectRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.refundRequestUsecase.MerchantRejectRefundProsess(user.Username, uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MERCHANT_REJECT_REFUND_REQUEST",
		Message: "Success merchant reject refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AdminAcceptRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	resBody, err := h.refundRequestUsecase.AdminAcceptRefundProcess(uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADMIN_ACCEPT_REFUND_REQUEST",
		Message: "Success admin accept refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AdminRejectRequestRefund(c *gin.Context) {
	idStr := c.Param("refund_id")
	refundId, err := strconv.Atoi(idStr)
	if err != nil {
		_ = c.Error(domain.ErrInvalidRefundId)
		return
	}

	resBody, err := h.refundRequestUsecase.AdminRejectRefundProcess(uint(refundId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADMIN_REJECT_REFUND_REQUEST",
		Message: "Success admin reject refund request",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
