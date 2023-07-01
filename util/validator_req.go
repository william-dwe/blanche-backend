package util

import (
	"fmt"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ShouldBindWithValidation(c *gin.Context, dto any) error {
	if err := c.ShouldBind(dto); err != nil {
		errTag, cvtErr := err.(validator.ValidationErrors)
		if !cvtErr {
			return domain.ErrRequestBodyJSONInvalid
		}
		if errTag != nil {
			return httperror.BadRequestError(fmt.Sprintf("check the requested param: %s must be %s %s",
				ToSnakeCase(errTag[0].StructField()),
				errTag[0].ActualTag(),
				errTag[0].Param()), "DATA_NOT_VALID")
		}
	}

	return nil
}
