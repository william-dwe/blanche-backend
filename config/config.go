package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	ENV_MODE_TEST    = "test"
	ENV_MODE_RELEASE = "release"
)

type authConfig struct {
	AccessTokenExpTimeMinutes     string
	RefreshTokenExpTimeMinutes    string
	AccessTokenSecretString       string
	RefreshTokenSecretString      string
	AdminAccessTokenSecretString  string
	AdminRefreshTokenSecretString string
}

type dbConfig struct {
	Host        string
	User        string
	Password    string
	DBName      string
	Port        string
	GormLogMode string
}

type rdbConfig struct {
	Host     string
	DB       string
	Password string
	Expires  string
}

type envConfig struct {
	Mode                string
	RequestTimeoutLimit int
}

type seaLabsPayConfig struct {
	Url          string
	MerchantCode string
	ApiKey       string
	RedirectUrl  string
	CallbackUrl  string
}

type walletpayConfig struct {
	RedirectUrl string
}

type gcsConfig struct {
	ProjectID  string
	Bucket     string
	UploadPath string
}

type oauthConfig struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	ConfigObj    *oauth2.Config
}

type rajaOngkirConfig struct {
	Url    string
	ApiKey string
}

type smtpConfig struct {
	Host                  string
	Port                  string
	Username              string
	Password              string
	EmailAddress          string
	ResetPasswordAddress  string
	ForgetPasswordAddress string
}

type cronConfig struct {
	IsEnableCron                     bool
	TrxQueueSizeWaitingToCanceled    int
	TrxBatchSizeWaitingToCanceled    int
	TrxQueueSieProcessedToCanceled   int
	TrxBatchSizeProcessedToCanceled  int
	TrxQueueSizeDeliveredToCompleted int
	TrxBatchSizeDeliveredToCompleted int
}

type AppConfig struct {
	AppName           string
	AppUrlUser        string
	AppUrlAdmin       string
	WebUrlUser        string
	WebUrlAdmin       string
	WebTransactionURL string
	DBConfig          dbConfig
	RDBConfig         rdbConfig
	AuthConfig        authConfig
	RajaOngkirConfig  rajaOngkirConfig
	ENVConfig         envConfig
	SeaLabsPayConfig  seaLabsPayConfig
	GCSConfig         gcsConfig
	OauthConfig       oauthConfig
	SmtpConfig        smtpConfig
	WalletpayConfig   walletpayConfig
	CronConfig        cronConfig
}

func getENV(key, defaultVal string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	return env
}

func getENVinteger(key string, defaultVal int) int {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(env)
	if err != nil {
		return defaultVal
	}
	return val
}

func getENVbool(key string, defaultVal bool) bool {
	env := os.Getenv(key)
	if env == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(env)
	if err != nil {
		return defaultVal
	}
	return val
}

func initLogConfig() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Error loading .env file")
	}
	log.Info().Msg(".env file Loaded")
}

func initTimezone() {
	offsetHour, err := strconv.Atoi(getENV("TIMEZONE_OFFSET_HOUR", "7"))
	if err != nil {
		log.Error().Msg("Error parsing timezone offset hour")
		offsetHour = 7
	}
	offsetMinute, err := strconv.Atoi(getENV("TIMEZONE_OFFSET_MINUTE", "0"))
	if err != nil {
		log.Error().Msg("Error parsing timezone offset minute")
		offsetMinute = 0
	}

	time.Local = time.FixedZone(getENV("TIMEZONE_LOCATION", "Asia/Jakarta"), offsetHour*60*60+offsetMinute*60)
}

