package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/rs/zerolog/log"
)

type TransactionUsecase interface {
	GetTransactionList(username string, req dto.TransactionReqParamDTO) (*dto.TransactionListResDTO, error)
	GetSellerTransactionList(username string, req dto.TransactionReqParamDTO) (*dto.TransactionSellerListResDTO, error)
	GetTransactionDetail(username string, invoiceCode string) (*dto.TransactionDetailResDTO, error)
	GetSellerTransactionDetail(username string, invoiceCode string) (*dto.TransactionSellerDetailResDTO, error)

	MakeTransaction(username string, req dto.MakeTransactionReqDTO) (*dto.MakeTransactionResDTO, error)
	HandleTransactionSlpRes(trxRecords []entity.TransactionPaymentRecord, isSuccess bool, paymentRec entity.PaymentRecord) error

	UpdateMerchantTransactionStatus(username string, req dto.UpdateMerchantTransactionStatusReqDTO) (*dto.UpdateMerchantTransactionStatusResDTO, error)
	UpdateUserTransactionStatus(username string, req dto.UpdateUserTransactionStatusReqDTO) (*dto.UpdateUserTransactionStatusResDTO, error)
}

type transactionUsecaseImpl struct {
	cartItemRepository                  repository.CartItemRepository
	transactionRepository               repository.TransactionRepository
	transactionDeliveryStatusRepository repository.TransactionDeliveryStatusRepository
	transactionStatusRepository         repository.TransactionStatusRepository
	userRepository                      repository.UserRepository
	merchantRepository                  repository.MerchantRepository
	transactionStatusUsecase            TransactionStatusUsecase
	transactionDeliveryStatusUsecase    TransactionDeliveryStatusUsecase
	orderItemUsecase                    OrderItemUsecase
	sealabspayRepository                repository.SealabspayRepository
	walletRepository                    repository.WalletRepository
	paymentMethodRepository             repository.PaymentMethodRepository
}

type TransactionUsecaseConfig struct {
	CartItemRepository                  repository.CartItemRepository
	TransactionRepository               repository.TransactionRepository
	TransactionDeliveryStatusRepository repository.TransactionDeliveryStatusRepository
	TransactionStatusRepository         repository.TransactionStatusRepository
	UserRepository                      repository.UserRepository
	MerchantRepository                  repository.MerchantRepository
	TransactionStatusUsecase            TransactionStatusUsecase
	TransactionDeliveryStatusUsecase    TransactionDeliveryStatusUsecase
	OrderItemUsecase                    OrderItemUsecase
	SealabspayRepository                repository.SealabspayRepository
	WalletRepository                    repository.WalletRepository
	PaymentMethodRepository             repository.PaymentMethodRepository
}

func NewTransactionUsecase(c TransactionUsecaseConfig) TransactionUsecase {
	return &transactionUsecaseImpl{
		cartItemRepository:                  c.CartItemRepository,
		transactionRepository:               c.TransactionRepository,
		userRepository:                      c.UserRepository,
		merchantRepository:                  c.MerchantRepository,
		transactionDeliveryStatusRepository: c.TransactionDeliveryStatusRepository,
		transactionStatusRepository:         c.TransactionStatusRepository,
		transactionStatusUsecase:            c.TransactionStatusUsecase,
		transactionDeliveryStatusUsecase:    c.TransactionDeliveryStatusUsecase,
		orderItemUsecase:                    c.OrderItemUsecase,
		sealabspayRepository:                c.SealabspayRepository,
		walletRepository:                    c.WalletRepository,
		paymentMethodRepository:             c.PaymentMethodRepository,
	}
}

