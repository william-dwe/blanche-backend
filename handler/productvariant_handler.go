package handler

import (
	"fmt"
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProductVariants(c *gin.Context) {
	domain := c.Param("domain")
	slug := c.Param("slug")

	res, err := h.productVariantUsecase.GetByProductSlug(fmt.Sprintf("%s/%s", domain, slug))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_VARIANTS",
		Message: "Success retrieve product variants",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetMerchantProductVariants(c *gin.Context) {
	productId := c.Param("product_id")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		_ = c.Error(domain.ErrProductIdNotValid)
		return
	}

	res, err := h.productVariantUsecase.GetByProductId(uint(productIdInt))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_PRODUCT_VARIANTS",
		Message: "Success retrieve product variants",
		Data:    res,
	}

	util.ResponseSuccessJSON(c, response)
}
