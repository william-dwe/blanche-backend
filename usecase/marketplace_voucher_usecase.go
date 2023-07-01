package usecase

import (
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

const MARKETPLACE_PREFIX = "BLANCHE"

type MarketplaceVoucherUsecase interface {
	GetMarketplaceVoucherList() ([]dto.MarketplaceVoucherResDTO, error)
	GetMarketplaceVoucherByCode(voucherCode string) (*dto.MarketplaceAdminVoucherResDTO, error)
	GetMarketplaceAdminVoucherList(req dto.MerchantVoucherListParamReqDTO) (*dto.MarketplaceAdminVoucherListResDTO, error)
	CreateMarketplaceVoucher(req dto.UpsertMarketplaceVoucherReqDTO) (*dto.MarketplaceAdminVoucherResDTO, error)
	UpdateMarketplaceVoucher(voucherCode string, req dto.UpsertMarketplaceVoucherReqDTO) (*dto.MarketplaceAdminVoucherResDTO, error)
	DeleteMarketplaceVoucher(voucherCode string) (*dto.MarketplaceAdminVoucherResDTO, error)
}

type MarketplaceVoucherUsecaseConfig struct {
	MarketplaceVoucherRepo repository.MarketplaceVoucherRepository
}

type marketplaceVoucherUsecaseImpl struct {
	mpVoucherRepo repository.MarketplaceVoucherRepository
}

func NewMarketplaceVoucherUsecase(c MarketplaceVoucherUsecaseConfig) MarketplaceVoucherUsecase {
	return &marketplaceVoucherUsecaseImpl{
		mpVoucherRepo: c.MarketplaceVoucherRepo,
	}
}

func (u *marketplaceVoucherUsecaseImpl) GetMarketplaceVoucherList() ([]dto.MarketplaceVoucherResDTO, error) {
	mpVouchers, err := u.mpVoucherRepo.GetMarketplaceVoucherList()
	if err != nil {
		return nil, err
	}

	if len(mpVouchers) == 0 {
		return []dto.MarketplaceVoucherResDTO{}, nil
	}

	var mpVoucherResDTOs []dto.MarketplaceVoucherResDTO
	for _, mpVoucher := range mpVouchers {
		mpVoucherResDTOs = append(mpVoucherResDTOs, dto.MarketplaceVoucherResDTO{
			ID:                 mpVoucher.ID,
			DiscountPercentage: mpVoucher.DiscountPercentage,
			ExpiredAt:          mpVoucher.ExpiredAt,
			Code:               mpVoucher.Code,
			MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
			MinOrderNominal:    mpVoucher.MinOrderNominal,
			Quota:              mpVoucher.Quota,
		})
	}

	return mpVoucherResDTOs, nil
}

func (u *marketplaceVoucherUsecaseImpl) GetMarketplaceVoucherByCode(voucherCode string) (*dto.MarketplaceAdminVoucherResDTO, error) {
	mpVoucher, err := u.mpVoucherRepo.GetMarketplaceVoucherByCode(voucherCode)
	if err != nil {
		return nil, err
	}

	mpVoucherResDTO := dto.MarketplaceAdminVoucherResDTO{
		ID:                 mpVoucher.ID,
		DiscountPercentage: mpVoucher.DiscountPercentage,
		Code:               mpVoucher.Code,
		MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
		MpDomain:           MARKETPLACE_PREFIX,
		CodeSuffix:         strings.TrimPrefix(mpVoucher.Code, MARKETPLACE_PREFIX),
		MinOrderNominal:    mpVoucher.MinOrderNominal,
		StartDate:          mpVoucher.StartDate,
		EndDate:            mpVoucher.ExpiredAt,
		Quota:              mpVoucher.Quantity,
		UsedQuota:          mpVoucher.Quantity - mpVoucher.Quota,
	}

	return &mpVoucherResDTO, nil
}

func (u *marketplaceVoucherUsecaseImpl) GetMarketplaceAdminVoucherList(req dto.MerchantVoucherListParamReqDTO) (*dto.MarketplaceAdminVoucherListResDTO, error) {
	mpVouchers, total, err := u.mpVoucherRepo.GetMarketplaceAdminVoucherList(req)
	if err != nil {
		return nil, err
	}

	if len(mpVouchers) == 0 {
		return &dto.MarketplaceAdminVoucherListResDTO{
			PaginationResponse: dto.PaginationResponse{
				CurrentPage: req.Page,
			},
			Vouchers: []dto.MarketplaceAdminVoucherResDTO{},
		}, nil
	}

	var mpVoucherResDTOs []dto.MarketplaceAdminVoucherResDTO
	for _, mpVoucher := range mpVouchers {
		mpVoucherResDTOs = append(mpVoucherResDTOs, dto.MarketplaceAdminVoucherResDTO{
			ID:                 mpVoucher.ID,
			Code:               mpVoucher.Code,
			DiscountPercentage: mpVoucher.DiscountPercentage,
			MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
			MinOrderNominal:    mpVoucher.MinOrderNominal,
			StartDate:          mpVoucher.StartDate,
			EndDate:            mpVoucher.ExpiredAt,
			Quota:              mpVoucher.Quantity,
			UsedQuota:          mpVoucher.Quantity - mpVoucher.Quota,
		})
	}

	return &dto.MarketplaceAdminVoucherListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   total,
			TotalPage:   (total + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
		Vouchers: mpVoucherResDTOs,
	}, nil
}

func (u *marketplaceVoucherUsecaseImpl) CreateMarketplaceVoucher(req dto.UpsertMarketplaceVoucherReqDTO) (*dto.MarketplaceAdminVoucherResDTO, error) {
	err := util.ValidateVoucherCode(req.Code, MARKETPLACE_PREFIX)
	if err != nil {
		return nil, err
	}

	if req.StartDate.After(req.EndDate) {
		return nil, domain.ErrInvalidVoucherDateRange
	}

	if req.EndDate.Before(time.Now()) || req.StartDate.Before(time.Now()) {
		return nil, domain.ErrInvalidVoucherDateBeforeNow
	}

	voucher := entity.MarketplaceVoucher{
		Code:               req.Code,
		DiscountPercentage: req.DiscountPercentage,
		StartDate:          req.StartDate,
		ExpiredAt:          req.EndDate,
		Quantity:           req.Quota,
		Quota:              req.Quota,
		MaxDiscountNominal: req.MaxDiscountNominal,
		MinOrderNominal:    req.MinOrderNominal,
	}

	mpVoucher, err := u.mpVoucherRepo.CreateMarketplaceVoucher(&voucher)
	if err != nil {
		return nil, err
	}

	mpVoucherResDTO := dto.MarketplaceAdminVoucherResDTO{
		ID:                 mpVoucher.ID,
		Code:               mpVoucher.Code,
		DiscountPercentage: mpVoucher.DiscountPercentage,
		MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
		MinOrderNominal:    mpVoucher.MinOrderNominal,
		StartDate:          mpVoucher.StartDate,
		EndDate:            mpVoucher.ExpiredAt,
		Quota:              mpVoucher.Quantity,
		UsedQuota:          mpVoucher.Quantity - mpVoucher.Quota,
	}

	return &mpVoucherResDTO, nil
}

func (u *marketplaceVoucherUsecaseImpl) UpdateMarketplaceVoucher(voucherCode string, req dto.UpsertMarketplaceVoucherReqDTO) (*dto.MarketplaceAdminVoucherResDTO, error) {
	err := util.ValidateVoucherCode(req.Code, MARKETPLACE_PREFIX)
	if err != nil {
		return nil, err
	}

	voucher, err := u.mpVoucherRepo.GetMarketplaceVoucherByCode(voucherCode)
	if err != nil {
		return nil, err
	}

	if req.StartDate.After(req.EndDate) {
		return nil, domain.ErrInvalidVoucherDateRange
	}

	if req.StartDate.Before(voucher.StartDate) {
		return nil, domain.ErrUpdateOngoingVoucher
	}

	voucher.DiscountPercentage = req.DiscountPercentage
	voucher.StartDate = req.StartDate
	voucher.Code = req.Code
	voucher.ExpiredAt = req.EndDate
	voucher.Quantity = req.Quota
	voucher.Quota = req.Quota
	voucher.MinOrderNominal = req.MinOrderNominal
	voucher.MaxDiscountNominal = req.MaxDiscountNominal
	if voucher.Code != req.Code || voucherCode != req.Code {
		return nil, domain.ErrUpdateVoucherCode
	}

	mpVoucher, err := u.mpVoucherRepo.UpdateMarketplaceVoucher(voucher)
	if err != nil {
		return nil, err
	}

	mpVoucherResDTO := dto.MarketplaceAdminVoucherResDTO{
		ID:                 mpVoucher.ID,
		Code:               mpVoucher.Code,
		DiscountPercentage: mpVoucher.DiscountPercentage,
		MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
		MinOrderNominal:    mpVoucher.MinOrderNominal,
		StartDate:          mpVoucher.StartDate,
		EndDate:            mpVoucher.ExpiredAt,
		Quota:              mpVoucher.Quantity,
		UsedQuota:          mpVoucher.Quantity - mpVoucher.Quota,
	}

	return &mpVoucherResDTO, nil
}

func (u *marketplaceVoucherUsecaseImpl) DeleteMarketplaceVoucher(voucherCode string) (*dto.MarketplaceAdminVoucherResDTO, error) {
	voucher, err := u.mpVoucherRepo.GetMarketplaceVoucherByCode(voucherCode)
	if err != nil {
		return nil, err
	}

	if voucher.StartDate.Before(time.Now()) && voucher.ExpiredAt.After(time.Now()) {
		return nil, domain.ErrDelVoucherIsOngoing
	}

	mpVoucher, err := u.mpVoucherRepo.DeleteMarketplaceVoucher(voucher)
	if err != nil {
		return nil, err
	}

	mpVoucherResDTO := dto.MarketplaceAdminVoucherResDTO{
		ID:                 mpVoucher.ID,
		Code:               mpVoucher.Code,
		DiscountPercentage: mpVoucher.DiscountPercentage,
		MaxDiscountNominal: mpVoucher.MaxDiscountNominal,
		MinOrderNominal:    mpVoucher.MinOrderNominal,
		StartDate:          mpVoucher.StartDate,
		EndDate:            mpVoucher.ExpiredAt,
		Quota:              mpVoucher.Quantity,
		UsedQuota:          mpVoucher.Quantity - mpVoucher.Quota,
	}

	return &mpVoucherResDTO, nil
}
