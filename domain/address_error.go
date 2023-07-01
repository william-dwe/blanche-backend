package domain

import "git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/httperror"

var ErrGetProvinces = httperror.InternalServerError("failed to get provinces")
var ErrGetProvinceById = httperror.InternalServerError("failed to get province by id")
var ErrGetProvinceByIdNotFound = httperror.BadRequestError("province not found, cannot get province data", "PROVINCE_NOT_FOUND")

var ErrGetCitiesByProvinceId = httperror.InternalServerError("failed to get cities by province id")
var ErrGetCities = httperror.InternalServerError("failed to get cities")
var ErrGetCityByIdNotFound = httperror.BadRequestError("city not found, cannot get city data", "CITY_NOT_FOUND")
var ErrGetCityById = httperror.InternalServerError("failed to get city by id")
var ErrGetCitiesByProvinceDeficient = httperror.BadRequestError("check your request, province_id param are not appropriate", "NOT_VALID_PARAM")

var ErrGetSubDistrictsByDistrictId = httperror.InternalServerError("failed to get sub districts by district id")
var ErrGetSubDistrictsByDistrictIdDeficient = httperror.BadRequestError("check your request, district_id param are not appropriate", "NOT_VALID_PARAM")
var ErrGetSubDistrictByIdNotFound = httperror.BadRequestError("sub district not found, cannot get sub district data", "SUB_DISTRICT_NOT_FOUND")
var ErrGetSubDistrictById = httperror.InternalServerError("failed to get sub district by id")

var ErrGetDistrictsByCityId = httperror.InternalServerError("failed to get districts by city id")
var ErrGetDistrictsByCityIdDeficient = httperror.BadRequestError("check your request, city_id param are not appropriate", "NOT_VALID_PARAM")
var ErrGetDistrictByIdNotFound = httperror.BadRequestError("district not found, cannot get district data", "DISTRICT_NOT_FOUND")
var ErrGetDistrictById = httperror.InternalServerError("failed to get district by id")

var ErrAddMerchantAddressNotFound = httperror.BadRequestError("address not found, cannot add merchant address", "ADDRESS_NOT_FOUND")
