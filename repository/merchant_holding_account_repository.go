package repository

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type MerchantHoldingAccountRepository interface {
	UpdateBalanceTx(tx *gorm.DB, merchantId uint, amount int) (*entity.MerchantHoldingAccount, int64, error)
	WithdrawBalance(userId uint, merchantId uint, amount int) (*entity.MerchantHoldingAccountHistory, error)
	GetByMerchantId(merchantId uint) (*entity.MerchantHoldingAccount, error)
	GetByUserId(userId uint) (*entity.MerchantHoldingAccount, error)
}

type MerchantHoldingAccountRepositoryImpl struct {
	db               *gorm.DB
	walletRepository WalletRepository
}

type MerchantHoldingAccountRepositoryConfig struct {
	DB               *gorm.DB
	WalletRepository WalletRepository
}

func NewMerchantHoldingAccountRepository(c MerchantHoldingAccountRepositoryConfig) MerchantHoldingAccountRepository {
	return &MerchantHoldingAccountRepositoryImpl{
		db:               c.DB,
		walletRepository: c.WalletRepository,
	}
}

func (r *MerchantHoldingAccountRepositoryImpl) UpdateBalanceTx(tx *gorm.DB, merchantId uint, amount int) (*entity.MerchantHoldingAccount, int64, error) {
	var merchantHoldingAcc entity.MerchantHoldingAccount
	res := tx.Raw(fmt.Sprintf(`
			UPDATE "merchant_holding_accounts" SET "balance"=balance + %d,"updated_at"=now() WHERE merchant_id = %d AND "merchant_holding_accounts"."deleted_at" IS null
			returning id, merchant_id, balance, created_at, updated_at
		`, amount, merchantId)).Scan(&merchantHoldingAcc)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	return &merchantHoldingAcc, res.RowsAffected, nil
}

func (r *MerchantHoldingAccountRepositoryImpl) GetByMerchantId(merchantId uint) (*entity.MerchantHoldingAccount, error) {
	var merchantHoldingAcc entity.MerchantHoldingAccount
	res := r.db.Where("merchant_id = ?", merchantId).First(&merchantHoldingAcc)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantHoldingAccountNotFound
		}

		return nil, domain.ErrMerchantHoldingAccount
	}

	return &merchantHoldingAcc, nil
}

func (r *MerchantHoldingAccountRepositoryImpl) GetByUserId(userId uint) (*entity.MerchantHoldingAccount, error) {
	var merchantHoldingAcc entity.MerchantHoldingAccount
	res := r.db.Model(&merchantHoldingAcc).Joins("left join merchants on merchants.id = merchant_id").Where("user_id = ?", userId).First(&merchantHoldingAcc)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantHoldingAccountNotFound
		}

		return nil, domain.ErrMerchantHoldingAccount
	}

	return &merchantHoldingAcc, nil
}

func (r *MerchantHoldingAccountRepositoryImpl) WithdrawBalance(userId uint, merchantId uint, amount int) (resWd *entity.MerchantHoldingAccountHistory, errWd error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in WithdrawBalance repo: %v", r)
			errWd = domain.ErrMerchantHoldingAccWithdraw
		}
	}()

	// reduce balance from merchant holding account
	updatedHoldingAcc, rowsAffected, err := r.UpdateBalanceTx(tx, merchantId, -amount)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Failed to reduce merchant holding account balance: %v", err)
		return nil, err
	}
	if updatedHoldingAcc.Balance < 0 || rowsAffected <= 0 {
		tx.Rollback()
		log.Error().Msgf("Failed to reduce merchant holding account balance unsufficient: %v", err)
		return nil, domain.ErrMerchantHoldingAccInsufficientBalance
	}

	// add history trx merchant holding account
	newHistoryMerchantAcc := entity.MerchantHoldingAccountHistory{
		MerchantHoldingAccountID: updatedHoldingAcc.ID,
		Amount:                   float64(amount),
		Type:                     dto.MERCHANT_HOLDING_ACC_DEBIT_CODE,
		Notes:                    "Withdrawal to wallet",
	}
	err = tx.Create(&newHistoryMerchantAcc).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Failed to add merchant holding account history: %v", err)
		return nil, domain.ErrMerchantHoldingAccWithdrawAddHistory
	}

	// add history trx wallet and update balance
	err = r.walletRepository.AddWalletHistoryMerchantWdTx(tx, userId, float64(amount))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Failed to update balance and add wallet history merchant withdrawal: %v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Failed to commit merchant holding account withdraw: %v", err)
		return nil, domain.ErrMerchantHoldingAccWithdraw
	}

	return &newHistoryMerchantAcc, nil
}
