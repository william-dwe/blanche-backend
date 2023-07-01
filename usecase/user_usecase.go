package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type UserUsecase interface {
	GetProfile(email string) (*dto.UserProfileResDTO, error)
	UpdateProfile(username string, input dto.UserUpdateProfileReqDTO) (*dto.UserProfileResDTO, error)
	UpdateProfileDetail(username string, input dto.UserUpdateProfileDetailFormReqDTO) (*dto.UserProfileResDTO, error)

	GetAllUserAddress(username string) ([]dto.UserAddressDTO, error)
	AddUserAddress(username string, userAddressReq dto.UserAddressReqDTO) (*dto.UserAddressDTO, error)
	SetDefaultUserAddress(username string, userAddressId uint) (*dto.UserAddressDTO, error)
	DeleteUserAddress(username string, userAddressId uint) (*dto.UserAddressDTO, error)
	UpdateUserAddress(username string, userAddressId uint, userAddressReq dto.UserAddressUpdateReqDTO) (*dto.UserAddressDTO, error)
}

type UserUsecaseConfig struct {
	UserRepository     repository.UserRepository
	AddressRepository  repository.AddressRepository
	MerchantRepository repository.MerchantRepository
	MediaUsecase       MediaUsecase
}

type userUsecaseImpl struct {
	userRepository     repository.UserRepository
	addressRepository  repository.AddressRepository
	merchantRepository repository.MerchantRepository
	mediaUsecase       MediaUsecase
}

func NewUserUsecase(c UserUsecaseConfig) UserUsecase {
	return &userUsecaseImpl{
		userRepository:     c.UserRepository,
		addressRepository:  c.AddressRepository,
		merchantRepository: c.MerchantRepository,
		mediaUsecase:       c.MediaUsecase,
	}
}

func (u *userUsecaseImpl) GetProfile(username string) (*dto.UserProfileResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userWithDetail, err := u.userRepository.GetUserAndUserDetailByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	resBody := dto.UserProfileResDTO{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Fullname:       userWithDetail.UserDetail.Fullname,
		Phone:          userWithDetail.UserDetail.Phone,
		Gender:         userWithDetail.UserDetail.Gender,
		BirthDate:      userWithDetail.UserDetail.BirthDate,
		ProfilePicture: userWithDetail.UserDetail.ProfilePicture,
		Role:           userWithDetail.Role.RoleName,
	}
	return &resBody, nil
}

func (u *userUsecaseImpl) UpdateProfile(username string, input dto.UserUpdateProfileReqDTO) (*dto.UserProfileResDTO, error) {
	if isEmailValid := util.IsValidEmail(input.Email); !isEmailValid {
		return nil, domain.ErrInvalidEmailFormat
	}
	err := u.userRepository.CheckEmailBlacklist(input.Email)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Email == input.Email {
		return nil, domain.ErrUserEmailAlreadyExist
	}

	newUser := entity.User{
		ID:    user.ID,
		Email: input.Email,
	}
	newUserEmailBlacklist := entity.UserEmailBlacklist{
		UserId: user.ID,
		Email:  user.Email,
	}

	err = u.userRepository.UpdateUser(newUser)
	if err != nil {
		return nil, err
	}

	userWithDetail, err := u.userRepository.GetUserAndUserDetailByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	err = u.userRepository.AddUserEmailBlacklist(newUserEmailBlacklist)
	if err != nil {
		return nil, err
	}

	resBody := dto.UserProfileResDTO{
		ID:             userWithDetail.ID,
		Username:       userWithDetail.Username,
		Email:          userWithDetail.Email,
		Fullname:       userWithDetail.UserDetail.Fullname,
		Phone:          userWithDetail.UserDetail.Phone,
		Gender:         userWithDetail.UserDetail.Gender,
		BirthDate:      userWithDetail.UserDetail.BirthDate,
		ProfilePicture: userWithDetail.UserDetail.ProfilePicture,
	}

	return &resBody, nil
}

