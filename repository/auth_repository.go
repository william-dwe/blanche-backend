package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cache"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckUserEmailExistence(email string) (bool, error)
	CheckUserUsernameExistence(username string) (bool, error)

	GetGoogleToken(googleConfig *oauth2.Config, code string) (*oauth2.Token, error)
	GetUserDataFromGoogle(code string) (*dto.GoogleUserData, error)

	AddUserLoginActivity(loginLog entity.UserLoginActivity) error
	GetUserIdByRefreshToken(refreshToken string) (*uint, error)
	DeleteUserActivityByRefreshToken(refreshToken string) error

	GetBlacklistedToken(cacheKey string) (bool, error)
	AddBlacklistedToken(token string) error

	BlockResetPasswordRequestTx(tx *gorm.DB, username string, action string) error
	CheckBlockResetPasswordRequest(username string, action string) error
	GetTTLBlockPasswordCode(username, action string) (int, error)

	GetChangePasswordCode(username string) (string, error)
	AddChangePasswordCode(hashedOTP string, username string) (attempt int, err error)
	GetTTLChangePasswordCodeTx(tx *gorm.DB, username, action string) (int, error)
	RemoveChangePasswordCode(username string) error
	GetForgetPasswordCode(uuid string) (string, error)
	AddForgetPasswordCode(uuid string, username string) (int, error)
	RemoveForgetPasswordCode(uuid string) error
	InvalidateRefreshTokenOnResetPasswordTx(tx *gorm.DB, userId uint) error
	UpdatePassword(user *entity.User, hashedPassword string) error

	BlockWalletGetAttempt(username string) (int, error)
	BlockWalletAddAttempt(username string) error
	BlockWalletResetAttempt(username string) error
}

type AuthRepositoryConfig struct {
	DB  *gorm.DB
	RDB *cache.RDBConnection
}

type authRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.RDBConnection
}

func NewAuthRepository(c AuthRepositoryConfig) AuthRepository {
	return &authRepositoryImpl{
		db:  c.DB,
		rdb: c.RDB,
	}
}

func (r *authRepositoryImpl) CheckUserEmailExistence(email string) (bool, error) {
	var user entity.User
	err := r.db.
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, domain.ErrCheckUserInternalServer
	}
	return false, nil
}

func (r *authRepositoryImpl) CheckUserUsernameExistence(username string) (bool, error) {
	var user entity.User
	err := r.db.
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, domain.ErrCheckUserInternalServer
	}
	return false, nil
}

func (r *authRepositoryImpl) GetGoogleToken(googleConfig *oauth2.Config, code string) (*oauth2.Token, error) {
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, domain.ErrGetUserDataFromGoogleInternalError
	}
	return token, nil
}

func (r *authRepositoryImpl) GetUserDataFromGoogle(code string) (*dto.GoogleUserData, error) {
	googleConfig := config.Config.OauthConfig.ConfigObj
	token, err := r.GetGoogleToken(googleConfig, code)
	if err != nil {
		return nil, err
	}

	client := googleConfig.Client(context.Background(), token)
	gApiRes, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, domain.ErrGetUserDataFromGoogleInternalError
	}
	defer gApiRes.Body.Close()

	gApiBody, err := io.ReadAll(gApiRes.Body)
	if err != nil {
		return nil, domain.ErrGetUserDataFromGoogleInternalError
	}

	var gApiUserData dto.GoogleUserData
	err = json.Unmarshal(gApiBody, &gApiUserData)
	if err != nil {
		return nil, domain.ErrGetUserDataFromGoogleInternalError
	}
	return &gApiUserData, nil
}

func (r *authRepositoryImpl) AddUserLoginActivity(loginLog entity.UserLoginActivity) error {
	err := r.db.Create(&loginLog).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_email_blacklists_user_id_fkey": domain.ErrUserNotFound,
				"user_email_blacklists_email_key":    domain.ErrUpdateProfileDuplicateEmailBlackList,
			},
			domain.ErrRegister,
		)
		return maskedErr
	}
	return err
}

