package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const (
	MERCHANT_VOUCHER_IS_VALID  = false
	MERCHANT_VOUCHER_MIN_QUOTA = 0
)

type MerchantRepository interface {
	GetByDomain(domain string) (*entity.Merchant, error)
	GetByStoreName(storeName string) (*entity.Merchant, error)
	GetByUserID(userId uint) (*entity.Merchant, error)
	GetByUsername(username string) (*entity.Merchant, error)
	AddMerchant(merchant *entity.Merchant) error
	UpdateMerchantDetails(merchant *entity.Merchant) (*entity.Merchant, error)
	SynchronizeMerchantCity(merchant *entity.Merchant) error
	UpdateMerchantAddress(merchant *entity.Merchant) error
	GetMerchantProductCategoryIds(merchantDomain string) ([]uint, error)

	GetMerchantVoucherList(merchantDomain string) ([]entity.MerchantVoucher, error)
	GetMerchantAdminVoucherList(merchantDomain string, req dto.MerchantVoucherListParamReqDTO) ([]entity.MerchantVoucher, int64, error)
	GetMerchantVoucher(merchantDomain string, voucherId uint) (*entity.MerchantVoucher, error)
	GetMerchantVoucherByCode(merchantDomain string, voucherCode string) (*entity.MerchantVoucher, error)
	CreateMerchantVoucher(voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error)
	UpdateMerchantVoucher(voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error)
	DeleteMerchantVoucher(merchantDomain string, voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error)

	DecreaseMerchantVoucherQuotaTx(tx *gorm.DB, merchantDomain string, voucherID uint) error
	IncreaseMerchantVoucherQuotaTx(tx *gorm.DB, merchantDomain string, voucherID uint) error
	IncreaseMerchantNumOfSaleTx(tx *gorm.DB, merchantId uint, delta uint) error
	UpdateMerchantRatingAndNumOfReviewTx(tx *gorm.DB, merchantId uint, rating float64) error
}

type MerchantRepositoryConfig struct {
	DB                 *gorm.DB
	UserRepository     UserRepository
	DeliveryRepository DeliveryRepository
}

type merchantRepositoryImpl struct {
	db                 *gorm.DB
	userRepository     UserRepository
	deliveryRepository DeliveryRepository
}

func NewMerchantRepository(c MerchantRepositoryConfig) MerchantRepository {
	return &merchantRepositoryImpl{
		db:                 c.DB,
		userRepository:     c.UserRepository,
		deliveryRepository: c.DeliveryRepository,
	}
}

func (r *merchantRepositoryImpl) GetByDomain(merchantDomain string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.Model(&merchant).Preload("MerchantAnalytical").Preload("UserAddress.Province").Preload("City").Where("domain = ?", merchantDomain).First(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantDomainNotFound
		}

		log.Error().Msgf("cannot get merchant record: %v", err)
		return nil, domain.ErrGetMerchant
	}

	return &merchant, nil
}

func (r *merchantRepositoryImpl) GetByStoreName(name string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.Model(&merchant).Where("name = ?", name).First(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantStoreNameNotFound
		}

		log.Error().Msgf("cannot get merchant record: %v", err)
		return nil, domain.ErrGetMerchant
	}

	return &merchant, nil
}

func (r *merchantRepositoryImpl) GetByUserID(userId uint) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.Model(&merchant).Preload("MerchantAnalytical").Preload("UserAddress.Province").Preload("City").Where("user_id = ?", userId).First(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantUserIDNotFound
		}

		log.Error().Msgf("cannot get merchant record: %v", err)
		return nil, domain.ErrGetMerchant
	}

	return &merchant, nil
}

func (r *merchantRepositoryImpl) GetByUsername(username string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	sq := r.db.Model(&entity.User{}).
		Select("id").
		Where("username = ?", username)
	err := r.db.Model(&merchant).
		Preload("MerchantAnalytical").
		Preload("UserAddress.Province").
		Preload("UserAddress.City").
		Where("user_id in (?)", sq).
		First(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantUsernameNotFound
		}

		log.Error().Msgf("cannot get merchant record: %v", err)
		return nil, domain.ErrGetMerchant
	}

	return &merchant, nil
}

func (r *merchantRepositoryImpl) UpdateMerchantDetails(merchant *entity.Merchant) (*entity.Merchant, error) {
	err := r.db.Model(&merchant).Updates(merchant).Error
	if err != nil {
		log.Error().Msgf("cannot update merchant record: %v", err)
		return nil, domain.ErrUpdateMerchant
	}

	return merchant, nil
}