func (u *userUsecaseImpl) UpdateProfileDetail(username string, input dto.UserUpdateProfileDetailFormReqDTO) (*dto.UserProfileResDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := util.ValidateFullname(input.Fullname); err != nil && input.Fullname != "" {
		return nil, err
	}
	if input.Phone != nil {
		if err := util.ValidatePhone(*input.Phone); err != nil {
			return nil, err
		}
	}

	newUserDetail := entity.UserDetail{
		UserId:    user.ID,
		Fullname:  input.Fullname,
		Phone:     input.Phone,
		Gender:    input.Gender,
		BirthDate: input.BirthDate,
	}
	if input.ProfilePicture != nil {
		url, err := u.mediaUsecase.UploadFileForBinding(*input.ProfilePicture, "profile_picture:"+username)
		if err != nil {
			return nil, err
		}
		newUserDetail.ProfilePicture = &url
	}

	err = u.userRepository.UpdateUserDetail(newUserDetail)
	if err != nil {
		return nil, err
	}

	userWithDetail, err := u.userRepository.GetUserAndUserDetailByUserId(user.ID)
	if err != nil {
		return nil, err
	}

	resBody := dto.UserProfileResDTO{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Fullname:       userWithDetail.UserDetail.Fullname,
		Phone:          userWithDetail.UserDetail.Phone,
		Gender:         userWithDetail.UserDetail.Gender,
		BirthDate:      userWithDetail.UserDetail.BirthDate,
		ProfilePicture: userWithDetail.UserDetail.ProfilePicture,
	}

	return &resBody, nil
}

func (u *userUsecaseImpl) GetAllUserAddress(username string) ([]dto.UserAddressDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userAddresses, err := u.userRepository.GetAllUserAddress(*user)
	if err != nil {
		return nil, err
	}

	merchant, err := u.merchantRepository.GetByUserID(user.ID)
	if err != nil && err != domain.ErrMerchantUserIDNotFound {
		return nil, err
	}

	var resBody []dto.UserAddressDTO
	for _, userAddress := range userAddresses {
		resBody = append(resBody, dto.UserAddressDTO{
			ID:              userAddress.ID,
			Phone:           userAddress.PhoneNumber,
			Name:            userAddress.Name,
			Details:         userAddress.Details,
			ProvinceName:    userAddress.Province.Name,
			ProvinceId:      userAddress.Province.ID,
			CityName:        userAddress.City.Name,
			CityId:          userAddress.City.ID,
			DistrictName:    userAddress.District.Name,
			DistrictId:      userAddress.District.ID,
			SubdistrictName: userAddress.Subdistrict.Name,
			SubdistrictId:   userAddress.Subdistrict.ID,
			ZipCode:         userAddress.Subdistrict.ZipCode,
			Label:           userAddress.Label,
			IsDefault:       userAddress.IsDefault,
		})

		if merchant != nil && merchant.UserAddressId == userAddress.ID {
			resBody[len(resBody)-1].IsMerchantAddress = true
		}
	}

	return resBody, nil
}

