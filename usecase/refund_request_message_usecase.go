package usecase

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
)

type RefundRequestMessageUsecase interface {
	GetListMessageByRefundRequestId(username string, refundRequestId uint) (*dto.RefundRequestMsgListResDTO, error)
	GetAdminListMessageByRefundRequestId(refundRequestId uint) (*dto.RefundRequestMsgListResDTO, error)

	AdminAddMessage(refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error)
	MerchantAddMessage(username string, refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error)
	BuyerAddMessage(username string, refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error)
}

type RefundRequestMessageUsecaseConfig struct {
	RefundRequestMessageRepository repository.RefundRequestMessageRepository
	RefundRequestRepository        repository.RefundRequestRepository
	UserRepository                 repository.UserRepository
	GcsUploader                    util.GCSUploader
}

type refundRequestMessageUsecaseImpl struct {
	refundRequestMessageRepository repository.RefundRequestMessageRepository
	refundRequestRepository        repository.RefundRequestRepository
	userRepository                 repository.UserRepository
	gcsUploader                    util.GCSUploader
}

func NewRefundRequestMessageUsecase(c RefundRequestMessageUsecaseConfig) RefundRequestMessageUsecase {
	return &refundRequestMessageUsecaseImpl{
		refundRequestMessageRepository: c.RefundRequestMessageRepository,
		refundRequestRepository:        c.RefundRequestRepository,
		userRepository:                 c.UserRepository,
		gcsUploader:                    c.GcsUploader,
	}
}

func (u *refundRequestMessageUsecaseImpl) MerchantAddMessage(username string, refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	refundReq, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil {
		return nil, err
	}

	if refundReq.RefundRequestStatuses[0].ClosedAt != nil {
		return nil, domain.ErrRefundRequestClosed
	}

	if refundReq.Transaction.Merchant.UserId != user.ID {
		return nil, domain.ErrGetRefundRequestNotFound
	}

	return u.addRefundReqMessage(refundId, user.ID, dto.REFUND_REQ_MSG_ROLE_MERCHANT_ID, message)
}

func (u *refundRequestMessageUsecaseImpl) BuyerAddMessage(username string, refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	refundReq, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil {
		return nil, err
	}

	if refundReq.RefundRequestStatuses[0].ClosedAt != nil {
		return nil, domain.ErrRefundRequestClosed
	}

	if refundReq.Transaction.UserId != user.ID {
		return nil, domain.ErrGetRefundRequestNotFound
	}

	return u.addRefundReqMessage(refundId, user.ID, dto.REFUND_REQ_MSG_ROLE_BUYER_ID, message)
}

func (u *refundRequestMessageUsecaseImpl) AdminAddMessage(refundId uint, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error) {
	refundReq, err := u.refundRequestRepository.GetRefundRequestById(refundId)
	if err != nil || refundReq == nil {
		return nil, err
	}

	if refundReq.RefundRequestStatuses[0].ClosedAt != nil {
		return nil, domain.ErrRefundRequestClosed
	}

	return u.addRefundReqMessage(refundId, 0, dto.REFUND_REQ_MSG_ROLE_ADMIN_ID, message)
}

func (u *refundRequestMessageUsecaseImpl) addRefundReqMessage(refundId uint, userId uint, roleId int, message dto.RefundRequestMsgFormReqDTO) (*dto.RefundRequestMsgResDTO, error) {
	newMsg := entity.RefundReqMessage{
		RefundRequestId:        refundId,
		RefundReqMessageRoleId: uint(roleId),
		Message:                message.Message,
	}

	//upload image if there is any
	if message.Image != nil {
		url, err := u.gcsUploader.UploadFileFromFileHeader(*message.Image,
			fmt.Sprintf("ref_req_%d_%d", userId, refundId))
		if err != nil {
			log.Error().Msgf("Failed to upload file: %v", err)
			return nil, domain.ErrUploadFile
		}
		newMsg.ImageUrl = &url
	}

	createdMsg, err := u.refundRequestMessageRepository.AddMessage(newMsg)
	if err != nil {
		return nil, err
	}

	res := dto.RefundRequestMsgResDTO{
		ID:        createdMsg.ID,
		Message:   createdMsg.Message,
		Image:     createdMsg.ImageUrl,
		CreatedAt: createdMsg.CreatedAt,
	}

	return &res, nil
}

