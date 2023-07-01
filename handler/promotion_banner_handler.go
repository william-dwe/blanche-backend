package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPromotionBannerList(c *gin.Context) {
	var bannerRequest dto.PaginationRequest
	if err := util.ShouldBindQueryWithValidation(c, &bannerRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.promotionBannerUsecase.GetPromotionBannerList(bannerRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PROMOTION_BANNER",
		Message: "Success retrieve promotion banner",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetPromotionBannerByID(c *gin.Context) {
	bannerId := c.Param("banner_id")
	bannerIdInt, err := strconv.Atoi(bannerId)
	if err != nil {
		_ = c.Error(domain.ErrPromotionBannerIdNotValid)
		return
	}

	res, err := h.promotionBannerUsecase.GetPromotionBannerByID(uint(bannerIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PROMOTION_BANNER",
		Message: "Success retrieve promotion banner",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreatePromotionBanner(c *gin.Context) {
	var bannerRequest dto.UpsertPromotionBannerReqDTO
	if err := util.ShouldBindWithValidation(c, &bannerRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.promotionBannerUsecase.CreatePromotionBanner(bannerRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_PROMOTION_BANNER",
		Message: "Success create promotion banner",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdatePromotionBanner(c *gin.Context) {
	var bannerRequest dto.UpsertPromotionBannerReqDTO
	if err := util.ShouldBindWithValidation(c, &bannerRequest); err != nil {
		_ = c.Error(err)
		return
	}

	bannerId := c.Param("banner_id")
	bannerIdInt, err := strconv.Atoi(bannerId)
	if err != nil {
		_ = c.Error(domain.ErrPromotionBannerIdNotValid)
		return
	}

	resBody, err := h.promotionBannerUsecase.UpdatePromotionBanner(uint(bannerIdInt), bannerRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PROMOTION_BANNER",
		Message: "Success update promotion banner",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeletePromotionBanner(c *gin.Context) {
	bannerId := c.Param("banner_id")
	bannerIdInt, err := strconv.Atoi(bannerId)
	if err != nil {
		_ = c.Error(domain.ErrPromotionBannerIdNotValid)
		return
	}

	resBody, err := h.promotionBannerUsecase.DeletePromotionBanner(uint(bannerIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_PROMOTION_BANNER",
		Message: "Success delete promotion banner",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
