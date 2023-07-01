package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cache"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetRoleByRoleName(roleName string) (*entity.Role, error)
	GetRoleByRoleId(roleId uint) (*entity.Role, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)

	GetUserByUserId(userId uint) (*entity.User, error)
	GetUserAndUserDetailByUserId(userId uint) (*entity.User, error)
	AddUser(user entity.User) (*entity.User, error)
	UpdateUser(user entity.User) error
	AddUserEmailBlacklist(UserEmailBlacklist entity.UserEmailBlacklist) error
	CheckEmailBlacklist(email string) error
	UpdateUserDetail(userDetail entity.UserDetail) error

	GetAllUserAddress(user entity.User) ([]entity.UserAddress, error)
	GetDefaultUserAddress(user entity.User) (*entity.UserAddress, error)
	GetUserAddressById(userId, userAddressId uint) (*entity.UserAddress, error)
	AddUserAddress(userAddress entity.UserAddress) (*entity.UserAddress, error)
	SetDefaultUserAddress(userId uint, userAddressId uint) error
	UpdateUserAddress(userAddress, newUserAddress entity.UserAddress) (*entity.UserAddress, error)
	DeleteUserAddress(userAddress entity.UserAddress) (*entity.UserAddress, error)

	UpdateUserRoleToMerchantTx(txGorm *gorm.DB, userId uint) error
}

type UserRepositoryConfig struct {
	DB  *gorm.DB
	RDB *cache.RDBConnection
}

type userRepositoryImpl struct {
	db  *gorm.DB
	rdb *cache.RDBConnection
}

func NewUserRepository(c UserRepositoryConfig) UserRepository {
	return &userRepositoryImpl{
		db:  c.DB,
		rdb: c.RDB,
	}
}

func (r *userRepositoryImpl) GetRoleByRoleName(roleName string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.
		Where("role_name = ?", roleName).
		First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetRoleInvalidInput
		}
		return nil, domain.ErrGetRoleInternalError
	}
	return &role, err
}

func (r *userRepositoryImpl) GetRoleByRoleId(roleId uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.
		Where("id = ?", roleId).
		First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetRoleInvalidInput
		}
		return nil, domain.ErrGetRoleInternalError
	}
	return &role, err
}

func (r *userRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCheckEmailInvalidInput
		}
		return nil, domain.ErrGetUserInternalError
	}
	return &user, err
}

func (r *userRepositoryImpl) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetUserNotFound
		}

		return nil, domain.ErrGetUser
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetUserByUserId(userId uint) (*entity.User, error) {
	var user entity.User
	err := r.db.
		Where("id = ?", userId).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCheckEmailInvalidInput
		}
		return nil, domain.ErrGetUserInternalError
	}
	return &user, err
}

func (r *userRepositoryImpl) GetUserAndUserDetailByUserId(userId uint) (*entity.User, error) {
	var user entity.User
	err := r.db.
		Where("id = ?", userId).
		Preload("Role").
		Preload("UserDetail").
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCheckEmailInvalidInput
		}
		return nil, domain.ErrGetUserInternalError
	}
	return &user, err
}

func (r *userRepositoryImpl) AddUser(user entity.User) (*entity.User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return nil, util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"users_role_id_fkey":          domain.ErrRoleNotFound,
				"user_details_user_id_fkey":   domain.ErrUserDetailNotFound,
				"users_username_check":        domain.ErrCheckUsernameInvalidInput,
				"users_username_key":          domain.ErrUsernameAlreadyExist,
				"users_email_key":             domain.ErrUserEmailAlreadyExist,
				"user_details_fullname_check": domain.ErrFullnameWrongFormat,
				"user_details_phone_check":    domain.ErrPhoneWrongFormat,
				"user_details_phone_key":      domain.ErrPhoneAlreadyExist,
			},
			domain.ErrRegister,
		)
	}

	return &user, err
}

func (r *userRepositoryImpl) UpdateUser(user entity.User) error {
	err := r.db.
		Where("id = ?", user.ID).
		Updates(&user).
		Error

	if err != nil {
		return util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"users_role_id_fkey":   domain.ErrRoleNotFound,
				"users_username_check": domain.ErrCheckUsernameInvalidInput,
				"users_username_key":   domain.ErrUsernameAlreadyExist,
				"users_email_key":      domain.ErrUserEmailAlreadyExist,
			},
			domain.ErrUpdateErrorInternalError,
		)
	}

	return nil
}

func (r *userRepositoryImpl) AddUserEmailBlacklist(UserEmailBlacklist entity.UserEmailBlacklist) error {
	err := r.db.
		Create(&UserEmailBlacklist).
		Error
	if err != nil {
		return util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_email_blacklists_email_key":    domain.ErrUserEmailBlacklistAlreadyExist,
				"user_email_blacklists_user_id_fkey": domain.ErrUserEmailBlacklistUserNotFound,
			},
			domain.ErrAddUserEmailBlacklistInternalError,
		)
	}
	return nil
}

