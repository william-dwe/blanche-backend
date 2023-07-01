package repository

import (
	"fmt"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByUserId(userId uint) (*entity.Wallet, error)
	GetById(walletId uint) (*entity.Wallet, error)

	Create(newWallet entity.Wallet) (*entity.Wallet, error)

	Update(wallet entity.Wallet) (*entity.Wallet, error)
	UpdateBalance(wallet entity.Wallet, amount int) (*entity.Wallet, error)
	UpdateBalanceTx(tx *gorm.DB, wallet entity.Wallet, amount int) error
	BalanceTopUpSuccess(wallet entity.Wallet, walletTrx entity.WalletTransactionRecord) error
	BalanceTopUpFailed(wallet entity.Wallet, walletTrx entity.WalletTransactionRecord) error

	GetTransactionsByWalletId(walletID uint, req dto.WalletTransactionReqParamDTO) ([]entity.WalletTransactionRecord, int64, error)
	AddTranscation(wallet entity.Wallet, transaction entity.WalletTransactionRecord) (*entity.WalletTransactionRecord, error)
	GetTransactionByPaymentId(paymentId string) (*entity.WalletTransactionRecord, error)

	AddWalletHistoryMerchantWdTx(tx *gorm.DB, userId uint, amount float64) error
	AddWalletHistoryRefundTx(tx *gorm.DB, userId uint, amount float64, transactionId uint) error

	MakePayment(walletId string, amount uint) (redirectUrl string, paymentId string, paymentRecordId uint, walletpayError error)
}

type WalletRepositoryConfig struct {
	DB                      *gorm.DB
	PaymentRecordRepository PaymentRecordRepository
}

type walletRepositoryImpl struct {
	db                      *gorm.DB
	paymentRecordRepository PaymentRecordRepository
}

func NewWalletRepository(c WalletRepositoryConfig) WalletRepository {
	return &walletRepositoryImpl{
		db:                      c.DB,
		paymentRecordRepository: c.PaymentRecordRepository,
	}
}

func (r *walletRepositoryImpl) GetById(walletId uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	res := r.db.Where("id = ?", walletId).First(&wallet)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletNotFound
		}
		return nil, domain.ErrGetWallet
	}

	return &wallet, nil
}

func (r *walletRepositoryImpl) Create(newWallet entity.Wallet) (*entity.Wallet, error) {
	res := r.db.Create(&newWallet)
	if res.Error != nil {
		pgErr := res.Error.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "wallets_user_id_fkey" {
				return nil, domain.ErrCreateWalletUserIdNotFound
			}
		}
		if pgErr.Code == "23505" {
			if pgErr.ConstraintName == "wallets_user_id_key" {
				return nil, domain.ErrCreateWalletUserDuplicate
			}
		}

		log.Error().Msgf("Error create wallet: %v", pgErr.Message)
		return nil, domain.ErrCreateWallet
	}

	return &newWallet, nil
}

func (r *walletRepositoryImpl) GetByUserId(userId uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	res := r.db.Where("user_id = ?", userId).First(&wallet)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletNotFound
		}
		return nil, domain.ErrGetWallet
	}

	return &wallet, nil
}

func (r *walletRepositoryImpl) Update(newWallet entity.Wallet) (*entity.Wallet, error) {
	res := r.db.Updates(&newWallet)
	if res.Error != nil {
		return nil, domain.ErrUpdateWalletPin
	}

	return &newWallet, nil
}

