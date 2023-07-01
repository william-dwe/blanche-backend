package usecase

import (
	"fmt"
	"strings"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
)

type MerchantUsecase interface {
	GetInfoByDomain(string) (*dto.MerchantInfoResDTO, error)
	GetInfoByUsername(username string) (*dto.MerchantInfoResDTO, error)
	GetProductCategories(string) ([]dto.CategoryResDTO, error)
	CheckMerchantDomain(req dto.CheckMerchantDomainReqDTO) (*dto.CheckMerchantDomainResDTO, error)
	CheckMerchantStoreName(req dto.CheckMerchantStoreNameReqDTO) (*dto.CheckMerchantStoreNameResDTO, error)
	RegisterMerchant(string, dto.RegisterMerchantReqDTO) (*dto.RegisterMerchantResDTO, error)
	UpdateMerchantProfile(username string, req dto.UpdateMerchantProfileFormReqDTO) (*dto.UpdateMerchantProfileResDTO, error)

	GetFundActivities(username string, reqParam dto.MerchantFundActivitiesReqParamDTO) (*dto.MerchantFundActivitiesResDTO, error)
	GetMerchantFundBalance(username string) (*dto.MerchantFundBalanceDTO, error)
	WithdrawToWallet(username string, req dto.MerchantWithdrawReqDTO) (*dto.MerchantWithdrawResDTO, error)

	CreateMerchantVoucher(username string, req dto.UpsertMerchantVoucherReqDTO) (*dto.UpsertMerchantVoucherResDTO, error)
	GetMerchantVoucherList(merchantDomain string) ([]dto.MerchantVoucherResDTO, error)
	GetMerchantAdminVoucherList(username string, req dto.MerchantVoucherListParamReqDTO) (*dto.MerchantAdminVoucherListResDTO, error)
	GetMerchantVoucherByCode(username string, voucherCode string) (*dto.MerchantAdminVoucherResDTO, error)
	UpdateMerchantVoucher(username string, voucherCode string, req dto.UpsertMerchantVoucherReqDTO) (*dto.UpsertMerchantVoucherResDTO, error)
	DeleteMerchantVoucher(username, voucherCode string) (*dto.MerchantAdminVoucherResDTO, error)
	UpdateMerchantAddress(username string, userAddressId uint) (*dto.MerchantAddressDTO, error)
}

type MerchantUsecaseConfig struct {
	MerchantRepository                      repository.MerchantRepository
	CategoryRepository                      repository.CategoryRepository
	UserRepository                          repository.UserRepository
	MerchantHoldingAccountHistoryRepository repository.MerchantHoldingAccountHistoryRepository
	MerchantHoldingAccountRepository        repository.MerchantHoldingAccountRepository
	GcsUploader                             util.GCSUploader
}

type merchantUsecaseImpl struct {
	merchantRepository                      repository.MerchantRepository
	categoryRepository                      repository.CategoryRepository
	userRepository                          repository.UserRepository
	merchantHoldingAccountHistoryRepository repository.MerchantHoldingAccountHistoryRepository
	merchantHoldingAccountRepository        repository.MerchantHoldingAccountRepository
	gcsUploader                             util.GCSUploader
}

func NewMerchantUsecase(c MerchantUsecaseConfig) MerchantUsecase {
	return &merchantUsecaseImpl{
		merchantRepository:                      c.MerchantRepository,
		categoryRepository:                      c.CategoryRepository,
		userRepository:                          c.UserRepository,
		merchantHoldingAccountHistoryRepository: c.MerchantHoldingAccountHistoryRepository,
		merchantHoldingAccountRepository:        c.MerchantHoldingAccountRepository,
		gcsUploader:                             c.GcsUploader,
	}
}

