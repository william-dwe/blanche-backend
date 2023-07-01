package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SealabspayResponseHandler(c *gin.Context) {
	var inputResponseSlp dto.SealabspayReqDTO
	err := util.ShouldBindJsonWithValidation(c, &inputResponseSlp)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.sealabspayUsecase.HandlePaymentResponse(inputResponseSlp)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_SEALABSPAY_RESPONSE",
		Message: "Success sealabspay response",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
