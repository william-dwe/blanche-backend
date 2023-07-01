package repository

import (
	"encoding/json"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	AddRefundRequest(req entity.RefundRequest) (*entity.RefundRequest, error)
	AddMessage(msg entity.RefundReqMessage) (*entity.RefundReqMessage, error)
	GetMessages(refundReqId uint) ([]entity.RefundReqMessage, error)

	GetRefundRequestList(req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error)
	GetRefundRequestListByMerchantDomain(merchantDomain string, req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error)
	GetRefundRequestListByUserId(userId uint, req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error)

	GetRefundRequestById(refundReqId uint) (*entity.RefundRequest, error)

	UserCancelRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, amountVoucherMp float64) (*entity.RefundRequestStatus, error)
	UserAcceptRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (resRefundReq *entity.RefundRequestStatus, errAcceptRefund error)
	UserRejectRefundRequest(refundReqId uint) (resRefundReq *entity.RefundRequestStatus, errAcceptRefund error)

	MerchantAcceptRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error)
	MerchantRejectRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error)

	AdminAcceptRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (*entity.RefundRequestStatus, error)
	AdminRejectRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error)
	AdminRejectRefundRequestClosed(refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (*entity.RefundRequestStatus, error)

	UpdateAllRefundRequestStatusToAcceptedBySeller() error
	UpdateAllRefundRequestStatusToAcceptedByBuyer() error
}

type RefundRequestRepositoryConfig struct {
	DB                    *gorm.DB
	TransactionRepository TransactionRepository
}

type refundRequestRepositoryImpl struct {
	db                    *gorm.DB
	transactionRepository TransactionRepository
}

func NewRefundRequestRepository(c RefundRequestRepositoryConfig) RefundRequestRepository {
	return &refundRequestRepositoryImpl{
		db:                    c.DB,
		transactionRepository: c.TransactionRepository,
	}
}

func (r *refundRequestRepositoryImpl) AddRefundRequest(req entity.RefundRequest) (resRefundReq *entity.RefundRequest, addRefundReqErr error) {
	// begin transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in AddRefundRequest repo: %v", r)
			addRefundReqErr = domain.ErrCreateRefundRequest
		}
	}()

	// create refund request
	err := tx.Create(&req).Error
	if err != nil {
		tx.Rollback()
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"unique_refund_requests_transaction_id": domain.ErrCreateRefundRequestDuplicate,
			},
			domain.ErrCreateRefundRequest,
		)
		return nil, maskedErr
	}

	// update transaction status to request refund
	err = tx.Model(&entity.TransactionStatus{}).Where("transaction_id = ?", req.TransactionID).Update("on_request_refund_at", "now()").Error
	if err != nil {
		tx.Rollback()
		return nil, domain.ErrCreateRefundRequestUpdateStatus
	}

	// commit transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, domain.ErrCreateRefundRequest
	}

	return &req, nil
}

func (r *refundRequestRepositoryImpl) AddMessage(msg entity.RefundReqMessage) (*entity.RefundReqMessage, error) {
	err := r.db.Create(&msg).Error
	if err != nil {
		return nil, domain.ErrAddRefundRequestMessage
	}

	return &msg, nil
}

func (r *refundRequestRepositoryImpl) GetMessages(refundReqId uint) ([]entity.RefundReqMessage, error) {
	panic("implement me")
}

func (r *refundRequestRepositoryImpl) GetRefundRequestList(req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error) {
	qFilter := r.subQueryFilterProcess(dto.RefundRequestFilter(req.Status))

	var refundRequests []entity.RefundRequest
	var count int64
	pageOffset := req.Limit * (req.Page - 1)

	err := r.db.
		Where("refund_requests.id in (?)", qFilter).
		Limit(req.Limit).
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Transaction.Merchant").
		Preload("RefundRequestStatuses", func(db *gorm.DB) *gorm.DB {
			return db.Order("refund_request_statuses.created_at DESC")
		}).
		Offset(pageOffset).
		Order("created_at DESC").
		Find(&refundRequests).
		Count(&count).
		Error

	if err != nil {
		return nil, 0, domain.ErrGetRefundRequestList
	}

	return refundRequests, count, nil
}

