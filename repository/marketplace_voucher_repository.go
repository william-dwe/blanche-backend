package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

const (
	MP_VOUCHER_IS_VALID  = false
	MP_VOUCHER_MIN_QUOTA = 0
)

type MarketplaceVoucherRepository interface {
	GetMarketplaceVoucherList() ([]entity.MarketplaceVoucher, error)
	GetMarketplaceVoucherByID(voucherID uint) (*entity.MarketplaceVoucher, error)
	GetMarketplaceVoucherByCode(voucherCode string) (*entity.MarketplaceVoucher, error)
	GetMarketplaceAdminVoucherList(req dto.MerchantVoucherListParamReqDTO) ([]entity.MarketplaceVoucher, int64, error)
	CreateMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error)
	UpdateMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error)
	DeleteMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error)

	DecreaseMarketplaceVoucherQuotaTx(tx *gorm.DB, voucherID uint) error
	IncreaseMarketplaceVoucherQuotaTx(tx *gorm.DB, voucherID uint) error
}

type MarketplaceVoucherRepositoryConfig struct {
	DB *gorm.DB
}

type mpVoucherRepositoryImpl struct {
	db *gorm.DB
}

func NewMarketplaceVoucherRepository(c MarketplaceVoucherRepositoryConfig) MarketplaceVoucherRepository {
	return &mpVoucherRepositoryImpl{
		db: c.DB,
	}
}

func (r *mpVoucherRepositoryImpl) GetMarketplaceVoucherList() ([]entity.MarketplaceVoucher, error) {
	var mpVouchers []entity.MarketplaceVoucher
	err := r.db.Model(&mpVouchers).
		Where("is_invalid = ? AND quota > ? AND start_date <= now() AND expired_at >= now()", MP_VOUCHER_IS_VALID, MP_VOUCHER_MIN_QUOTA).
		Find(&mpVouchers).Error
	if err != nil {
		return nil, domain.ErrGetMarketplaceVoucherList
	}

	return mpVouchers, nil
}

func (r *mpVoucherRepositoryImpl) GetMarketplaceVoucherByID(voucherID uint) (*entity.MarketplaceVoucher, error) {
	var mpVoucher entity.MarketplaceVoucher
	err := r.db.Model(&mpVoucher).
		Where("is_invalid = ? AND quota > ? AND start_date <= now() AND expired_at >= now()", MP_VOUCHER_IS_VALID, MP_VOUCHER_MIN_QUOTA).
		Where("id = ?", voucherID).
		First(&mpVoucher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMarketplaceVoucherNotFound
		}

		return nil, domain.ErrGetMarketplaceVoucher
	}

	return &mpVoucher, nil
}

func (r *mpVoucherRepositoryImpl) GetMarketplaceVoucherByCode(voucherCode string) (*entity.MarketplaceVoucher, error) {
	var mpVoucher entity.MarketplaceVoucher
	err := r.db.Model(&entity.MarketplaceVoucher{}).
		Where("code = ?", voucherCode).
		First(&mpVoucher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrMarketplaceVoucherNotFound
		}

		return nil, domain.ErrGetMarketplaceVoucher
	}

	return &mpVoucher, nil
}

func (r *mpVoucherRepositoryImpl) GetMarketplaceAdminVoucherList(req dto.MerchantVoucherListParamReqDTO) ([]entity.MarketplaceVoucher, int64, error) {
	var marketplaceVouchers []entity.MarketplaceVoucher
	var total int64
	PageOffset := req.Limit * (req.Page - 1)
	baseQuery := r.db.Model(&entity.MarketplaceVoucher{}).
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

	err := baseQuery.Find(&marketplaceVouchers).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, total, domain.ErrGetMarketplaceVoucherList
	}

	return marketplaceVouchers, total, nil
}

func (r *mpVoucherRepositoryImpl) CreateMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error) {
	err := r.db.Create(voucher).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"marketplace_vouchers_code_unique": domain.ErrMarketplaceVoucherCodeAlreadyExist,
			},
			domain.ErrCreateMarketplaceVoucher,
		)
		return nil, maskedErr
	}

	return voucher, nil
}

func (r *mpVoucherRepositoryImpl) UpdateMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error) {
	err := r.db.Model(&entity.MarketplaceVoucher{}).
		Where("id = ?", voucher.ID).
		Updates(voucher).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"marketplace_vouchers_code_unique": domain.ErrMarketplaceVoucherCodeAlreadyExist,
			},
			domain.ErrUpdateMarketplaceVoucher,
		)
		return nil, maskedErr
	}

	return voucher, nil
}

func (r *mpVoucherRepositoryImpl) DeleteMarketplaceVoucher(voucher *entity.MarketplaceVoucher) (*entity.MarketplaceVoucher, error) {
	err := r.db.Model(&voucher).
		Delete(voucher).Error
	if err != nil {
		return nil, domain.ErrDeleteMarketplaceVoucher
	}

	return voucher, nil
}

func (r *mpVoucherRepositoryImpl) DecreaseMarketplaceVoucherQuotaTx(tx *gorm.DB, voucherID uint) error {
	err := tx.Model(&entity.MarketplaceVoucher{}).
		Where("id = ?", voucherID).
		Update("quota", gorm.Expr("quota - ?", 1)).Error
	if err != nil {
		return domain.ErrDecreaseMarketplaceVoucherQuota
	}

	return nil
}

func (r *mpVoucherRepositoryImpl) IncreaseMarketplaceVoucherQuotaTx(tx *gorm.DB, voucherID uint) error {
	err := tx.Model(&entity.MarketplaceVoucher{}).
		Where("id = ?", voucherID).
		Update("quota", gorm.Expr("quota + ?", 1)).Error
	if err != nil {
		return domain.ErrIncreaseMarketplaceVoucherQuota
	}

	return nil
}
