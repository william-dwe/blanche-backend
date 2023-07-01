package usecase

import (
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type SlpAccountUsecase interface {
	GetSlpAccountListByUsername(username string) ([]dto.SlpAccountResDTO, error)
	GetUserSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error)
	RegisterSlpAccount(username string, req dto.SlpAccountReqDTO) (*dto.SlpAccountResDTO, error)
	DeleteSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error)
	SetDefaultSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error)
}

type SlpAccountUsecaseConfig struct {
	SlpAccountsRepository repository.SlpAccountsRepository
	UserRepository        repository.UserRepository
}

type slpAccountUsecaseImpl struct {
	slpAccountsRepository repository.SlpAccountsRepository
	userRepository        repository.UserRepository
}

func NewSlpAccountUsecase(c SlpAccountUsecaseConfig) SlpAccountUsecase {
	return &slpAccountUsecaseImpl{
		slpAccountsRepository: c.SlpAccountsRepository,
		userRepository:        c.UserRepository,
	}
}

func (u *slpAccountUsecaseImpl) GetSlpAccountListByUsername(username string) ([]dto.SlpAccountResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	slpAccounts, err := u.slpAccountsRepository.GetSlpAccountListByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	var slpAccountsResDTO []dto.SlpAccountResDTO
	for _, slpAccount := range slpAccounts {
		slpAccountsResDTO = append(slpAccountsResDTO, dto.SlpAccountResDTO{
			ID:         slpAccount.ID,
			CardNumber: slpAccount.CardNumber,
			NameOnCard: slpAccount.NameOnCard,
			ActiveDate: slpAccount.ActiveDate,
			IsDefault:  slpAccount.IsDefault,
		})
	}

	return slpAccountsResDTO, nil
}

func (u *slpAccountUsecaseImpl) GetUserSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	slpAccount, err := u.slpAccountsRepository.GetUserSlpAccountByID(user.ID, uint(slpAccountId))
	if err != nil {
		return nil, err
	}

	if slpAccount.UserID != user.ID {
		return nil, domain.ErrUnauthorizedSlpAccount
	}

	return &dto.SlpAccountResDTO{
		ID:         slpAccount.ID,
		CardNumber: slpAccount.CardNumber,
		NameOnCard: slpAccount.NameOnCard,
		ActiveDate: slpAccount.ActiveDate,
		IsDefault:  slpAccount.IsDefault,
	}, nil
}

func (u *slpAccountUsecaseImpl) RegisterSlpAccount(username string, req dto.SlpAccountReqDTO) (*dto.SlpAccountResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if len(req.CardNumber) != 16 {
		return nil, domain.ErrInvalidCardNumber
	}
	if req.ActiveDate.Before(time.Now()) {
		return nil, domain.ErrInvalidActiveDate
	}

	slpAccounts, err := u.slpAccountsRepository.GetSlpAccountListByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	if len(slpAccounts) == 0 {
		req.IsDefault = true
	}

	newSlpAccount := entity.SlpAccount{
		CardNumber: req.CardNumber,
		NameOnCard: req.NameOnCard,
		ActiveDate: req.ActiveDate,
		UserID:     user.ID,
		IsDefault:  req.IsDefault,
	}

	slpAccount, err := u.slpAccountsRepository.RegisterSlpAccount(newSlpAccount)
	if err != nil {
		return nil, err
	}

	return &dto.SlpAccountResDTO{
		ID:         slpAccount.ID,
		CardNumber: slpAccount.CardNumber,
		NameOnCard: slpAccount.NameOnCard,
		ActiveDate: slpAccount.ActiveDate,
		IsDefault:  slpAccount.IsDefault,
	}, nil
}

func (u *slpAccountUsecaseImpl) SetDefaultSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	slpAccount, err := u.slpAccountsRepository.GetUserSlpAccountByID(user.ID, uint(slpAccountId))
	if err != nil {
		return nil, err
	}

	if slpAccount.UserID != user.ID {
		return nil, domain.ErrUnauthorizedSlpAccount
	}

	err = u.slpAccountsRepository.SetDefaultSlpAccount(user.ID, slpAccount.ID)
	if err != nil {
		return nil, err
	}

	return &dto.SlpAccountResDTO{
		ID:         slpAccount.ID,
		CardNumber: slpAccount.CardNumber,
		NameOnCard: slpAccount.NameOnCard,
		ActiveDate: slpAccount.ActiveDate,
		IsDefault:  true,
	}, nil
}

func (u *slpAccountUsecaseImpl) DeleteSlpAccount(username string, slpAccountId int) (*dto.SlpAccountResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	slpAccount, err := u.slpAccountsRepository.GetUserSlpAccountByID(user.ID, uint(slpAccountId))
	if err != nil {
		return nil, err
	}
	slpAccounts, err := u.slpAccountsRepository.GetSlpAccountListByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	if slpAccount.UserID != user.ID {
		return nil, domain.ErrUnauthorizedSlpAccount
	}
	if slpAccount.IsDefault && len(slpAccounts) > 1 {
		return nil, domain.ErrDeleteDefaultSlpAccount
	}

	slpAccount, err = u.slpAccountsRepository.DeleteUserSlpAccount(slpAccount)
	if err != nil {
		return nil, err
	}

	return &dto.SlpAccountResDTO{
		ID:         slpAccount.ID,
		CardNumber: slpAccount.CardNumber,
		NameOnCard: slpAccount.NameOnCard,
		ActiveDate: slpAccount.ActiveDate,
		IsDefault:  slpAccount.IsDefault,
	}, nil
}
