package usecase

import (
	"encoding/json"
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/cronjob"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type RefundRequestUsecase interface {
	RequestRefundProcess(username string, req dto.RefundRequestFormReqDTO) (*dto.RefundRequestFormResDTO, error)

	UserCancelRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error)
	UserAcceptRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error)
	UserRejectRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error)

	MerchantRejectRefundProsess(username string, refundId uint) (*dto.RefundRequestDTO, error)
	MerchantAcceptRefundProsess(username string, refundId uint) (*dto.RefundRequestDTO, error)

	AdminRejectRefundProcess(refundId uint) (*dto.RefundRequestDTO, error)
	AdminAcceptRefundProcess(refundId uint) (*dto.RefundRequestDTO, error)

	GetUserRefundRequestList(username string, req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error)
	GetMerchantRefundRequestList(username string, req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error)
	GetAdminRefundRequestList(req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error)

	CronRefundRequestStatusToAcceptedBySeller()
	CronRefundRequestStatusToAcceptedByBuyer()
}

type RefundRequestUsecaseConfig struct {
	RefundRequestRepository repository.RefundRequestRepository
	TransactionRepository   repository.TransactionRepository
	UserRepository          repository.UserRepository
	GCSUploader             util.GCSUploader
	MerchantRepository      repository.MerchantRepository
	WalletRepository        repository.WalletRepository
	Cron                    *cronjob.CronJob
}

type refundRequestUsecaseImpl struct {
	refundRequestRepository repository.RefundRequestRepository
	transactionRepository   repository.TransactionRepository
	userRepository          repository.UserRepository
	gCSUploader             util.GCSUploader
	merchantRepository      repository.MerchantRepository
	walletRepository        repository.WalletRepository
	cron                    *cronjob.CronJob
}

func NewRefundRequestUsecase(c RefundRequestUsecaseConfig) RefundRequestUsecase {
	refundRequestUsecase := &refundRequestUsecaseImpl{
		refundRequestRepository: c.RefundRequestRepository,
		transactionRepository:   c.TransactionRepository,
		userRepository:          c.UserRepository,
		gCSUploader:             c.GCSUploader,
		merchantRepository:      c.MerchantRepository,
		walletRepository:        c.WalletRepository,
		cron:                    c.Cron,
	}

	c.Cron.AddJob("* * * * *", refundRequestUsecase.CronRefundRequestStatusToAcceptedBySeller)
	c.Cron.AddJob("* * * * *", refundRequestUsecase.CronRefundRequestStatusToAcceptedByBuyer)

	return refundRequestUsecase
}

func (u *refundRequestUsecaseImpl) RequestRefundProcess(username string, req dto.RefundRequestFormReqDTO) (*dto.RefundRequestFormResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// check if user already has wallet
	wallet, err := u.walletRepository.GetByUserId(user.ID)
	if err != nil || wallet == nil {
		return nil, domain.ErrUserWalletNotActivated
	}

	// check transaction status
	transaction, err := u.transactionRepository.GetTransactionDetailByInvoiceCode(user.ID, req.InvoiceCode)
	if err != nil {
		return nil, err
	}

	// check if transaction status is eligible
	if transaction.TransactionStatus.OnRequestRefundAt != nil {
		return nil, domain.ErrRefundTransactionAlreadyRequested
	}

	if transaction.TransactionStatus.OnDeliveredAt == nil ||
		transaction.TransactionStatus.OnCompletedAt != nil ||
		transaction.TransactionStatus.OnCanceledAt != nil {
		return nil, domain.ErrRefundTransactionNotEligible
	}

	//upload image
	imageUrl, err := u.gCSUploader.UploadFileFromFileHeader(*req.Image, fmt.Sprintf("refund_request-%d-%d", user.ID, transaction.ID))
	if err != nil {
		return nil, domain.ErrUploadFile
	}

	// create refund request
	reqfundReqStatuses := make([]entity.RefundRequestStatus, 1)
	refundRequest := entity.RefundRequest{
		TransactionID:         transaction.ID,
		Reason:                req.Reason,
		ImageUrl:              imageUrl,
		RefundRequestStatuses: reqfundReqStatuses,
	}
	createdRefundRequest, err := u.refundRequestRepository.AddRefundRequest(refundRequest)
	if err != nil {
		return nil, err
	}

	return &dto.RefundRequestFormResDTO{
		ID:            createdRefundRequest.ID,
		TransactionId: createdRefundRequest.TransactionID,
		Reason:        createdRefundRequest.Reason,
		ImageUrl:      createdRefundRequest.ImageUrl,
		CreatedAt:     createdRefundRequest.CreatedAt,
	}, nil
}

