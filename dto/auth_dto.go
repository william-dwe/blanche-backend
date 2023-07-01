package dto

const MAX_RETRY_WALLET = 3
const TIME_LIMIT_BLACKLISTED_TOKEN = 10
const TIME_LIMIT_RESET_PASSWORD_OTP = 10
const TIME_LIMIT_FORGET_PASSWORD_UUID = 10
const TIME_LIMIT_BLOCK_RESET_PASSWORD_REQUEST = 1
const TIME_LIMIT_OAUTH2_STATE = 2

const (
	MANUAL_REGISTER = false
	OAUTH_REGISTER  = true
)

const ACTION_FORGET_PASSWORD = "forget_password"
const ACTION_CHANGE_PASSWORD = "change_password"

const (
	ROLE_USER            = "user"
	ROLE_MERCHANT        = "merchant"
	ROLE_ADMIN           = "admin"
	SCOPE_PIN            = "pin"
	SCOPE_PASSWORD       = "pass"
	SCOPE_OTP            = "otp"
	SCOPE_RESET_PASSWORD = "reset_password"
)

const (
	SMTP_CHANGE_PASS_SENDER_NAME = "Blanche"
	SMTP_CHANGE_PASS_SUBJECT     = "Blanche - Change Password Verification Code"
	SMTP_CHANGE_PASS_HTML_PATH   = "template/email/change_password.html"
	SMTP_FORGOT_PASS_SENDER_NAME = "Blanche"
	SMTP_FORGOT_PASS_SUBJECT     = "Blanche - Forgot Password"
	SMTP_FORGOT_PASS_HTML_PATH   = "template/email/forget_password.html"
)

type StepUpTokenScopeWithPinReqDTO struct {
	Pin string `json:"pin" binding:"required"`
}

type StepUpTokenScopeWithPassReqDTO struct {
	Password string `json:"password" binding:"required"`
}

type UserStepUpTokenScopeResDTO struct {
	AccessToken string `json:"access_token"`
}

type ChangePasswordRequestVerificationCodeResDTO struct {
	IsEmailSent bool   `json:"is_email_sent"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	RetryIn     int    `json:"retry_in"`
}
type ChangePasswordVerificationCodeReqDTO struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}
type ChangePasswordVerificationCodeResDTO struct {
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type ForgetPasswordRequestVerificationCodeReqDTO struct {
	Email string `json:"email"`
}
type ForgetPasswordRequestVerificationCodeResDTO struct {
	IsEmailSent bool   `json:"is_email_sent"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	RetryIn     int    `json:"retry_in"`
}
type ForgetPasswordVerificationCodeReqDTO struct {
	VerificationCode string `json:"verification_code" binding:"required"`
}
type ForgetPasswordVerificationCodeResDTO struct {
	AccessToken string `json:"access_token"`
}

type ResetPasswordReqBody struct {
	Password string `json:"password" binding:"required"`
}

type ResetPasswordResBody struct {
	Username string `json:"username"`
}

type GoogleRequestLoginResDTO struct {
	Url string `json:"url"`
}

type GoogleCallbackReqDTO struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

type UserOAuthLoginReqDTO struct {
	Code string `json:"code"`
}

type UserOAuthLoginResDTO struct {
	Email        string  `json:"email"`
	IsRegistered bool    `json:"is_registered"`
	Name         string  `json:"name"`
	Picture      string  `json:"picture"`
	AccessToken  *string `json:"access_token"`
}

type GoogleUserData struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

type AdminLoginReqDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResDTO struct {
	AccessToken string `json:"access_token"`
}
