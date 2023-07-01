package handler

import (
	"fmt"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProductList(c *gin.Context) {
	var productRequest dto.ProductListReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if productRequest.SellerCityId != "" {
		productRequest.SellerCityIdList = util.StringToArrInt(productRequest.SellerCityId, ",")
	}

	res, err := h.productUsecase.GetProductList(productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_LIST",
		Message: "Success retrieve product list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantProductList(c *gin.Context) {
	var productRequest dto.ProductListReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productUsecase.GetMerchantProductList(user.Username, productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_LIST_BY_MERCHANT",
		Message: "Success retrieve product list by merchant",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetRecommendationProductList(c *gin.Context) {
	var paginationRequest dto.PaginationRequest
	err := c.ShouldBindQuery(&paginationRequest)
	if err != nil {
		_ = c.Error(domain.ErrProductQueryParam)
		return
	}

	res, err := h.productUsecase.GetRecommendationProductList(paginationRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_RECOMMENDATION_PRODUCT_LIST",
		Message: "Success retrieve recommendation product list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetProductDetails(c *gin.Context) {
	domain := c.Param("domain")
	slug := c.Param("slug")

	//check is logged in
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		userJwt = nil
	}

	res, err := h.productUsecase.GetProductDetailsBySlug(userJwt, fmt.Sprintf("%s/%s", domain, slug))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_DETAILS",
		Message: "Success retrieve product details",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantProductDetails(c *gin.Context) {
	productId := c.Param("product_id")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		_ = c.Error(domain.ErrProductIdNotValid)
		return
	}
	//check is logged in
	userJwt, err := util.GetUserJWTContext(c)
	if err != nil {
		userJwt = nil
	}

	res, err := h.productUsecase.GetAdminProductDetailByProductID(userJwt, uint(productIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_DETAILS",
		Message: "Success retrieve product details",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UploadProductImage(c *gin.Context) {
	var productImageRequest dto.UploadImageReqDTO
	if err := util.ShouldBindWithValidation(c, &productImageRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productUsecase.UploadProductFiles(productImageRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPLOAD_PRODUCT_IMAGE",
		Message: "Success upload product image",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var productRequest dto.CreateProductReqDTO
	err := util.ShouldBindJsonWithValidation(c, &productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = productRequest.Validate()
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productUsecase.CreateProduct(user.Username, productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_PRODUCT",
		Message: "Success create product",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CheckMerchantProductName(c *gin.Context) {
	var productRequest dto.CheckMerchantProductNameReqDTO
	err := util.ShouldBindJsonWithValidation(c, &productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productUsecase.CheckMerchantProductName(user.Username, productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CHECK_MERCHANT_PRODUCT_NAME",
		Message: "Success check merchant product name",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteMerchantProduct(c *gin.Context) {
	productIdList := c.Param("product_id")
	productIdIntList := util.StringToArrInt(productIdList, ",")

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.productUsecase.DeleteMerchantProduct(user.Username, productIdIntList)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_MERCHANT_PRODUCT",
		Message: "Success delete merchant product",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantProduct(c *gin.Context) {
	productId := c.Param("product_id")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		_ = c.Error(domain.ErrProductIdNotValid)
		return
	}
	var productRequest dto.CreateProductReqDTO
	err = util.ShouldBindJsonWithValidation(c, &productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = productRequest.Validate()
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productUsecase.UpdateMerchantProduct(user.Username, uint(productIdInt), productRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PRODUCT",
		Message: "Success update product",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateProductAvailability(c *gin.Context) {
	productIdList := c.Param("product_id")
	productIdIntList := util.StringToArrInt(productIdList, ",")

	var updateRequest dto.UpdateProductAvailabilityReqDTO
	err := util.ShouldBindJsonWithValidation(c, &updateRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.productUsecase.UpdateMerchantProductAvailability(user.Username, productIdIntList, updateRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PRODUCT_AVAILABILITY",
		Message: "Success update product availability",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