func (u *transactionUsecaseImpl) GetSellerTransactionList(username string, req dto.TransactionReqParamDTO) (*dto.TransactionSellerListResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	transactions, totalTransactions, err := u.transactionRepository.GetTransactionListByMerchant(merchant.Domain, req)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return &dto.TransactionSellerListResDTO{
			PaginationResponse: dto.PaginationResponse{
				TotalData:   totalTransactions,
				TotalPage:   (totalTransactions + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
				CurrentPage: req.Pagination.Page,
			},
			Transactions: []dto.TransactionSellerResDTO{},
		}, nil
	}

	var transactionResDTOs []dto.TransactionSellerResDTO
	for _, transaction := range transactions {
		transactionResDTOs = append(transactionResDTOs, dto.TransactionSellerResDTO{
			InvoiceCode: transaction.InvoiceCode,
		})

		var cartItems []entity.TransactionCartItem
		err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONCartItem
		}
		cartItemsRes := dto.TransactionListCartResDTO{
			Product:      cartItems[0],
			TotalProduct: len(cartItems),
		}
		transactionResDTOs[len(transactionResDTOs)-1].ProductOverview = cartItemsRes

		var paymentDetails entity.TransactionPaymentDetails
		err = json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONPaymentDetails
		}
		subtotal := paymentDetails.Subtotal
		delivery_fee := paymentDetails.DeliveryFee
		merchant_voucher_nominal := paymentDetails.MerchantVoucherNominal
		transactionResDTOs[len(transactionResDTOs)-1].Total = subtotal + delivery_fee - merchant_voucher_nominal

		transactionResDTOs[len(transactionResDTOs)-1].Username = transaction.User.Username

		transactionResDTOs[len(transactionResDTOs)-1].TransactionDate = transaction.CreatedAt

		transactionStatus, err := u.transactionStatusUsecase.GetTransactionStatusByTransactionID(transaction.ID)
		if err != nil {
			return nil, err
		}

		transactionResDTOs[len(transactionResDTOs)-1].TransactionStatus = *transactionStatus

		var addressDetails entity.TransactionAddress
		err = json.Unmarshal([]byte(transaction.Address.Bytes), &addressDetails)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONTransactionAddress
		}
		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.Address = addressDetails

		var deliveryOption entity.TransactionDeliveryOption
		err = json.Unmarshal([]byte(transaction.DeliveryOption.Bytes), &deliveryOption)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONDeliveryOption
		}
		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.DeliveryOption.CourierName = deliveryOption.CourierName

		transactionDeliveryStatus, deliveryReceiptNumber, err := u.transactionDeliveryStatusUsecase.GetTransactionDeliveryStatusByTransactionID(transaction.ID)
		if err != nil {
			return nil, err
		}

		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.TransactionDeliveryStatus = *transactionDeliveryStatus
		if deliveryReceiptNumber != nil {
			transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.DeliveryOption.ReceiptNumber = *deliveryReceiptNumber
		}
	}

	return &dto.TransactionSellerListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalTransactions,
			TotalPage:   (totalTransactions + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
			CurrentPage: req.Pagination.Page,
		},
		Transactions: transactionResDTOs,
	}, nil
}

func (u *transactionUsecaseImpl) GetTransactionList(username string, req dto.TransactionReqParamDTO) (*dto.TransactionListResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	transactions, totalTransactions, err := u.transactionRepository.GetTransactionListByUserId(user.ID, req)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return &dto.TransactionListResDTO{
			PaginationResponse: dto.PaginationResponse{
				TotalData:   totalTransactions,
				TotalPage:   (totalTransactions + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
				CurrentPage: req.Pagination.Page,
			},
			Transactions: []dto.TransactionResDTO{},
		}, nil
	}

	var transactionResDTOs []dto.TransactionResDTO
	for _, transaction := range transactions {
		transactionResDTOs = append(transactionResDTOs, dto.TransactionResDTO{
			InvoiceCode: transaction.InvoiceCode,
		})

		var cartItems []entity.TransactionCartItem
		err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONCartItem
		}
		cartItemsRes := dto.TransactionListCartResDTO{
			Product:      cartItems[0],
			TotalProduct: len(cartItems),
		}
		transactionResDTOs[len(transactionResDTOs)-1].ProductOverview = cartItemsRes

		var paymentDetails entity.TransactionPaymentDetails
		err = json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONPaymentDetails
		}
		subtotal := paymentDetails.Subtotal
		delivery_fee := paymentDetails.DeliveryFee
		merchant_voucher_nominal := paymentDetails.MerchantVoucherNominal
		marketplace_voucher_nominal := paymentDetails.MarketplaceVoucherNominal
		transactionResDTOs[len(transactionResDTOs)-1].Total = subtotal + delivery_fee - merchant_voucher_nominal - marketplace_voucher_nominal

		merchant, err := u.merchantRepository.GetByDomain(transaction.MerchantDomain)
		if err != nil {
			return nil, err
		}

		transactionResDTOs[len(transactionResDTOs)-1].Merchant = dto.TransactionMerchantResDTO{
			Domain: merchant.Domain,
			Name:   merchant.Name,
		}

		transactionStatus, err := u.transactionStatusUsecase.GetTransactionStatusByTransactionID(transaction.ID)
		if err != nil {
			return nil, err
		}

		transactionResDTOs[len(transactionResDTOs)-1].TransactionStatus = *transactionStatus

		var addressDetails entity.TransactionAddress
		err = json.Unmarshal([]byte(transaction.Address.Bytes), &addressDetails)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONTransactionAddress
		}
		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.Address = addressDetails

		var deliveryOption entity.TransactionDeliveryOption
		err = json.Unmarshal([]byte(transaction.DeliveryOption.Bytes), &deliveryOption)
		if err != nil {
			return nil, domain.ErrUnmarshalJSONDeliveryOption
		}
		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.DeliveryOption.CourierName = deliveryOption.CourierName

		transactionDeliveryStatus, deliveryReceiptNumber, err := u.transactionDeliveryStatusUsecase.GetTransactionDeliveryStatusByTransactionID(transaction.ID)
		if err != nil {
			return nil, err
		}

		transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.TransactionDeliveryStatus = *transactionDeliveryStatus
		if deliveryReceiptNumber != nil {
			transactionResDTOs[len(transactionResDTOs)-1].ShippingDetails.DeliveryOption.ReceiptNumber = *deliveryReceiptNumber
		}
	}

	return &dto.TransactionListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   totalTransactions,
			TotalPage:   (totalTransactions + int64(req.Pagination.Limit) - 1) / int64(req.Pagination.Limit),
			CurrentPage: req.Pagination.Page,
		},
		Transactions: transactionResDTOs,
	}, nil
}

