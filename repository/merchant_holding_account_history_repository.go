package repository

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type MerchantHoldingAccountHistoryRepository interface {
	AddHistoryTx(tx *gorm.DB, merchantHoldingAccId uint, amount float64, notes string) error
	GetHistoryByMerchantHoldingAccId(merchantHoldingAccId uint, reqParam dto.MerchantFundActivitiesReqParamDTO) ([]entity.MerchantHoldingAccountHistory, int64, error)
}

type MerchantHoldingAccountHistoryRepositoryConfig struct {
	DB *gorm.DB
}

type merchantHoldingAccountHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewMerchantHoldingAccountHistoryRepository(c MerchantHoldingAccountHistoryRepositoryConfig) MerchantHoldingAccountHistoryRepository {
	return &merchantHoldingAccountHistoryRepositoryImpl{
		db: c.DB,
	}
}

func (m *merchantHoldingAccountHistoryRepositoryImpl) AddHistoryTx(tx *gorm.DB, merchantHoldingAccId uint, amount float64, notes string) error {
	codeHistory := dto.MERCHANT_HOLDING_ACC_CREDIT_CODE
	if amount < 0 {
		codeHistory = dto.MERCHANT_HOLDING_ACC_DEBIT_CODE
		amount = -amount
	}

	newHistoryMerchantAcc := entity.MerchantHoldingAccountHistory{
		MerchantHoldingAccountID: merchantHoldingAccId,
		Amount:                   float64(amount),
		Type:                     codeHistory,
		Notes:                    notes,
	}

	err := tx.Create(&newHistoryMerchantAcc).Error

	return err
}

func (m *merchantHoldingAccountHistoryRepositoryImpl) GetHistoryByMerchantHoldingAccId(merchantHoldingAccId uint, req dto.MerchantFundActivitiesReqParamDTO) ([]entity.MerchantHoldingAccountHistory, int64, error) {
	qDateTo := ""
	qDateFrom := ""
	if req.StartDate != "" {
		qDateFrom = fmt.Sprintf("created_at >= '%s'", req.StartDate)
	}
	if req.EndDate != "" {
		qDateTo = fmt.Sprintf("created_at <= '%s'", req.EndDate)
	}

	PageOffset := req.Limit * (req.Page - 1)
	var countData int64

	var history []entity.MerchantHoldingAccountHistory
	res := m.db.
		Where("merchant_holding_account_id = ?", merchantHoldingAccId).
		Where(qDateTo).
		Where(qDateFrom).
		Order("created_at desc").
		Limit(req.Limit).
		Offset(PageOffset).
		Find(&history).
		Limit(-1).
		Offset(-1).
		Count(&countData)

	if res.Error != nil {
		return nil, 0, domain.ErrGetMerchantFundActivities
	}

	return history, countData, res.Error
}
