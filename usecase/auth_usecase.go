package usecase

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type AuthUsecase interface {
	CheckUserEmail(dto.UserRegisterCheckEmailReqDTO) (*dto.UserRegisterCheckEmailResDTO, error)
	CheckUserUsername(input dto.UserRegisterCheckUsernameReqDTO) (*dto.UserRegisterCheckUsernameResDTO, error)

	Register(input dto.UserRegisterReqDTO) (string, *dto.UserRegisterResDTO, error)
	Login(input dto.UserLoginReqDTO) (string, *dto.UserLoginResDTO, error)
	ValidateOAuthLoginRequest(oauthState string, input *dto.GoogleCallbackReqDTO) error
	OAuthLogin(input *dto.UserOAuthLoginReqDTO) (string, *dto.UserOAuthLoginResDTO, error)
	AdminLogin(input dto.AdminLoginReqDTO) (string, *dto.AdminLoginResDTO, error)
	Refresh(refreshToken string) (*dto.UserRefreshResDTO, error)
	Logout(refreshToken string) error

	StepUpTokenScopeWithPin(currentTokenScope string, payload *dto.AccessTokenPayload, input dto.StepUpTokenScopeWithPinReqDTO) (*dto.UserStepUpTokenScopeResDTO, error)
	StepUpTokenScopeWithPass(currentTokenScope string, payload *dto.AccessTokenPayload, input dto.StepUpTokenScopeWithPassReqDTO) (*dto.UserStepUpTokenScopeResDTO, error)

	CheckBlacklistToken(token string) (bool, error)
	BlacklistToken(token string) error

	SendChangePasswordRequest(username string) (*dto.ChangePasswordRequestVerificationCodeResDTO, error)
	VerifyChangePasswordRequest(currentTokenScope string, payload *dto.AccessTokenPayload, input *dto.ChangePasswordVerificationCodeReqDTO) (*dto.ChangePasswordVerificationCodeResDTO, error)
	SendForgetPasswordRequest(input dto.ForgetPasswordRequestVerificationCodeReqDTO) (*dto.ForgetPasswordRequestVerificationCodeResDTO, error)
	VerifyForgetPasswordRequest(input *dto.ForgetPasswordVerificationCodeReqDTO) (*dto.ForgetPasswordVerificationCodeResDTO, error)
	ResetPassword(username string, input *dto.ResetPasswordReqBody) (*dto.ResetPasswordResBody, error)
}

type AuthUsecaseConfig struct {
	AuthRepository   repository.AuthRepository
	UserRepository   repository.UserRepository
	WalletRepository repository.WalletRepository
	AuthUtil         util.AuthUtil
}

type authUsecaseImpl struct {
	authRepository   repository.AuthRepository
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	authUtil         util.AuthUtil
}

func NewAuthUsecase(c AuthUsecaseConfig) AuthUsecase {
	return &authUsecaseImpl{
		authRepository:   c.AuthRepository,
		userRepository:   c.UserRepository,
		walletRepository: c.WalletRepository,
		authUtil:         c.AuthUtil,
	}
}

func (u *authUsecaseImpl) CheckUserEmail(input dto.UserRegisterCheckEmailReqDTO) (*dto.UserRegisterCheckEmailResDTO, error) {
	if isEmailValid := util.IsValidEmail(input.Email); !isEmailValid {
		return nil, domain.ErrInvalidEmailFormat
	}
	err := u.userRepository.CheckEmailBlacklist(input.Email)
	if err != nil {
		return &dto.UserRegisterCheckEmailResDTO{
			Email:       input.Email,
			IsAvailable: false,
		}, nil
	}
	isAvailable, err := u.authRepository.CheckUserEmailExistence(input.Email)
	if err != nil {
		return nil, err
	}
	return &dto.UserRegisterCheckEmailResDTO{
		Email:       input.Email,
		IsAvailable: isAvailable,
	}, nil
}

func (u *authUsecaseImpl) CheckUserUsername(input dto.UserRegisterCheckUsernameReqDTO) (*dto.UserRegisterCheckUsernameResDTO, error) {
	isAvailable, err := u.authRepository.CheckUserUsernameExistence(input.Username)
	if err != nil {
		return nil, err
	}
	return &dto.UserRegisterCheckUsernameResDTO{
		Username:    input.Username,
		IsAvailable: isAvailable,
	}, nil
}

