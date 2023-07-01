package handler

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateUserFavoriteProduct(c *gin.Context) {
	var inputRequest dto.UserFavoriteProductReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	inputRequest.Username = userJWT.Username

	resBody, err := h.userFavoriteProductUsecase.UpdateFavoriteProduct(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_FAVORITE_PRODUCT",
		Message: "Success update favorite product",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetUserFavoriteProducts(c *gin.Context) {
	var queryParam dto.UserFavoriteProductReqParamDTO
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		_ = c.Error(domain.ErrUserFavoriteProductQueryDeficient)
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_USER_FAVORITE_PRODUCT_LIST",
		Message: "Success retrieve user favorite product list",
		Data: dto.ProductListResDTO{
			PaginationResponse: dto.PaginationResponse{
				CurrentPage: 1,
			},
			Products: []dto.ProductResDTO{},
		},
	}

	search, exist := c.GetQuery("q")
	if exist {
		if strings.TrimSpace(search) == "" {
			util.ResponseSuccessJSON(c, response)
			return
		}
	}

	userJWT, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userFavoriteProductUsecase.GetFavoriteProducts(userJWT.Username, queryParam)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Data = resBody

	util.ResponseSuccessJSON(c, response)
}
