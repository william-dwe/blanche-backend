package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type WalletpayUsecase interface {
	HandleWalletpaySuccessRequest(username string, req dto.WalletpayReqDTO) (*dto.WalletpayResDTO, error)
	HandleWalletpayCancelRequest(username string, req dto.WalletpayReqDTO) (*dto.WalletpayResDTO, error)
}

type WalletpayUsecaseConfig struct {
	WalletRepository     repository.WalletRepository
	UserRepository       repository.UserRepository
	PaymentRecordUsecase PaymentRecordUsecase
}

type walletpayUsecaseImpl struct {
	walletRepository     repository.WalletRepository
	userRepository       repository.UserRepository
	paymentRecordUsecase PaymentRecordUsecase
}

func NewWalletpayUsecase(c WalletpayUsecaseConfig) WalletpayUsecase {
	return &walletpayUsecaseImpl{
		walletRepository:     c.WalletRepository,
		userRepository:       c.UserRepository,
		paymentRecordUsecase: c.PaymentRecordUsecase,
	}
}

func (u *walletpayUsecaseImpl) HandleWalletpaySuccessRequest(username string, req dto.WalletpayReqDTO) (*dto.WalletpayResDTO, error) {
	err := u.checkWalletBalance(username, req.Amount)
	if err != nil {
		return nil, err
	}

	err = u.paymentRecordUsecase.UpdatePaymentRecordStatus(req.PaymentId, req.Amount, true)
	if err != nil {
		return nil, err
	}

	return &dto.WalletpayResDTO{
		Amount:    req.Amount,
		PaymentId: req.PaymentId,
		Status:    "TXN_PAID",
	}, nil
}

func (u *walletpayUsecaseImpl) HandleWalletpayCancelRequest(username string, req dto.WalletpayReqDTO) (*dto.WalletpayResDTO, error) {
	err := u.paymentRecordUsecase.UpdatePaymentRecordStatus(req.PaymentId, req.Amount, false)
	if err != nil {
		return nil, err
	}

	return &dto.WalletpayResDTO{
		Amount:    req.Amount,
		PaymentId: req.PaymentId,
		Status:    "TXN_FAILED",
	}, nil
}

func (u *walletpayUsecaseImpl) checkWalletBalance(username string, amount uint) error {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return err
	}

	wallet, err := u.walletRepository.GetByUserId(user.ID)
	if err != nil {
		return err
	}

	if uint(wallet.Balance) < amount {
		return domain.ErrWalletBalanceNotSufficient
	}

	return nil
}
