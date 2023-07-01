package handler

import (
	"net/http"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) UserRegisterCheckEmailHandler(c *gin.Context) {
	var inputRequest dto.UserRegisterCheckEmailReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.CheckUserEmail(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_EMAIL_AVAILABILITY_INFO",
		Message: "Success retrieve email availability info",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserRegisterCheckUsernameHandler(c *gin.Context) {
	var inputRequest dto.UserRegisterCheckUsernameReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.CheckUserUsername(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_EMAIL_AVAILABILITY_INFO",
		Message: "Success retrieve username availability info",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func setAuthCookies(c *gin.Context, refreshToken *string, accessToken *string, isUserLoggedIn *bool, isAdminLoggedIn *bool, appUrl string) {
	refreshTokenExpLimit, _ := strconv.Atoi(config.Config.AuthConfig.RefreshTokenExpTimeMinutes)
	accessTokenExpLimit, _ := strconv.Atoi(config.Config.AuthConfig.AccessTokenExpTimeMinutes)
	isRelease := config.Config.ENVConfig.Mode == config.ENV_MODE_RELEASE
	if isRelease {
		c.SetSameSite(http.SameSiteNoneMode)
	}
	if refreshToken != nil {
		c.SetCookie(
			"refresh_token",
			*refreshToken,
			refreshTokenExpLimit*60,
			"/",
			appUrl,
			isRelease,
			true,
		)
		if *refreshToken == "" {
			c.SetCookie(
				"refresh_token",
				"",
				-1,
				"/",
				appUrl,
				isRelease,
				true,
			)
		}
	}
	if accessToken != nil {
		c.SetCookie(
			"access_token",
			*accessToken,
			accessTokenExpLimit*60,
			"/",
			appUrl,
			isRelease,
			true,
		)
		if *accessToken == "" {
			c.SetCookie(
				"access_token",
				"",
				-1,
				"/",
				appUrl,
				isRelease,
				true,
			)
		}
	}
	if isUserLoggedIn != nil {
		c.SetCookie(
			"is_user_logged_in",
			strconv.FormatBool(*isUserLoggedIn),
			refreshTokenExpLimit*60,
			"/",
			appUrl,
			true,
			false,
		)
		if !*isUserLoggedIn {
			c.SetCookie(
				"is_user_logged_in",
				"",
				-1,
				"/",
				appUrl,
				true,
				false,
			)
		}
	}

	if isAdminLoggedIn != nil {
		c.SetCookie(
			"is_admin_logged_in",
			strconv.FormatBool(*isAdminLoggedIn),
			refreshTokenExpLimit*60,
			"/",
			appUrl,
			true,
			false,
		)
		if !*isAdminLoggedIn {
			c.SetCookie(
				"is_admin_logged_in",
				"",
				-1,
				"/",
				appUrl,
				true,
				false,
			)
		}
	}
}

func (h *Handler) UserRegisterHandler(c *gin.Context) {
	var inputRequest dto.UserRegisterReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	refreshToken, resBody, err := h.authUsecase.Register(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_REGISTERED_USER",
		Message: "Success add user to database",
		Data:    resBody,
	}

	var isUserLoggedIn = true
	var isAdminLoggedIn = false
	setAuthCookies(c, &refreshToken, &resBody.AccessToken, &isUserLoggedIn, &isAdminLoggedIn, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserLoginHandler(c *gin.Context) {
	var inputRequest dto.UserLoginReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	refreshToken, resBody, err := h.authUsecase.Login(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_LOGIN_USER",
		Message: "Success authenticate user",
		Data:    resBody,
	}

	var isUserLoggedIn = true
	var isAdminLoggedIn = false
	setAuthCookies(c, &refreshToken, &resBody.AccessToken, &isUserLoggedIn, &isAdminLoggedIn, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) AdminLoginHandler(c *gin.Context) {
	var inputRequest dto.AdminLoginReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	refreshToken, resBody, err := h.authUsecase.AdminLogin(inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_LOGIN_ADMIN",
		Message: "Success authenticate admin",
		Data:    resBody,
	}

	var isUserLoggedIn = false
	var isAdminLoggedIn = true
	setAuthCookies(c, &refreshToken, &resBody.AccessToken, &isUserLoggedIn, &isAdminLoggedIn, config.Config.AppUrlAdmin)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserRefreshHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	appUrl := config.Config.AppUrlUser
	if c.Request.Header.Get("Origin") == config.Config.WebUrlAdmin {
		appUrl = config.Config.AppUrlAdmin
	}
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		var accessToken, refreshToken, isUserLoggedIn = "", "", false
		setAuthCookies(c, &accessToken, &refreshToken, &isUserLoggedIn, &isUserLoggedIn, appUrl)
		return
	}

	resBody, err := h.authUsecase.Refresh(refreshToken)
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		var accessToken, refreshToken, isLoggedIn = "", "", false
		setAuthCookies(c, &accessToken, &refreshToken, &isLoggedIn, &isLoggedIn, appUrl)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_REFRESH_TOKEN",
		Message: "Success generate new access token",
		Data:    resBody,
	}

	var isUserLoggedIn = true
	var isAdminLoggedIn = false
	if resBody.Role == dto.ROLE_ADMIN {
		isAdminLoggedIn = true
		isUserLoggedIn = false
	}
	setAuthCookies(c, nil, &resBody.AccessToken, &isUserLoggedIn, &isAdminLoggedIn, appUrl)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UserLogoutHandler(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	appUrl := config.Config.AppUrlUser
	if c.Request.Header.Get("Origin") == config.Config.WebUrlAdmin {
		appUrl = config.Config.AppUrlAdmin
	}
	if refreshToken != "" {
		err := h.authUsecase.Logout(refreshToken)
		if err != nil {
			_ = c.Error(httperror.UnauthorizedError())
			return
		}
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_LOGOUT_USER",
		Message: "Success logout user",
		Data:    nil,
	}

	accessToken, refreshToken, isLoggedIn := "", "", false
	setAuthCookies(c, &accessToken, &refreshToken, &isLoggedIn, &isLoggedIn, appUrl)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GoogleRequestLogin(c *gin.Context) {
	googleConfig := config.Config.OauthConfig.ConfigObj

	oauthState := util.GenerateStateOauthCookie(c)
	url := googleConfig.AuthCodeURL(oauthState)

	resBody := dto.GoogleRequestLoginResDTO{
		Url: url,
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_LOGIN_USER",
		Message: "Success authenticate user",
		Data:    resBody,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GoogleRequestCallback(c *gin.Context) {
	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.GoogleCallbackReqDTO
	err = util.ShouldBindQueryWithValidation(c, &inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.authUsecase.ValidateOAuthLoginRequest(oauthState, &inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	oauthReq := dto.UserOAuthLoginReqDTO{Code: inputRequest.Code}
	refreshToken, resBody, err := h.authUsecase.OAuthLogin(&oauthReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if resBody.IsRegistered {
		var isLoggedIn = true
		var isAdminLoggedIn = false
		setAuthCookies(c, &refreshToken, resBody.AccessToken, &isLoggedIn, &isAdminLoggedIn, config.Config.AppUrlUser)
		c.Redirect(http.StatusTemporaryRedirect, config.Config.WebUrlUser)
		return
	}

	accessToken, refreshToken, isLoggedIn := "", "", false
	setAuthCookies(c, &accessToken, &refreshToken, &isLoggedIn, &isLoggedIn, config.Config.AppUrlUser)
	registerUrl := config.Config.WebUrlUser + "/register?email=" + resBody.Email + "&name=" + resBody.Name + "&picture=" + resBody.Picture
	c.Redirect(http.StatusTemporaryRedirect, registerUrl)
}

func (h *Handler) StepUpTokenScopeWithPin(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	scope, err := util.GetScopeJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.StepUpTokenScopeWithPinReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.StepUpTokenScopeWithPin(scope, user, inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_STEP_UP_TOKEN_SCOPE",
		Message: "Success step up token scope",
		Data:    resBody,
	}

	setAuthCookies(c, nil, &resBody.AccessToken, nil, nil, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) StepUpTokenScopeWithPass(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	scope, err := util.GetScopeJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var inputRequest dto.StepUpTokenScopeWithPassReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &inputRequest); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.StepUpTokenScopeWithPass(scope, user, inputRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_STEP_UP_TOKEN_SCOPE",
		Message: "Success step up token scope",
		Data:    resBody,
	}

	setAuthCookies(c, nil, &resBody.AccessToken, nil, nil, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) BlacklistTokenChecker(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}
	isValid, err := h.authUsecase.CheckBlacklistToken(accessToken)
	if err != nil {
		_ = c.Error(err)
		return
	}
	if !isValid {
		util.AbortWithError(c, httperror.UnauthorizedError())
		return
	}
}

func (h *Handler) BlacklistToken(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		_ = c.Error(httperror.UnauthorizedError())
		return
	}

	err = h.authUsecase.BlacklistToken(accessToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	accessToken = ""
	setAuthCookies(c, nil, &accessToken, nil, nil, config.Config.AppUrlUser)
}

func (h *Handler) ChangePasswordRequestVerificationCode(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.SendChangePasswordRequest(user.Username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_RESET_PASSWORD_REQUEST",
		Message: "Success change password request",
		Data:    resBody,
	}
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ChangePasswordVerification(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	scope, err := util.GetScopeJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.ChangePasswordVerificationCodeReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.VerifyChangePasswordRequest(scope, user, &reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_RESET_PASSWORD_REQUEST",
		Message: "Success verified change password request",
		Data:    resBody,
	}

	setAuthCookies(c, nil, &resBody.AccessToken, nil, nil, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ForgetPasswordRequestVerificationUrl(c *gin.Context) {
	var reqBody dto.ForgetPasswordRequestVerificationCodeReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.SendForgetPasswordRequest(reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_RESET_PASSWORD_REQUEST",
		Message: "Success forget password request",
		Data:    resBody,
	}
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ForgetPasswordVerification(c *gin.Context) {
	var reqBody dto.ForgetPasswordVerificationCodeReqDTO
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.VerifyForgetPasswordRequest(&reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_RESET_PASSWORD_REQUEST",
		Message: "Success verified forget password request",
		Data:    resBody,
	}

	setAuthCookies(c, nil, &resBody.AccessToken, nil, nil, config.Config.AppUrlUser)
	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var reqBody dto.ResetPasswordReqBody
	if err := util.ShouldBindJsonWithValidation(c, &reqBody); err != nil {
		_ = c.Error(err)
		return
	}

	resBody, err := h.authUsecase.ResetPassword(user.Username, &reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_RESET_PASSWORD",
		Message: "Success change password",
		Data:    resBody,
	}

	h.BlacklistToken(c)

	accessToken := ""
	setAuthCookies(c, nil, &accessToken, nil, nil, config.Config.AppUrlUser)

	util.ResponseSuccessJSON(c, response)
}
