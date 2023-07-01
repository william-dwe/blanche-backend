package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMarketplaceDashboardActiveUserStatistics(c *gin.Context) {
	var reqBody dto.MarketplaceAnalyticsActiveUserReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.marketplaceAnalyticsUsecase.GetMarketplaceDashboardActiveUserStatistics(reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ACTIVE_USER_STATISTICS",
		Message: "Success get active user statistics",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMarketplaceDashboardUserConversionStatistics(c *gin.Context) {
	var reqBody dto.MarketplaceAnalyticsUserConversionReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.marketplaceAnalyticsUsecase.GetMarketplaceDashboardUserConversionStatistics(reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_USER_CONVERSION_STATISTICS",
		Message: "Success get user conversion statistics",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMarketplaceDashboardSalesStatistics(c *gin.Context) {
	var reqBody dto.MarketplaceAnalyticsSalesReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.marketplaceAnalyticsUsecase.GetMarketplaceDashboardSalesStatistics(reqBody)
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

func (h *Handler) GetMarketplaceDashboardCustomerSatisfactionStatistics(c *gin.Context) {
	var reqBody dto.MarketplaceAnalyticsCustomerSatisfactionReqBody
	if err := util.ShouldBindQueryWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.marketplaceAnalyticsUsecase.GetMarketplaceDashboardCustomerSatisfactionStatistics(reqBody)
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

func (h *Handler) UpdateMarketplaceDashboard(c *gin.Context) {
	var reqBody dto.MarketplaceAnalyticsUpdateReqBody
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.marketplaceAnalyticsUsecase.UpdateMarketplaceDashboard(&reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MARKETPLACE_DASHBOARD",
		Message: "Success update marketplace dashboard",
	}

	util.ResponseSuccessJSON(c, response)
}
