package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var transactionStatusMap = map[uint]string{
	1: "on_waited_at",
	2: "on_processed_at",
	3: "on_delivery_at",
	4: "tds.on_delivered_at",
	5: "on_request_refund_at",
	6: "on_refunded_at",
	7: "on_completed_at",
	8: "on_canceled_at",
}

var transactionSortDateOrder = map[uint]string{
	1: "desc",
	2: "asc",
}

var transactionStatusOrderMap = map[uint]uint{
	1: 1,
	2: 2,
	3: 3,
	4: 4,
	5: 7,
	6: 8,
	7: 5,
	8: 6,
}

type TransactionRepository interface {
	GetTransactionListByMerchant(merchantDomain string, req dto.TransactionReqParamDTO) ([]entity.Transaction, int64, error)
	GetTransactionListByUserId(userId uint, req dto.TransactionReqParamDTO) ([]entity.Transaction, int64, error)
	GetTransactionDetailByInvoiceCode(userId uint, invoiceCode string) (*entity.Transaction, error)
	GetTransactionByInvoiceCode(invoiceCode string) (*entity.Transaction, error)
	GetMerchantTransactionDetailByInvoiceCode(merchantDomain string, invoiceCode string) (*entity.Transaction, error)

	MakeTransaction(transactions []entity.Transaction, orderSummary dto.PostOrderSummaryResDTO, paymentId string) error
	UpdateTransactionPaymentSuccess(transactions []entity.Transaction, paymentRec entity.PaymentRecord) error
	UpdateTransactionPaymentFailed(transactions []entity.Transaction, cartItems []entity.TransactionCartItem, paymentRec entity.PaymentRecord) error

	UpdateTransactionStatusCanceled(transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trxStatus *entity.TransactionStatus, cancelTrx error)
	UpdateTransactionStatusCanceledTx(tx *gorm.DB, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trxStatus *entity.TransactionStatus, cancelTrx error)
	UpdateTransactionStatusCompleted(transaction entity.Transaction, amount float64, amountPromotionMarketplace float64) (trsStatus *entity.TransactionStatus, cancelTrx error)
	UpdateTransactionStatusCompletedTx(tx *gorm.DB, transaction entity.Transaction, amount float64, amountPromotionMarketplace float64) (trsStatus *entity.TransactionStatus, cancelTrx error)
	UpdateTransactionStatusRefundedTx(tx *gorm.DB, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trsStatus *entity.TransactionStatus, cancelTrx error)
}

type transactionRepositoryImpl struct {
	db                                      *gorm.DB
	marketplaceVoucherRepository            MarketplaceVoucherRepository
	merchantRepository                      MerchantRepository
	productRepository                       ProductRepository
	paymentRecordRepository                 PaymentRecordRepository
	transactionDeliveryStatusRepository     TransactionDeliveryStatusRepository
	transactionStatusRepository             TransactionStatusRepository
	transactionPaymentRecordRepository      TransactionPaymentRecordRepository
	walletRepository                        WalletRepository
	merchantHoldingAccountRepository        MerchantHoldingAccountRepository
	merchantHoldingAccountHistoryRepository MerchantHoldingAccountHistoryRepository
}

type TransactionRepositoryConfig struct {
	DB                                      *gorm.DB
	MarketplaceVoucherRepository            MarketplaceVoucherRepository
	MerchantRepository                      MerchantRepository
	ProductRepository                       ProductRepository
	PaymentRecordRepository                 PaymentRecordRepository
	TransactionDeliveryStatusRepository     TransactionDeliveryStatusRepository
	TransactionStatusRepository             TransactionStatusRepository
	TransactionPaymentRecordRepository      TransactionPaymentRecordRepository
	WalletRepository                        WalletRepository
	MerchantHoldingAccountRepository        MerchantHoldingAccountRepository
	MerchantHoldingAccountHistoryRepository MerchantHoldingAccountHistoryRepository
}

