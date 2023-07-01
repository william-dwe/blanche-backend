package dto

import (
	"mime/multipart"
	"time"
)

const (
	TIME_LIMIT_BLOCK_WALLET = 30
)

type UserRegisterCheckEmailReqDTO struct {
	Email string `json:"email" binding:"required"`
}

type UserRegisterCheckEmailResDTO struct {
	Email       string `json:"email" binding:"required"`
	IsAvailable bool   `json:"is_available" binding:"required"`
}

type UserRegisterCheckUsernameReqDTO struct {
	Username string `json:"username" binding:"required"`
}

type UserRegisterCheckUsernameResDTO struct {
	Username    string `json:"username" binding:"required"`
	IsAvailable bool   `json:"is_available" binding:"required"`
}

type UserRegisterReqDTO struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegisterResDTO struct {
	AccessToken string `json:"access_token"`
}

type UserLoginReqDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResDTO struct {
	AccessToken string `json:"access_token"`
}

type UserRefreshResDTO struct {
	AccessToken string `json:"access_token"`
	Role        string `json:"-"`
}

type AccessTokenPayload struct {
	Username string `json:"username"`
}

type UserProfileResDTO struct {
	ID             uint       `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Fullname       string     `json:"fullname"`
	Phone          *string    `json:"phone"`
	Gender         *string    `json:"gender"`
	BirthDate      *time.Time `json:"birth_date"`
	ProfilePicture *string    `json:"profile_picture"`
	Role           string     `json:"role"`
}

type UserUpdateProfileReqDTO struct {
	Email string `json:"email"`
}

type UserUpdateProfileDetailFormReqDTO struct {
	Fullname       string                `form:"fullname,omitempty"`
	Phone          *string               `form:"phone,omitempty"`
	Gender         *string               `form:"gender,omitempty"`
	BirthDate      *time.Time            `form:"birth_date,omitempty"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture,omitempty"`
}

type UserAddressDTO struct {
	ID                uint   `json:"id"`
	Phone             string `json:"phone"`
	Name              string `json:"name"`
	Details           string `json:"details"`
	ProvinceName      string `json:"province_name"`
	ProvinceId        uint   `json:"province_id"`
	CityName          string `json:"city_name"`
	CityId            uint   `json:"city_id"`
	DistrictName      string `json:"district_name"`
	DistrictId        uint   `json:"district_id"`
	SubdistrictName   string `json:"subdistrict_name"`
	SubdistrictId     uint   `json:"subdistrict_id"`
	ZipCode           string `json:"zip_code"`
	Label             string `json:"label"`
	IsDefault         bool   `json:"is_default"`
	IsMerchantAddress bool   `json:"is_merchant_address"`
}

type UserAddressReqDTO struct {
	Phone         string `json:"phone" binding:"required"`
	Name          string `json:"name" binding:"required"`
	ProvinceId    uint   `json:"province_id" binding:"required"`
	CityId        uint   `json:"city_id" binding:"required"`
	DistrictId    uint   `json:"district_id" binding:"required"`
	SubdistrictId uint   `json:"subdistrict_id" binding:"required"`
	Label         string `json:"label" binding:"required"`
	Details       string `json:"details"`
	IsDefault     bool   `json:"is_default"`
}

type UserAddressUpdateReqDTO struct {
	Phone         string `json:"phone"`
	Name          string `json:"name"`
	ProvinceId    uint   `json:"province_id"`
	CityId        uint   `json:"city_id"`
	DistrictId    uint   `json:"district_id"`
	SubdistrictId uint   `json:"subdistrict_id"`
	Label         string `json:"label"`
	Details       string `json:"details"`
	IsDefault     bool   `json:"is_default"`
}