func (u *merchantUsecaseImpl) GetInfoByDomain(merchantDomain string) (*dto.MerchantInfoResDTO, error) {
	merchantInfo, err := u.merchantRepository.GetByDomain(merchantDomain)
	if err != nil {
		return nil, err
	}

	merchantInfoDTO := dto.MerchantInfoResDTO{
		ID:     merchantInfo.ID,
		Domain: merchantInfo.Domain,
		Name:   merchantInfo.Name,
		Address: dto.MerchantAddressDTO{
			Province: merchantInfo.UserAddress.Province.Name,
			City:     merchantInfo.City.Name,
		},
		AvgRating:    merchantInfo.MerchantAnalytical.AvgRating,
		JoinDate:     merchantInfo.JoinDate.String(),
		NumOfProduct: merchantInfo.MerchantAnalytical.NumOfProduct,
		NumOfSale:    merchantInfo.MerchantAnalytical.NumOfSale,
		NumOfReview:  merchantInfo.MerchantAnalytical.NumOfReview,
		Image:        merchantInfo.ImageUrl,
	}

	return &merchantInfoDTO, nil
}

func (u *merchantUsecaseImpl) GetInfoByUsername(username string) (*dto.MerchantInfoResDTO, error) {
	merchantInfo, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	merchantInfoDTO := dto.MerchantInfoResDTO{
		ID:     merchantInfo.ID,
		Domain: merchantInfo.Domain,
		Name:   merchantInfo.Name,
		Address: dto.MerchantAddressDTO{
			Province: merchantInfo.UserAddress.Province.Name,
			City:     merchantInfo.City.Name,
		},
		AvgRating:    merchantInfo.MerchantAnalytical.AvgRating,
		JoinDate:     merchantInfo.JoinDate.String(),
		NumOfProduct: merchantInfo.MerchantAnalytical.NumOfProduct,
		NumOfSale:    merchantInfo.MerchantAnalytical.NumOfSale,
		NumOfReview:  merchantInfo.MerchantAnalytical.NumOfReview,
		Image:        merchantInfo.ImageUrl,
	}

	return &merchantInfoDTO, nil
}

func (u *merchantUsecaseImpl) GetProductCategories(merchantDomain string) ([]dto.CategoryResDTO, error) {
	merchantInfo, err := u.merchantRepository.GetByDomain(merchantDomain)
	if err != nil {
		return nil, err
	}

	productCategoriesIds, err := u.merchantRepository.GetMerchantProductCategoryIds(merchantInfo.Domain)
	if err != nil {
		return nil, err
	}

	productCategories, err := u.categoryRepository.GetCategoryTreeByListId(productCategoriesIds)
	if err != nil {
		return nil, err
	}

	categoriesTreeMap := make(map[*entity.Category]map[*entity.Category][]entity.Category)
	for _, category := range productCategories {
		if category.GrandparentId != 0 && category.ParentId != 0 {
			if _, ok := categoriesTreeMap[category.Grandparent]; !ok {
				categoriesTreeMap[category.Grandparent] = make(map[*entity.Category][]entity.Category)
			}
			if _, ok := categoriesTreeMap[category.Grandparent][category.Parent]; !ok {
				categoriesTreeMap[category.Grandparent][category.Parent] = []entity.Category{}
			}
			categoriesTreeMap[category.Grandparent][category.Parent] = append(categoriesTreeMap[category.Grandparent][category.Parent], category)
		}
	}

	categoriesTree := make([]dto.CategoryResDTO, 0)
	for grandparent, parentMap := range categoriesTreeMap {
		var parentCategories []dto.CategoryResDTO
		for parent, children := range parentMap {
			parentCategories = append(parentCategories, u.buildCategoryTO(*parent, children))
		}

		grandparentDTO := u.buildCategoryTO(*grandparent, nil)
		grandparentDTO.Children = parentCategories

		categoriesTree = append(categoriesTree, grandparentDTO)
	}

	return categoriesTree, nil
}

func (u *merchantUsecaseImpl) buildCategoryTO(parentCategory entity.Category, childrenCategory []entity.Category) dto.CategoryResDTO {
	var children []dto.CategoryResDTO
	for _, child := range childrenCategory {
		children = append(children, dto.CategoryResDTO{
			ID:       child.ID,
			Name:     child.Name,
			Slug:     child.Slug,
			ImageUrl: child.ImageUrl,
			Children: nil,
		})
	}
	return dto.CategoryResDTO{
		ID:       parentCategory.ID,
		Name:     parentCategory.Name,
		Slug:     parentCategory.Slug,
		ImageUrl: parentCategory.ImageUrl,
		Children: children,
	}
}