func (r *walletRepositoryImpl) GetTransactionsByWalletId(walletID uint, req dto.WalletTransactionReqParamDTO) ([]entity.WalletTransactionRecord, int64, error) {
	qDate := ""
	if req.StartDate != "" && req.EndDate != "" {
		qDate = fmt.Sprintf("payment_records.paid_at BETWEEN '%s' AND '%s'", req.StartDate, req.EndDate)
	} else if req.StartDate != "" {
		qDate = fmt.Sprintf("payment_records.paid_at >= '%s'", req.StartDate)
	} else if req.EndDate != "" {
		qDate = fmt.Sprintf("payment_records.paid_at <= '%s'", req.EndDate)
	}

	var totalRows int64
	var walletTrx []entity.WalletTransactionRecord
	res := r.db.Debug().Model(&entity.WalletTransactionRecord{}).
		Where("wallet_id = ?", walletID).
		Joins("LEFT JOIN payment_records ON wallet_transaction_records.payment_id = payment_records.payment_id").
		Preload("WalletTransactionType").
		Preload("PaymentRecord").
		Preload("PaymentRecord.Transactions").
		Where("payment_records.payment_id = wallet_transaction_records.payment_id").
		Where("payment_records.paid_at IS NOT NULL").
		Where(qDate).
		Order("created_at DESC").
		Limit(req.Limit).
		Offset(req.Limit * (req.Page - 1)).
		Find(&walletTrx).Limit(-1).Offset(-1).Count(&totalRows)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, 0, domain.ErrWalletTransactionNotFound
		}
		return nil, 0, domain.ErrGetWalletTransaction
	}

	return walletTrx, totalRows, nil
}

func (r *walletRepositoryImpl) UpdateBalance(newWallet entity.Wallet, amount int) (*entity.Wallet, error) {
	res := r.db.Model(&newWallet).Update("balance", gorm.Expr("balance + ?", amount))
	if res.Error != nil {
		return nil, domain.ErrUpdateWalletBalance
	}

	return &newWallet, nil
}

func (r *walletRepositoryImpl) UpdateBalanceTx(tx *gorm.DB, wallet entity.Wallet, amount int) error {
	qUserId := "id = 0"
	if wallet.UserId != 0 {
		qUserId = fmt.Sprintf("user_id = %d", wallet.UserId)
	}
	if wallet.ID != 0 {
		qUserId = fmt.Sprintf("id = %d", wallet.ID)
	}

	var amountBalance float64
	res := tx.
		Raw(fmt.Sprintf(`
		UPDATE "wallets" SET "balance"=balance + %d,"updated_at"=now() WHERE %s AND "wallets"."deleted_at" IS null
		returning balance
	`, amount, qUserId)).Scan(&amountBalance)

	if res.Error != nil {
		return domain.ErrUpdateWalletBalance
	}

	if amountBalance < 0 {
		return domain.ErrWalletBalanceNotSufficient
	}

	return nil
}

func (r *walletRepositoryImpl) AddTranscation(wallet entity.Wallet, transaction entity.WalletTransactionRecord) (resWalletTrx *entity.WalletTransactionRecord, addTrxErr error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			addTrxErr = domain.ErrAddWalletTransaction
			resWalletTrx = nil
		}
	}()

	var err error
	createdPaymentRecord, err := r.paymentRecordRepository.CreateTx(tx, transaction.PaymentRecord)
	if err != nil {
		tx.Rollback()

		log.Error().Msgf("Error create payment record: %v", err)
		return nil, err
	}

	transaction.PaymentRecord = *createdPaymentRecord
	err = tx.Model(&wallet).Association("WalletTransactionRecords").Append(&transaction)
	if err != nil {
		tx.Rollback()

		log.Error().Msgf("Error create wallet transaction: %v", err)
		return nil, domain.ErrAddWalletTransaction
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, domain.ErrAddWalletTransaction
	}

	return &transaction, nil
}

func (r *walletRepositoryImpl) GetTransactionByPaymentId(paymentId string) (*entity.WalletTransactionRecord, error) {
	var transaction entity.WalletTransactionRecord
	res := r.db.
		Preload("PaymentRecord").
		Where("payment_id = ?", paymentId).
		First(&transaction)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrWalletTransactionNotFound
		}
		return nil, domain.ErrGetWallet
	}

	return &transaction, nil
}