func (u *transactionUsecaseImpl) GetSellerTransactionDetail(username string, invoiceCode string) (*dto.TransactionSellerDetailResDTO, error) {
	merchant, err := u.merchantRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	transaction, err := u.transactionRepository.GetMerchantTransactionDetailByInvoiceCode(merchant.Domain, invoiceCode)
	if err != nil {
		return nil, err
	}

	transactionResDTO := dto.TransactionSellerDetailResDTO{
		InvoiceCode: transaction.InvoiceCode,
	}

	var cartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItem
	}
	transactionResDTO.ProductDetails.CartItems = cartItems

	var paymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	subtotal := paymentDetails.Subtotal
	delivery_fee := paymentDetails.DeliveryFee
	merchant_voucher_nominal := paymentDetails.MerchantVoucherNominal
	transactionResDTO.PaymentDetails.PaymentDetails = paymentDetails
	transactionResDTO.PaymentDetails.PaymentDetails.Total = subtotal + delivery_fee - merchant_voucher_nominal

	var addressDetails entity.TransactionAddress
	err = json.Unmarshal([]byte(transaction.Address.Bytes), &addressDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONTransactionAddress
	}
	transactionResDTO.ShippingDetails.Address = addressDetails

	var paymentMethod entity.TransactionPaymentMethod
	err = json.Unmarshal([]byte(transaction.PaymentMethod.Bytes), &paymentMethod)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentMethod
	}
	transactionResDTO.PaymentDetails.PaymentMethod = paymentMethod

	transactionResDTO.ProductDetails.User.Username = transaction.User.Username
	if transaction.User.UserDetail.ProfilePicture != nil {
		transactionResDTO.ProductDetails.User.ProfilePicture = *transaction.User.UserDetail.ProfilePicture
	}

	transactionStatus, err := u.transactionStatusUsecase.GetTransactionStatusByTransactionID(transaction.ID)
	if err != nil {
		return nil, err
	}

	transactionResDTO.TransactionStatus = *transactionStatus

	var deliveryOption entity.TransactionDeliveryOption
	err = json.Unmarshal([]byte(transaction.DeliveryOption.Bytes), &deliveryOption)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONDeliveryOption
	}
	transactionResDTO.ShippingDetails.DeliveryOption.CourierName = deliveryOption.CourierName

	transactionDeliveryStatus, deliveryReceiptNumber, err := u.transactionDeliveryStatusUsecase.GetTransactionDeliveryStatusByTransactionID(transaction.ID)
	if err != nil {
		return nil, err
	}

	transactionResDTO.ShippingDetails.TransactionDeliveryStatus = *transactionDeliveryStatus
	if deliveryReceiptNumber != nil {
		transactionResDTO.ShippingDetails.DeliveryOption.ReceiptNumber = *deliveryReceiptNumber
	}

	return &transactionResDTO, nil
}

