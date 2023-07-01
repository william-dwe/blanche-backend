package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllDeliveryOption(c *gin.Context) {
	resBody, err := h.deliveryUsecase.GetAllDeliveryOption()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_DELIVERY_OPTIONS",
		Message: "Success get all delivery options",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetDeliveryOptionByMerchantDomain(c *gin.Context) {
	resBody, err := h.deliveryUsecase.GetDeliveryOptionByMerchantDomain(c.Param("domain"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_DELIVERY_OPTIONS",
		Message: "Success get all delivery options",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantUserDeliveryOption(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.deliveryUsecase.GetMerchantDeliveryOption(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_DELIVERY_OPTIONS",
		Message: "Success get all delivery options",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ChangeMerchantUserDeliveryOption(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest []dto.DeliveryUpdateMerchantOptionReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.deliveryUsecase.UpdateMerchantDeliveryOption(user.Username, inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_DELIVERY_OPTIONS",
		Message: "Success update delivery options",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