func (r *walletRepositoryImpl) BalanceTopUpSuccess(wallet entity.Wallet, transaction entity.WalletTransactionRecord) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = domain.ErrTopUpWallet
		}
	}()

	err = tx.Model(&transaction.PaymentRecord).Where("id = ?", transaction.PaymentRecord.ID).Update("paid_at", time.Now()).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrTopUpWalletPaymentRecord
	}

	err = tx.Model(&wallet).Where("id = ?", wallet.ID).Update("balance", gorm.Expr("balance + ?", transaction.PaymentRecord.Amount)).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrUpdateWalletBalance
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrTopUpWallet
	}

	return nil
}

func (r *walletRepositoryImpl) BalanceTopUpFailed(wallet entity.Wallet, transaction entity.WalletTransactionRecord) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = domain.ErrTopUpWallet
		}
	}()

	err = tx.Model(&transaction.PaymentRecord).Where("id = ?", transaction.PaymentRecord.ID).Update("canceled_at", time.Now()).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrTopUpWalletPaymentRecord
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrTopUpWallet
	}

	return nil
}

func (r *walletRepositoryImpl) MakePayment(walletId string, amount uint) (redirectUrl string, paymentId string, paymentRecordId uint, walletpayError error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			walletpayError = domain.ErrCreateWalletPaymentRecord
		}
	}()

	var wallet entity.Wallet
	err := tx.Where("id = ?", walletId).First(&wallet).Error
	if err != nil {
		tx.Rollback()
		return "", "", 0, domain.ErrWalletNotFound
	}
	if uint(wallet.Balance) < amount {
		tx.Rollback()
		return "", "", 0, domain.ErrWalletBalanceNotSufficient
	}

	var newPaymentRecord entity.PaymentRecord
	err = tx.Model(&newPaymentRecord).Raw(`insert into payment_records (payment_id, payment_method_id, amount)
		values( concat('WAL', nextval('wallet_payment_id_seq')), ?, ?) 
		returning id, payment_id , amount, payment_method_id, created_at, updated_at, deleted_at , payment_url`, dto.PAYMENT_METHOD_ID_WALLET, amount).Scan(&newPaymentRecord).Error
	if err != nil {
		tx.Rollback()
		return "", "", 0, domain.ErrCreateWalletPaymentRecord
	}

	// update payment Url
	redirectPayUrl := fmt.Sprintf("%s?payment_id=%s&amount=%d", config.Config.WalletpayConfig.RedirectUrl, newPaymentRecord.PaymentId, amount)
	err = tx.Model(&newPaymentRecord).Where("id = ?", newPaymentRecord.ID).Update("payment_url", redirectPayUrl).Error
	if err != nil {
		tx.Rollback()
		return "", "", 0, domain.ErrCreateWalletPaymentRecord
	}

	// //create wallet transaction records
	newWalletTransactionRecord := entity.WalletTransactionRecord{
		WalletId:                wallet.ID,
		WalletTransactionTypeId: dto.WALLET_TRANSACTION_TYPE_ID_TRANSACTION,
		PaymentId:               newPaymentRecord.PaymentId,
	}
	err = tx.Model(&entity.WalletTransactionRecord{}).Create(&newWalletTransactionRecord).Error
	if err != nil {
		tx.Rollback()
		return "", "", 0, domain.ErrCreateWalletTransactionRecord
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return "", "", 0, domain.ErrCreateWalletPaymentRecord
	}

	return redirectPayUrl, newPaymentRecord.PaymentId, newPaymentRecord.ID, nil
}