func (u *authUsecaseImpl) Register(input dto.UserRegisterReqDTO) (string, *dto.UserRegisterResDTO, error) {
	user, err := u.validateAndAddUserRegisterRequest(input)
	if err != nil {
		return "", nil, err
	}
	role, err := u.userRepository.GetRoleByRoleName(dto.ROLE_USER)
	if err != nil {
		return "", nil, err
	}

	scope := generateRoleScope(role.RoleName)
	accessTokenStr, err := u.authUtil.GenerateAccessToken(user.Username, scope)
	if err != nil {
		return "", nil, domain.ErrFailedToGenerateAccessToken
	}

	refreshTokenStr, err := u.authUtil.GenerateRefreshToken()
	if err != nil {
		return "", nil, domain.ErrFailedToGenerateRefreshToken
	}
	timeLimit, err := strconv.Atoi(config.Config.AuthConfig.AccessTokenExpTimeMinutes)
	if err != nil {
		return "", nil, domain.ErrFailedToGenerateRefreshToken
	}
	loginLog := entity.UserLoginActivity{
		UserId:       user.ID,
		RefreshToken: refreshTokenStr,
		ExpiredAt:    time.Now().Add(time.Minute * time.Duration(timeLimit)),
	}

	err = u.authRepository.AddUserLoginActivity(loginLog)
	if err != nil {
		return "", nil, err
	}

	respBody := dto.UserRegisterResDTO{
		AccessToken: accessTokenStr,
	}
	return refreshTokenStr, &respBody, nil
}

func (u *authUsecaseImpl) validateAndAddUserRegisterRequest(input dto.UserRegisterReqDTO) (*entity.User, error) {
	input.Email = strings.ToLower(input.Email)
	if isEmailValid := util.IsValidEmail(input.Email); !isEmailValid {
		return nil, domain.ErrInvalidEmailFormat
	}
	err := u.userRepository.CheckEmailBlacklist(input.Email)
	if err != nil {
		return nil, err
	}
	if err := util.ValidateFullname(input.Fullname); err != nil {
		return nil, err
	}
	if err := util.ValidateUsername(input.Username); err != nil {
		return nil, err
	}

	role, err := u.userRepository.GetRoleByRoleName(dto.ROLE_USER)
	if err != nil {
		return nil, err
	}
	newUser := entity.User{
		RoleId:   role.ID,
		Email:    input.Email,
		Username: input.Username,
		UserDetail: entity.UserDetail{
			Fullname: input.Fullname,
		},
	}

	if err := util.ValidatePassword(input.Password, input.Username); err != nil {
		return nil, err
	}
	hashedPass, _ := util.HashAndSalt(input.Password)
	newUser.Password = hashedPass

	user, err := u.userRepository.AddUser(newUser)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *authUsecaseImpl) Login(input dto.UserLoginReqDTO) (string, *dto.UserLoginResDTO, error) {
	input.Email = strings.ToLower(input.Email)
	user, err := u.userRepository.GetUserByEmail(input.Email)
	if err != nil {
		return "", nil, httperror.UnauthorizedErrorLogin()
	}

	if !util.ValidateHash(user.Password, input.Password) {
		return "", nil, httperror.UnauthorizedErrorLogin()
	}

	refreshToken, accessToken, err := u.preparingLogin(user, false)
	if err != nil {
		return "", nil, err
	}

	resBody := dto.UserLoginResDTO{
		AccessToken: accessToken,
	}

	return refreshToken, &resBody, nil
}

func (u *authUsecaseImpl) preparingLogin(user *entity.User, isAdminLogin bool) (string, string, error) {
	role, err := u.userRepository.GetRoleByRoleId(user.RoleId)
	if err != nil {
		return "", "", domain.ErrFailedToRetrieveUserRole
	}
	var scope = generateRoleScope(role.RoleName)

	var accessTokenStr, refreshTokenStr string
	if isAdminLogin {
		if role.RoleName != dto.ROLE_ADMIN {
			return "", "", httperror.UnauthorizedErrorLogin()
		}
		accessTokenStr, err = u.authUtil.GenerateAdminAccessToken(user.Username, scope)
		if err != nil {
			return "", "", domain.ErrFailedToGenerateAccessToken
		}
		refreshTokenStr, err = u.authUtil.GenerateAdminRefreshToken()
		if err != nil {
			return "", "", domain.ErrFailedToGenerateRefreshToken
		}
	} else {
		if role.RoleName == dto.ROLE_ADMIN {
			return "", "", httperror.UnauthorizedErrorLogin()
		}
		accessTokenStr, err = u.authUtil.GenerateAccessToken(user.Username, scope)
		if err != nil {
			return "", "", domain.ErrFailedToGenerateAccessToken
		}
		refreshTokenStr, err = u.authUtil.GenerateRefreshToken()
		if err != nil {
			return "", "", domain.ErrFailedToGenerateRefreshToken
		}
	}

	timeLimit, err := strconv.Atoi(config.Config.AuthConfig.RefreshTokenExpTimeMinutes)
	if err != nil {
		return "", "", domain.ErrFailedToGenerateRefreshToken
	}
	loginLog := entity.UserLoginActivity{
		UserId:       user.ID,
		RefreshToken: refreshTokenStr,
		ExpiredAt:    time.Now().Add(time.Minute * time.Duration(timeLimit)),
	}

	err = u.authRepository.AddUserLoginActivity(loginLog)
	if err != nil {
		return "", "", domain.ErrFailedToGenerateRefreshToken
	}

	return refreshTokenStr, accessTokenStr, nil
}

