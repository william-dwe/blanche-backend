package util

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"github.com/gin-gonic/gin"
)

func AbortWithError(c *gin.Context, err httperror.AppError) {
	c.AbortWithStatusJSON(err.StatusCode, err)
}
