package handler

import (
	"strconv"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllProvinces(c *gin.Context) {
	provinces, err := h.addressUsecase.GetProvinces()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_PROVINCES",
		Message: "Success get all provinces",
		Data:    provinces,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetAllCities(c *gin.Context) {
	cities, err := h.addressUsecase.GetCities()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_ALL_CITIES",
		Message: "Success get all cities",
		Data:    cities,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetCitiesByProvinceID(c *gin.Context) {
	provinceId, err := strconv.Atoi(c.Param("provinceId"))
	if err != nil {
		_ = c.Error(domain.ErrGetCitiesByProvinceDeficient)
		return
	}

	cities, err := h.addressUsecase.GetCitiesByProvinceID(uint(provinceId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_CITIES_BY_PROVINCE_ID",
		Message: "Success get cities by province id",
		Data:    cities,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetDistrictsByCityID(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("cityId"))
	if err != nil {
		_ = c.Error(domain.ErrGetDistrictsByCityIdDeficient)
		return
	}

	districtDTO, err := h.addressUsecase.GetDistrictsByCityID(uint(cityId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_DISTRICTS_BY_CITY_ID",
		Message: "Success get districts by city id",
		Data:    districtDTO,
	}

	util.ResponseSuccessJSON(c, response)
}

func (h *Handler) GetSubDistrictsByDistrictID(c *gin.Context) {
	districtId, err := strconv.Atoi(c.Param("districtId"))
	if err != nil {
		_ = c.Error(domain.ErrGetSubDistrictsByDistrictIdDeficient)
		return
	}

	subDistrictDTO, err := h.addressUsecase.GetSubDistrictsByDistrictID(uint(districtId))
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := util.ResponseStruct{
		Code:    "SUCCESS_GET_SUB_DISTRICTS_BY_DISTRICT_ID",
		Message: "Success get sub districts by district id",
		Data:    subDistrictDTO,
	}

	util.ResponseSuccessJSON(c, response)
}