func (r *merchantRepositoryImpl) AddMerchant(merchant *entity.Merchant) (errAddMerchant error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			errAddMerchant = domain.ErrRegisterMerchant
		}
	}()

	err := tx.Create(merchant).Error
	if err != nil {
		tx.Rollback()
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == "23503" {
			if pgErr.ConstraintName == "user_id_fkey" {
				return domain.ErrUserNotFound
			}
			if pgErr.ConstraintName == "user_address_id_fkey" {
				return domain.ErrAddMerchantAddressNotFound
			}
		}
		if pgErr.Code == "23505" {
			if pgErr.ConstraintName == "merchants_domain_key" {
				return domain.ErrMerchantDomainUnique
			}
			if pgErr.ConstraintName == "merchants_name_key" {
				return domain.ErrMerchantNameUnique
			}
		}

		log.Error().Msgf("Error create merchant account: %v", pgErr.Message)
		return domain.ErrRegisterMerchant
	}

	newMerchantHoldingAcc := entity.MerchantHoldingAccount{
		MerchantID: merchant.ID,
		Balance:    0,
	}
	err = tx.Create(&newMerchantHoldingAcc).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrRegisterMerchantAccount
	}

	newMerchantAnalytics := entity.MerchantAnalytical{
		MerchantId:   merchant.ID,
		AvgRating:    0,
		NumOfSale:    0,
		NumOfProduct: 0,
		NumOfReview:  0,
	}

	err = tx.Create(&newMerchantAnalytics).Error
	if err != nil {
		tx.Rollback()
		return domain.ErrRegisterMerchantAccount
	}

	allDeliveryOptions, err := r.deliveryRepository.GetAllDeliveryOption()
	if err != nil {
		tx.Rollback()
		return domain.ErrRegisterMerchantAccount
	}

	var newMerchantDeliveryOption []entity.MerchantDeliveryOption
	for _, deliveryOption := range allDeliveryOptions {
		newMerchantDeliveryOption = append(newMerchantDeliveryOption,
			entity.MerchantDeliveryOption{
				MerchantId:       merchant.ID,
				DeliveryOptionId: deliveryOption.ID,
			})
	}

	_, err = r.deliveryRepository.AddMerchantDeliveryOptionsTx(tx, newMerchantDeliveryOption)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.userRepository.UpdateUserRoleToMerchantTx(tx, merchant.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return domain.ErrRegisterMerchant
	}

	return nil
}

func (r *merchantRepositoryImpl) SynchronizeMerchantCity(merchant *entity.Merchant) error {
	err := r.db.Model(&entity.Merchant{}).
		Where("id = ?", merchant.ID).
		Update("city_id", merchant.UserAddress.CityId).Error
	if err != nil {
		return domain.ErrUpdateMerchantCity
	}

	return nil
}

func (r *merchantRepositoryImpl) UpdateMerchantAddress(merchant *entity.Merchant) error {
	err := r.db.Model(&entity.Merchant{}).
		Where("id = ?", merchant.ID).
		Update("user_address_id", merchant.UserAddressId).Error
	if err != nil {
		return domain.ErrUpdateMerchantAddress
	}

	return nil
}

func (r *merchantRepositoryImpl) GetMerchantProductCategoryIds(merchantDomain string) ([]uint, error) {
	var productCategoryIds []uint
	err := r.db.Model(&entity.Product{}).
		Select("distinct category_id").
		Where("merchant_domain = ?", merchantDomain).
		Find(&productCategoryIds).Error
	if err != nil {
		return nil, domain.ErrGetMerchantProductCategoryIds
	}

	return productCategoryIds, nil
}

func (r *merchantRepositoryImpl) GetMerchantVoucherList(merchantDomain string) ([]entity.MerchantVoucher, error) {
	var merchantVouchers []entity.MerchantVoucher
	err := r.db.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ?", merchantDomain).
		Where("is_invalid = ? AND quota > ? AND start_date <= now() AND expired_at >= now()", MERCHANT_VOUCHER_IS_VALID, MERCHANT_VOUCHER_MIN_QUOTA).
		Find(&merchantVouchers).Error
	if err != nil {
		return nil, domain.ErrGetMerchantVoucherList
	}

	return merchantVouchers, nil
}

func (r *merchantRepositoryImpl) GetMerchantAdminVoucherList(merchantDomain string, req dto.MerchantVoucherListParamReqDTO) ([]entity.MerchantVoucher, int64, error) {
	var merchantVouchers []entity.MerchantVoucher
	var total int64
	PageOffset := req.Limit * (req.Page - 1)
	baseQuery := r.db.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ?", merchantDomain).
		Limit(req.Limit).
		Offset(PageOffset)

	if req.Status == dto.VoucherStatusOngoing {
		baseQuery = baseQuery.Where("start_date <= now() AND expired_at >= now()")
	}
	if req.Status == dto.VoucherStatusExpired {
		baseQuery = baseQuery.Where("expired_at < now()")
	}
	if req.Status == dto.VoucherStatusIncoming {
		baseQuery = baseQuery.Where("start_date > now()")
	}

	err := baseQuery.Find(&merchantVouchers).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, total, domain.ErrGetMerchantVoucherList
	}

	return merchantVouchers, total, nil
}