func (u *merchantUsecaseImpl) CheckMerchantDomain(req dto.CheckMerchantDomainReqDTO) (*dto.CheckMerchantDomainResDTO, error) {
	merchant, err := u.merchantRepository.GetByDomain(req.Domain)
	if merchant != nil {
		return &dto.CheckMerchantDomainResDTO{
			Domain:      req.Domain,
			IsAvailable: false,
		}, nil
	}
	if err != nil {
		if err != domain.ErrMerchantDomainNotFound {
			return nil, err
		}
	}
	return &dto.CheckMerchantDomainResDTO{
		Domain:      req.Domain,
		IsAvailable: true,
	}, nil
}

func (u *merchantUsecaseImpl) CheckMerchantStoreName(req dto.CheckMerchantStoreNameReqDTO) (*dto.CheckMerchantStoreNameResDTO, error) {
	merchant, err := u.merchantRepository.GetByStoreName(req.Name)
	if merchant != nil {
		return &dto.CheckMerchantStoreNameResDTO{
			Name:        req.Name,
			IsAvailable: false,
		}, nil
	}
	if err != nil {
		if err != domain.ErrMerchantStoreNameNotFound {
			return nil, err
		}
	}
	return &dto.CheckMerchantStoreNameResDTO{
		Name:        req.Name,
		IsAvailable: true,
	}, nil
}

func (u *merchantUsecaseImpl) RegisterMerchant(username string, req dto.RegisterMerchantReqDTO) (*dto.RegisterMerchantResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, domain.ErrGetUserNotFound
	}

	merchant, err := u.merchantRepository.GetByUserID(user.ID)
	if merchant != nil {
		return nil, domain.ErrMerchantAlreadyRegistered
	}
	if err != nil {
		if err != domain.ErrMerchantUserIDNotFound {
			return nil, err
		}
	}

	if req.AddressId == 0 {
		defaultAddress, err := u.userRepository.GetDefaultUserAddress(*user)
		if err != nil {
			return nil, err
		}
		req.AddressId = defaultAddress.ID
	}

	address, err := u.userRepository.GetUserAddressById(user.ID, req.AddressId)
	if err != nil {
		return nil, err
	}

	newMerchant := entity.Merchant{
		UserId:        user.ID,
		UserAddressId: req.AddressId,
		Domain:        req.Domain,
		Name:          req.Name,
		City:          address.City,
		JoinDate:      time.Now(),
	}

	err = u.merchantRepository.AddMerchant(&newMerchant)
	if err != nil {
		return nil, err
	}

	merchantInfoDTO := dto.RegisterMerchantResDTO{
		Name:   req.Name,
		Domain: req.Domain,
		Address: dto.MerchantAddressDTO{
			Province: address.Province.Name,
			City:     address.City.Name,
		},
		JoinDate: newMerchant.JoinDate,
	}

	return &merchantInfoDTO, nil
}