func (u *refundRequestUsecaseImpl) GetAdminRefundRequestList(req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error) {
	refundRequests, totalData, err := u.refundRequestRepository.GetRefundRequestList(req)
	if err != nil {
		return nil, err
	}

	refundRequestsDTO := u.buildRefundRequestListResDTO(refundRequests)

	return &dto.RefundRequestListResDTO{
		RefundRequests: refundRequestsDTO,
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalData,
			TotalPage:   (totalData + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
	}, nil
}

func (u *refundRequestUsecaseImpl) GetMerchantRefundRequestList(username string, req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	refundRequests, totalData, err := u.refundRequestRepository.GetRefundRequestListByMerchantDomain(merchant.Domain, req)
	if err != nil {
		return nil, err
	}

	refundRequestsDTO := u.buildRefundRequestListResDTO(refundRequests)

	return &dto.RefundRequestListResDTO{
		RefundRequests: refundRequestsDTO,
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalData,
			TotalPage:   (totalData + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
	}, nil
}

func (u *refundRequestUsecaseImpl) GetUserRefundRequestList(username string, req dto.RefundRequestListReqParamDTO) (*dto.RefundRequestListResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	refundRequests, totalData, err := u.refundRequestRepository.GetRefundRequestListByUserId(user.ID, req)
	if err != nil {
		return nil, err
	}

	refundRequestsDTO := u.buildRefundRequestListResDTO(refundRequests)

	return &dto.RefundRequestListResDTO{
		RefundRequests: refundRequestsDTO,
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalData,
			TotalPage:   (totalData + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
	}, nil

}

func (u *refundRequestUsecaseImpl) buildRefundRequestListResDTO(refundRequests []entity.RefundRequest) []dto.RefundRequestDTO {
	refundRequestsDTO := make([]dto.RefundRequestDTO, 0)
	for _, refundRequest := range refundRequests {
		refundStatuses := make([]dto.RefundRequestStatusDTO, 0)
		for _, refundStatus := range refundRequest.RefundRequestStatuses {
			refundStatuses = append(refundStatuses, dto.RefundRequestStatusDTO{
				CanceledByBuyerAt:  refundStatus.CanceledByBuyerAt,
				AcceptedByBuyerAt:  refundStatus.AcceptedByBuyerAt,
				RejectedByBuyerAt:  refundStatus.RejectedByBuyerAt,
				AcceptedBySellerAt: refundStatus.AcceptedBySellerAt,
				RejectedBySellerAt: refundStatus.RejectedBySellerAt,
				AcceptedByAdminAt:  refundStatus.AcceptedByAdminAt,
				RejectedByAdminAt:  refundStatus.RejectedByAdminAt,
				ClosedAt:           refundStatus.ClosedAt,
			})
		}

		refundRequestsDTO = append(refundRequestsDTO, dto.RefundRequestDTO{
			ID:             refundRequest.ID,
			TransactionId:  refundRequest.TransactionID,
			InvoiceCode:    refundRequest.Transaction.InvoiceCode,
			Username:       refundRequest.Transaction.User.Username,
			MerchantDomain: refundRequest.Transaction.Merchant.Domain,

			Reason:    refundRequest.Reason,
			ImageUrl:  refundRequest.ImageUrl,
			CreatedAt: refundRequest.CreatedAt,

			RefundRequestStatusesDTO: refundStatuses,
		})
	}

	return refundRequestsDTO
}

func (u *refundRequestUsecaseImpl) buildRefundRequestResDTO(refundRequest entity.RefundRequest, refReqStatusRes entity.RefundRequestStatus) *dto.RefundRequestDTO {
	return &dto.RefundRequestDTO{
		ID:             refundRequest.ID,
		TransactionId:  refundRequest.TransactionID,
		InvoiceCode:    refundRequest.Transaction.InvoiceCode,
		Username:       refundRequest.Transaction.User.Username,
		MerchantDomain: refundRequest.Transaction.Merchant.Domain,
		Reason:         refundRequest.Reason,
		ImageUrl:       refundRequest.ImageUrl,
		CreatedAt:      refundRequest.CreatedAt,
		RefundRequestStatusesDTO: []dto.RefundRequestStatusDTO{
			{
				CanceledByBuyerAt:  refReqStatusRes.CanceledByBuyerAt,
				AcceptedByBuyerAt:  refReqStatusRes.AcceptedByBuyerAt,
				RejectedByBuyerAt:  refReqStatusRes.RejectedByBuyerAt,
				AcceptedBySellerAt: refReqStatusRes.AcceptedBySellerAt,
				RejectedBySellerAt: refReqStatusRes.RejectedBySellerAt,
				AcceptedByAdminAt:  refReqStatusRes.AcceptedByAdminAt,
				RejectedByAdminAt:  refReqStatusRes.RejectedByAdminAt,
				ClosedAt:           refReqStatusRes.ClosedAt,
			},
		},
	}
}

func (u *refundRequestUsecaseImpl) getUserRefundRequest(username string, refundId uint) (*entity.RefundRequest, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	refundRequest, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil {
		return nil, err
	}

	if refundRequest.Transaction.UserId != user.ID {
		return nil, domain.ErrGetRefundRequestNotFound
	}

	if refundRequest.RefundRequestStatuses[0].CanceledByBuyerAt != nil {
		return nil, domain.ErrRefundRequestAlreadyCanceledOrProcessed
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedByBuyerAt != nil ||
		refundRequest.RefundRequestStatuses[0].RejectedByBuyerAt != nil {
		return nil, domain.ErrRefundRequestUserAlreadyAcceptedOrRejected
	}

	return refundRequest, nil
}

func (u *refundRequestUsecaseImpl) UserCancelRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error) {
	refundRequest, err := u.getUserRefundRequest(username, refundId)
	if err != nil {
		return nil, err
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedByAdminAt != nil ||
		refundRequest.RefundRequestStatuses[0].RejectedByAdminAt != nil {
		return nil, domain.ErrRefundRequestAlreadyCanceledOrProcessed
	}

	var amount float64
	var trxPaymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(refundRequest.Transaction.PaymentDetails.Bytes), &trxPaymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	amount = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MerchantVoucherNominal

	refReqStatusRes, err := u.refundRequestRepository.UserCancelRefundRequest(refundId, refundRequest.Transaction, amount, trxPaymentDetails.MarketplaceVoucherNominal)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundRequest, *refReqStatusRes), nil
}

// it means refund is closed with declined refund request
// fund will be forwarded to seller
// transaction status will be completed
func (u *refundRequestUsecaseImpl) UserAcceptRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error) {
	refundRequest, err := u.getUserRefundRequest(username, refundId)
	if err != nil {
		return nil, err
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedByAdminAt == nil &&
		refundRequest.RefundRequestStatuses[0].RejectedByAdminAt == nil {
		return nil, domain.ErrRefundRequestNotYetProcessedByAdmin
	}

	var amount float64
	var trxPaymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(refundRequest.Transaction.PaymentDetails.Bytes), &trxPaymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	amount = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MerchantVoucherNominal

	refReqStatusRes, err := u.refundRequestRepository.UserAcceptRefundRequest(refundId, refundRequest.Transaction, amount, trxPaymentDetails.MarketplaceVoucherNominal)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundRequest, *refReqStatusRes), nil
}

// it means refund will be continued, because user not agree with admin decision
// new status will be created
// but will be checked first if user already rejected three times
// if yes, then refund will be closed with declined refund request
// fund will be forwarded to seller
func (u *refundRequestUsecaseImpl) UserRejectRefundProcess(username string, refundId uint) (*dto.RefundRequestDTO, error) {
	refundRequest, err := u.getUserRefundRequest(username, refundId)
	if err != nil {
		return nil, err
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedByAdminAt == nil &&
		refundRequest.RefundRequestStatuses[0].RejectedByAdminAt == nil {
		return nil, domain.ErrRefundRequestNotYetProcessedByAdmin
	}

	if len(refundRequest.RefundRequestStatuses) >= 3 {
		return nil, domain.ErrRefundRequestUserAlreadyRejectedThreeTimes
	}

	refReqStatusRes, err := u.refundRequestRepository.UserRejectRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundRequest, *refReqStatusRes), nil
}

func (u *refundRequestUsecaseImpl) getMerchantRefundRequest(username string, refundId uint) (*entity.User, *entity.RefundRequest, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	refundRequest, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil {
		return nil, nil, err
	}

	if refundRequest.Transaction.Merchant.UserId != user.ID {
		return nil, nil, domain.ErrGetRefundRequestNotFound
	}

	if refundRequest.RefundRequestStatuses[0].RejectedBySellerAt != nil ||
		refundRequest.RefundRequestStatuses[0].AcceptedBySellerAt != nil ||
		refundRequest.RefundRequestStatuses[0].CanceledByBuyerAt != nil {
		return nil, nil, domain.ErrRefundRequestAlreadyCanceledOrProcessed
	}

	return user, refundRequest, nil
}

func (u *refundRequestUsecaseImpl) MerchantRejectRefundProsess(username string, refundId uint) (*dto.RefundRequestDTO, error) {
	_, refundReq, err := u.getMerchantRefundRequest(username, refundId)
	if err != nil {
		return nil, err
	}

	refReqStatusRes, err := u.refundRequestRepository.MerchantRejectRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundReq, *refReqStatusRes), nil
}

func (u *refundRequestUsecaseImpl) MerchantAcceptRefundProsess(username string, refundId uint) (*dto.RefundRequestDTO, error) {
	_, refundReq, err := u.getMerchantRefundRequest(username, refundId)
	if err != nil {
		return nil, err
	}

	refReqStatusRes, err := u.refundRequestRepository.MerchantAcceptRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundReq, *refReqStatusRes), nil
}

