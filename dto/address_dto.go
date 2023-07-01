package dto

type ProvinceResDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AllProvincesResDTO struct {
	Provinces []ProvinceResDTO `json:"provinces"`
}

type CityResDTO struct {
	ID   uint   `json:"id"`
	RoId uint   `json:"ro_id"`
	Name string `json:"name"`
}

type AllCitiesResDTO struct {
	Cities []CityResDTO `json:"cities"`
}

type CitiesResDTO struct {
	Province ProvinceResDTO `json:"province"`
	Cities   []CityResDTO   `json:"cities"`
}

type DistrictResDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AllDistrictsResDTO struct {
	City      CityResDTO       `json:"city"`
	Districts []DistrictResDTO `json:"districts"`
}

type SubDistrictResDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	ZipCode string `json:"zip_code"`
}

type AllSubDistrictsResDTO struct {
	District     DistrictResDTO      `json:"district"`
	SubDistricts []SubDistrictResDTO `json:"sub_districts"`
}