func NewTransactionRepository(c TransactionRepositoryConfig) TransactionRepository {
	return &transactionRepositoryImpl{
		db:                                      c.DB,
		marketplaceVoucherRepository:            c.MarketplaceVoucherRepository,
		merchantRepository:                      c.MerchantRepository,
		productRepository:                       c.ProductRepository,
		paymentRecordRepository:                 c.PaymentRecordRepository,
		transactionDeliveryStatusRepository:     c.TransactionDeliveryStatusRepository,
		transactionStatusRepository:             c.TransactionStatusRepository,
		transactionPaymentRecordRepository:      c.TransactionPaymentRecordRepository,
		walletRepository:                        c.WalletRepository,
		merchantHoldingAccountRepository:        c.MerchantHoldingAccountRepository,
		merchantHoldingAccountHistoryRepository: c.MerchantHoldingAccountHistoryRepository,
	}
}

func (r *transactionRepositoryImpl) GetTransactionByInvoiceCode(invoiceCode string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Unscoped().Model(&transaction).Where("invoice_code = ?", invoiceCode).First(&transaction).Error
	if err != nil {
		return nil, domain.ErrTransactionNotFound
	}
	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetTransactionListByMerchant(merchantDomain string, req dto.TransactionReqParamDTO) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64
	PageOffset := req.Pagination.Limit * (req.Pagination.Page - 1)
	query := r.db.Model(&transactions).Where("merchant_domain = ?", merchantDomain).
		Preload("User").
		Joins("left join transaction_statuses ts on transactions.id = ts.transaction_id").
		Joins("left join transaction_delivery_statuses tds on transactions.id = tds.transaction_id").
		Where("ts.on_waited_at IS NOT NULL").
		Order("created_at "+transactionSortDateOrder[req.Sort]).
		Limit(req.Pagination.Limit).
		Offset(PageOffset).
		Where("EXISTS (SELECT * FROM jsonb_array_elements(cart_items) f(x) WHERE x->>'name' ILIKE ?)", "%"+req.Search+"%")

	fieldName := transactionStatusMap[transactionStatusOrderMap[req.Status]]
	maxLength := len(transactionStatusOrderMap)
	if req.Status != 0 {
		if int(transactionStatusOrderMap[req.Status]) == maxLength || int(req.Status) == maxLength {
			query = query.Where(fieldName + " IS NOT NULL")
		} else {
			nextFieldName := transactionStatusMap[transactionStatusOrderMap[req.Status+1]]
			query = query.Where(fieldName + " IS NOT NULL AND " + nextFieldName + " IS NULL AND on_canceled_at IS NULL AND on_refunded_at IS NULL")
		}
	}

	err := query.Find(&transactions).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, total, domain.ErrTransactionsNotFound
		}

		return nil, total, domain.ErrGetTransactions
	}

	return transactions, total, nil
}

func (r *transactionRepositoryImpl) GetTransactionListByUserId(userId uint, req dto.TransactionReqParamDTO) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64
	PageOffset := req.Pagination.Limit * (req.Pagination.Page - 1)
	query := r.db.Model(&transactions).Where("transactions.user_id = ?", userId).
		Joins("Merchant").
		Joins("left join transaction_statuses ts on transactions.id = ts.transaction_id").
		Joins("left join transaction_delivery_statuses tds on transactions.id = tds.transaction_id").
		Where("ts.on_waited_at IS NOT NULL").
		Order("created_at "+transactionSortDateOrder[req.Sort]).
		Limit(req.Pagination.Limit).
		Offset(PageOffset).
		Where("EXISTS (SELECT * FROM jsonb_array_elements(cart_items) f(x) WHERE x->>'name' ILIKE ?) OR name ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")

	fieldName := transactionStatusMap[transactionStatusOrderMap[req.Status]]
	maxLength := len(transactionStatusOrderMap)
	if req.Status != 0 {
		if int(transactionStatusOrderMap[req.Status]) == maxLength || int(req.Status) == maxLength {
			query = query.Where(fieldName + " IS NOT NULL")
		} else {
			nextFieldName := transactionStatusMap[transactionStatusOrderMap[req.Status+1]]
			query = query.Where(fieldName + " IS NOT NULL AND " + nextFieldName + " IS NULL AND on_canceled_at IS NULL AND on_refunded_at IS NULL")
		}
	}

	err := query.Find(&transactions).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, total, domain.ErrTransactionsNotFound
		}

		return nil, total, domain.ErrGetTransactions
	}

	return transactions, total, nil
}