func (u *merchantUsecaseImpl) CreateMerchantVoucher(username string, req dto.UpsertMerchantVoucherReqDTO) (*dto.UpsertMerchantVoucherResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	err = util.ValidateVoucherCode(req.Code, merchant.Domain)
	if err != nil {
		return nil, err
	}

	if req.StartDate.After(req.EndDate) {
		return nil, domain.ErrInvalidVoucherDateRange
	}

	if req.EndDate.Before(time.Now()) || req.StartDate.Before(time.Now()) {
		return nil, domain.ErrInvalidVoucherDateBeforeNow
	}

	voucher := entity.MerchantVoucher{
		MerchantDomain:     merchant.Domain,
		DiscountNominal:    req.DiscountNominal,
		Code:               req.Code,
		StartDate:          req.StartDate,
		ExpiredAt:          req.EndDate,
		Quantity:           req.Quota,
		Quota:              req.Quota,
		MinOrderNominal:    req.MinOrderNominal,
		MaxDiscountNominal: req.DiscountNominal,
	}

	voucherRes, err := u.merchantRepository.CreateMerchantVoucher(&voucher)
	if err != nil {
		return nil, err
	}

	voucherDTO := dto.UpsertMerchantVoucherResDTO{
		ID:              voucherRes.ID,
		MerchantDomain:  voucherRes.MerchantDomain,
		DiscountNominal: voucherRes.DiscountNominal,
		Code:            voucherRes.Code,
		MinOrderNominal: voucherRes.MinOrderNominal,
		StartDate:       voucherRes.StartDate,
		ExpiredAt:       voucherRes.ExpiredAt,
		Quantity:        voucherRes.Quantity,
		Quota:           voucherRes.Quota,
	}

	return &voucherDTO, nil
}

func (u *merchantUsecaseImpl) UpdateMerchantVoucher(username string, voucherCode string, req dto.UpsertMerchantVoucherReqDTO) (*dto.UpsertMerchantVoucherResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	voucher, err := u.merchantRepository.GetMerchantVoucherByCode(merchant.Domain, voucherCode)
	if err != nil {
		return nil, err
	}

	err = util.ValidateVoucherCode(req.Code, merchant.Domain)
	if err != nil {
		return nil, err
	}

	if req.StartDate.After(req.EndDate) {
		return nil, domain.ErrInvalidVoucherDateRange
	}

	if req.StartDate.Before(voucher.StartDate) {
		return nil, domain.ErrUpdateOngoingVoucher
	}

	voucher.DiscountNominal = req.DiscountNominal
	voucher.StartDate = req.StartDate
	voucher.Code = req.Code
	voucher.ExpiredAt = req.EndDate
	voucher.Quantity = req.Quota
	voucher.Quota = req.Quota
	voucher.MinOrderNominal = req.MinOrderNominal
	voucher.MaxDiscountNominal = req.DiscountNominal
	if voucher.Code != req.Code || voucherCode != req.Code {
		return nil, domain.ErrUpdateVoucherCode
	}

	voucherRes, err := u.merchantRepository.UpdateMerchantVoucher(voucher)
	if err != nil {
		return nil, err
	}

	voucherDTO := dto.UpsertMerchantVoucherResDTO{
		ID:              voucherRes.ID,
		MerchantDomain:  voucherRes.MerchantDomain,
		DiscountNominal: voucherRes.DiscountNominal,
		Code:            voucherRes.Code,
		MinOrderNominal: voucherRes.MinOrderNominal,
		StartDate:       voucherRes.StartDate,
		ExpiredAt:       voucherRes.ExpiredAt,
		Quantity:        voucherRes.Quantity,
		Quota:           voucherRes.Quota,
	}

	return &voucherDTO, nil
}

func (u *merchantUsecaseImpl) GetMerchantAdminVoucherList(username string, req dto.MerchantVoucherListParamReqDTO) (*dto.MerchantAdminVoucherListResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	vouchers, totalVouchers, err := u.merchantRepository.GetMerchantAdminVoucherList(merchant.Domain, req)
	if err != nil {
		return nil, err
	}

	if len(vouchers) == 0 {
		return &dto.MerchantAdminVoucherListResDTO{
			PaginationResponse: dto.PaginationResponse{
				CurrentPage: req.Page,
			},
			Vouchers: []dto.MerchantAdminVoucherResDTO{},
		}, nil
	}

	voucherDTOs := make([]dto.MerchantAdminVoucherResDTO, 0)
	for _, voucher := range vouchers {
		voucherDTOs = append(voucherDTOs, dto.MerchantAdminVoucherResDTO{
			ID:              voucher.ID,
			Code:            voucher.Code,
			StartDate:       voucher.StartDate,
			ExpiredAt:       voucher.ExpiredAt,
			DiscountNominal: voucher.DiscountNominal,
			MinOrderNominal: voucher.MinOrderNominal,
			Quota:           voucher.Quantity,
			UsedQuota:       voucher.Quantity - voucher.Quota,
		})
	}

	return &dto.MerchantAdminVoucherListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalVouchers,
			TotalPage:   (totalVouchers + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
		Vouchers: voucherDTOs,
	}, nil
}