func (u *refundRequestMessageUsecaseImpl) getRefundReqMessageList(refundRequest entity.RefundRequest) (*dto.RefundRequestMsgListResDTO, error) {
	refundReq, err := u.refundRequestRepository.GetRefundRequestById(refundRequest.ID)
	if err != nil {
		return nil, err
	}

	msgList, err := u.refundRequestMessageRepository.GetAllMessageByRefundRequestId(refundRequest.ID)
	if err != nil {
		return nil, err
	}

	res := dto.RefundRequestMsgListResDTO{
		Messages:            make([]dto.RefundRequestMsgResDTO, 0),
		RefundRequestStatus: make([]dto.RefundRequestStatusDTO, 0),
		Details: dto.RefundRequestMsgDetailsDTO{
			RefundId:       refundRequest.ID,
			TransactionId:  refundRequest.Transaction.ID,
			InvoiceCode:    refundRequest.Transaction.InvoiceCode,
			BuyerUsername:  refundRequest.Transaction.User.Username,
			MerchantDomain: refundRequest.Transaction.MerchantDomain,
			Reason:         refundRequest.Reason,
			ImageUrl:       refundRequest.ImageUrl,
			ClosedAt:       refundRequest.RefundRequestStatuses[0].ClosedAt,
		},
	}

	for _, status := range refundReq.RefundRequestStatuses {
		res.RefundRequestStatus = append(res.RefundRequestStatus, dto.RefundRequestStatusDTO{
			CanceledByBuyerAt:  status.CanceledByBuyerAt,
			AcceptedByBuyerAt:  status.AcceptedByBuyerAt,
			RejectedByBuyerAt:  status.RejectedByBuyerAt,
			AcceptedBySellerAt: status.AcceptedBySellerAt,
			RejectedBySellerAt: status.RejectedBySellerAt,
			AcceptedByAdminAt:  status.AcceptedByAdminAt,
			RejectedByAdminAt:  status.RejectedByAdminAt,
			ClosedAt:           status.ClosedAt,
		})
	}

	for _, msg := range msgList {
		name := ""
		if msg.RefundReqMessageRoleId == dto.REFUND_REQ_MSG_ROLE_ADMIN_ID {
			name = "blanche_admin"
		} else if msg.RefundReqMessageRoleId == dto.REFUND_REQ_MSG_ROLE_MERCHANT_ID {
			name = refundRequest.Transaction.MerchantDomain
		} else if msg.RefundReqMessageRoleId == dto.REFUND_REQ_MSG_ROLE_BUYER_ID {
			name = refundRequest.Transaction.User.Username
		}

		res.Messages = append(res.Messages, dto.RefundRequestMsgResDTO{
			ID:         msg.ID,
			Message:    msg.Message,
			Image:      msg.ImageUrl,
			CreatedAt:  msg.CreatedAt,
			SenderName: name,
			Role: dto.RefundRequestMsgRoleDTO{
				ID:   msg.RefundReqMessageRoleId,
				Name: msg.RefundReqMessageRole.RoleName,
			},
		})
	}

	return &res, nil
}

func (u *refundRequestMessageUsecaseImpl) GetListMessageByRefundRequestId(username string, refundRequestId uint) (*dto.RefundRequestMsgListResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	refundReq, err := u.refundRequestRepository.GetRefundRequestById(refundRequestId)
	if err != nil {
		return nil, err
	}

	if refundReq.Transaction.Merchant.UserId != user.ID && refundReq.Transaction.UserId != user.ID {
		return nil, domain.ErrGetRefundRequestNotFound
	}

	return u.getRefundReqMessageList(*refundReq)
}

func (u *refundRequestMessageUsecaseImpl) GetAdminListMessageByRefundRequestId(refundRequestId uint) (*dto.RefundRequestMsgListResDTO, error) {
	refundRequest, err := u.refundRequestRepository.GetRefundRequestById(refundRequestId)
	if err != nil {
		return nil, err
	}

	return u.getRefundReqMessageList(*refundRequest)
}