func (u *authUsecaseImpl) ValidateOAuthLoginRequest(oauthState string, input *dto.GoogleCallbackReqDTO) error {
	if input.State != oauthState {
		return domain.ErrInvalidOAuthState
	}

	return nil
}

func (u *authUsecaseImpl) OAuthLogin(input *dto.UserOAuthLoginReqDTO) (string, *dto.UserOAuthLoginResDTO, error) {
	gApiUserData, err := u.authRepository.GetUserDataFromGoogle(input.Code)
	if err != nil {
		return "", nil, err
	}

	user, err := u.userRepository.GetUserByEmail(gApiUserData.Email)
	if err == domain.ErrCheckEmailInvalidInput {
		return "", &dto.UserOAuthLoginResDTO{
			Email:        gApiUserData.Email,
			IsRegistered: false,
			Name:         gApiUserData.Name,
			Picture:      gApiUserData.Picture,
			AccessToken:  nil,
		}, nil
	}

	if err != nil {
		return "", nil, err
	}

	refreshToken, accessToken, err := u.preparingLogin(user, false)
	if err != nil {
		return "", nil, err
	}
	return refreshToken, &dto.UserOAuthLoginResDTO{
		Email:        gApiUserData.Email,
		IsRegistered: true,
		Name:         gApiUserData.Name,
		Picture:      gApiUserData.Picture,
		AccessToken:  &accessToken,
	}, nil
}

func (u *authUsecaseImpl) AdminLogin(input dto.AdminLoginReqDTO) (string, *dto.AdminLoginResDTO, error) {
	user, err := u.userRepository.GetUserByEmail(input.Email)
	if err != nil {
		return "", nil, httperror.UnauthorizedErrorLogin()
	}

	if !util.ValidateHash(user.Password, input.Password) {
		return "", nil, httperror.UnauthorizedErrorLogin()
	}

	refreshToken, accessToken, err := u.preparingLogin(user, true)
	if err != nil {
		return "", nil, err
	}

	resBody := dto.AdminLoginResDTO{
		AccessToken: accessToken,
	}

	return refreshToken, &resBody, nil
}

func (u *authUsecaseImpl) Refresh(refreshToken string) (*dto.UserRefreshResDTO, error) {
	_, err := u.authUtil.ValidateToken(refreshToken, config.Config.AuthConfig.RefreshTokenSecretString)
	if err != nil {
		return nil, err
	}

	userId, err := u.authRepository.GetUserIdByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByUserId(*userId)
	if err != nil {
		return nil, err
	}

	role, err := u.userRepository.GetRoleByRoleId(user.RoleId)
	if err != nil {
		return nil, err
	}

	var accessTokenStr string
	var scope = generateRoleScope(role.RoleName)
	if role.RoleName != dto.ROLE_ADMIN {
		accessTokenStr, err = u.authUtil.GenerateAccessToken(user.Username, scope)
	} else {
		accessTokenStr, err = u.authUtil.GenerateAdminAccessToken(user.Username, scope)
	}
	if err != nil {
		return nil, err
	}

	resBody := dto.UserRefreshResDTO{
		AccessToken: accessTokenStr,
		Role:        role.RoleName,
	}
	return &resBody, nil
}

func generateRoleScope(roleName string) string {
	var scope = roleName
	if roleName == dto.ROLE_MERCHANT {
		scope += " " + dto.ROLE_USER
	}
	return scope
}