func (r *userRepositoryImpl) CheckEmailBlacklist(email string) error {
	var UserEmailBlacklist entity.UserEmailBlacklist
	err := r.db.
		Where("email = ?", email).
		First(&UserEmailBlacklist).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return domain.ErrCheckEmailBlacklistInternalError
	}
	return domain.ErrCheckEmailBlacklistEmailBlacklisted
}

func (r *userRepositoryImpl) UpdateUserDetail(userDetail entity.UserDetail) error {
	err := r.db.
		Where("user_id = ?", userDetail.UserId).
		Updates(&userDetail).
		Error

	if err != nil {
		return util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_details_user_id_fkey":   domain.ErrUserDetailNotFound,
				"user_details_fullname_check": domain.ErrFullnameWrongFormat,
				"user_details_phone_check":    domain.ErrPhoneWrongFormat,
				"user_details_phone_key":      domain.ErrPhoneAlreadyExist,
			},
			domain.ErrUpdateErrorInternalError,
		)
	}

	return nil
}

func (r *userRepositoryImpl) GetAllUserAddress(user entity.User) ([]entity.UserAddress, error) {
	var userAddress []entity.UserAddress
	err := r.db.Model(&user).
		Preload("Subdistrict").
		Preload("City").
		Preload("District").
		Preload("Province").
		Association("UserAddress").
		Find(&userAddress)
	if err != nil {
		return nil, domain.ErrUserAddressNotFound
	}
	return userAddress, nil
}

func (r *userRepositoryImpl) GetDefaultUserAddress(user entity.User) (*entity.UserAddress, error) {
	var userAddress entity.UserAddress
	err := r.db.Model(&user).
		Where("is_default = ?", true).
		Preload("Subdistrict").
		Preload("City").
		Preload("District").
		Preload("Province").
		Association("UserAddress").
		Find(&userAddress)
	if err != nil || userAddress.ID == 0 {
		return nil, domain.ErrUserAddressNotFound
	}

	return &userAddress, nil
}

func (r *userRepositoryImpl) GetUserAddressById(userId, userAddressId uint) (*entity.UserAddress, error) {
	var userAddress entity.UserAddress
	err := r.db.
		Where("user_id = ? AND id = ?", userId, userAddressId).
		Preload("Subdistrict").
		Preload("City").
		Preload("District").
		Preload("Province").
		First(&userAddress).Error
	if err != nil {
		return nil, domain.ErrUserAddressNotFound
	}
	return &userAddress, nil
}

func (r *userRepositoryImpl) AddUserAddress(userAddress entity.UserAddress) (*entity.UserAddress, error) {
	err := r.db.Create(&userAddress).Error
	if err != nil {
		return nil, util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_addresses_user_id_fkey": domain.ErrUserNotFound,
				"user_addresses_phone_check":  domain.ErrAddUserAddressInvalidPhone,
			},
			domain.ErrAddUserAddress,
		)
	}
	return &userAddress, nil
}

func (r *userRepositoryImpl) SetDefaultUserAddress(userId uint, userAddressId uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.UserAddress{}).
			Where("user_id = ?", userId).
			Update("is_default", false).Error
		if err != nil {
			return err
		}

		err = tx.Model(&entity.UserAddress{}).
			Where("id = ?", userAddressId).
			Updates(map[string]interface{}{
				"is_default": true,
			}).Error
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (r *userRepositoryImpl) UpdateUserAddress(userAddress, newUserAddress entity.UserAddress) (*entity.UserAddress, error) {
	err := r.db.Model(&userAddress).Updates(&newUserAddress).Error
	if err != nil {
		return nil, util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_addresses_user_id_fkey": domain.ErrUserNotFound,
				"user_addresses_phone_check":  domain.ErrAddUserAddressInvalidPhone,
			},
			domain.ErrUpdateUserAddress,
		)
	}

	return &userAddress, nil
}

func (r *userRepositoryImpl) DeleteUserAddress(userAddress entity.UserAddress) (*entity.UserAddress, error) {
	err := r.db.Delete(&userAddress).Error
	if err != nil {
		return nil, util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"user_addresses_user_id_fkey": domain.ErrUserNotFound,
			},
			domain.ErrDeleteUserAddress,
		)
	}
	return &userAddress, nil
}

func (r *userRepositoryImpl) UpdateUserRoleToMerchantTx(txGorm *gorm.DB, userId uint) error {
	sq := txGorm.Model(&entity.Role{}).Where("role_name = ?", "merchant").Select("id")
	err := txGorm.Model(&entity.User{}).Where("id = ?", userId).Update("role_id", sq).Error

	if err != nil {
		return util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"users_role_id_fkey": domain.ErrRoleNotFound,
			},
			domain.ErrUpdateErrorInternalError,
		)
	}

	return nil
}
