package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllPaymentMethod(c *gin.Context) {
	paymentMethods, err := h.paymentMethodUsecase.GetAll()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_PAYMENT_METHOD",
		Message: "Success retrieve all payment method",
		Data:    paymentMethods,
	}

	util.ResponseSuccessJSON(c, response)
}