func (u *authUsecaseImpl) Logout(refreshToken string) error {
	_, err := u.authUtil.ValidateToken(refreshToken, config.Config.AuthConfig.RefreshTokenSecretString)
	if err != nil {
		return err
	}

	err = u.authRepository.DeleteUserActivityByRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecaseImpl) StepUpTokenScopeWithPin(currentTokenScope string, payload *dto.AccessTokenPayload, input dto.StepUpTokenScopeWithPinReqDTO) (*dto.UserStepUpTokenScopeResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(payload.Username)
	if err != nil {
		return nil, err
	}
	wallet, err := u.walletRepository.GetByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	if isBlocked, err := u.CheckBlockWallet(payload.Username); isBlocked {
		return nil, err
	}

	if !util.ValidateHash(wallet.Pin, input.Pin) {
		u.authRepository.BlockWalletAddAttempt(payload.Username)
		return nil, domain.ErrInvalidPin
	}
	u.authRepository.BlockWalletResetAttempt(payload.Username)

	var newScope = util.ScopeAddTag(currentTokenScope, dto.SCOPE_PIN)
	accessTokenStr, err := u.authUtil.GenerateAccessToken(user.Username, newScope)
	if err != nil {
		return nil, domain.ErrFailedToGenerateAccessToken
	}

	resBody := dto.UserStepUpTokenScopeResDTO{
		AccessToken: accessTokenStr,
	}
	return &resBody, nil
}

func (u *authUsecaseImpl) StepUpTokenScopeWithPass(currentTokenScope string, payload *dto.AccessTokenPayload, input dto.StepUpTokenScopeWithPassReqDTO) (*dto.UserStepUpTokenScopeResDTO, error) {
	if isBlocked, err := u.CheckBlockWallet(payload.Username); isBlocked {
		return nil, err
	}

	user, err := u.userRepository.GetUserByUsername(payload.Username)
	if err != nil {
		return nil, err
	}

	if !util.ValidateHash(user.Password, input.Password) {
		u.authRepository.BlockWalletAddAttempt(payload.Username)
		return nil, domain.ErrInvalidPass
	}
	u.authRepository.BlockWalletResetAttempt(payload.Username)

	var newScope = util.ScopeAddTag(currentTokenScope, dto.SCOPE_PASSWORD)
	accessTokenStr, err := u.authUtil.GenerateAccessToken(user.Username, newScope)
	if err != nil {
		return nil, domain.ErrFailedToGenerateAccessToken
	}

	resBody := dto.UserStepUpTokenScopeResDTO{
		AccessToken: accessTokenStr,
	}
	return &resBody, nil
}

func (u *authUsecaseImpl) CheckBlockWallet(username string) (bool, error) {
	attempt, err := u.authRepository.BlockWalletGetAttempt(username)
	if err != nil {
		return false, err
	}
	if attempt >= dto.MAX_RETRY_WALLET {
		return true, domain.ErrWalletBlocked
	}
	return false, nil
}

func (u *authUsecaseImpl) CheckBlacklistToken(token string) (bool, error) {
	cacheKey := "blacklist_token:" + token
	isValid, err := u.authRepository.GetBlacklistedToken(cacheKey)
	if err != nil {
		return isValid, err
	}

	return isValid, nil
}

func (u *authUsecaseImpl) BlacklistToken(token string) error {
	err := u.authRepository.AddBlacklistedToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecaseImpl) SendChangePasswordRequest(username string) (*dto.ChangePasswordRequestVerificationCodeResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	err = u.authRepository.CheckBlockResetPasswordRequest(username, dto.ACTION_CHANGE_PASSWORD)
	if err != nil {
		if err == domain.ErrResetPasswordWait {
			ttl, err_block := u.authRepository.GetTTLBlockPasswordCode(username, dto.ACTION_CHANGE_PASSWORD)
			if err_block != nil {
				return nil, err_block
			}
			return &dto.ChangePasswordRequestVerificationCodeResDTO{IsEmailSent: false, Email: user.Email, Username: user.Username, RetryIn: ttl}, nil
		}
		return nil, err
	}

	code, err := u.authUtil.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	hashedCode, err := util.HashAndSalt(code)
	if err != nil {
		return nil, domain.ErrChangePasswordInternalError
	}

	ttl, err := u.authRepository.AddChangePasswordCode(hashedCode, username)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(dto.SMTP_CHANGE_PASS_HTML_PATH)
	if err != nil {
		return nil, domain.ErrChangePasswordInternalError
	}

	body := strings.Replace(string(b), "{{OTP}}", code, 1)

	mailStruct := util.Mail{
		SenderAddress: config.Config.SmtpConfig.ResetPasswordAddress,
		SenderName:    dto.SMTP_CHANGE_PASS_SENDER_NAME,
		ToAddress:     user.Email,
		Subject:       dto.SMTP_CHANGE_PASS_SUBJECT,
		Body:          body,
	}

	err = util.SMTPSendMail(mailStruct)
	if err != nil {
		return nil, domain.ErrChangePasswordInternalError
	}

	return &dto.ChangePasswordRequestVerificationCodeResDTO{IsEmailSent: true, Email: user.Email, Username: user.Username, RetryIn: ttl}, nil
}

