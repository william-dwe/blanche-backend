package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type AddressRepository interface {
	GetProvinces() ([]entity.Province, error)
	GetProvinceById(provinceId uint) (*entity.Province, error)
	GetCities() ([]entity.City, error)
	GetCityById(cityId uint) (*entity.City, error)
	GetCitiesByProvinceId(provinceId uint) ([]entity.City, error)
	GetDistrictsByCityId(cityId uint) ([]entity.District, error)
	GetDistrictById(uint) (*entity.District, error)
	GetSubDistrictsByDistrictId(districtId uint) ([]entity.Subdistrict, error)
	GetSubDistrictById(subDistrictId uint) (*entity.Subdistrict, error)
}

type AddressRepositoryConfig struct {
	DB *gorm.DB
}

type addressRepositoryImpl struct {
	db *gorm.DB
}

func NewAddressRepository(c AddressRepositoryConfig) AddressRepository {
	return &addressRepositoryImpl{
		db: c.DB,
	}
}

func (r *addressRepositoryImpl) GetProvinces() ([]entity.Province, error) {
	var provinces []entity.Province
	res := r.db.Find(&provinces)
	if res.Error != nil {
		return nil, domain.ErrGetProvinces
	}

	return provinces, nil
}

func (r *addressRepositoryImpl) GetProvinceById(provinceId uint) (*entity.Province, error) {
	var province entity.Province
	res := r.db.Where("id = ?", provinceId).First(&province)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetProvinceByIdNotFound
		}
		return nil, domain.ErrGetProvinceById
	}

	return &province, nil
}

func (r *addressRepositoryImpl) GetCities() ([]entity.City, error) {
	var cities []entity.City
	res := r.db.Order("ro_id ASC").Find(&cities)
	if res.Error != nil {
		return nil, domain.ErrGetCities
	}

	return cities, nil
}

func (r *addressRepositoryImpl) GetCityById(cityId uint) (*entity.City, error) {
	var city entity.City
	res := r.db.Where("id = ?", cityId).Order("ro_id ASC").First(&city)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetCityByIdNotFound
		}
		return nil, domain.ErrGetCityById
	}

	return &city, nil
}

func (r *addressRepositoryImpl) GetCitiesByProvinceId(provinceId uint) ([]entity.City, error) {
	var cities []entity.City
	res := r.db.Where("province_id = ?", provinceId).Order("ro_id ASC").Find(&cities)
	if res.Error != nil {
		return nil, domain.ErrGetCitiesByProvinceId
	}

	return cities, nil
}

func (r *addressRepositoryImpl) GetDistrictsByCityId(cityId uint) ([]entity.District, error) {
	var districts []entity.District
	res := r.db.Where("city_id = ?", cityId).Find(&districts)
	if res.Error != nil {
		return nil, domain.ErrGetDistrictsByCityId
	}

	return districts, nil
}

func (r *addressRepositoryImpl) GetDistrictById(districtId uint) (*entity.District, error) {
	var district entity.District
	res := r.db.Where("id = ?", districtId).First(&district)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetDistrictByIdNotFound
		}
		return nil, domain.ErrGetDistrictById
	}

	return &district, nil
}

func (r *addressRepositoryImpl) GetSubDistrictsByDistrictId(districtId uint) ([]entity.Subdistrict, error) {
	var subDistricts []entity.Subdistrict
	res := r.db.Distinct("id", "name", "zip_code").Where("district_id = ?", districtId).Find(&subDistricts)
	if res.Error != nil {
		return nil, domain.ErrGetSubDistrictsByDistrictId
	}

	return subDistricts, nil
}

func (r *addressRepositoryImpl) GetSubDistrictById(subDistrictId uint) (*entity.Subdistrict, error) {
	var subDistrict entity.Subdistrict
	res := r.db.Where("id = ?", subDistrictId).First(&subDistrict)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrGetSubDistrictByIdNotFound
		}
		return nil, domain.ErrGetSubDistrictById
	}

	return &subDistrict, nil
}