func (u *userUsecaseImpl) AddUserAddress(username string, userAddressReq dto.UserAddressReqDTO) (*dto.UserAddressDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := util.ValidatePhone(userAddressReq.Phone); err != nil {
		return nil, err
	}

	province, err := u.addressRepository.GetProvinceById(userAddressReq.ProvinceId)
	if err != nil {
		return nil, err
	}
	if province.ID != userAddressReq.ProvinceId {
		return nil, domain.ErrInvalidProvince
	}

	city, err := u.addressRepository.GetCityById(userAddressReq.CityId)
	if err != nil {
		return nil, err
	}
	if city.ProvinceID != userAddressReq.ProvinceId {
		return nil, domain.ErrInvalidCity
	}

	district, err := u.addressRepository.GetDistrictById(userAddressReq.DistrictId)
	if err != nil {
		return nil, err
	}
	if district.CityId != userAddressReq.CityId {
		return nil, domain.ErrInvalidDistrict
	}

	subdistrict, err := u.addressRepository.GetSubDistrictById(userAddressReq.SubdistrictId)
	if err != nil {
		return nil, err
	}
	if subdistrict.DistrictId != userAddressReq.DistrictId {
		return nil, domain.ErrInvalidSubdistrict
	}

	userAddresses, err := u.userRepository.GetAllUserAddress(*user)
	if err != nil {
		return nil, err
	}
	if len(userAddresses) == 0 {
		userAddressReq.IsDefault = true
	}

	userAddress := entity.UserAddress{
		UserID:        user.ID,
		PhoneNumber:   userAddressReq.Phone,
		ProvinceId:    userAddressReq.ProvinceId,
		CityId:        userAddressReq.CityId,
		DistrictId:    userAddressReq.DistrictId,
		SubdistrictId: userAddressReq.SubdistrictId,
		Name:          userAddressReq.Name,
		Label:         userAddressReq.Label,
		Details:       userAddressReq.Details,
		IsDefault:     userAddressReq.IsDefault,
	}

	userAddressRes, err := u.userRepository.AddUserAddress(userAddress)
	if err != nil {
		return nil, err
	}

	return &dto.UserAddressDTO{
		ID:                userAddressRes.ID,
		Phone:             userAddressRes.PhoneNumber,
		Name:              userAddressRes.Name,
		Details:           userAddressRes.Details,
		ProvinceName:      province.Name,
		ProvinceId:        province.ID,
		CityName:          city.Name,
		CityId:            city.ID,
		DistrictName:      district.Name,
		DistrictId:        district.ID,
		SubdistrictName:   subdistrict.Name,
		SubdistrictId:     subdistrict.ID,
		ZipCode:           subdistrict.ZipCode,
		Label:             userAddressRes.Label,
		IsDefault:         userAddressRes.IsDefault,
		IsMerchantAddress: false,
	}, nil
}

