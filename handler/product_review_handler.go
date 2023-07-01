package handler

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProductReviewFromTransaction(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.productReviewUsecase.GetProductReviewByInvoiceCode(user.Username, invoiceCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_REVIEW",
		Message: "Success retrieve product review",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AddProductReviewFromTransaction(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.ReviewProductFormReqDTO
	if err := util.ShouldBindWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.productReviewUsecase.AddProductReview(user.Username, reqBody, invoiceCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADD_PRODUCT_REVIEW",
		Message: "Success add product review",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetProductReviewByProductSlug(c *gin.Context) {
	productSlug := c.Param("slug")
	merchantDomain := c.Param("domain")

	var prodReviewParam dto.ProductReviewReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &prodReviewParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.productReviewUsecase.GetProductReviewByProductSlug(fmt.Sprintf("%s/%s", merchantDomain, productSlug), prodReviewParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_REVIEW_DETAILS",
		Message: "Success retrieve product review details",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
