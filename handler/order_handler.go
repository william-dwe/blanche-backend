package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrderSummary(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.PostOrderSummaryReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.orderItemUsecase.GetOrderCheckoutSummary(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ORDER_SUMMARY",
		Message: "Success get order summary",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MakeOrderCheckout(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody []dto.MakeOrderCheckoutProductDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.orderItemUsecase.MakeOrderCheckout(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MAKE_ORDER_CHECKOUT",
		Message: "Success make order checkout",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetWaitingForPayments(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentRecordUsecase.GetWaitingForPayment(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_WAITING_FOR_PAYMENT",
		Message: "Success get waiting for payment",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetWaitingForPaymentDetails(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	paymentId := c.Param("payment_id")

	res, err := h.paymentRecordUsecase.GetWaitingForPaymentDetail(user.Username, paymentId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_WAITING_FOR_PAYMENT_DETAILS",
		Message: "Success get waiting for payment details",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
