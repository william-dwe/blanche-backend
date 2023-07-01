package usecase

import (
	"fmt"
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type WalletUsecase interface {
	CreateWalletPin(input dto.CreateWalletPinReq) (*dto.CreateWalletPinRes, error)
	GetWalletDetails(username string) (*dto.WalletDetails, error)
	UpdateWalletPin(input dto.WalletUpdatePin) (*dto.WalletDetails, error)
	GetWalletTransactions(username string, query dto.WalletTransactionReqParamDTO) (*dto.WalletTransactionRecordDataDTO, error)

	MakeTopUpWalletUsingSlpReq(username string, input dto.TopUpWalletUsingSlpReqDTO) (*dto.TopUpWalletUsingSlpResDTO, error)
	HandleTopUpWalletUsingSlpRes(walletTrx entity.WalletTransactionRecord, isSuccess bool) error
}

type WalletUsecaseConfig struct {
	UserRepository       repository.UserRepository
	WalletRepository     repository.WalletRepository
	SealabspayRepository repository.SealabspayRepository
}

type walletUsecaseImpl struct {
	userRepository       repository.UserRepository
	walletRepository     repository.WalletRepository
	sealabspayRepository repository.SealabspayRepository
}

func NewWalletUsecase(c WalletUsecaseConfig) WalletUsecase {
	return &walletUsecaseImpl{
		userRepository:       c.UserRepository,
		walletRepository:     c.WalletRepository,
		sealabspayRepository: c.SealabspayRepository,
	}
}

func (u *walletUsecaseImpl) CreateWalletPin(walletInputDTO dto.CreateWalletPinReq) (*dto.CreateWalletPinRes, error) {
	hashedPin, err := util.HashAndSalt(walletInputDTO.Pin)
	if err != nil {
		return nil, domain.ErrCreateWalletPinHash
	}

	foundUser, err := u.userRepository.GetUserByUsername(walletInputDTO.Username)
	if err != nil {
		return nil, err
	}

	newWallet := entity.Wallet{
		User: *foundUser,
		Pin:  hashedPin,
	}
	createdWallet, err := u.walletRepository.Create(newWallet)
	if err != nil {
		return nil, err
	}

	createdWalletDTO := dto.CreateWalletPinRes{
		WalletId: createdWallet.ID,
	}

	return &createdWalletDTO, nil
}

func (u *walletUsecaseImpl) GetWalletDetails(username string) (*dto.WalletDetails, error) {
	foundUser, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	foundWallet, err := u.walletRepository.GetByUserId(foundUser.ID)
	if err != nil {
		return nil, err
	}

	walletDetailsDTO := dto.WalletDetails{
		ID:      foundWallet.ID,
		Balance: foundWallet.Balance,
	}

	return &walletDetailsDTO, nil
}

func (u *walletUsecaseImpl) UpdateWalletPin(walletUpdatePinDTO dto.WalletUpdatePin) (*dto.WalletDetails, error) {
	foundUser, err := u.userRepository.GetUserByUsername(walletUpdatePinDTO.Username)
	if err != nil {
		return nil, err
	}

	hashedPin, err := util.HashAndSalt(walletUpdatePinDTO.NewPin)
	if err != nil {
		return nil, domain.ErrCreateWalletPinHash
	}

	foundWallet, err := u.walletRepository.GetByUserId(foundUser.ID)
	if err != nil {
		return nil, domain.ErrUpdateWalletUserNotFound
	}
	foundWallet.Pin = hashedPin

	updatedWallet, err := u.walletRepository.Update(*foundWallet)
	if err != nil {
		return nil, err
	}

	walletDetailsDTO := dto.WalletDetails{
		ID:      updatedWallet.ID,
		Balance: updatedWallet.Balance,
	}

	return &walletDetailsDTO, nil
}

func (u *walletUsecaseImpl) GetWalletTransactions(username string, query dto.WalletTransactionReqParamDTO) (*dto.WalletTransactionRecordDataDTO, error) {
	foundUser, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	foundWallet, err := u.walletRepository.GetByUserId(foundUser.ID)
	if err != nil {
		return nil, err
	}

	walletTransactionsDTO := make([]dto.WalletTransactionRecordDTO, 0)
	resDTO := dto.WalletTransactionRecordDataDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   0,
			TotalPage:   0,
			CurrentPage: 1,
		},
		Transactions: walletTransactionsDTO,
	}

	foundWalletTrx, trxCount, err := u.walletRepository.GetTransactionsByWalletId(foundWallet.ID, query)
	if err != nil {
		if err == domain.ErrWalletTransactionNotFound {
			return &resDTO, nil
		}
		return nil, err
	}

	for _, transaction := range foundWalletTrx {
		amount := transaction.PaymentRecord.Amount
		if transaction.WalletTransactionType.Code == dto.WALLET_TRANSACTION_DEBIT_CODE {
			amount = -amount
		}
		payId := transaction.PaymentId
		notes := ""
		title := ""

		if transaction.WalletTransactionTypeId == dto.WALLET_TRANSACTION_TYPE_ID_TOP_UP_SLP {
			notes = "Top Up from SealabsPay"
			title = "Top Up"
		} else if transaction.WalletTransactionTypeId == dto.WALLET_TRANSACTION_TYPE_ID_TRANSACTION {
			invoiceCodes := "["
			for _, trx := range transaction.PaymentRecord.Transactions {
				invoiceCodes += trx.InvoiceCode
				invoiceCodes += ","
			}
			invoiceCodes = strings.TrimSuffix(invoiceCodes, ",")
			invoiceCodes += "]"

			notes = fmt.Sprintf("Pay Transaction for %s", invoiceCodes)
			title = "Transaction"
		} else if transaction.WalletTransactionTypeId == dto.WALLET_TRANSACTION_TYPE_ID_REFUND {
			notes = "Refund from Transaction"
			title = "Refund"
			if len(transaction.PaymentRecord.Transactions) > 0 {
				notes = fmt.Sprintf("Refund from Transaction [%s]", transaction.PaymentRecord.Transactions[0].InvoiceCode)
			}
		} else if transaction.WalletTransactionTypeId == dto.WALLET_TRANSACTION_TYPE_ID_MERCHANT_WITHDRAWAL {
			notes = "Merchant income withdraw"
			title = "Withdraw"
		}

		walletTransactionDTO := dto.WalletTransactionRecordDTO{
			WalletTransactionTypeDTO: dto.WalletTransactionTypeDTO{
				Name: transaction.WalletTransactionType.Name,
				Code: transaction.WalletTransactionType.Code,
			},
			Amount:    amount,
			PaymentId: &payId,
			Notes:     notes,
			Title:     title,
			IssuedAt:  *transaction.PaymentRecord.PaidAt,
		}
		walletTransactionsDTO = append(walletTransactionsDTO, walletTransactionDTO)
	}

	resDTO.TotalData = trxCount
	resDTO.TotalPage = (trxCount + int64(query.Limit) - 1) / int64(query.Limit)
	resDTO.CurrentPage = query.Page
	resDTO.Transactions = walletTransactionsDTO

	return &resDTO, nil
}