func (r *authRepositoryImpl) GetUserIdByRefreshToken(refreshToken string) (*uint, error) {
	var userLog entity.UserLoginActivity
	err := r.db.
		Where("refresh_token = ?", refreshToken).
		First(&userLog).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrCheckUserInternalServer
	}
	return &userLog.UserId, err
}

func (r *authRepositoryImpl) DeleteUserActivityByRefreshToken(refreshToken string) error {
	err := r.db.
		Where("refresh_token = ?", refreshToken).
		Delete(&entity.UserLoginActivity{}).Error
	return err
}

func (r *authRepositoryImpl) GetBlacklistedToken(cacheKey string) (bool, error) {
	var tokenStatus = true
	err := r.rdb.GetCache(cacheKey, &tokenStatus)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return tokenStatus, nil
		}
		log.Error().Msgf("Error get blacklist token: %v", err)
		return tokenStatus, domain.ErrBlacklistToken
	}
	return tokenStatus, nil
}

func (r *authRepositoryImpl) AddBlacklistedToken(token string) error {
	err := r.rdb.SetCache("blacklist_token:"+token, true, dto.TIME_LIMIT_BLACKLISTED_TOKEN)
	if err != nil {
		log.Error().Msgf("Error blacklist token: %v", err)
		return domain.ErrBlacklistToken
	}
	return nil
}

func (r *authRepositoryImpl) BlockResetPasswordRequestTx(tx *gorm.DB, username string, action string) error {
	err := r.rdb.SetCache("block-"+action+":"+username, false, dto.TIME_LIMIT_BLOCK_RESET_PASSWORD_REQUEST)
	if err != nil {
		log.Error().Msgf("Error block reset password: %v", err)
		return domain.ErrChangePasswordInternalError
	}
	return nil
}

func (r *authRepositoryImpl) CheckBlockResetPasswordRequest(username string, action string) error {
	var isValid bool
	err := r.rdb.GetCache("block-"+action+":"+username, &isValid)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		log.Error().Msgf("Error check block reset password: %v", err)
		return domain.ErrChangePasswordInternalError
	}

	return domain.ErrResetPasswordWait
}

func (r *authRepositoryImpl) GetTTLBlockPasswordCode(username, action string) (int, error) {
	ttl, err := r.rdb.GetTTLCache("block-" + action + ":" + username)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, domain.ErrResetPasswordCodeExpired
		}
		log.Error().Msgf("Error get ttl reset password: %v", err)
		return 0, domain.ErrChangePasswordInternalError
	}
	return ttl, nil
}

func (r *authRepositoryImpl) GetChangePasswordCode(username string) (string, error) {
	var hashedOTP string
	err := r.rdb.GetCache(dto.ACTION_CHANGE_PASSWORD+":"+username, &hashedOTP)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrResetPasswordCodeExpired
		}
		log.Error().Msgf("Error get reset password: %v", err)
		return "", domain.ErrChangePasswordInternalError
	}
	return hashedOTP, nil
}

func (r *authRepositoryImpl) AddChangePasswordCode(hashedOTP string, username string) (attempt int, errChangePw error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errChangePw = domain.ErrChangePasswordInternalError
		}
	}()

	err := r.rdb.SetCache(dto.ACTION_CHANGE_PASSWORD+":"+username, hashedOTP, dto.TIME_LIMIT_RESET_PASSWORD_OTP)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error change password: %v", err)
		return 0, domain.ErrChangePasswordInternalError
	}

	err = r.BlockResetPasswordRequestTx(tx, username, dto.ACTION_CHANGE_PASSWORD)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	ttl, err := r.GetTTLBlockPasswordCode(username, dto.ACTION_CHANGE_PASSWORD)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit().Error
	if err != nil {
		return 0, domain.ErrChangePasswordInternalError
	}

	return ttl, nil
}