func (u *authUsecaseImpl) VerifyChangePasswordRequest(currentTokenScope string, payload *dto.AccessTokenPayload, input *dto.ChangePasswordVerificationCodeReqDTO) (*dto.ChangePasswordVerificationCodeResDTO, error) {
	hashedCode, err := u.authRepository.GetChangePasswordCode(payload.Username)
	if err != nil {
		return nil, err
	}

	if !util.ValidateHash(hashedCode, input.VerificationCode) {
		return nil, domain.ErrInvalidVerificationCode
	}

	var newScope = util.ScopeAddTag(currentTokenScope, dto.SCOPE_RESET_PASSWORD)
	accessTokenStr, err := u.authUtil.GenerateAccessToken(payload.Username, newScope)
	if err != nil {
		return nil, domain.ErrFailedToGenerateAccessToken
	}

	err = u.authRepository.RemoveChangePasswordCode(payload.Username)
	if err != nil {
		return nil, err
	}

	return &dto.ChangePasswordVerificationCodeResDTO{Username: payload.Username, AccessToken: accessTokenStr}, nil
}

func (u *authUsecaseImpl) SendForgetPasswordRequest(input dto.ForgetPasswordRequestVerificationCodeReqDTO) (*dto.ForgetPasswordRequestVerificationCodeResDTO, error) {
	user, err := u.userRepository.GetUserByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	err = u.authRepository.CheckBlockResetPasswordRequest(user.Username, dto.ACTION_FORGET_PASSWORD)
	if err != nil {
		if err == domain.ErrResetPasswordWait {
			ttl, err_block := u.authRepository.GetTTLBlockPasswordCode(user.Username, dto.ACTION_FORGET_PASSWORD)
			if err_block != nil {
				return nil, err_block
			}
			return &dto.ForgetPasswordRequestVerificationCodeResDTO{IsEmailSent: false, Email: user.Email, Username: user.Username, RetryIn: ttl}, nil
		}
		return nil, err
	}

	uuid := util.GenerateUUIDWithDate()
	ttl, err := u.authRepository.AddForgetPasswordCode(uuid, user.Username)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(dto.SMTP_FORGOT_PASS_HTML_PATH)
	if err != nil {
		return nil, domain.ErrChangePasswordInternalError
	}
	url := config.Config.WebUrlUser + "/reset-password/" + uuid
	body := strings.Replace(string(b), "{{URL}}", url, 1)

	mailStruct := util.Mail{
		SenderAddress: config.Config.SmtpConfig.ForgetPasswordAddress,
		SenderName:    dto.SMTP_FORGOT_PASS_SENDER_NAME,
		ToAddress:     user.Email,
		Subject:       dto.SMTP_FORGOT_PASS_SUBJECT,
		Body:          body,
	}

	err = util.SMTPSendMail(mailStruct)
	if err != nil {
		return nil, domain.ErrChangePasswordInternalError
	}

	return &dto.ForgetPasswordRequestVerificationCodeResDTO{IsEmailSent: true, Email: user.Email, Username: user.Username, RetryIn: ttl}, nil
}

func (u *authUsecaseImpl) VerifyForgetPasswordRequest(input *dto.ForgetPasswordVerificationCodeReqDTO) (*dto.ForgetPasswordVerificationCodeResDTO, error) {
	username, err := u.authRepository.GetForgetPasswordCode(input.VerificationCode)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	var newScope = util.ScopeAddTag("", dto.SCOPE_RESET_PASSWORD)
	accessTokenStr, err := u.authUtil.GenerateAccessToken(user.Username, newScope)
	if err != nil {
		return nil, domain.ErrFailedToGenerateAccessToken
	}

	err = u.authRepository.RemoveForgetPasswordCode(input.VerificationCode)
	if err != nil {
		return nil, err
	}

	return &dto.ForgetPasswordVerificationCodeResDTO{AccessToken: accessTokenStr}, nil
}

func (u *authUsecaseImpl) ResetPassword(username string, input *dto.ResetPasswordReqBody) (*dto.ResetPasswordResBody, error) {
	if err := util.ValidatePassword(input.Password, username); err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if util.ValidateHash(user.Password, input.Password) {
		return nil, domain.ErrSamePassword
	}

	hashedPassword, err := util.HashAndSalt(input.Password)
	if err != nil {
		return nil, err
	}

	err = u.authRepository.UpdatePassword(user, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &dto.ResetPasswordResBody{Username: username}, nil
}
