package handler

import (
	"fmt"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) AddItemToCart(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var cartReq dto.AddItemToCartReqDTO
	if err := c.ShouldBindJSON(&cartReq); err != nil {
		errTag, cvtErr := err.(validator.ValidationErrors)
		if !cvtErr {
			_ = c.Error(domain.ErrAddCartItemInvalidInput)
			return
		}
		if errTag != nil {
			_ = c.Error(
				httperror.BadRequestError(
					fmt.Sprintf(
						"check the requested param: %s must be %s %s",
						errTag[0].StructField(),
						errTag[0].ActualTag(),
						errTag[0].Param(),
					),
					"DATA_NOT_VALID",
				),
			)
			return
		}
	}

	res, err := h.cartItemUsecase.AddCartItem(user.Username, cartReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADD_ITEM_TO_CART",
		Message: "Success add item to cart",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCart(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.cartItemUsecase.GetCartItems(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CART",
		Message: "Success get cart",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetHomeCart(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.cartItemUsecase.GetHomeCartItems(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CART",
		Message: "Success get cart",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteCartItem(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	cartItemID := c.Param("cart_item_id")
	cartItemIDInt, err := strconv.Atoi(cartItemID)
	if err != nil {
		_ = c.Error(domain.ErrInvalidCartItemID)
		return
	}
	err = h.cartItemUsecase.DeleteCartItemByCartId(user.Username, uint(cartItemIDInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_CART_ITEM",
		Message: "Success delete cart item",
		Data:    nil,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteSelectedCartItem(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.cartItemUsecase.DeleteSelectedCartItem(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_CART_ITEM",
		Message: "Success delete cart item",
		Data:    nil,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateCart(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	cartItemID := c.Param("cart_item_id")
	cartItemIDInt, err := strconv.Atoi(cartItemID)
	if err != nil {
		_ = c.Error(domain.ErrInvalidCartItemID)
		return
	}

	var cartReq dto.UpdateCartItemDTO
	if err := c.ShouldBindJSON(&cartReq); err != nil {
		errTag, cvtErr := err.(validator.ValidationErrors)
		if !cvtErr {
			_ = c.Error(domain.ErrUpdateCartItemInvalidInput)
			return
		}
		if errTag != nil {
			_ = c.Error(
				httperror.BadRequestError(
					fmt.Sprintf(
						"check the requested param: %s must be %s %s",
						errTag[0].StructField(),
						errTag[0].ActualTag(),
						errTag[0].Param(),
					),
					"DATA_NOT_VALID",
				),
			)
			return
		}
	}

	res, err := h.cartItemUsecase.UpdateCartItem(user.Username, uint(cartItemIDInt), cartReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_CART_ITEM",
		Message: "Success update cart item",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateAllCart(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var cartReq []dto.UpdateAllCartCheckStatusDTO
	if err := c.ShouldBindJSON(&cartReq); err != nil {
		errTag, cvtErr := err.(validator.ValidationErrors)
		if !cvtErr {
			_ = c.Error(domain.ErrUpdateCartItemInvalidInput)
			return
		}
		if errTag != nil {
			_ = c.Error(
				httperror.BadRequestError(
					fmt.Sprintf(
						"check the requested param: %s must be %s %s",
						errTag[0].StructField(),
						errTag[0].ActualTag(),
						errTag[0].Param(),
					),
					"DATA_NOT_VALID",
				),
			)
			return
		}
	}

	res, err := h.cartItemUsecase.UpdateAllCartCheckStatus(user.Username, cartReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_CART_ITEM",
		Message: "Success update cart item",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