func (r *authRepositoryImpl) GetTTLChangePasswordCodeTx(tx *gorm.DB, username string, action string) (int, error) {
	ttl, err := r.rdb.GetTTLCache(action + ":" + username)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, domain.ErrResetPasswordCodeExpired
		}
		log.Error().Msgf("Error get ttl reset password: %v", err)
		return 0, domain.ErrChangePasswordInternalError
	}
	return ttl, nil
}

func (r *authRepositoryImpl) RemoveChangePasswordCode(username string) error {
	err := r.rdb.DeleteCache(dto.ACTION_CHANGE_PASSWORD + ":" + username)
	if err != nil {
		log.Error().Msgf("Error remove reset password: %v", err)
		return domain.ErrChangePasswordInternalError
	}
	return nil
}

func (r *authRepositoryImpl) GetForgetPasswordCode(uuid string) (string, error) {
	var username string
	err := r.rdb.GetCache(dto.ACTION_FORGET_PASSWORD+":"+uuid, &username)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrResetPasswordCodeExpired
		}
		log.Error().Msgf("Error get reset password: %v", err)
		return "", domain.ErrChangePasswordInternalError
	}
	return username, nil
}

func (r *authRepositoryImpl) AddForgetPasswordCode(uuid string, username string) (att int, errForgetPw error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errForgetPw = domain.ErrChangePasswordInternalError
		}
	}()

	err := r.rdb.SetCache(dto.ACTION_FORGET_PASSWORD+":"+uuid, username, dto.TIME_LIMIT_FORGET_PASSWORD_UUID)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error change password: %v", err)
		return 0, domain.ErrChangePasswordInternalError
	}

	err = r.BlockResetPasswordRequestTx(tx, username, dto.ACTION_FORGET_PASSWORD)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	ttl, err := r.GetTTLBlockPasswordCode(username, dto.ACTION_FORGET_PASSWORD)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit().Error
	if err != nil {
		return 0, domain.ErrChangePasswordInternalError
	}

	return ttl, nil
}

func (r *authRepositoryImpl) RemoveForgetPasswordCode(uuid string) error {
	err := r.rdb.DeleteCache(dto.ACTION_FORGET_PASSWORD + ":" + uuid)
	if err != nil {
		log.Error().Msgf("Error remove reset password: %v", err)
		return domain.ErrChangePasswordInternalError
	}
	return nil
}

func (r *authRepositoryImpl) InvalidateRefreshTokenOnResetPasswordTx(tx *gorm.DB, userId uint) error {
	err := r.db.
		Where("user_id = ?", userId).
		Delete(&entity.UserLoginActivity{}).Error
	if err != nil {
		log.Error().Msgf("Error delete user login activity: %v", err)
		return domain.ErrChangePasswordInternalError
	}
	return nil
}

func (r *authRepositoryImpl) UpdatePassword(user *entity.User, hashedPassword string) (errUpdatePw error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errUpdatePw = domain.ErrChangePasswordInternalError
		}
	}()
	err := r.db.Model(&entity.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password": hashedPassword,
		}).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrChangePasswordInternalError
	}

	err = r.InvalidateRefreshTokenOnResetPasswordTx(tx, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrRegister
	}

	return nil
}

func (r *authRepositoryImpl) BlockWalletGetAttempt(username string) (int, error) {
	var attempt int
	err := r.rdb.GetCache("wallet_attempt"+username, &attempt)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return attempt, domain.ErrPinInternalError
	}
	return attempt, nil
}

func (r *authRepositoryImpl) BlockWalletAddAttempt(username string) error {
	attempt, err := r.BlockWalletGetAttempt(username)
	if err != nil {
		return err
	}

	err = r.rdb.SetCache("wallet_attempt"+username, attempt+1, dto.TIME_LIMIT_BLOCK_WALLET)
	if err != nil {
		return domain.ErrPinInternalError
	}

	return nil
}

func (r *authRepositoryImpl) BlockWalletResetAttempt(username string) error {
	err := r.rdb.DeleteCache("wallet_attempt" + username)
	if err != nil {
		return domain.ErrPinInternalError
	}

	return nil
}