func (r *transactionRepositoryImpl) GetTransactionDetailByInvoiceCode(userId uint, invoiceCode string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Where("user_id = ?", userId).
		Where("invoice_code = ?", invoiceCode).
		Preload("Merchant").
		Preload("TransactionStatus").
		Preload("TransactionDeliveryStatus").
		First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransactionNotFound
		}

		return nil, domain.ErrGetTransaction
	}

	return &transaction, nil
}

func (r *transactionRepositoryImpl) GetMerchantTransactionDetailByInvoiceCode(merchantDomain string, invoiceCode string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Model(&transaction).Where("merchant_domain = ?", merchantDomain).
		Where("invoice_code = ?", invoiceCode).
		Preload("TransactionStatus").
		Preload("TransactionDeliveryStatus").
		Preload("User.UserDetail").
		First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrTransactionNotFound
		}

		return nil, domain.ErrGetTransaction
	}

	return &transaction, nil
}

func (r *transactionRepositoryImpl) MakeTransaction(transactions []entity.Transaction, orderSummary dto.PostOrderSummaryResDTO, paymentId string) (errMakeTrx error) {
	//begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in MakeTransaction repo: %v", r)
			errMakeTrx = domain.ErrMakeTransaction
		}
	}()

	// make payment records with the transaction
	err := tx.Create(&transactions).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error create transaction record: %v", err)
		return domain.ErrCreateTransaction
	}

	// decrease voucher marketplace quota
	if transactions[0].MarketplaceVoucherId != nil {
		err = r.marketplaceVoucherRepository.DecreaseMarketplaceVoucherQuotaTx(tx, *transactions[0].MarketplaceVoucherId)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error decrease marketplace voucher quota: %v", err)
			return domain.ErrCreateTransaction
		}
	}

	// decrease product stock and product promotion and merchant voucher quota
	for _, order := range orderSummary.Orders {
		if order.MerchantVoucherId != nil {
			err = r.merchantRepository.DecreaseMerchantVoucherQuotaTx(tx, order.Merchant.MerchantDomain, *order.MerchantVoucherId)
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error decrease merchant voucher quota: %v", err)
				return domain.ErrCreateTransaction
			}
		}

		//decrease stock and promotion
		for _, item := range order.Items {
			//decrease product stock
			err = r.productRepository.DecreaseProductStockTx(tx, item.ProductId, *item.VariantItemId, uint(item.Quantity))
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error decrease product stock: %v", err)
				return domain.ErrCreateTransaction
			}
			//decrease product promotion
			if item.RealPrice != item.DiscountPrice {
				err = r.productRepository.DecreaseProductPromotionTx(tx, item.ProductId, uint(item.Quantity))
				if err != nil {
					tx.Rollback()
					log.Error().Msgf("Error decrease product promotion: %v", err)
					return domain.ErrCreateTransaction
				}
			}

			//increase pending product pending sale
			err = r.productRepository.ChangeNumOfPendingSaleTx(tx, item.ProductId, 1)
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error increase product pending sale: %v", err)
				return domain.ErrCreateTransaction
			}
		}
	}

	// update transaction payment record order code
	err = r.transactionPaymentRecordRepository.UpdateRecordOrderCodebyPaymentIdTx(tx, orderSummary.OrderCode, paymentId)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction payment record order code: %v", err)
		return domain.ErrCreateTransaction
	}

	//delete order summary
	var userOder entity.UserOrder
	err = tx.Where("order_code = ?", orderSummary.OrderCode).Delete(&userOder).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error delete order summary: %v", err)
		return domain.ErrCreateTransaction
	}

	//delete userorder item
	// var orderItem entity.OrderItem
	// err = tx.Where("order_code = ?", orderSummary.OrderCode).Delete(&orderItem).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	log.Error().Msgf("Error delete order item: %v", err)
	// 	return domain.ErrCreateTransaction
	// }

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrCreateTransaction
	}

	return nil
}

