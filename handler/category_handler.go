package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCategoryTree(c *gin.Context) {
	var categoryRequest dto.CategoryListReqParamDTO
	err := c.ShouldBindQuery(&categoryRequest)
	if err != nil {
		_ = c.Error(domain.ErrCategoryQueryParam)
		return
	}

	res, err := h.categoryUsecase.GetCategoryTree(categoryRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CATEGORY_TREE",
		Message: "Success retrieve category Tree",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCategoryAncestorsById(c *gin.Context) {
	res, err := h.categoryUsecase.GetCategoryAncestorsBySlug(c.Param("category_slug"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CATEGORY_ANCESTORS",
		Message: "Success retrieve category ancestors",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCategoryBySlug(c *gin.Context) {
	res, err := h.categoryUsecase.GetCategoryBySlug(c.Param("category_slug"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CATEGORY",
		Message: "Success retrieve category",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var categoryRequest dto.UpsertCategoryReqDTO
	if err := util.ShouldBindWithValidation(c, &categoryRequest); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.categoryUsecase.CreateCategory(categoryRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_CREATE_CATEGORY",
		Message: "Success create category",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCategoryDetailByID(c *gin.Context) {
	categoryId := c.Param("category_id")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		_ = c.Error(domain.ErrCategoryQueryParam)
		return
	}

	res, err := h.categoryUsecase.GetCategoryByID(uint(categoryIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CATEGORY_DETAIL",
		Message: "Success retrieve category detail",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var categoryRequest dto.UpsertCategoryReqDTO
	if err := util.ShouldBindWithValidation(c, &categoryRequest); err != nil {
		_ = c.Error(err)
		return
	}

	categoryId := c.Param("category_id")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		_ = c.Error(domain.ErrCategoryQueryParam)
		return
	}

	res, err := h.categoryUsecase.UpdateCategory(uint(categoryIdInt), categoryRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_CATEGORY",
		Message: "Success update category",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	categoryId := c.Param("category_id")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		_ = c.Error(domain.ErrCategoryQueryParam)
		return
	}

	resBody, err := h.categoryUsecase.DeleteCategory(uint(categoryIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_CATEGORY",
		Message: "Success delete category",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCategoryList(c *gin.Context) {
	var categoryRequest dto.PaginationRequest
	if err := util.ShouldBindQueryWithValidation(c, &categoryRequest); err != nil {
		_ = c.Error(domain.ErrCategoryQueryParam)
		return
	}

	res, err := h.categoryUsecase.GetCategoryList(categoryRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CATEGORY_LIST",
		Message: "Success retrieve category list",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