func (r *merchantRepositoryImpl) GetMerchantVoucher(merchantDomain string, voucherId uint) (*entity.MerchantVoucher, error) {
	var merchantVoucher entity.MerchantVoucher
	err := r.db.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ? AND id = ?", merchantDomain, voucherId).
		Where("is_invalid = ? AND quota > ? AND start_date <= now() AND expired_at >= now()", MERCHANT_VOUCHER_IS_VALID, MERCHANT_VOUCHER_MIN_QUOTA).
		First(&merchantVoucher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantVoucherNotFound
		}

		return nil, domain.ErrGetMerchantVoucher
	}

	return &merchantVoucher, nil
}

func (r *merchantRepositoryImpl) GetMerchantVoucherByCode(merchantDomain string, voucherCode string) (*entity.MerchantVoucher, error) {
	var merchantVoucher entity.MerchantVoucher
	err := r.db.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ? AND code = ?", merchantDomain, voucherCode).
		First(&merchantVoucher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMerchantVoucherNotFound
		}

		return nil, domain.ErrGetMerchantVoucher
	}

	return &merchantVoucher, nil
}

func (r *merchantRepositoryImpl) CreateMerchantVoucher(voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error) {
	err := r.db.Model(&entity.MerchantVoucher{}).
		Create(voucher).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"merchant_voucher_code_unique": domain.ErrMerchantVoucherCodeAlreadyExist,
			},
			domain.ErrCreateMerchantVoucher,
		)
		return nil, maskedErr
	}

	return voucher, nil
}

func (r *merchantRepositoryImpl) UpdateMerchantVoucher(voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error) {
	err := r.db.Model(&entity.MerchantVoucher{}).
		Where("id = ?", voucher.ID).
		Save(voucher).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"merchant_voucher_code_unique": domain.ErrMerchantVoucherCodeAlreadyExist,
			},
			domain.ErrUpdateMerchantVoucher,
		)
		return nil, maskedErr
	}

	return voucher, nil
}

func (r *merchantRepositoryImpl) DecreaseMerchantVoucherQuotaTx(tx *gorm.DB, merchantDomain string, voucherId uint) error {
	err := tx.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ? AND id = ?", merchantDomain, voucherId).
		Where("is_invalid = ? AND quota > ?", MERCHANT_VOUCHER_IS_VALID, MERCHANT_VOUCHER_MIN_QUOTA).
		Update("quota", gorm.Expr("quota - ?", 1)).Error
	if err != nil {
		return domain.ErrDecreaseMerchantVoucherQuota
	}

	return nil
}

func (r *merchantRepositoryImpl) IncreaseMerchantVoucherQuotaTx(tx *gorm.DB, merchantDomain string, voucherId uint) error {
	err := tx.Model(&entity.MerchantVoucher{}).
		Where("merchant_domain = ? AND id = ?", merchantDomain, voucherId).
		Update("quota", gorm.Expr("quota + ?", 1)).Error
	if err != nil {
		return domain.ErrIncreaseMarketplaceVoucherQuota
	}

	return nil
}

func (r *merchantRepositoryImpl) DeleteMerchantVoucher(merchantDomain string, voucher *entity.MerchantVoucher) (*entity.MerchantVoucher, error) {
	err := r.db.Model(&voucher).
		Where("merchant_domain = ?", merchantDomain).
		Delete(voucher).Error
	if err != nil {
		return nil, domain.ErrDeleteMerchantVoucher
	}

	return voucher, nil
}

func (r *merchantRepositoryImpl) IncreaseMerchantNumOfSaleTx(tx *gorm.DB, merchantId uint, delta uint) error {
	err := tx.Model(&entity.MerchantAnalytical{}).
		Where("merchant_id = ?", merchantId).
		Update("num_of_sale", gorm.Expr("num_of_sale + ?", delta)).Error
	if err != nil {
		return domain.ErrIncreaseMerchantNumOfSale
	}

	return nil
}

func (r *merchantRepositoryImpl) UpdateMerchantRatingAndNumOfReviewTx(tx *gorm.DB, merchantId uint, rating float64) error {
	err := tx.Model(&entity.MerchantAnalytical{}).
		Where("merchant_id = ?", merchantId).
		Updates(map[string]interface{}{
			"avg_rating":    gorm.Expr("((avg_rating * num_of_review) + ?)/(num_of_review + 1)", rating),
			"num_of_review": gorm.Expr("num_of_review + ?", 1)},
		).Error
	if err != nil {
		return domain.ErrUpdateMerchantRatingAndNumOfReview
	}

	return nil
}
