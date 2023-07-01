package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type AddressUsecase interface {
	GetProvinces() (*dto.AllProvincesResDTO, error)
	GetCities() (*dto.AllCitiesResDTO, error)
	GetCitiesByProvinceID(provinceID uint) (*dto.CitiesResDTO, error)
	GetDistrictsByCityID(cityID uint) (*dto.AllDistrictsResDTO, error)
	GetSubDistrictsByDistrictID(districtID uint) (*dto.AllSubDistrictsResDTO, error)
}

type AddressUsecaseConfig struct {
	AddressRepository repository.AddressRepository
}

type addressUsecaseImpl struct {
	addressRepository repository.AddressRepository
}

func NewAddressUsecase(c AddressUsecaseConfig) AddressUsecase {
	return &addressUsecaseImpl{
		addressRepository: c.AddressRepository,
	}
}

func (u *addressUsecaseImpl) GetProvinces() (*dto.AllProvincesResDTO, error) {
	provinces, err := u.addressRepository.GetProvinces()
	if err != nil {
		return nil, err
	}

	provincesDTO := dto.AllProvincesResDTO{
		Provinces: make([]dto.ProvinceResDTO, 0),
	}

	for _, province := range provinces {
		provinceDTO := dto.ProvinceResDTO{
			ID:   province.ID,
			Name: province.Name,
		}
		provincesDTO.Provinces = append(provincesDTO.Provinces, provinceDTO)
	}

	return &provincesDTO, nil
}

func (u *addressUsecaseImpl) GetCities() (*dto.AllCitiesResDTO, error) {
	cities, err := u.addressRepository.GetCities()
	if err != nil {
		return nil, err
	}

	citiesDTO := dto.AllCitiesResDTO{
		Cities: make([]dto.CityResDTO, 0),
	}

	for _, city := range cities {
		cityDTO := dto.CityResDTO{
			ID:   city.ID,
			Name: city.Name,
			RoId: city.RoId,
		}
		citiesDTO.Cities = append(citiesDTO.Cities, cityDTO)
	}

	return &citiesDTO, nil
}

func (u *addressUsecaseImpl) GetCitiesByProvinceID(provinceID uint) (*dto.CitiesResDTO, error) {
	province, err := u.addressRepository.GetProvinceById(provinceID)
	if err != nil {
		return nil, err
	}

	cities, err := u.addressRepository.GetCitiesByProvinceId(provinceID)
	if err != nil {
		return nil, err
	}

	citiesDTO := dto.CitiesResDTO{
		Province: dto.ProvinceResDTO{
			ID:   province.ID,
			Name: province.Name,
		},
		Cities: make([]dto.CityResDTO, 0),
	}

	for _, city := range cities {
		cityDTO := dto.CityResDTO{
			ID:   city.ID,
			Name: city.Name,
			RoId: city.RoId,
		}
		citiesDTO.Cities = append(citiesDTO.Cities, cityDTO)
	}

	return &citiesDTO, nil
}

func (u *addressUsecaseImpl) GetDistrictsByCityID(cityID uint) (*dto.AllDistrictsResDTO, error) {
	city, err := u.addressRepository.GetCityById(cityID)
	if err != nil {
		return nil, err
	}

	districts, err := u.addressRepository.GetDistrictsByCityId(cityID)
	if err != nil {
		return nil, err
	}

	districtsDTO := dto.AllDistrictsResDTO{
		City: dto.CityResDTO{
			ID:   city.ID,
			Name: city.Name,
			RoId: city.RoId,
		},
		Districts: make([]dto.DistrictResDTO, 0),
	}

	for _, district := range districts {
		districtDTO := dto.DistrictResDTO{
			ID:   district.ID,
			Name: district.Name,
		}
		districtsDTO.Districts = append(districtsDTO.Districts, districtDTO)
	}

	return &districtsDTO, nil
}

func (u *addressUsecaseImpl) GetSubDistrictsByDistrictID(districtID uint) (*dto.AllSubDistrictsResDTO, error) {
	district, err := u.addressRepository.GetDistrictById(districtID)
	if err != nil {
		return nil, err
	}

	subDistricts, err := u.addressRepository.GetSubDistrictsByDistrictId(districtID)
	if err != nil {
		return nil, err
	}

	subDistrictsDTO := dto.AllSubDistrictsResDTO{
		District: dto.DistrictResDTO{
			ID:   district.ID,
			Name: district.Name,
		},
		SubDistricts: make([]dto.SubDistrictResDTO, 0),
	}

	for _, subDistrict := range subDistricts {
		subDistrictDTO := dto.SubDistrictResDTO{
			ID:      subDistrict.ID,
			Name:    subDistrict.Name,
			ZipCode: subDistrict.ZipCode,
		}
		subDistrictsDTO.SubDistricts = append(subDistrictsDTO.SubDistricts, subDistrictDTO)
	}

	return &subDistrictsDTO, nil
}