func (r *refundRequestRepositoryImpl) GetRefundRequestListByMerchantDomain(merchantDomain string, req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error) {
	qFilter := r.subQueryFilterProcessMerchantDomain(dto.RefundRequestFilter(req.Status), merchantDomain)

	var refundRequests []entity.RefundRequest
	var count int64

	pageOffset := req.Limit * (req.Page - 1)

	err := r.db.
		Where("refund_requests.id in (?)", qFilter).
		Limit(req.Limit).
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Transaction.Merchant").
		Preload("RefundRequestStatuses", func(db *gorm.DB) *gorm.DB {
			return db.Order("refund_request_statuses.created_at DESC")
		}).
		Offset(pageOffset).
		Order("created_at DESC").
		Find(&refundRequests).
		Count(&count).
		Error

	if err != nil {
		return nil, 0, domain.ErrGetRefundRequestList
	}

	return refundRequests, count, nil
}

func (r *refundRequestRepositoryImpl) GetRefundRequestListByUserId(userId uint, req dto.RefundRequestListReqParamDTO) ([]entity.RefundRequest, int64, error) {
	var refundRequests []entity.RefundRequest
	var count int64

	pageOffset := req.Limit * (req.Page - 1)

	err := r.db.
		Limit(req.Limit).
		Where("refund_requests.id in (?)", r.subQueryFilterProcessUserId(dto.RefundRequestFilter(req.Status), userId)).
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Transaction.Merchant").
		Preload("RefundRequestStatuses", func(db *gorm.DB) *gorm.DB {
			return r.db.Order("refund_request_statuses.created_at DESC")
		}).
		Offset(pageOffset).
		Order("created_at DESC").
		Find(&refundRequests).
		Count(&count).
		Error

	if err != nil {
		return nil, 0, domain.ErrGetRefundRequestList
	}

	return refundRequests, count, nil
}

func (r *refundRequestRepositoryImpl) GetRefundRequestById(refundReqId uint) (*entity.RefundRequest, error) {
	var refundRequest entity.RefundRequest
	err := r.db.
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Transaction.Merchant").
		Preload("Transaction.TransactionStatus").
		Preload("RefundRequestStatuses", func(db *gorm.DB) *gorm.DB {
			return db.Order("refund_request_statuses.created_at DESC")
		}).
		Where("id = ?", refundReqId).
		First(&refundRequest).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetRefundRequestNotFound
		}

		return nil, domain.ErrGetRefundRequest
	}

	return &refundRequest, nil
}

func (r *refundRequestRepositoryImpl) UserCancelRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (refundReq *entity.RefundRequestStatus, errCancelRefund error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in CancelRefundRequest repo: %v", r)
			errCancelRefund = domain.ErrCancelRefundRequest
		}
	}()

	var refundRequestStatus entity.RefundRequestStatus
	err := tx.Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Updates(map[string]interface{}{"canceled_by_buyer_at": "now()", "closed_at": "now()"}).
		Find(&refundRequestStatus).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update refund request canceled_by_buyer_at: %v", err)
		return nil, domain.ErrCancelRefundRequest
	}

	err = tx.
		Model(&entity.TransactionStatus{}).
		Where("transaction_id = ?", transaction.ID).
		Update("on_completed_at", "now()").Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status on_completed_at: %v", err)
		return nil, domain.ErrCancelRefundRequestUpdateTransactionStatus
	}

	// update transaction status like completed, use transction repo
	_, err = r.transactionRepository.UpdateTransactionStatusCompletedTx(tx, transaction, amount, amountPromotionMp)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status completed: %v", err)
		return nil, domain.ErrCancelRefundRequestUpdateTransactionStatus
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error commit refund cancel: %v", err)
		return nil, domain.ErrCancelRefundRequest
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) MerchantAcceptRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error) {
	var refundRequestStatus entity.RefundRequestStatus
	res := r.db.
		Model(&refundRequestStatus).
		Where("id = (?)", r.db.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Update("accepted_by_seller_at", "now()").Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, domain.ErrMerchantAcceptRefundRequest
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) MerchantRejectRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error) {
	var refundRequestStatus entity.RefundRequestStatus
	res := r.db.
		Model(&refundRequestStatus).
		Where("id = (?)", r.db.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Update("rejected_by_seller_at", "now()").Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, domain.ErrMerchantRejectRefundRequest
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) AdminAcceptRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, cartItems []entity.TransactionCartItem) (resRefReq *entity.RefundRequestStatus, errAcceptRefund error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in AdminAcceptRefundRequest repo: %v", r)
			errAcceptRefund = domain.ErrAdminRejectRefundRequest
		}
	}()

	var refundRequestStatus entity.RefundRequestStatus
	res := tx.
		Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Updates(map[string]interface{}{"accepted_by_admin_at": "now()", "closed_at": "now()"}).
		Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, domain.ErrAdminAcceptRefundRequest
	}

	// update transaction status refunded and return amount to user wallet
	_, err := r.transactionRepository.UpdateTransactionStatusRefundedTx(tx, transaction, amount, cartItems)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status refunded: %v", err)
		return nil, domain.ErrAdminAcceptRefundRequestUpdateTransactionStatus
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error commit refund accept: %v", err)
		return nil, domain.ErrAdminAcceptRefundRequestCommit
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) AdminRejectRefundRequestClosed(refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (refReqRes *entity.RefundRequestStatus, errRejectRefund error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in AdminRejectRefundRequestClosed repo: %v", r)
			errRejectRefund = domain.ErrAdminRejectRefundRequest
		}
	}()

	var refundRequestStatus entity.RefundRequestStatus
	res := tx.
		Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Updates(map[string]interface{}{"rejected_by_admin_at": "now()", "closed_at": "now()"}).
		Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		log.Error().Msgf("Error update refund request rejected by admin closed: %v", res.Error)
		return nil, domain.ErrAdminRejectRefundRequest
	}

	// update transaction status like completed, use transction repo
	// forward fund to merchant
	// update transaction status like completed, use transction repo
	_, err := r.transactionRepository.UpdateTransactionStatusCompletedTx(tx, transaction, amount, amountPromotionMp)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status completed: %v", err)
		return nil, domain.ErrAdminRejectRefundRequestUpdateTransactionStatus
	}

	// commit
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error commit admin reject refund request closed: %v", err)
		return nil, domain.ErrAdminRejectRefundRequestCommit
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) AdminRejectRefundRequest(refundReqId uint) (*entity.RefundRequestStatus, error) {
	var refundRequestStatus entity.RefundRequestStatus
	res := r.db.
		Model(&refundRequestStatus).
		Where("id = (?)", r.db.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Update("rejected_by_admin_at", "now()").Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, domain.ErrAdminRejectRefundRequest
	}

	return &refundRequestStatus, nil

}

