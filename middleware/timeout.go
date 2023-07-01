package middleware

import (
	"net/http"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func testResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusRequestTimeout, httperror.TimeoutError().Error())
	util.AbortWithError(c, httperror.TimeoutError())
}

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Duration(20)*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(testResponse),
	)
}