func (u *transactionUsecaseImpl) GetTransactionDetail(username string, invoiceCode string) (*dto.TransactionDetailResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	transaction, err := u.transactionRepository.GetTransactionDetailByInvoiceCode(user.ID, invoiceCode)
	if err != nil {
		return nil, err
	}

	transactionResDTO := dto.TransactionDetailResDTO{
		InvoiceCode: transaction.InvoiceCode,
	}

	var cartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONCartItem
	}
	transactionResDTO.ProductDetails.CartItems = cartItems

	var paymentDetails entity.TransactionPaymentDetails
	err = json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}
	subtotal := paymentDetails.Subtotal
	delivery_fee := paymentDetails.DeliveryFee
	merchant_voucher_nominal := paymentDetails.MerchantVoucherNominal
	marketplace_voucher_nominal := paymentDetails.MarketplaceVoucherNominal
	transactionResDTO.PaymentDetails.PaymentDetails = paymentDetails
	transactionResDTO.PaymentDetails.PaymentDetails.Total = subtotal + delivery_fee - merchant_voucher_nominal - marketplace_voucher_nominal

	var addressDetails entity.TransactionAddress
	err = json.Unmarshal([]byte(transaction.Address.Bytes), &addressDetails)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONTransactionAddress
	}
	transactionResDTO.ShippingDetails.Address = addressDetails

	var paymentMethod entity.TransactionPaymentMethod
	err = json.Unmarshal([]byte(transaction.PaymentMethod.Bytes), &paymentMethod)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONPaymentMethod
	}
	transactionResDTO.PaymentDetails.PaymentMethod = paymentMethod

	merchant, err := u.merchantRepository.GetByDomain(transaction.MerchantDomain)
	if err != nil {
		return nil, err
	}

	transactionResDTO.ProductDetails.Merchant = dto.TransactionMerchantResDTO{
		Domain: merchant.Domain,
		Name:   merchant.Name,
	}

	transactionStatus, err := u.transactionStatusUsecase.GetTransactionStatusByTransactionID(transaction.ID)
	if err != nil {
		return nil, err
	}

	transactionResDTO.TransactionStatus = *transactionStatus

	var deliveryOption entity.TransactionDeliveryOption
	err = json.Unmarshal([]byte(transaction.DeliveryOption.Bytes), &deliveryOption)
	if err != nil {
		return nil, domain.ErrUnmarshalJSONDeliveryOption
	}
	transactionResDTO.ShippingDetails.DeliveryOption.CourierName = deliveryOption.CourierName

	transactionDeliveryStatus, deliveryReceiptNumber, err := u.transactionDeliveryStatusUsecase.GetTransactionDeliveryStatusByTransactionID(transaction.ID)
	if err != nil {
		return nil, err
	}

	transactionResDTO.ShippingDetails.TransactionDeliveryStatus = *transactionDeliveryStatus
	if deliveryReceiptNumber != nil {
		transactionResDTO.ShippingDetails.DeliveryOption.ReceiptNumber = *deliveryReceiptNumber
	}

	return &transactionResDTO, nil
}

func (u *transactionUsecaseImpl) MakeTransaction(username string, req dto.MakeTransactionReqDTO) (*dto.MakeTransactionResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	//get paymentId from paymentCode
	redirectUrl, paymentId, payRecId, paymentMethod, err := u.getPaymentDetails(req.PaymentMethodCode, req.PaymentAccountNumber, req.PaymentTotal, req.OrderCode)
	if err != nil {
		return nil, err
	}

	//check order summary, is the promo is available and the total price is correct
	orderValidated, err := u.validateOrderRequest(username, req)
	if err != nil {
		return nil, err
	}

	if orderValidated.SubTotal < 0 || orderValidated.Total < 0 {
		return nil, domain.ErrPaymentTotalNotMatch
	}

	//make payment record and the transaction record
	paymentRec := entity.PaymentRecord{
		ID:              payRecId,
		PaymentId:       paymentId,
		PaymentMethodId: paymentMethod.ID,
		Amount:          orderValidated.Total,
		PaymentUrl:      redirectUrl,
	}
	transactionRecords, err := u.makeTransactionsEntity(user.ID, paymentRec, *orderValidated, *paymentMethod, req.PaymentAccountNumber)
	if err != nil {
		return nil, err
	}

	err = u.transactionRepository.MakeTransaction(transactionRecords, *orderValidated, paymentId)
	if err != nil {
		log.Error().Msgf("error: in usecase make transaction %v", err)
		return nil, err
	}

	return &dto.MakeTransactionResDTO{
		PaymentId:          paymentId,
		OrderCode:          orderValidated.OrderCode,
		Amount:             orderValidated.Total,
		PaymentRedirectUrl: redirectUrl,
	}, nil
}