func (r *walletRepositoryImpl) AddWalletHistoryMerchantWdTx(tx *gorm.DB, userId uint, amount float64) error {
	//make payment record
	var newPaymentRecord entity.PaymentRecord
	err := tx.
		Raw(`insert into payment_records (payment_id, payment_method_id, amount, paid_at)
		values( concat('MEW', nextval('withdraw_merchant_payment_id_seq')), ?, ?, now()) 
		returning id, payment_id , amount, payment_method_id, created_at, updated_at, deleted_at , payment_url`, dto.PAYMENT_METHOD_ID_MERCHANT_WITHDRAW, amount).
		Scan(&newPaymentRecord).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for withdraw in making payment record %v", err)
		return domain.ErrCreateWalletTransactionRecordWithdrawPaymentRec
	}

	//update balance
	qUserId := fmt.Sprintf("user_id = %d", userId)
	var walletId uint
	err = tx.
		Raw(fmt.Sprintf(`
			UPDATE "wallets" SET "balance"=balance + %d,"updated_at"=now() WHERE %s AND "wallets"."deleted_at" IS null
			returning id
		`, int(amount), qUserId)).Scan(&walletId).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for withdraw in updating balance %v", err)
		return domain.ErrCreateWalletTransactionRecordWithdrawUpdateBalance
	}

	newWalletTrxRecord := entity.WalletTransactionRecord{
		WalletId:                walletId,
		WalletTransactionTypeId: dto.WALLET_TRANSACTION_TYPE_ID_MERCHANT_WITHDRAWAL,
		PaymentId:               newPaymentRecord.PaymentId,
	}
	err = tx.Model(&entity.WalletTransactionRecord{}).Create(&newWalletTrxRecord).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for withdraw %v", err)
		return domain.ErrCreateWalletTransactionRecordWithdraw
	}

	return nil
}

func (r *walletRepositoryImpl) AddWalletHistoryRefundTx(tx *gorm.DB, userId uint, amount float64, transctionId uint) error {
	//make payment record
	var newPaymentRecord entity.PaymentRecord
	err := tx.
		Raw(`insert into payment_records (payment_id, payment_method_id, amount, paid_at)
		values( concat('REF', nextval('refund_transaction_payment_id_seq')), ?, ?, now()) 
		returning id, payment_id , amount, payment_method_id, created_at, updated_at, deleted_at , payment_url`, dto.PAYMENT_METHOD_ID_TRANSACTION_REFUND, amount).
		Scan(&newPaymentRecord).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for trx refund in making payment record %v", err)
		return domain.ErrCreateWalletTransactionRecordRefundPaymentRec
	}

	//update balance
	qUserId := fmt.Sprintf("user_id = %d", userId)
	var walletId uint
	err = tx.
		Raw(fmt.Sprintf(`
			UPDATE "wallets" SET "balance"=balance + %d,"updated_at"=now() WHERE %s AND "wallets"."deleted_at" IS null
			returning id
		`, int(amount), qUserId)).Scan(&walletId).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for refund in updating balance %v", err)
		return domain.ErrCreateWalletTransactionRecordRefundUpdateBalance
	}

	//TODO make transaction record payment for refund
	newTransactionPaymentRecord := entity.TransactionPaymentRecord{
		TransactionId: transctionId,
		PaymentId:     newPaymentRecord.PaymentId,
		OrderCode:     fmt.Sprintf("ref-%s", util.GenerateUUID()),
	}
	err = tx.Create(&newTransactionPaymentRecord).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for refund in making transaction payment record %v", err)
		return domain.ErrCreateWalletTransactionRecordRefundTransactionPaymentRec
	}

	newWalletTrxRecord := entity.WalletTransactionRecord{
		WalletId:                walletId,
		WalletTransactionTypeId: dto.WALLET_TRANSACTION_TYPE_ID_REFUND,
		PaymentId:               newPaymentRecord.PaymentId,
	}
	err = tx.Model(&entity.WalletTransactionRecord{}).Create(&newWalletTrxRecord).Error
	if err != nil {
		log.Error().Msgf("cannot create wallet transaction record for refund %v", err)
		return domain.ErrCreateWalletTransactionRecordRefund
	}

	return nil
}