func (u *merchantUsecaseImpl) GetMerchantVoucherList(merchantDomain string) ([]dto.MerchantVoucherResDTO, error) {
	vouchers, err := u.merchantRepository.GetMerchantVoucherList(merchantDomain)
	if err != nil {
		return nil, err
	}

	if len(vouchers) == 0 {
		return []dto.MerchantVoucherResDTO{}, nil
	}

	voucherDTOs := make([]dto.MerchantVoucherResDTO, 0)
	for _, voucher := range vouchers {
		voucherDTOs = append(voucherDTOs, dto.MerchantVoucherResDTO{
			ID:              voucher.ID,
			Code:            voucher.Code,
			ExpiredAt:       voucher.ExpiredAt,
			DiscountNominal: voucher.DiscountNominal,
			Quota:           voucher.Quota,
			MinOrderNominal: voucher.MinOrderNominal,
		})
	}

	return voucherDTOs, nil
}

func (u *merchantUsecaseImpl) GetMerchantVoucherByCode(username string, voucherCode string) (*dto.MerchantAdminVoucherResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	voucher, err := u.merchantRepository.GetMerchantVoucherByCode(merchant.Domain, voucherCode)
	if err != nil {
		return nil, err
	}

	codeSuffix := strings.TrimPrefix(voucher.Code, strings.ToUpper(merchant.Domain))
	voucherDTO := dto.MerchantAdminVoucherResDTO{
		ID:              voucher.ID,
		Code:            voucher.Code,
		MerchantDomain:  strings.ToUpper(voucher.MerchantDomain),
		CodeSuffix:      codeSuffix,
		DiscountNominal: voucher.DiscountNominal,
		MinOrderNominal: voucher.MinOrderNominal,
		StartDate:       voucher.StartDate,
		ExpiredAt:       voucher.ExpiredAt,
		Quota:           voucher.Quantity,
		UsedQuota:       voucher.Quantity - voucher.Quota,
	}

	return &voucherDTO, nil
}

func (u *merchantUsecaseImpl) GetFundActivities(username string, reqParam dto.MerchantFundActivitiesReqParamDTO) (*dto.MerchantFundActivitiesResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	merchantHoldingAcc, err := u.merchantHoldingAccountRepository.GetByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	fundActivities, countData, err := u.merchantHoldingAccountHistoryRepository.GetHistoryByMerchantHoldingAccId(merchantHoldingAcc.ID, reqParam)
	if err != nil {
		return nil, err
	}

	fundActivitiesDTO := make([]dto.MerchantFundActivitiesDTO, 0)

	for _, fundActivity := range fundActivities {
		fundActivitiesDTO = append(fundActivitiesDTO, dto.MerchantFundActivitiesDTO{
			ID:       fundActivity.ID,
			Notes:    fundActivity.Notes,
			Amount:   fundActivity.Amount,
			IssuedAt: fundActivity.CreatedAt,
			Type:     fundActivity.Type,
		})
	}

	return &dto.MerchantFundActivitiesResDTO{
		FundActivities: fundActivitiesDTO,
		PaginationResponse: dto.PaginationResponse{
			TotalData:   countData,
			TotalPage:   (countData + int64(reqParam.Limit) - 1) / int64(reqParam.Limit),
			CurrentPage: reqParam.Page,
		},
	}, nil
}

