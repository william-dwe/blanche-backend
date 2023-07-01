package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CachedExampleHandler(c *gin.Context) {
	var inputRequest dto.ExampleReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		util.ResponseErrorJSON(c, err)
		return
	}

	res, err := h.exampleUsecase.CachedExampleProcess(inputRequest)
	if err != nil {
		util.ResponseErrorJSON(c, err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCEED_CACHED_EXAMPLE_HANDLER",
		Message: "Success",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ExampleHandler(c *gin.Context) {
	var inputRequest dto.ExampleReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		util.ResponseErrorJSON(c, err)
		return
	}

	res, err := h.exampleUsecase.ExampleProcess(inputRequest)
	if err != nil {
		util.ResponseErrorJSON(c, err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCEED_EXAMPLE_HANDLER",
		Message: "Success",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ExampleHandlerErrorMiddleware(c *gin.Context) {
	_ = c.Error(domain.ErrExampleUnexpected)
}