func (u *walletUsecaseImpl) MakeTopUpWalletUsingSlpReq(username string, input dto.TopUpWalletUsingSlpReqDTO) (*dto.TopUpWalletUsingSlpResDTO, error) {
	foundUser, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	foundWallet, err := u.walletRepository.GetByUserId(foundUser.ID)
	if err != nil {
		return nil, err
	}

	redirectUrl, paymentId, err := u.sealabspayRepository.MakePayment(input.SlpCardNumber, uint(input.Amount))
	if err != nil {
		return nil, err
	}

	newTrxRecord := entity.WalletTransactionRecord{
		WalletId:  foundWallet.ID,
		PaymentId: paymentId,
		PaymentRecord: entity.PaymentRecord{
			Amount:          input.Amount,
			PaymentId:       paymentId,
			PaymentMethodId: dto.PAYMENT_METHOD_ID_SLP,
		},
		WalletTransactionTypeId: dto.WALLET_TRANSACTION_TYPE_ID_TOP_UP_SLP,
	}

	_, err = u.walletRepository.AddTranscation(*foundWallet, newTrxRecord)
	if err != nil {
		return nil, err
	}

	res := dto.TopUpWalletUsingSlpResDTO{
		PaymentId:      paymentId,
		Amount:         input.Amount,
		WalletId:       foundWallet.ID,
		SlpCardNumber:  input.SlpCardNumber,
		SlpRedirectUrl: redirectUrl,
	}

	return &res, nil
}

func (u *walletUsecaseImpl) HandleTopUpWalletUsingSlpRes(walletTrx entity.WalletTransactionRecord, isSuccess bool) error {
	foundWallet, err := u.walletRepository.GetById(walletTrx.WalletId)
	if err != nil {
		return err
	}

	if isSuccess {
		err = u.walletRepository.BalanceTopUpSuccess(*foundWallet, walletTrx)
		if err != nil {
			return err
		}
	} else {
		err = u.walletRepository.BalanceTopUpFailed(*foundWallet, walletTrx)
		if err != nil {
			return err
		}
	}

	return nil
}
