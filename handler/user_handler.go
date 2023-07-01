package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UserGetProfileHandler(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.GetProfile(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PROFILE",
		Message: "Success retrieve user profile",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserUpdateProfileHandler(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.UserUpdateProfileReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.UpdateProfile(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PROFILE",
		Message: "Success retrieve user profile",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserUpdateProfileDetailHandler(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.UserUpdateProfileDetailFormReqDTO
	if err := util.ShouldBindWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.UpdateProfileDetail(user.Username, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_PROFILE",
		Message: "Success retrieve user profile",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetAllUserAddress(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.GetAllUserAddress(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_ADDRESS",
		Message: "Success retrieve all user address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AddUserAddress(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.UserAddressReqDTO
	err = util.ShouldBindJsonWithValidation(c, &inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.AddUserAddress(user.Username, inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_ADD_ADDRESS",
		Message: "Success add user address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) SetDefaultUserAddress(c *gin.Context) {
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

	resBody, err := h.userUsecase.SetDefaultUserAddress(user.Username, uint(addressIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_SET_DEFAULT_ADDRESS",
		Message: "Success set default user address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) DeleteUserAddress(c *gin.Context) {
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

	resBody, err := h.userUsecase.DeleteUserAddress(user.Username, uint(addressIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_DELETE_ADDRESS",
		Message: "Success delete user address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateUserAddress(c *gin.Context) {
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

	var inputRequest dto.UserAddressUpdateReqDTO
	err = util.ShouldBindJsonWithValidation(c, &inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.userUsecase.UpdateUserAddress(user.Username, uint(addressIdInt), inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_ADDRESS",
		Message: "Success update user address",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}
