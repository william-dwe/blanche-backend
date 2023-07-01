package handler

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTransactionList(c *gin.Context) {
	var transactionRequest dto.TransactionReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transactions, err := h.transactionUsecase.GetTransactionList(user.Username, transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_TRANSACTIONS",
		Message: "Success retrieve transactions",
		Data:    transactions,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetSellerTransactionList(c *gin.Context) {
	var transactionRequest dto.TransactionReqParamDTO
	err := util.ShouldBindQueryWithValidation(c, &transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transactions, err := h.transactionUsecase.GetSellerTransactionList(user.Username, transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_TRANSACTIONS",
		Message: "Success retrieve seller's transactions",
		Data:    transactions,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetTransactionDetail(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionUsecase.GetTransactionDetail(user.Username, invoiceCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_TRANSACTION",
		Message: "Success retrieve transaction",
		Data:    transaction,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetSellerTransactionDetail(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionUsecase.GetSellerTransactionDetail(user.Username, invoiceCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_TRANSACTION",
		Message: "Success retrieve transaction",
		Data:    transaction,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) MakeTransaction(c *gin.Context) {
	var transactionRequest dto.MakeTransactionReqDTO
	err := util.ShouldBindJsonWithValidation(c, &transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionUsecase.MakeTransaction(user.Username, transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_MAKE_TRANSACTION",
		Message: "Success make transaction",
		Data:    transaction,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateMerchantTransactionStatus(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	var transactionRequest dto.UpdateMerchantTransactionStatusReqDTO
	err := util.ShouldBindJsonWithValidation(c, &transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}
	transactionRequest.InvoiceCode = invoiceCode

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionUsecase.UpdateMerchantTransactionStatus(user.Username, transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_MERCHANT_TRANSACTION_STATUS",
		Message: "Success update merchant transaction status",
		Data:    transaction,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) UpdateUserTransactionStatus(c *gin.Context) {
	invoiceCode := c.Param("invoice_code")

	var transactionRequest dto.UpdateUserTransactionStatusReqDTO
	err := util.ShouldBindJsonWithValidation(c, &transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}
	transactionRequest.InvoiceCode = invoiceCode

	user, err := util.GetUserJWTContext(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := h.transactionUsecase.UpdateUserTransactionStatus(user.Username, transactionRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_UPDATE_USER_TRANSACTION_STATUS",
		Message: "Success update user transaction status",
		Data:    transaction,
	}

	util.ResponseSuccessJSON(c, response)
}
