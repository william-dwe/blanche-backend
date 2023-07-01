package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMarketplaceVoucherList(c *gin.Context) {
	res, err := h.mpVoucherUsecase.GetMarketplaceVoucherList()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MARKETPLACE_VOUCHER_LIST",
		Message: "Success retrieve marketplace voucher list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMarketplaceVoucherDetails(c *gin.Context) {
	voucherCode := c.Param("voucher_code")
	res, err := h.mpVoucherUsecase.GetMarketplaceVoucherByCode(voucherCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MARKETPLACE_VOUCHER",
		Message: "Success retrieve marketplace voucher",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreateMarketplaceVoucher(c *gin.Context) {
	var req dto.UpsertMarketplaceVoucherReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.mpVoucherUsecase.CreateMarketplaceVoucher(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_MARKETPLACE_VOUCHER",
		Message: "Success create marketplace voucher",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMarketplaceAdminVoucherList(c *gin.Context) {
	var req dto.MerchantVoucherListParamReqDTO
	if err := util.ShouldBindQueryWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.mpVoucherUsecase.GetMarketplaceAdminVoucherList(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MARKETPLACE_ADMIN_VOUCHER_LIST",
		Message: "Success retrieve marketplace admin voucher list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMarketplaceVoucher(c *gin.Context) {
	var req dto.UpsertMarketplaceVoucherReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.mpVoucherUsecase.UpdateMarketplaceVoucher(c.Param("voucher_code"), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MARKETPLACE_VOUCHER",
		Message: "Success update marketplace voucher",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteMarketplaceVoucher(c *gin.Context) {
	voucherCode := c.Param("voucher_code")
	resBody, err := h.mpVoucherUsecase.DeleteMarketplaceVoucher(voucherCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_MARKETPLACE_VOUCHER",
		Message: "Success delete marketplace voucher",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