func (r *refundRequestRepositoryImpl) UserRejectRefundRequest(refundReqId uint) (resRefundReq *entity.RefundRequestStatus, errRejectRefund error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UserRejectRefundRequest repo: %v", r)
			errRejectRefund = domain.ErrUserAcceptRefundRequest
		}
	}()

	var refundRequestStatus entity.RefundRequestStatus
	res := tx.
		Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Update("rejected_by_buyer_at", "now()").Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, domain.ErrUserRejectRefundRequestUpdateRefundRequestStatus
	}

	// update make new refund request status
	err := tx.Create(&entity.RefundRequestStatus{
		RefundRequestId: refundReqId,
	}).Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error create new refund request status: %v", err)
		return nil, domain.ErrUserRejectRefundRequestNewRefundRequestStatus
	}

	// commit
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error commit user reject refund request: %v", err)
		return nil, domain.ErrUserRejectRefundRequestCommit
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) UserAcceptRefundRequest(refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (resRefundReq *entity.RefundRequestStatus, errAcceptRefund error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UserAcceptRefundRequest repo: %v", r)
			errAcceptRefund = domain.ErrUserAcceptRefundRequest
		}
	}()

	var refundRequestStatus entity.RefundRequestStatus
	res := tx.
		Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Updates(map[string]interface{}{"accepted_by_buyer_at": "now()", "closed_at": "now()"}).
		Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		log.Error().Msgf("Error update refund request status when user accept stop refund request: %v", res.Error)
		return nil, domain.ErrUserAcceptRefundRequestUpdateRefundRequestStatus
	}

	// update transcation status to completed
	// so refund is rejected (because user accept closing refund request)
	_, err := r.transactionRepository.UpdateTransactionStatusCompletedTx(tx, transaction, amount, amountPromotionMp)
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update transaction status completed when user accept stop refund request: %v", err)
		return nil, domain.ErrUserAcceptRefundRequestUpdateTransactionStatus
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error commit transaction when user accept stop refund request: %v", err)
		return nil, domain.ErrUserAcceptRefundRequestCommitTransaction
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) UserAcceptRefundRequestTx(tx *gorm.DB, refundReqId uint, transaction entity.Transaction, amount float64, amountPromotionMp float64) (resRefundReq *entity.RefundRequestStatus, errAcceptRefund error) {
	var refundRequestStatus entity.RefundRequestStatus
	res := tx.
		Model(&refundRequestStatus).
		Where("id = (?)", tx.Model(&entity.RefundRequestStatus{}).
			Select("id").Where("refund_request_id = ?", refundReqId).Order("created_at DESC").Limit(1)).
		Updates(map[string]interface{}{"accepted_by_buyer_at": "now()", "closed_at": "now()"}).
		Find(&refundRequestStatus)

	if res.Error != nil || res.RowsAffected == 0 {
		log.Error().Msgf("Error update refund request status when user accept stop refund request: %v", res.Error)
		return nil, domain.ErrUserAcceptRefundRequestUpdateRefundRequestStatus
	}

	// update transcation status to completed
	// so refund is rejected (because user accept closing refund request)
	_, err := r.transactionRepository.UpdateTransactionStatusCompletedTx(tx, transaction, amount, amountPromotionMp)
	if err != nil {
		log.Error().Msgf("Error update transaction status completed when user accept stop refund request: %v", err)
		return nil, domain.ErrUserAcceptRefundRequestUpdateTransactionStatus
	}

	return &refundRequestStatus, nil
}

