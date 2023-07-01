package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMerchantDashboardMerchantResponsivenessStatistics(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.MerchantAnalyticsMerchantResponsivenessReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantAnalyticsUsecase.GetMerchantDashboardMerchantResponsivenessStatistics(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_RESPONSIVENESS_STATISTICS",
		Message: "Success get merchant responsiveness statistics",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantDashboardSalesStatistics(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.MerchantAnalyticsSalesReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantAnalyticsUsecase.GetMerchantDashboardSalesStatistics(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_SALES_STATISTICS",
		Message: "Success get sales statistics",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantDashboardCustomerSatisfactionStatistics(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.MerchantAnalyticsCustomerSatisfactionReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantAnalyticsUsecase.GetMerchantDashboardCustomerSatisfactionStatistics(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CUSTOMER_SATISFACTION_STATISTICS",
		Message: "Success get customer satisfaction statistics",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantDashboard(c *gin.Context) {
	var reqBody dto.MerchantAnalyticsUpdateReqBody
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.merchantAnalyticsUsecase.UpdateMerchantDashboard(&reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MARKETPLACE_DASHBOARD",
		Message: "Success update merchant dashboard",
	}

	util.ResponseSuccessJSON(c, response)
}