func (u *userUsecaseImpl) SetDefaultUserAddress(username string, userAddressId uint) (*dto.UserAddressDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userAddress, err := u.userRepository.GetUserAddressById(user.ID, userAddressId)
	if err != nil {
		return nil, err
	}

	if userAddress.UserID != user.ID {
		return nil, domain.ErrInvalidCredentialUserAddress
	}

	err = u.userRepository.SetDefaultUserAddress(user.ID, userAddress.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserAddressDTO{
		ID:              userAddress.ID,
		Phone:           userAddress.PhoneNumber,
		Name:            userAddress.Name,
		Details:         userAddress.Details,
		ProvinceName:    userAddress.Province.Name,
		ProvinceId:      userAddress.Province.ID,
		CityName:        userAddress.City.Name,
		CityId:          userAddress.City.ID,
		DistrictName:    userAddress.District.Name,
		DistrictId:      userAddress.District.ID,
		SubdistrictName: userAddress.Subdistrict.Name,
		SubdistrictId:   userAddress.Subdistrict.ID,
		ZipCode:         userAddress.Subdistrict.ZipCode,
		Label:           userAddress.Label,
		IsDefault:       true,
	}, nil
}

func (u *userUsecaseImpl) DeleteUserAddress(username string, userAddressId uint) (*dto.UserAddressDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userAddress, err := u.userRepository.GetUserAddressById(user.ID, userAddressId)
	if err != nil {
		return nil, err
	}

	merchant, err := u.merchantRepository.GetByUserID(user.ID)
	if err != nil && err != domain.ErrMerchantUserIDNotFound {
		return nil, err
	}

	if userAddress.UserID != user.ID {
		return nil, domain.ErrInvalidCredentialUserAddress
	}
	if userAddress.IsDefault {
		return nil, domain.ErrDeleteDefaultUserAddress
	}

	if merchant != nil {
		if merchant.UserAddressId == userAddress.ID {
			return nil, domain.ErrDeleteDefaultMerchantAddress
		}
	}

	userAddress, err = u.userRepository.DeleteUserAddress(*userAddress)
	if err != nil {
		return nil, err
	}

	return &dto.UserAddressDTO{
		ID:                userAddress.ID,
		Phone:             userAddress.PhoneNumber,
		Name:              userAddress.Name,
		Details:           userAddress.Details,
		ProvinceName:      userAddress.Province.Name,
		ProvinceId:        userAddress.Province.ID,
		CityName:          userAddress.City.Name,
		CityId:            userAddress.City.ID,
		DistrictName:      userAddress.District.Name,
		DistrictId:        userAddress.District.ID,
		SubdistrictName:   userAddress.Subdistrict.Name,
		SubdistrictId:     userAddress.Subdistrict.ID,
		ZipCode:           userAddress.Subdistrict.ZipCode,
		Label:             userAddress.Label,
		IsDefault:         userAddress.IsDefault,
		IsMerchantAddress: merchant != nil && merchant.UserAddressId == userAddress.ID,
	}, nil
}

func (u *userUsecaseImpl) UpdateUserAddress(username string, userAddressId uint, userAddressReq dto.UserAddressUpdateReqDTO) (*dto.UserAddressDTO, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userAddress, err := u.userRepository.GetUserAddressById(user.ID, userAddressId)
	if err != nil {
		return nil, err
	}

	if userAddress.UserID != user.ID {
		return nil, domain.ErrInvalidCredentialUserAddress
	}

	if userAddressReq.ProvinceId != 0 {
		province, err := u.addressRepository.GetProvinceById(userAddressReq.ProvinceId)
		if err != nil {
			return nil, err
		}
		if province.ID != userAddressReq.ProvinceId {
			return nil, domain.ErrInvalidProvince
		}
	}

	if userAddressReq.CityId != 0 {
		city, err := u.addressRepository.GetCityById(userAddressReq.CityId)
		if err != nil {
			return nil, err
		}
		if city.ProvinceID != userAddressReq.ProvinceId {
			return nil, domain.ErrInvalidCity
		}
	}

	if userAddressReq.DistrictId != 0 {
		district, err := u.addressRepository.GetDistrictById(userAddressReq.DistrictId)
		if err != nil {
			return nil, err
		}
		if district.CityId != userAddressReq.CityId {
			return nil, domain.ErrInvalidDistrict
		}
	}

	if userAddressReq.SubdistrictId != 0 {
		subdistrict, err := u.addressRepository.GetSubDistrictById(userAddressReq.SubdistrictId)
		if err != nil {
			return nil, err
		}
		if subdistrict.DistrictId != userAddressReq.DistrictId {
			return nil, domain.ErrInvalidSubdistrict
		}
	}

	newUserAddress := entity.UserAddress{
		UserID:        user.ID,
		PhoneNumber:   userAddressReq.Phone,
		ProvinceId:    userAddressReq.ProvinceId,
		CityId:        userAddressReq.CityId,
		DistrictId:    userAddressReq.DistrictId,
		SubdistrictId: userAddressReq.SubdistrictId,
		Name:          userAddressReq.Name,
		Label:         userAddressReq.Label,
		Details:       userAddressReq.Details,
	}

	userAddressRes, err := u.userRepository.UpdateUserAddress(*userAddress, newUserAddress)
	if err != nil {
		return nil, err
	}

	merchant, err := u.merchantRepository.GetByUserID(user.ID)
	if err != nil && err != domain.ErrMerchantUserIDNotFound {
		return nil, err
	}
	isMerchantAddress := merchant != nil && merchant.UserAddressId == userAddress.ID
	if isMerchantAddress {
		err = u.merchantRepository.SynchronizeMerchantCity(merchant)
		if err != nil {
			return nil, err
		}
	}

	return &dto.UserAddressDTO{
		ID:                userAddress.ID,
		Phone:             userAddressRes.PhoneNumber,
		Name:              userAddressRes.Name,
		Details:           userAddressRes.Details,
		ProvinceName:      userAddressRes.Province.Name,
		ProvinceId:        userAddressRes.Province.ID,
		CityName:          userAddressRes.City.Name,
		CityId:            userAddressRes.City.ID,
		DistrictName:      userAddressRes.District.Name,
		DistrictId:        userAddressRes.District.ID,
		SubdistrictName:   userAddressRes.Subdistrict.Name,
		SubdistrictId:     userAddressRes.Subdistrict.ID,
		ZipCode:           userAddressRes.Subdistrict.ZipCode,
		Label:             userAddressRes.Label,
		IsDefault:         userAddressRes.IsDefault,
		IsMerchantAddress: isMerchantAddress,
	}, nil
}