func (u *refundRequestUsecaseImpl) getAdminRefundRequest(refundId uint) (*entity.RefundRequest, error) {
	refundRequest, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil {
		return nil, err
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedBySellerAt == nil &&
		refundRequest.RefundRequestStatuses[0].RejectedBySellerAt == nil {
		return nil, domain.ErrRefundRequestNotYetProcessedBySeller
	}

	if refundRequest.RefundRequestStatuses[0].AcceptedByAdminAt != nil ||
		refundRequest.RefundRequestStatuses[0].RejectedByAdminAt != nil ||
		refundRequest.RefundRequestStatuses[0].CanceledByBuyerAt != nil {
		return nil, domain.ErrRefundRequestAlreadyCanceledOrProcessed
	}

	return refundRequest, nil
}

func (u *refundRequestUsecaseImpl) AdminAcceptRefundProcess(refundId uint) (*dto.RefundRequestDTO, error) {
	refundReq, err := u.getAdminRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	var trxCartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(refundReq.Transaction.CartItems.Bytes), &trxCartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItems
	}

	var amount float64
	var trxPaymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(refundReq.Transaction.PaymentDetails.Bytes), &trxPaymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	amount = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MerchantVoucherNominal

	refReqStatusRes, err := u.refundRequestRepository.AdminAcceptRefundRequest(refundId, refundReq.Transaction, amount, trxCartItems)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundReq, *refReqStatusRes), nil
}

