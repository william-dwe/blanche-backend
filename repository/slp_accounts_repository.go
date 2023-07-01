package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type SlpAccountsRepository interface {
	GetSlpAccountListByUserId(userId uint) ([]entity.SlpAccount, error)
	GetUserSlpAccountByID(userId, slpAccountId uint) (*entity.SlpAccount, error)
	RegisterSlpAccount(slpAccount entity.SlpAccount) (*entity.SlpAccount, error)
	SetDefaultSlpAccount(userId, slpAccountId uint) error
	DeleteUserSlpAccount(slpAccount *entity.SlpAccount) (*entity.SlpAccount, error)
}

type SlpAccountsRepositoryConfig struct {
	DB *gorm.DB
}

type slpAccountsRepositoryImpl struct {
	db *gorm.DB
}

func NewSlpAccountsRepository(c SlpAccountsRepositoryConfig) SlpAccountsRepository {
	return &slpAccountsRepositoryImpl{
		db: c.DB,
	}
}

func (r *slpAccountsRepositoryImpl) GetSlpAccountListByUserId(userId uint) ([]entity.SlpAccount, error) {
	var slpAccounts []entity.SlpAccount
	err := r.db.Model(&slpAccounts).Where("user_id = ?", userId).Find(&slpAccounts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSlpAccountsNotFound
		}

		return nil, domain.ErrGetSlpAccounts
	}

	return slpAccounts, nil
}

func (r *slpAccountsRepositoryImpl) GetUserSlpAccountByID(userId, slpAccountId uint) (*entity.SlpAccount, error) {
	var slpAccount entity.SlpAccount
	err := r.db.Model(&slpAccount).Where("user_id = ? AND id = ?", userId, slpAccountId).First(&slpAccount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrSlpAccountNotFound
		}

		return nil, domain.ErrGetSlpAccount
	}

	return &slpAccount, nil
}

func (r *slpAccountsRepositoryImpl) RegisterSlpAccount(slpAccount entity.SlpAccount) (*entity.SlpAccount, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.SlpAccount{}).
			Where("user_id = ? AND card_number = ?", slpAccount.UserID, slpAccount.CardNumber).
			First(&entity.SlpAccount{}).Error

		if err == nil {
			return domain.ErrDuplicateSlpAccount
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return domain.ErrGetSlpAccount
		}

		if err == gorm.ErrRecordNotFound {
			err = tx.Model(&entity.SlpAccount{}).Create(&slpAccount).Error
			if err != nil {
				maskedErr := util.PgConsErrMasker(
					err,
					entity.ConstraintErrMaskerMap{
						"slp_accounts_cardnum_check": domain.ErrInvalidCardNumber,
						"slp_accounts_pkey":          domain.ErrDuplicateSlpAccount,
					},
					domain.ErrRegisterSlpAccount,
				)
				return maskedErr
			}
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &slpAccount, nil
}

func (r *slpAccountsRepositoryImpl) SetDefaultSlpAccount(userId, slpAccountId uint) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity.SlpAccount{}).
			Where("user_id = ?", userId).
			Update("is_default", false).Error
		if err != nil {
			return err
		}

		err = tx.Model(&entity.SlpAccount{}).
			Where("id = ?", slpAccountId).
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

func (r *slpAccountsRepositoryImpl) DeleteUserSlpAccount(slpAccount *entity.SlpAccount) (*entity.SlpAccount, error) {
	err := r.db.Model(&slpAccount).Delete(&slpAccount).Error
	if err != nil {
		return nil, domain.ErrDeleteSlpAccount
	}

	return slpAccount, nil
}
