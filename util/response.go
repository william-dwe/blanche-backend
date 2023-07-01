package util

import (
	"net/http"
	"reflect"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ResponseStruct struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseSuccessJSONData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func ResponseSuccessJSON(c *gin.Context, data ResponseStruct) {
	c.JSON(http.StatusOK, data)
}

func ResponseSuccessNoContent(c *gin.Context) {
	c.Status(http.StatusOK)
}

func ResponseErrorJSON(c *gin.Context, err any) {
	log.Error().Msgf("Internal Error: %v", err)

	if appErr, isAppError := err.(httperror.AppError); isAppError {
		c.AbortWithStatusJSON(appErr.StatusCode, appErr)
		return
	}

	if reflect.TypeOf(err).String() == "error" {
		otherError := err.(error)
		serverErr := httperror.InternalServerError(otherError.Error())
		c.AbortWithStatusJSON(serverErr.StatusCode, serverErr)
		return
	}

	if reflect.TypeOf(err).String() == "string" {
		serverErr := httperror.InternalServerError(err.(string))
		c.AbortWithStatusJSON(serverErr.StatusCode, serverErr)
		return
	}

	serverErr := httperror.InternalServerError("Internal Error")
	c.AbortWithStatusJSON(serverErr.StatusCode, serverErr)
}