func (u *transactionUsecaseImpl) HandleTransactionSlpRes(trxRecords []entity.TransactionPaymentRecord, isSuccess bool, paymentRec entity.PaymentRecord) error {
	transactions := u.getTransactionsFromTransactionPaymentRecords(trxRecords)

	if len(transactions) <= 0 {
		return domain.ErrTransactionNotFound
	}

	if isSuccess {
		err := u.transactionRepository.UpdateTransactionPaymentSuccess(transactions, paymentRec)
		if err != nil {
			return err
		}
	} else {
		trxCartItems, err := u.getTransactionCartItemsFromTransactions(transactions)
		if err != nil {
			return err
		}

		err = u.transactionRepository.UpdateTransactionPaymentFailed(transactions, trxCartItems, paymentRec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *transactionUsecaseImpl) UpdateUserTransactionStatus(username string, req dto.UpdateUserTransactionStatusReqDTO) (*dto.UpdateUserTransactionStatusResDTO, error) {
	if (req.Status != dto.TransactionStatusCompleted) ==
		(req.Status != dto.TransactionStatusCanceled) ==
		(req.Status != dto.TransactionStatusRequestRefund) {
		return nil, domain.ErrInvalidTransactionStatusForbidden
	}

	validatedTransaction, err := u.validateUpdateUserTransactionStatus(username, req)
	if err != nil {
		return nil, err
	}

	//update transaction status
	updatedTrxStatus, err := u.updateUserTransactionStatusImpl(*validatedTransaction.TransactionStatus, *validatedTransaction, req)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateUserTransactionStatusResDTO{
		InvoiceCode: req.InvoiceCode,
		TransactionStatus: dto.TransactionStatusResDTO{
			OnWaitedAt:    updatedTrxStatus.OnWaitedAt,
			OnProcessedAt: updatedTrxStatus.OnProcessedAt,
			OnDeliveredAt: updatedTrxStatus.OnDeliveredAt,
			OnCompletedAt: updatedTrxStatus.OnCompletedAt,
			OnCanceledAt:  updatedTrxStatus.OnCanceledAt,
			OnRefundedAt:  updatedTrxStatus.OnRefundedAt,
		},
		UpdatedAt: updatedTrxStatus.UpdatedAt,
	}, nil
}

func (u *transactionUsecaseImpl) updateUserTransactionStatusImpl(trxStatus entity.TransactionStatus, transaction entity.Transaction, req dto.UpdateUserTransactionStatusReqDTO) (*entity.TransactionStatus, error) {
	updatedStatus, err := u.changeUserTransactionStatus(trxStatus, req)
	if err != nil {
		return nil, err
	}

	transaction.TransactionStatus = updatedStatus

	if req.Status == dto.TransactionStatusCanceled {
		updatedTrxStatus, err := u.updateTransactionCanceledProcess(transaction)
		if err != nil {
			return nil, err
		}

		return updatedTrxStatus, nil
	}

	updatedTrxStatus, err := u.updateTransactionCompletedProcess(transaction)
	if err != nil {
		return nil, err
	}

	return updatedTrxStatus, nil
}

func (u *transactionUsecaseImpl) updateTransactionCanceledProcess(transaction entity.Transaction) (*entity.TransactionStatus, error) {
	//count amount payment
	var paymentDetails entity.TransactionPaymentDetails
	err := json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
	if err != nil {
		log.Error().Msgf("error: in unmarshal transaction payment details items %v", err)
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}

	var trxCartItems []entity.TransactionCartItem
	err = json.Unmarshal([]byte(transaction.CartItems.Bytes), &trxCartItems)
	if err != nil {
		log.Error().Msgf("error: in unmarshal transaction cart items %v", err)
		return nil, domain.ErrUnmarshalJSONCartItems
	}

	amountPayment := paymentDetails.Subtotal + paymentDetails.DeliveryFee - paymentDetails.MarketplaceVoucherNominal - paymentDetails.MerchantVoucherNominal
	updatedStatus, err := u.transactionRepository.UpdateTransactionStatusCanceled(transaction, amountPayment, trxCartItems)
	if err != nil {
		return nil, err
	}

	return updatedStatus, nil
}

func (u *transactionUsecaseImpl) updateTransactionCompletedProcess(transaction entity.Transaction) (*entity.TransactionStatus, error) {
	//count amount payment
	var paymentDetails entity.TransactionPaymentDetails
	err := json.Unmarshal([]byte(transaction.PaymentDetails.Bytes), &paymentDetails)
	if err != nil {
		log.Error().Msgf("error: in unmarshal transaction payment details items %v", err)
		return nil, domain.ErrUnmarshalJSONPaymentDetails
	}

	amountPayment := paymentDetails.Subtotal + paymentDetails.DeliveryFee - paymentDetails.MerchantVoucherNominal
	updatedStatus, err := u.transactionRepository.UpdateTransactionStatusCompleted(transaction, amountPayment, paymentDetails.MarketplaceVoucherNominal)
	if err != nil {
		return nil, err
	}

	return updatedStatus, nil
}

func (u *transactionUsecaseImpl) getTransactionCartItemsFromTransactions(transactions []entity.Transaction) ([]entity.TransactionCartItem, error) {
	var transactionCartItems []entity.TransactionCartItem
	for _, transaction := range transactions {
		var cartItems []entity.TransactionCartItem
		err := json.Unmarshal([]byte(transaction.CartItems.Bytes), &cartItems)
		if err != nil {
			log.Error().Msgf("error: in unmarshal transaction cart items %v", err)
			return nil, domain.ErrUnmarshalJSONCartItem
		}

		transactionCartItems = append(transactionCartItems, cartItems...)
	}

	return transactionCartItems, nil
}

func (u *transactionUsecaseImpl) getTransactionsFromTransactionPaymentRecords(trxRecords []entity.TransactionPaymentRecord) []entity.Transaction {
	var transactions []entity.Transaction
	for _, trxRecord := range trxRecords {
		transactions = append(transactions, trxRecord.Transaction)
	}

	return transactions
}

func (u *transactionUsecaseImpl) getPaymentDetails(paymentMethodCode string, paymentAccNumber string, amount float64, orderCode string) (payUrl string, paymentId string, paymentRecId uint, payMethod *entity.PaymentMethod, payErr error) {
	//get payment method id from paymentCode
	paymentMethod, err := u.paymentMethodRepository.GetByCode(paymentMethodCode)
	if err != nil {
		return "", "", 0, nil, err
	}

	if paymentMethod.ID == dto.PAYMENT_METHOD_ID_SLP {
		redirectTrxUrl := fmt.Sprintf("%s/%s/payment-status", config.Config.WebTransactionURL, orderCode)
		redirectUrl, paymentId, err := u.sealabspayRepository.MakePaymentCustomRedirect(paymentAccNumber, uint(amount), redirectTrxUrl)
		if err != nil {
			return "", "", 0, nil, err
		}

		return redirectUrl, paymentId, 0, paymentMethod, nil
	}

	if paymentMethod.ID == dto.PAYMENT_METHOD_ID_WALLET {
		redirectUrl, paymentId, paymentRecordId, err := u.walletRepository.MakePayment(paymentAccNumber, uint(amount))
		if err != nil {
			return "", "", 0, nil, err
		}

		return redirectUrl, paymentId, paymentRecordId, paymentMethod, nil
	}

	return "", "", 0, nil, domain.ErrCreateTransactionPayment
}

func (u *transactionUsecaseImpl) validateOrderRequest(username string, req dto.MakeTransactionReqDTO) (*dto.PostOrderSummaryResDTO, error) {
	orderSummary, err := u.orderItemUsecase.GetOrderCheckoutSummary(username, dto.PostOrderSummaryReqDTO{
		OrderCode:          req.OrderCode,
		AddressId:          req.AddressId,
		Merchants:          req.Merchants,
		VoucherMarketplace: req.VoucherMarketplace,
	})
	if err != nil {
		return nil, err
	}

	//check orderValidated total with req total
	if orderSummary.Total != req.PaymentTotal ||
		!orderSummary.IsOrderValid ||
		!orderSummary.IsVouchervalid ||
		!orderSummary.IsOrderEligible {
		return nil, domain.ErrMakeTransaction
	}

	return orderSummary, nil
}

func (u *transactionUsecaseImpl) deleteCartItems(trxCartItems []entity.TransactionCartItem, userId uint) error {
	var produVarIds []uint

	for _, trxCartItem := range trxCartItems {
		produVarIds = append(produVarIds, trxCartItem.ProductVariantId)
	}

	err := u.cartItemRepository.DeleteCartItemsByProductVariantIds(userId, produVarIds)

	return err
}

func (u *transactionUsecaseImpl) makeTransactionsEntity(userId uint, paymentRecord entity.PaymentRecord, orderSummary dto.PostOrderSummaryResDTO, paymentMethod entity.PaymentMethod, paymentAccRelated string) ([]entity.Transaction, error) {

	var transactions []entity.Transaction
	numOfTrx := len(orderSummary.Orders)
	for _, orderItem := range orderSummary.Orders {
		newInvoiceCode := u.generateInvoiceCode(userId, orderItem.Merchant.MerchantId, paymentRecord.PaymentId)
		for i := 0; i < 11; i++ {
			if i == 10 {
				return nil, domain.ErrCreateTransaction
			}

			_, err := u.transactionRepository.GetTransactionByInvoiceCode(newInvoiceCode)
			if err != nil {
				break
			}
			newInvoiceCode = u.generateInvoiceCode(userId, orderItem.Merchant.MerchantId, paymentRecord.PaymentId)
		}
		transaction := entity.Transaction{
			MerchantDomain:            orderItem.Merchant.MerchantDomain,
			UserId:                    userId,
			InvoiceCode:               newInvoiceCode,
			MerchantVoucherId:         orderItem.MerchantVoucherId,
			MarketplaceVoucherId:      orderSummary.MarketplaceVoucherId,
			TransactionStatus:         &entity.TransactionStatus{},
			TransactionDeliveryStatus: &entity.TransactionDeliveryStatus{},
			PaymentRecords:            []entity.PaymentRecord{paymentRecord},
		}

		transaction.PaymentMethod.Set(entity.TransactionPaymentMethod{
			ID:                   paymentMethod.ID,
			Name:                 paymentMethod.Name,
			Code:                 paymentMethod.Code,
			AccountRelatedNumber: paymentAccRelated,
		})
		transaction.PaymentDetails.Set(entity.TransactionPaymentDetails{
			Subtotal:                  orderItem.SubTotal,
			DeliveryFee:               orderItem.DeliveryCost,
			MarketplaceVoucherNominal: orderSummary.DiscountMarketplace / float64(numOfTrx),
			MerchantVoucherNominal:    orderItem.Discount,
		})
		transaction.Address.Set(entity.TransactionAddress{
			Name:            orderSummary.Address.Name,
			Phone:           orderSummary.Address.PhoneNumber,
			Label:           orderSummary.Address.Label,
			Details:         orderSummary.Address.Details,
			CityName:        orderSummary.Address.City.Name,
			DistrictName:    orderSummary.Address.District.Name,
			SubdistrictName: orderSummary.Address.Subdistrict.Name,
			ProvinceName:    orderSummary.Address.Province.Name,
			ZipCode:         orderSummary.Address.Subdistrict.ZipCode,
		})
		transaction.DeliveryOption.Set(entity.TransactionDeliveryOption{
			CourierName: orderItem.DeliveryService.Name,
		})
		trxCartItems := u.makeTransactionCartItemEntity(orderItem.Items)
		transaction.CartItems.Set(trxCartItems)

		u.deleteCartItems(trxCartItems, userId)

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (u *transactionUsecaseImpl) makeTransactionCartItemEntity(items []dto.OrderItemDTO) []entity.TransactionCartItem {
	var cartItems []entity.TransactionCartItem
	for _, item := range items {
		cartItem := entity.TransactionCartItem{
			ProductId:        item.ProductId,
			ProductVariantId: *item.VariantItemId,
			Name:             item.Name,
			Image:            item.Image,
			RealPrice:        item.RealPrice,
			DiscountPrice:    item.DiscountPrice,
			Notes:            item.Notes,
			ProductSlug:      item.ProductSlug,
			VariantName:      item.VariantName,
			Quantity:         item.Quantity,
		}

		cartItems = append(cartItems, cartItem)
	}
	return cartItems
}

func (u *transactionUsecaseImpl) generateInvoiceCode(userId uint, merchantId uint, paymentId string) string {
	invCodeDate := fmt.Sprintf("INV%s", time.Now().Format("01022006"))
	invCodeNumber := fmt.Sprintf("%d%d%s%d", 123, 123123, "WAL828282", util.GenerateRandomNumber(9999, 1000))
	invCodeNumberHashed := fmt.Sprintf("%.10d%d", util.HashFnv(invCodeNumber), util.GenerateRandomNumber(9999, 1000))
	invCode := fmt.Sprintf("%s-%s", invCodeDate, invCodeNumberHashed)

	return invCode
}

func (u *transactionUsecaseImpl) checkUpdateStatusIsValid(transaction entity.Transaction, req dto.UpdateMerchantTransactionStatusReqDTO) error {
	currect_status := u.parseTransactionStatusToId(*transaction.TransactionStatus, *transaction.TransactionDeliveryStatus)

	if currect_status >= req.Status || currect_status == dto.TransactionStatusCanceled {
		return domain.ErrUpdateTransactionStatusCannotReverse
	}

	stepDiff := req.Status - currect_status
	if stepDiff > 2 {
		return domain.ErrUpdateTransactionStatusCannotSkip
	}

	if stepDiff == 2 &&
		(req.Status != dto.TransactionStatusOnDelivery &&
			req.Status != dto.TransactionStatusCompleted &&
			req.Status != dto.TransactionStatusCanceled &&
			req.Status != dto.TransactionStatusRefunded) {
		return domain.ErrUpdateTransactionStatusCannotSkip
	}

	return nil
}

func (u *transactionUsecaseImpl) validateUpdateUserTransactionStatus(username string, req dto.UpdateUserTransactionStatusReqDTO) (*entity.Transaction, error) {
	//get userId from username
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	//get transaction by invoiceCode
	transaction, err := u.transactionRepository.GetTransactionDetailByInvoiceCode(user.ID, req.InvoiceCode)
	if err != nil {
		return nil, err
	}

	// if transaction already processed, user cannot cancel it
	if transaction.TransactionStatus.OnProcessedAt != nil &&
		req.Status == dto.TransactionStatusCanceled {
		return nil, domain.ErrUpdateTransactionStatusCannotReverse
	}

	//check if status is valid
	err = u.checkUpdateStatusIsValid(*transaction, dto.UpdateMerchantTransactionStatusReqDTO{
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *transactionUsecaseImpl) parseTransactionStatusToId(statusTrx entity.TransactionStatus, statusDelivery entity.TransactionDeliveryStatus) int {
	switch {
	case statusTrx.OnRefundedAt != nil:
		return dto.TransactionStatusRefunded
	case statusTrx.OnCompletedAt != nil:
		return dto.TransactionStatusCompleted
	case statusTrx.OnRequestRefundAt != nil:
		return dto.TransactionStatusRequestRefund
	case statusTrx.OnCanceledAt != nil:
		return dto.TransactionStatusCanceled
	case statusTrx.OnDeliveredAt != nil || statusDelivery.OnDeliveredAt != nil:
		return dto.TransactionStatusDelivered
	case statusDelivery.OnDeliveryAt != nil:
		return dto.TransactionStatusOnDelivery
	case statusTrx.OnProcessedAt != nil:
		return dto.TransactionStatusProcessed
	default:
		return dto.TransactionStatusWaited
	}
}

func (u *transactionUsecaseImpl) changeUserTransactionStatus(transactionStatus entity.TransactionStatus, req dto.UpdateUserTransactionStatusReqDTO) (*entity.TransactionStatus, error) {
	timeNow := time.Now()
	switch req.Status {
	case dto.TransactionStatusCompleted:
		transactionStatus.OnCompletedAt = &timeNow
	case dto.TransactionStatusCanceled:
		transactionStatus.OnCanceledAt = &timeNow
	case dto.TransactionStatusRequestRefund:
		transactionStatus.OnRequestRefundAt = &timeNow
	default:
		return nil, domain.ErrUpdateTransactionStatusCannotSkip
	}

	return &transactionStatus, nil
}
