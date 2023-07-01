package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllPromotions(c *gin.Context) {
	var req dto.PromotionListReqParamDTO
	if err := util.ShouldBindQueryWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.promotionUsecase.GetAllPromotions(user.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PROMOTION_LIST",
		Message: "Success retrieve promotion list",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreateNewPromotion(c *gin.Context) {
	var req dto.UpsertPromotionReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.promotionUsecase.CreateNewPromotion(user.Username, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_NEW_PROMOTION",
		Message: "Success create new promotion",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetPromotionDetails(c *gin.Context) {
	promotionID := c.Param("promotion_id")
	promotionIDInt, err := strconv.Atoi(promotionID)
	if err != nil {
		_ = c.Error(domain.ErrPromotionIDNotValid)
		return
	}

	resBody, err := h.promotionUsecase.GetPromotionByID(uint(promotionIDInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PROMOTION_DETAILS",
		Message: "Success retrieve promotion details",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdatePromotion(c *gin.Context) {
	promotionID := c.Param("promotion_id")
	promotionIDInt, err := strconv.Atoi(promotionID)
	if err != nil {
		_ = c.Error(domain.ErrPromotionIDNotValid)
		return
	}

	var req dto.UpsertPromotionReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &req); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.promotionUsecase.UpdatePromotion(user.Username, promotionIDInt, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PROMOTION",
		Message: "Success update promotion",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeletePromotion(c *gin.Context) {
	promotionID := c.Param("promotion_id")
	promotionIDInt, err := strconv.Atoi(promotionID)
	if err != nil {
		_ = c.Error(domain.ErrPromotionIDNotValid)
		return
	}

	resBody, err := h.promotionUsecase.DeletePromotion(promotionIDInt)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_PROMOTION",
		Message: "Success delete promotion",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