func (u *merchantUsecaseImpl) DeleteMerchantVoucher(username, voucherCode string) (*dto.MerchantAdminVoucherResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	voucher, err := u.merchantRepository.GetMerchantVoucherByCode(merchant.Domain, voucherCode)
	if err != nil {
		return nil, err
	}

	if voucher.StartDate.Before(time.Now()) && voucher.ExpiredAt.After(time.Now()) {
		return nil, domain.ErrDelVoucherIsOngoing
	}

	voucherRes, err := u.merchantRepository.DeleteMerchantVoucher(merchant.Domain, voucher)
	if err != nil {
		return nil, err
	}

	voucherDTO := dto.MerchantAdminVoucherResDTO{
		ID:              voucherRes.ID,
		Code:            voucherRes.Code,
		DiscountNominal: voucherRes.DiscountNominal,
		StartDate:       voucherRes.StartDate,
		ExpiredAt:       voucherRes.ExpiredAt,
		Quota:           voucherRes.Quantity,
		UsedQuota:       voucherRes.Quantity - voucherRes.Quota,
	}

	return &voucherDTO, nil
}

func (u *merchantUsecaseImpl) UpdateMerchantAddress(username string, userAddressId uint) (*dto.MerchantAddressDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	merchant.UserAddressId = userAddressId
	err = u.merchantRepository.UpdateMerchantAddress(merchant)
	if err != nil {
		return nil, err
	}

	merchantNew, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	if merchantNew.UserAddressId == userAddressId {
		err := u.merchantRepository.SynchronizeMerchantCity(merchantNew)
		if err != nil {
			return nil, err
		}
	}

	merchantAddressDTO := dto.MerchantAddressDTO{
		Province: merchantNew.UserAddress.Province.Name,
		City:     merchantNew.UserAddress.City.Name,
	}

	return &merchantAddressDTO, nil
}

func (u *merchantUsecaseImpl) GetMerchantFundBalance(username string) (*dto.MerchantFundBalanceDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	merchantHoldingAcc, err := u.merchantHoldingAccountRepository.GetByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	merchantFundBalanceDTO := dto.MerchantFundBalanceDTO{
		TotalBalance: merchantHoldingAcc.Balance,
	}

	return &merchantFundBalanceDTO, nil
}

func (u *merchantUsecaseImpl) WithdrawToWallet(username string, req dto.MerchantWithdrawReqDTO) (*dto.MerchantWithdrawResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	merchantHoldingAcc, err := u.merchantHoldingAccountRepository.GetByUserId(user.ID)
	if err != nil {
		return nil, err
	}
	if merchantHoldingAcc.Balance < float64(req.Amount) {
		return nil, domain.ErrMerchantHoldingAccInsufficientBalance
	}

	res, err := u.merchantHoldingAccountRepository.WithdrawBalance(user.ID, merchantHoldingAcc.MerchantID, int(req.Amount))
	if err != nil {
		return nil, err
	}

	resDTO := dto.MerchantWithdrawResDTO{
		ID:     res.ID,
		Amount: res.Amount,
		Notes:  res.Notes,
	}

	return &resDTO, nil
}

func (u *merchantUsecaseImpl) UpdateMerchantProfile(username string, req dto.UpdateMerchantProfileFormReqDTO) (*dto.UpdateMerchantProfileResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, domain.ErrMerchantUsernameNotFound
	}

	if req.Name != nil && *req.Name != "" {
		merchant.Name = *req.Name
	}

	if req.Description != nil && *req.Description != "" {
		merchant.Description = *req.Description
	}

	if req.Image != nil {
		url, err := u.gcsUploader.UploadFileFromFileHeader(*req.Image,
			fmt.Sprintf("merchant_profile_%d", merchant.ID))
		if err != nil {
			log.Error().Msgf("Failed to upload file: %v", err)
			return nil, domain.ErrUploadFile
		}
		merchant.ImageUrl = url
	}

	merchantNew, err := u.merchantRepository.UpdateMerchantDetails(merchant)
	if err != nil {
		return nil, err
	}

	merchantProfileDTO := dto.UpdateMerchantProfileResDTO{
		Name:        merchantNew.Name,
		Description: merchantNew.Description,
		Image:       merchantNew.ImageUrl,
	}

	return &merchantProfileDTO, nil
}