func (r *transactionRepositoryImpl) UpdateTransactionPaymentSuccess(transactions []entity.Transaction, paymentRec entity.PaymentRecord) (updateTrx error) {
	//begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UpdateTransactionPaymentSuccess repo: %v", r)
			updateTrx = domain.ErrUpdateTransactionPayment
		}
	}()

	trxIds := r.parseTransactionstoTransacionIds(transactions)
	timeNow := time.Now()

	//update transaction status
	err := r.transactionStatusRepository.UpdateTransactionStatusWaitedAtTx(tx, trxIds, timeNow)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status: %v", err)
		return domain.ErrUpdateTransactionStatusPayment
	}

	//update payment record status
	paymentRec.PaidAt = &timeNow
	_, err = r.paymentRecordRepository.UpdateTx(tx, paymentRec)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update payment record status: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//update balance marketplace (add amount)
	mpWallet := entity.Wallet{
		ID: dto.WALLET_ID_ADMIN,
	}
	err = r.walletRepository.UpdateBalanceTx(tx, mpWallet, int(paymentRec.Amount))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update balance marketplace: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//update balance user (deduct) if using wallet as payment
	if paymentRec.PaymentMethodId == dto.PAYMENT_METHOD_ID_WALLET {
		userWallet := entity.Wallet{
			UserId: transactions[0].UserId,
		}
		err = r.walletRepository.UpdateBalanceTx(tx, userWallet, -1*int(paymentRec.Amount))
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error update balance user: %v", err)
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrUpdateTransactionPayment
	}

	return nil
}

