package usecase

import (
	"encoding/json"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type PaymentRecordUsecase interface {
	UpdatePaymentRecordStatus(paymentId string, paymentAmount uint, isSuccess bool) error

	GetWaitingForPayment(username string) ([]dto.WaitingForPaymentDTO, error)
	GetWaitingForPaymentDetail(username string, paymentId string) (*dto.WaitingForPaymentDetailDTO, error)
}

type PaymentRecordUsecaseConfig struct {
	PaymentRecordRepository            repository.PaymentRecordRepository
	TransactionPaymentRecordRepository repository.TransactionPaymentRecordRepository
	WalletRepository                   repository.WalletRepository
	UserRepository                     repository.UserRepository

	WalletUsecase      WalletUsecase
	TransactionUsecase TransactionUsecase
}

type paymentRecordUsecaseImpl struct {
	paymentRecordRepository            repository.PaymentRecordRepository
	transactionPaymentRecordRepository repository.TransactionPaymentRecordRepository
	walletRepository                   repository.WalletRepository
	userRepository                     repository.UserRepository

	walletUsecase      WalletUsecase
	transactionUsecase TransactionUsecase
}

func NewPaymentRecordUsecase(c PaymentRecordUsecaseConfig) PaymentRecordUsecase {
	return &paymentRecordUsecaseImpl{
		paymentRecordRepository:            c.PaymentRecordRepository,
		transactionPaymentRecordRepository: c.TransactionPaymentRecordRepository,
		walletRepository:                   c.WalletRepository,
		userRepository:                     c.UserRepository,

		walletUsecase:      c.WalletUsecase,
		transactionUsecase: c.TransactionUsecase,
	}
}

func (u *paymentRecordUsecaseImpl) GetWaitingForPaymentDetail(username string, paymentId string) (*dto.WaitingForPaymentDetailDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	paymentRecord, err := u.paymentRecordRepository.GetDetailByPaymentId(paymentId)
	if err != nil {
		return nil, err
	}

	if len(paymentRecord.Transactions) == 0 {
		return nil, domain.ErrPaymentIdNotFound
	}

	if len(paymentRecord.Transactions) > 0 && paymentRecord.Transactions[0].UserId != user.ID {
		return nil, domain.ErrPaymentIdNotFound
	}

	paymentMethodName := ""
	paymentRelatedAccount := ""
	if paymentRecord.Transactions != nil && len(paymentRecord.Transactions) > 0 {
		var paymentMethod entity.TransactionPaymentMethod
		err = json.Unmarshal([]byte(paymentRecord.Transactions[0].PaymentMethod.Bytes), &paymentMethod)
		if err == nil {
			paymentRelatedAccount = paymentMethod.AccountRelatedNumber
			paymentMethodName = paymentMethod.Name
		}
	}

	waitingForPaymentDTO := dto.WaitingForPaymentDetailDTO{
		PaymentId:             paymentRecord.PaymentId,
		OrderCode:             paymentRecord.TransactionPaymentRecord.OrderCode,
		Amount:                uint(paymentRecord.Amount),
		RedirectUrl:           paymentRecord.PaymentUrl,
		CreatedAt:             paymentRecord.CreatedAt,
		PayBefore:             paymentRecord.CreatedAt.Add(24 * time.Hour),
		PaymentMethod:         paymentMethodName,
		PaymentRelatedAccount: paymentRelatedAccount,
		Transactions:          make([]dto.WaitingForPaymentTransactions, 0),
	}

	for _, transaction := range paymentRecord.Transactions {
		cartItems := make([]entity.TransactionCartItem, 0)

		err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONCartItem
		}

		var trxPaymentDetails entity.TransactionPaymentDetails
		err = json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &trxPaymentDetails)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONPaymentDetails
		}

		trxPaymentDetails.Total = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MarketplaceVoucherNominal - trxPaymentDetails.MarketplaceVoucherNominal

		waitingForPaymentDTO.Transactions = append(waitingForPaymentDTO.Transactions, dto.WaitingForPaymentTransactions{
			TransactionDetailProductResDTO: dto.TransactionDetailProductResDTO{
				Merchant: dto.TransactionMerchantResDTO{
					Domain: transaction.Merchant.Domain,
					Name:   transaction.Merchant.Name,
				},
				CartItems: cartItems,
			},
			PaymentDetails: trxPaymentDetails,
		})
	}
	return &waitingForPaymentDTO, nil
}

