package middleware

import (
	"encoding/json"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/handler"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Authorize(h *handler.Handler, scopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.BlacklistTokenChecker(c)
		if len(c.Errors) != 0 {
			util.AbortWithError(c, httperror.UnauthorizedError())
			return
		}

		currentScope := c.MustGet("scope").(string)
		if !util.ScopeShouldContain(scopes, currentScope) {
			util.AbortWithError(c, httperror.ForbiddenError())
			return
		}
	}
}

func AuthorizeAndBlacklist(h *handler.Handler, scopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.BlacklistTokenChecker(c)
		if c.Errors != nil {
			util.AbortWithError(c, httperror.UnauthorizedError())
			return
		}

		currentScope := c.MustGet("scope").(string)
		if !util.ScopeShouldContain(scopes, currentScope) {
			util.AbortWithError(c, httperror.ForbiddenError())
			return
		}
		h.BlacklistToken(c)
	}
}

func Authenticate(c *gin.Context) {
	conf := config.Config.AuthConfig
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	a := util.NewAuthUtil()
	token, err := a.ValidateToken(accessToken, conf.AccessTokenSecretString)
	if err != nil || !token.Valid {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	userJson, _ := json.Marshal(claims["user"])
	var userPayload dto.AccessTokenPayload
	err = json.Unmarshal(userJson, &userPayload)
	if err != nil {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	c.Set("user", userPayload)
	c.Set("scope", claims["scope"])
}

func AuthenticateAdmin(c *gin.Context) {
	conf := config.Config.AuthConfig
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	a := util.NewAuthUtil()
	token, err := a.ValidateToken(accessToken, conf.AdminAccessTokenSecretString)
	if err != nil || !token.Valid {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	userJson, _ := json.Marshal(claims["user"])
	var userPayload dto.AccessTokenPayload
	err = json.Unmarshal(userJson, &userPayload)
	if err != nil {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}

	c.Set("user", userPayload)
	c.Set("scope", claims["scope"])
}

func AuthenticateWithByPass(c *gin.Context) {
	conf := config.Config.AuthConfig
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		return
	}

	a := util.NewAuthUtil()
	token, err := a.ValidateToken(accessToken, conf.AccessTokenSecretString)
	if err != nil || !token.Valid {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}

	userJson, _ := json.Marshal(claims["user"])
	var userPayload dto.AccessTokenPayload
	err = json.Unmarshal(userJson, &userPayload)
	if err != nil {
		return
	}

	c.Set("user", userPayload)
}