func (r *refundRequestRepositoryImpl) subQueryFilterProcess(filter dto.RefundRequestFilter) *gorm.DB {
	subQuery := r.db.
		Model(&entity.RefundRequestStatus{}).
		Joins("LEFT JOIN refund_requests ON refund_requests.id = refund_request_statuses.refund_request_id").
		Joins("LEFT JOIN transactions ON transactions.id = refund_requests.transaction_id").
		Select("refund_requests.id")

	return r.subQueryFilterStatusProcess(filter, subQuery)
}

func (r *refundRequestRepositoryImpl) subQueryFilterProcessMerchantDomain(filter dto.RefundRequestFilter, merchantDomain string) *gorm.DB {
	subQuery := r.subQueryFilterProcess(filter).
		Where("transactions.merchant_domain = ?", merchantDomain)

	return r.subQueryFilterStatusProcess(filter, subQuery)
}

func (r *refundRequestRepositoryImpl) subQueryFilterProcessUserId(filter dto.RefundRequestFilter, userId uint) *gorm.DB {
	subQuery := r.subQueryFilterProcess(filter).
		Where("transactions.user_id = ?", userId)

	return r.subQueryFilterStatusProcess(filter, subQuery)
}

func (r *refundRequestRepositoryImpl) subQueryFilterStatusProcess(filter dto.RefundRequestFilter, subQuery *gorm.DB) *gorm.DB {
	subQueryNewestId := r.db.Model(&entity.RefundRequestStatus{}).Group("refund_request_id").Select("MAX(id) AS newest_id")
	subQuery = subQuery.Where("refund_request_statuses.id IN (?)", subQueryNewestId)

	switch filter {
	case dto.RefundRequestFilterClosed:
		subQuery = subQuery.Where("refund_request_statuses.closed_at IS NOT NULL")
	case dto.RefundRequestFilterRefunded:
		subQuery = subQuery.Where("refund_request_statuses.accepted_by_admin_at IS NOT NULL")
	case dto.RefundRequestFilterCanceled:
		subQuery = subQuery.Where("refund_request_statuses.canceled_by_buyer_at IS NOT NULL")
	case dto.RefundRequestFilterRejected:
		subQuery = subQuery.Where("refund_request_statuses.rejected_by_admin_at IS NOT NULL").
			Where("refund_request_statuses.closed_at IS NOT NULL")
	case dto.RefundRequestFilterWaitingAdminAproval:
		subQuery = subQuery.
			Where("closed_at is NULL").
			Where("refund_request_statuses.accepted_by_admin_at IS NULL AND refund_request_statuses.rejected_by_admin_at IS NULL").
			Where("refund_request_statuses.accepted_by_seller_at IS NOT NULL OR refund_request_statuses.rejected_by_seller_at IS NOT NULL")
	case dto.RefundRequestFilterWaitingMerchantAproval:
		subQuery = subQuery.
			Where("closed_at is NULL").
			Where("refund_request_statuses.accepted_by_admin_at IS NULL AND refund_request_statuses.rejected_by_admin_at IS NULL").
			Where("refund_request_statuses.accepted_by_seller_at IS NULL AND refund_request_statuses.rejected_by_seller_at IS NULL")
	case dto.RefundRequestFilterWaitingBuyerAproval:
		subQuery = subQuery.
			Where("closed_at is NULL").
			Where("refund_request_statuses.accepted_by_buyer_at IS NULL AND refund_request_statuses.rejected_by_buyer_at IS NULL AND refund_request_statuses.canceled_by_buyer_at IS NULL").
			Where("refund_request_statuses.accepted_by_admin_at IS NULL AND refund_request_statuses.rejected_by_admin_at IS NOT NULL").
			Where("refund_request_statuses.accepted_by_seller_at IS NOT NULL OR refund_request_statuses.rejected_by_seller_at IS NOT NULL")
	}

	return subQuery
}