func (r *transactionRepositoryImpl) UpdateTransactionPaymentFailed(transactions []entity.Transaction, cartItems []entity.TransactionCartItem, paymentRec entity.PaymentRecord) (updateTrx error) {
	//begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UpdateTransactionPaymentFailed repo: %v", r)
			updateTrx = domain.ErrUpdateTransactionPayment
		}
	}()

	trxIds := r.parseTransactionstoTransacionIds(transactions)
	timeNow := time.Now()
	createdTransactionTime := transactions[0].CreatedAt

	//delete transaction
	err := tx.Where("id in ?", trxIds).Delete(&entity.Transaction{}).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error delete transaction: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//delete transaction status
	err = r.transactionStatusRepository.DeleteTransactionStatusTx(tx, trxIds)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error delete transaction status: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//delete transaction delivery status
	err = r.transactionDeliveryStatusRepository.DeleteTransactionDeliveryStatusTx(tx, trxIds)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error delete transaction delivery status: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//update payment record status
	paymentRec.CanceledAt = &timeNow
	_, err = r.paymentRecordRepository.UpdateTx(tx, paymentRec)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update payment record status: %v", err)
		return domain.ErrUpdateTransactionPayment
	}

	//return marketplace voucher if exist
	if transactions[0].MarketplaceVoucherId != nil {
		err = r.marketplaceVoucherRepository.IncreaseMarketplaceVoucherQuotaTx(tx, *transactions[0].MarketplaceVoucherId)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase marketplace voucher quota: %v", err)
			return domain.ErrUpdateTransactionPayment
		}
	}

	//return merchant voucher if exist
	for _, transaction := range transactions {
		if transaction.MerchantVoucherId != nil {
			err = r.merchantRepository.IncreaseMerchantVoucherQuotaTx(tx, transaction.MerchantDomain, *transaction.MerchantVoucherId)
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error increase merchant voucher quota: %v", err)
				return domain.ErrUpdateTransactionPayment
			}
		}
	}

	//return stock and promotions
	//return product stock
	for _, cartItem := range cartItems {
		//return product stock
		err = r.productRepository.IncreaseProductStockTx(tx, cartItem.ProductId, cartItem.ProductVariantId, uint(cartItem.Quantity))
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase product stock: %v", err)
			return domain.ErrUpdateTransactionPayment
		}
		//return product promotion
		if cartItem.RealPrice != cartItem.DiscountPrice {
			err = r.productRepository.IncreaseProductPromotionTx(tx, cartItem.ProductId, uint(cartItem.Quantity), createdTransactionTime)
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error increase product promotion: %v", err)
				return domain.ErrCreateTransaction
			}
		}

		//decrease pending product pending sale
		err = r.productRepository.ChangeNumOfPendingSaleTx(tx, cartItem.ProductId, -1)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase product pending sale: %v", err)
			return domain.ErrCreateTransaction
		}
	}

	//undelete order summary order code
	var trxPayRec entity.TransactionPaymentRecord
	err = tx.Where("payment_id = ?", paymentRec.PaymentId).First(&trxPayRec).Error
	if err == nil {
		err := tx.Unscoped().Model(&entity.UserOrder{}).Where("order_code = ?", trxPayRec.OrderCode).Update("deleted_at", nil).Error
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error delete order summary: %v", err)
			return domain.ErrCreateTransaction
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrUpdateTransactionPayment
	}

	return nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatusCanceledTx(tx *gorm.DB, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trxStatus *entity.TransactionStatus, cancelTrx error) {
	//update status transaction
	trxNewStatus, rowsAffected, err := r.transactionStatusRepository.UpdateTransactionStatusTx(tx, *transaction.TransactionStatus)
	if err != nil || rowsAffected <= 0 {
		log.Error().Msgf("Error update transaction status: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCancel
	}

	//deduct money from marketplace wallet
	err = r.walletRepository.UpdateBalanceTx(tx, entity.Wallet{ID: dto.WALLET_ID_ADMIN}, -1*int(amount))
	if err != nil {
		log.Error().Msgf("Error update balance marketplace wallet: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCancel
	}

	// add wallet history to user and add money to user wallet
	err = r.walletRepository.AddWalletHistoryRefundTx(tx, transaction.UserId, amount, transaction.ID)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error add wallet history to user: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToRefundedAddWalletHistory
	}

	//return marketplace voucher if exist
	if transaction.MarketplaceVoucherId != nil {
		err = r.marketplaceVoucherRepository.IncreaseMarketplaceVoucherQuotaTx(tx, *transaction.MarketplaceVoucherId)
		if err != nil {
			log.Error().Msgf("Error increase marketplace voucher quota: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	//return merchant voucher if exist
	if transaction.MerchantVoucherId != nil {
		err = r.merchantRepository.IncreaseMerchantVoucherQuotaTx(tx, transaction.MerchantDomain, *transaction.MerchantVoucherId)
		if err != nil {
			log.Error().Msgf("Error increase merchant voucher quota: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	//return stock and promotions
	//return product stock
	for _, cartItem := range cartItems {
		//return product stock
		err = r.productRepository.IncreaseProductStockTx(tx, cartItem.ProductId, cartItem.ProductVariantId, uint(cartItem.Quantity))
		if err != nil {
			log.Error().Msgf("Error increase product stock: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
		//return product promotion
		if cartItem.RealPrice != cartItem.DiscountPrice {
			err = r.productRepository.IncreaseProductPromotionTx(tx, cartItem.ProductId, uint(cartItem.Quantity), transaction.CreatedAt)
			if err != nil {
				log.Error().Msgf("Error increase product promotion: %v", err)
				return nil, domain.ErrUpdateTransactionStatusToCancel
			}
		}

		//decrease pending product pending sale
		err = r.productRepository.ChangeNumOfPendingSaleTx(tx, cartItem.ProductId, -1)
		if err != nil {
			log.Error().Msgf("Error decrease product pending sale: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	return trxNewStatus, nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatusCanceled(transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trxStatus *entity.TransactionStatus, cancelTrx error) {
	//begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UpdateTransactionStatusCanceled repo: %v", r)
			cancelTrx = domain.ErrUpdateTransactionStatusToCancel
			trxStatus = nil
		}
	}()

	//update status transaction
	trxNewStatus, rowsAffected, err := r.transactionStatusRepository.UpdateTransactionStatusTx(tx, *transaction.TransactionStatus)
	if err != nil || rowsAffected <= 0 {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCancel
	}

	//deduct money from marketplace wallet
	err = r.walletRepository.UpdateBalanceTx(tx, entity.Wallet{ID: dto.WALLET_ID_ADMIN}, -1*int(amount))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update balance marketplace wallet: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCancel
	}

	// add wallet history to user and add money to user wallet
	err = r.walletRepository.AddWalletHistoryRefundTx(tx, transaction.UserId, amount, transaction.ID)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error add wallet history to user: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToRefundedAddWalletHistory
	}

	//return marketplace voucher if exist
	if transaction.MarketplaceVoucherId != nil {
		err = r.marketplaceVoucherRepository.IncreaseMarketplaceVoucherQuotaTx(tx, *transaction.MarketplaceVoucherId)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase marketplace voucher quota: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	//return merchant voucher if exist
	if transaction.MerchantVoucherId != nil {
		err = r.merchantRepository.IncreaseMerchantVoucherQuotaTx(tx, transaction.MerchantDomain, *transaction.MerchantVoucherId)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase merchant voucher quota: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	//return stock and promotions
	//return product stock
	for _, cartItem := range cartItems {
		//return product stock
		err = r.productRepository.IncreaseProductStockTx(tx, cartItem.ProductId, cartItem.ProductVariantId, uint(cartItem.Quantity))
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase product stock: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
		//return product promotion
		if cartItem.RealPrice != cartItem.DiscountPrice {
			err = r.productRepository.IncreaseProductPromotionTx(tx, cartItem.ProductId, uint(cartItem.Quantity), transaction.CreatedAt)
			if err != nil {
				tx.Rollback()
				log.Error().Msgf("Error increase product promotion: %v", err)
				return nil, domain.ErrUpdateTransactionStatusToCancel
			}
		}

		//decrease pending product pending sale
		err = r.productRepository.ChangeNumOfPendingSaleTx(tx, cartItem.ProductId, -1)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error decrease product pending sale: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCancel
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, domain.ErrUpdateTransactionStatusToCancel
	}

	return trxNewStatus, nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatusCompletedTx(tx *gorm.DB, transaction entity.Transaction, amount float64, amountPromotionMarketplace float64) (trsStatus *entity.TransactionStatus, cancelTrx error) {
	// update status transaction
	timeNow := time.Now()
	transaction.TransactionStatus.OnCompletedAt = &timeNow
	trxNewStatus, rowsAffected, err := r.transactionStatusRepository.UpdateTransactionStatusTx(tx, *transaction.TransactionStatus)
	if err != nil || rowsAffected <= 0 {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCompleted
	}

	// deduct money from marketplace wallet
	err = r.walletRepository.UpdateBalanceTx(tx, entity.Wallet{ID: dto.WALLET_ID_ADMIN}, -1*int(amount))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update balance marketplace wallet: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCompleted
	}

	// if there is a marketplace voucher used, deduct the wallet promotion wallet
	if amountPromotionMarketplace > 0 {
		err = r.walletRepository.UpdateBalanceTx(tx, entity.Wallet{ID: dto.WALLET_ID_ADMIN_PROMOTION}, -1*int(amountPromotionMarketplace))
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error update balance marketplace wallet promotion: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCompleted
		}
	}

	// add money to merchants wallet
	merchantHoldingAcc, rowAffected, err := r.merchantHoldingAccountRepository.UpdateBalanceTx(tx, transaction.Merchant.ID, int(amount+amountPromotionMarketplace))
	if err != nil || rowAffected <= 0 {
		tx.Rollback()
		log.Error().Msgf("Error update balance to merchant wallet: %v rows affect: %d", err, rowAffected)
		return nil, domain.ErrUpdateTransactionStatusToCompletedFundActivities
	}

	//add history merchant holding account
	err = r.merchantHoldingAccountHistoryRepository.AddHistoryTx(tx, merchantHoldingAcc.ID, float64(amount+amountPromotionMarketplace), fmt.Sprintf("Income from %s", transaction.InvoiceCode))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error create history merchant holding acc: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCompletedFundActivities
	}

	var cartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItem
	}

	var totalProductSold uint = 0
	// decrease pending product pending sale and increase num of sale
	for _, cartItem := range cartItems {
		//decrease pending product pending sale
		err = r.productRepository.ChangeNumOfPendingSaleTx(tx, cartItem.ProductId, -1)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error decrease product pending sale: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCompleted
		}

		err = r.productRepository.IncreaseNumOfSaleTx(tx, cartItem.ProductId, cartItem.Quantity)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error increase product num of sale: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToCompleted
		}

		totalProductSold += uint(cartItem.Quantity)
	}

	// add merchant analytics num of sale
	err = r.merchantRepository.IncreaseMerchantNumOfSaleTx(tx, transaction.Merchant.ID, totalProductSold)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error increase merchant num of sale: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToCompleted
	}

	return trxNewStatus, nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatusCompleted(transaction entity.Transaction, amount float64, amountPromotionMarketplace float64) (trsStatus *entity.TransactionStatus, cancelTrx error) {
	//begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UpdateTransactionStatusCanceled repo: %v", r)
			cancelTrx = domain.ErrUpdateTransactionStatusToCompleted
			trsStatus = nil
		}
	}()

	trxNewStatus, err := r.UpdateTransactionStatusCompletedTx(tx, transaction, amount, amountPromotionMarketplace)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status completed: %v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, domain.ErrUpdateTransactionStatusToCompleted
	}

	return trxNewStatus, nil
}

func (r *transactionRepositoryImpl) UpdateTransactionStatusRefundedTx(tx *gorm.DB, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (trxStatus *entity.TransactionStatus, cancelTrx error) {
	//update status transaction
	timeNow := time.Now()
	transaction.TransactionStatus.OnRefundedAt = &timeNow
	trxNewStatus, rowsAffected, err := r.transactionStatusRepository.UpdateTransactionStatusTx(tx, *transaction.TransactionStatus)
	if err != nil || rowsAffected <= 0 {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status to refunded: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToRefunded
	}

	//deduct money from marketplace wallet
	err = r.walletRepository.UpdateBalanceTx(tx, entity.Wallet{ID: dto.WALLET_ID_ADMIN}, -1*int(amount))
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update balance marketplace wallet: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToRefundedDeductMpBalance
	}

	// add wallet history to user and add money to user wallet
	err = r.walletRepository.AddWalletHistoryRefundTx(tx, transaction.UserId, amount, transaction.ID)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error add wallet history to user: %v", err)
		return nil, domain.ErrUpdateTransactionStatusToRefundedAddWalletHistory
	}

	//decrease pending product pending sale
	for _, cartItem := range cartItems {
		err = r.productRepository.ChangeNumOfPendingSaleTx(tx, cartItem.ProductId, -1)
		if err != nil {
			tx.Rollback()
			log.Error().Msgf("Error decrease product pending sale: %v", err)
			return nil, domain.ErrUpdateTransactionStatusToRefundedPendingTransaction
		}
	}

	return trxNewStatus, nil
}

func (r *transactionRepositoryImpl) parseTransactionstoTransacionIds(transactions []entity.Transaction) []uint {
	var transactionIds []uint
	for _, transaction := range transactions {
		transactionIds = append(transactionIds, transaction.ID)
	}

	return transactionIds
}