func (u *paymentRecordUsecaseImpl) GetWaitingForPayment(username string) ([]dto.WaitingForPaymentDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	paymentRecords, err := u.paymentRecordRepository.GetWaitingForPayment(user.ID)
	if err != nil {
		return nil, err
	}

	waitingForPaymentDTOs := make([]dto.WaitingForPaymentDTO, 0)

	for _, paymentRecord := range paymentRecords {
		paymentMethodName := ""
		paymentRelatedAccount := ""
		if paymentRecord.Transactions != nil && len(paymentRecord.Transactions) > 0 {
			var paymentMethod entity.TransactionPaymentMethod
			err = json.Unmarshal([]byte(paymentRecord.Transactions[0].PaymentMethod.Bytes), &paymentMethod)
			if err == nil {
				paymentRelatedAccount = paymentMethod.AccountRelatedNumber
				paymentMethodName = paymentMethod.Name
			}
		}

		payBeforeTime := paymentRecord.CreatedAt.Add(24 * time.Hour)
		if paymentRecord.PaymentMethodId == dto.PAYMENT_METHOD_ID_SLP {
			payBeforeTime = paymentRecord.CreatedAt.Add(10 * time.Minute)
		}

		waitingForPaymentDTOs = append(waitingForPaymentDTOs, dto.WaitingForPaymentDTO{
			PaymentId:   paymentRecord.PaymentId,
			OrderCode:   paymentRecord.TransactionPaymentRecord.OrderCode,
			Amount:      uint(paymentRecord.Amount),
			RedirectUrl: paymentRecord.PaymentUrl,
			CreatedAt:   paymentRecord.CreatedAt,
			PayBefore:   payBeforeTime,

			PaymentMethod:         paymentMethodName,
			PaymentRelatedAccount: paymentRelatedAccount,
		})
	}

	return waitingForPaymentDTOs, nil
}

func (u *paymentRecordUsecaseImpl) UpdatePaymentRecordStatus(paymentId string, paymentAmount uint, isSuccess bool) error {
	paymentRecord, errPaymentRecord := u.paymentRecordRepository.GetByPaymentId(paymentId)
	if errPaymentRecord != nil {
		return errPaymentRecord
	}
	if paymentRecord.PaidAt != nil || paymentRecord.CanceledAt != nil {
		return domain.ErrPaymentIdExpired
	}
	if float64(paymentAmount) != paymentRecord.Amount {
		return domain.ErrPaymentAmountNotMatch
	}

	trxPayRecords, errTrxProduct := u.transactionPaymentRecordRepository.GetByPaymentId(paymentId)
	foundTrxWallet, errTrxWallet := u.walletRepository.GetTransactionByPaymentId(paymentId)

	if errTrxWallet != nil && errTrxProduct != nil {
		return domain.ErrPaymentIdNotFound
	}

	if foundTrxWallet != nil && errTrxWallet == nil && foundTrxWallet.WalletTransactionTypeId == dto.WALLET_TRANSACTION_TYPE_ID_TOP_UP_SLP {
		err := u.walletUsecase.HandleTopUpWalletUsingSlpRes(*foundTrxWallet, isSuccess)
		if err != nil {
			return err
		}
	}

	if len(trxPayRecords) > 0 && errTrxProduct == nil {
		err := u.transactionUsecase.HandleTransactionSlpRes(trxPayRecords, isSuccess, *paymentRecord)
		if err != nil {
			return err
		}
	}

	return nil
}