func (r *refundRequestRepositoryImpl) UpdateAllRefundRequestStatusToAcceptedBySeller() error {
	var refundRequestStatuses []entity.RefundRequestStatus
	// update which created_at is 24 hours ago and not accepted by seller yet
	err := r.db.
		Model(&refundRequestStatuses).
		Where("created_at <= ?", time.Now().Add(-24*time.Hour)).
		Where("accepted_by_seller_at IS NULL").
		Where("rejected_by_seller_at IS NULL").
		Where("accepted_by_admin_at IS NOT NULL").
		Where("rejected_by_admin_at IS NULL").
		Where("canceled_by_buyer_at IS NULL").
		Where("rejected_by_buyer_at IS NULL").
		Where("accepted_by_buyer_at IS NULL").
		Where("closed_at IS NULL").
		Update("accepted_by_seller_at", time.Now()).
		Error

	if err != nil {
		log.Error().Msgf("Error update refund request status to accepted by seller: %v", err)
		return domain.ErrCronRefundRequestStatusToAcceptedBySeller
	}

	return nil
}

func (r *refundRequestRepositoryImpl) UpdateAllRefundRequestStatusToAcceptedByBuyer() (errUpdateAll error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Error().Msgf("Recovered in UpdateAllRefundRequestStatusToAcceptedByBuyer repo: %v", r)
			errUpdateAll = domain.ErrUpdateAllRefundRequestStatus
		}
	}()
	// get all refund request status which created_at is 24 hours ago and not accepted by buyer yet
	var refundRequestStatuses []entity.RefundRequestStatus
	err := tx.
		Model(&refundRequestStatuses).
		Preload("RefundRequest").
		Preload("RefundRequest.Transaction").
		Preload("RefundRequest.Transaction.Merchant").
		Preload("RefundRequest.Transaction.TransactionStatus").
		Where("rejected_by_admin_at <= ?", time.Now().Add(-24*time.Hour)).
		Where("accepted_by_buyer_at IS NULL").
		Where("rejected_by_buyer_at IS NULL").
		Where("accepted_by_admin_at IS NULL").
		Where("rejected_by_admin_at IS NOT NULL").
		Where("canceled_by_buyer_at IS NULL").
		Where("rejected_by_seller_at IS NOT NULL OR accepted_by_seller_at IS NOT NULL").
		Where("closed_at IS NULL").
		Find(&refundRequestStatuses).
		Error

	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update refund request status to accepted by buyer: %v", err)
		return domain.ErrCronRefundRequestStatusToAcceptedByBuyer
	}

	//update as accepted by buyer
	for _, refundRequestStatus := range refundRequestStatuses {
		amount, mpAmount, err := r.countAmountAndPromotionTrx(refundRequestStatus.RefundRequest.Transaction)
		if err != nil {
			log.Error().Msgf("Error update refund request status to accepted by buyer: %v", err)
			continue
		}
		_, err = r.UserAcceptRefundRequestTx(tx, refundRequestStatus.RefundRequestId, refundRequestStatus.RefundRequest.Transaction, amount, mpAmount)
		if err != nil {
			log.Error().Msgf("Error update refund request status to accepted by buyer update each: %v", err)
			continue
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		log.Error().Msgf("Error update refund request status to accepted by buyer commit: %v", err)
		return domain.ErrCronRefundRequestStatusToAcceptedByBuyer
	}

	return nil
}

func (r *refundRequestRepositoryImpl) countAmountAndPromotionTrx(transaction entity.Transaction) (float64, float64, error) {
	var amount float64
	var promotionMarketplace float64

	var transactionPaymentDetails entity.TransactionPaymentDetails
	err := json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &transactionPaymentDetails)
	if err != nil {
		return 0, 0, domain.ErrUnmarshalJSONPaymentDetails
	}

	amount = transactionPaymentDetails.Subtotal + transactionPaymentDetails.DeliveryFee - transactionPaymentDetails.MerchantVoucherNominal
	promotionMarketplace = transactionPaymentDetails.MarketplaceVoucherNominal

	return amount, promotionMarketplace, nil
}
