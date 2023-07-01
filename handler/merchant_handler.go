package handler

import (
	"net/http"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserMerchantInfo(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.GetInfoByUsername(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_INFO",
		Message: "Success retrieve merchant info",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantInfo(c *gin.Context) {
	res, err := h.merchantUsecase.GetInfoByDomain(c.Param("domain"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_INFO",
		Message: "Success retrieve merchant info",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CheckMerchantStoreName(c *gin.Context) {
	_, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.CheckMerchantStoreNameReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.CheckMerchantStoreName(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CHECK_MERCHANT_DOMAIN",
		Message: "Success check merchant domain",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CheckMerchantDomain(c *gin.Context) {
	_, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.CheckMerchantDomainReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.CheckMerchantDomain(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CHECK_MERCHANT_DOMAIN",
		Message: "Success check merchant domain",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantProfile(c *gin.Context) {
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.UpdateMerchantProfileFormReqDTO
	if err := util.ShouldBindWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.UpdateMerchantProfile(userJwt.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MERCHANT_PROFILE",
		Message: "Success update merchant profile",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) RegisterMerchant(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.RegisterMerchantReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.RegisterMerchant(user.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_REGISTER_MERCHANT",
		Message: "Success register merchant",
		Data:    res,
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		return
	}

	resBody, err := h.authUsecase.Refresh(refreshToken)
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		return
	}

	appUrl := config.Config.AppUrlUser
	accessTokenExpLimit, _ := strconv.Atoi(config.Config.AuthConfig.AccessTokenExpTimeMinutes)
	isRelease := config.Config.ENVConfig.Mode == config.ENV_MODE_RELEASE
	if isRelease {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	c.SetCookie(
		"access_token",
		resBody.AccessToken,
		accessTokenExpLimit*60,
		"/",
		appUrl,
		isRelease,
		true,
	)

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantProductCategories(c *gin.Context) {
	res, err := h.merchantUsecase.GetProductCategories(c.Param("domain"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_PRODUCTS_CATEGORIES",
		Message: "Success retrieve merchant products categories",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreateMerchantVoucher(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.UpsertMerchantVoucherReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantUsecase.CreateMerchantVoucher(user.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_MERCHANT_VOUCHER",
		Message: "Success create merchant voucher",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantAdminVoucherList(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var voucherRequest dto.MerchantVoucherListParamReqDTO
	if err := util.ShouldBindQueryWithValidation(c, &voucherRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.GetMerchantAdminVoucherList(user.Username, voucherRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_ADMIN_VOUCHER_LIST",
		Message: "Success retrieve merchant admin voucher list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantVoucherList(c *gin.Context) {
	res, err := h.merchantUsecase.GetMerchantVoucherList(c.Param("domain"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_VOUCHER_LIST",
		Message: "Success retrieve merchant voucher list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantAdminVoucherDetails(c *gin.Context) {
	voucherCode := c.Param("voucher_code")
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.GetMerchantVoucherByCode(user.Username, voucherCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_VOUCHER",
		Message: "Success retrieve merchant voucher details",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantFundBalance(c *gin.Context) {
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.GetMerchantFundBalance(userJwt.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_FUND_BALANCE",
		Message: "Success retrieve merchant fund balance",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) WithdrawMerchantFundBalance(c *gin.Context) {
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var withdrawReq dto.MerchantWithdrawReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &withdrawReq); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.WithdrawToWallet(userJwt.Username, withdrawReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_WITHDRAW_MERCHANT_FUND_BALANCE",
		Message: "Success withdraw merchant fund balance",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantFundActivities(c *gin.Context) {
	var fundAccReqParam dto.MerchantFundActivitiesReqParamDTO
	err := c.ShouldBindQuery(&fundAccReqParam)
	if err != nil {
		_ = c.Error(domain.ErrMerchantFundActivitiesReqParam)
		return
	}

	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.merchantUsecase.GetFundActivities(userJwt.Username, fundAccReqParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_MERCHANT_FUND_ACTIVITIES",
		Message: "Success retrieve merchant fund activities",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantVoucher(c *gin.Context) {
	voucherCode := c.Param("voucher_code")
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var req dto.UpsertMerchantVoucherReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantUsecase.UpdateMerchantVoucher(user.Username, voucherCode, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MERCHANT_VOUCHER",
		Message: "Success update merchant voucher",
		Data:    resBody,
	}
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteMerchantAdminVoucher(c *gin.Context) {
	voucherCode := c.Param("voucher_code")
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.merchantUsecase.DeleteMerchantVoucher(user.Username, voucherCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_MERCHANT_VOUCHER",
		Message: "Success delete merchant voucher",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantAddress(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	addressId := c.Param("address_id")
	addressIdInt, err := strconv.Atoi(addressId)
	if err != nil {
		_ = c.Error(domain.ErrInvalidAddressID)
		return
	}

	resBody, err := h.merchantUsecase.UpdateMerchantAddress(user.Username, uint(addressIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MERCHANT_ADDRESS",
		Message: "Success update merchant address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