func initAllConfig() AppConfig {
	initLogConfig()
	initEnv()

	return AppConfig{
		AppName:           getENV("APP_NAME", "blanche"),
		AppUrlUser:        getENV("APP_URL_USER", ""),
		AppUrlAdmin:       getENV("APP_URL_ADMIN", ""),
		WebUrlUser:        getENV("WEB_URL_USER", "*"),
		WebUrlAdmin:       getENV("WEB_URL_ADMIN", "*"),
		WebTransactionURL: getENV("PAYMENT_TRANSACTION_REDIRECT", ""),

		AuthConfig: authConfig{
			AccessTokenExpTimeMinutes:     getENV("ACCESS_TOKEN_EXP_MINUTES", ""),
			RefreshTokenExpTimeMinutes:    getENV("REFRESH_TOKEN_EXP_MINUTES", ""),
			AccessTokenSecretString:       getENV("ACCESS_TOKEN_SECRET", ""),
			RefreshTokenSecretString:      getENV("REFRESH_TOKEN_SECRET", ""),
			AdminAccessTokenSecretString:  getENV("ADMIN_ACCESS_TOKEN_SECRET", ""),
			AdminRefreshTokenSecretString: getENV("ADMIN_REFRESH_TOKEN_SECRET", ""),
		},

		DBConfig: dbConfig{
			Host:        getENV("CONF_DB_HOST", ""),
			User:        getENV("CONF_DB_USER", ""),
			Password:    getENV("CONF_DB_PASSWORD", ""),
			DBName:      getENV("CONF_DB_NAME", ""),
			Port:        getENV("CONF_DB_PORT", ""),
			GormLogMode: getENV("GORM_LOG_MODE", ""),
		},

		RDBConfig: rdbConfig{
			Host:     getENV("RDB_HOST", ""),
			DB:       getENV("RDB_DB", ""),
			Password: getENV("RDB_PASSWORD", ""),
			Expires:  getENV("RDB_EXPIRES", ""),
		},

		ENVConfig: envConfig{
			Mode:                getENV("APP_ENV_MODE", ENV_MODE_RELEASE),
			RequestTimeoutLimit: getENVinteger("REQUEST_TIMEOUT_LIMIT", 60000),
		},

		SeaLabsPayConfig: seaLabsPayConfig{
			Url:          getENV("SEALABSPAY_URL", ""),
			MerchantCode: getENV("SEALABSPAY_MERCHANT", ""),
			ApiKey:       getENV("SEALABSPAY_API_KEY", ""),
			RedirectUrl:  getENV("SEALABSPAY_REDIRECT", ""),
			CallbackUrl:  getENV("SEALABSPAY_CALLBACK", ""),
		},

		RajaOngkirConfig: rajaOngkirConfig{
			Url:    getENV("RAJAONGKIR_URL", ""),
			ApiKey: getENV("RAJAONGKIR_API_KEY", ""),
		},

		SmtpConfig: smtpConfig{
			Host:                  getENV("SMTP_HOST", ""),
			Port:                  getENV("SMTP_PORT", ""),
			Username:              getENV("SMTP_USERNAME", ""),
			Password:              getENV("SMTP_PASSWORD", ""),
			EmailAddress:          getENV("SMTP_OTP_EMAIL_ADDRESS", ""),
			ResetPasswordAddress:  getENV("SMTP_RESET_PASSWORD_ADDRESS", ""),
			ForgetPasswordAddress: getENV("SMTP_FORGET_PASSWORD_ADDRESS", ""),
		},

		WalletpayConfig: walletpayConfig{
			RedirectUrl: getENV("WALLETPAY_REDIRECT", ""),
		},

		GCSConfig: gcsConfig{
			ProjectID:  getENV("GCS_PROJECT_ID", ""),
			Bucket:     getENV("GCS_BUCKET", ""),
			UploadPath: getENV("GCS_UPLOAD_PATH", ""),
		},

		OauthConfig: oauthConfig{
			RedirectURL:  getENV("OAUTH_REDIRECT_URL", ""),
			ClientID:     getENV("OAUTH_CLIENT_ID", ""),
			ClientSecret: getENV("OAUTH_CLIENT_SECRET", ""),
			ConfigObj: &oauth2.Config{
				RedirectURL:  getENV("OAUTH_REDIRECT_URL", ""),
				ClientID:     getENV("OAUTH_CLIENT_ID", ""),
				ClientSecret: getENV("OAUTH_CLIENT_SECRET", ""),
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint: google.Endpoint,
			},
		},

		CronConfig: cronConfig{
			IsEnableCron:                     getENVbool("IS_ENABLE_CRON", false),
			TrxQueueSizeWaitingToCanceled:    getENVinteger("TRX_QUEUE_SIZE_WAITING_TO_CANCELED", 100),
			TrxBatchSizeWaitingToCanceled:    getENVinteger("TRX_BATCH_SIZE_WAITING_TO_CANCELED", 25),
			TrxQueueSieProcessedToCanceled:   getENVinteger("TRX_QUEUE_SIZE_PROCESSED_TO_CANCELED", 100),
			TrxBatchSizeProcessedToCanceled:  getENVinteger("TRX_BATCH_SIZE_PROCESSED_TO_CANCELED", 25),
			TrxQueueSizeDeliveredToCompleted: getENVinteger("TRX_QUEUE_SIZE_DELIVERED_TO_COMPLETED", 100),
			TrxBatchSizeDeliveredToCompleted: getENVinteger("TRX_BATCH_SIZE_DELIVERED_TO_COMPLETED", 25),
		},
	}
}

var Config = initAllConfig()