// check if there are three times rejected by user
// if yes, then refund request will be closed
// fund will be forwarded to seller
// transaction status will be changed to completed
// if not, then refund request will be continued
func (u *refundRequestUsecaseImpl) AdminRejectRefundProcess(refundId uint) (*dto.RefundRequestDTO, error) {
	refundReq, err := u.getAdminRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	var amount float64
	var trxPaymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(refundReq.Transaction.PaymentDetails.Bytes), &trxPaymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	amount = trxPaymentDetails.Subtotal + trxPaymentDetails.DeliveryFee - trxPaymentDetails.MerchantVoucherNominal

	if len(refundReq.RefundRequestStatuses) >= 3 {
		refReqStatusRes, err := u.refundRequestRepository.AdminRejectRefundRequestClosed(refundId, refundReq.Transaction, amount, trxPaymentDetails.MarketplaceVoucherNominal)
		if err != nil {
			return nil, err
		}

		return u.buildRefundRequestResDTO(*refundReq, *refReqStatusRes), nil
	}

	refReqStatusRes, err := u.refundRequestRepository.AdminRejectRefundRequest(refundId)
	if err != nil {
		return nil, err
	}

	return u.buildRefundRequestResDTO(*refundReq, *refReqStatusRes), nil
}

func (u *refundRequestUsecaseImpl) CronRefundRequestStatusToAcceptedBySeller() {
	// log.Info().Msg("CronRefundRequestStatusToAcceptedBySeller Running")
	u.refundRequestRepository.UpdateAllRefundRequestStatusToAcceptedBySeller()
}

func (u *refundRequestUsecaseImpl) CronRefundRequestStatusToAcceptedByBuyer() {
	// log.Info().Msg("CronRefundRequestStatusToAcceptedByBuyer Running")
	u.refundRequestRepository.UpdateAllRefundRequestStatusToAcceptedByBuyer()
}
