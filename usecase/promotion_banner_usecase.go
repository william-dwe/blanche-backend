package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type PromotionBannerUsecase interface {
	GetPromotionBannerList(req dto.PaginationRequest) (*dto.PromotionBannerListResDTO, error)
	GetPromotionBannerByID(id uint) (*dto.PromotionBannerResDTO, error)
	CreatePromotionBanner(promotionBannerReqDTO dto.UpsertPromotionBannerReqDTO) (*dto.PromotionBannerResDTO, error)
	UpdatePromotionBanner(id uint, promotionBannerReqDTO dto.UpsertPromotionBannerReqDTO) (*dto.PromotionBannerResDTO, error)
	DeletePromotionBanner(id uint) (*dto.PromotionBannerResDTO, error)
}

type promotionBannerUsecaseImpl struct {
	promotionBannerRepo repository.PromotionBannerRepository
	mediaUsecase        MediaUsecase
}

type PromotionBannerUsecaseConfig struct {
	PromotionBannerRepo repository.PromotionBannerRepository
	MediaUsecase        MediaUsecase
}

func NewPromotionBannerUsecase(c PromotionBannerUsecaseConfig) PromotionBannerUsecase {
	return &promotionBannerUsecaseImpl{
		promotionBannerRepo: c.PromotionBannerRepo,
		mediaUsecase:        c.MediaUsecase,
	}
}

func (u *promotionBannerUsecaseImpl) GetPromotionBannerList(req dto.PaginationRequest) (*dto.PromotionBannerListResDTO, error) {
	promotionBanners, total, err := u.promotionBannerRepo.GetPromotionBannerList(req)
	if err != nil {
		return nil, err
	}

	promotionBannersDTOs := make([]dto.PromotionBannerResDTO, len(promotionBanners))
	for i, promotionBanner := range promotionBanners {
		promotionBannersDTOs[i] = dto.PromotionBannerResDTO{
			ID:          promotionBanner.ID,
			Name:        promotionBanner.Name,
			Description: promotionBanner.Description,
			ImageUrl:    promotionBanner.ImageUrl,
		}
	}

	return &dto.PromotionBannerListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   total,
			TotalPage:   (total + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
		PromotionBanners: promotionBannersDTOs,
	}, nil
}

func (u *promotionBannerUsecaseImpl) GetPromotionBannerByID(id uint) (*dto.PromotionBannerResDTO, error) {
	promotionBanner, err := u.promotionBannerRepo.GetPromotionBannerByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.PromotionBannerResDTO{
		ID:          promotionBanner.ID,
		Name:        promotionBanner.Name,
		Description: promotionBanner.Description,
		ImageUrl:    promotionBanner.ImageUrl,
	}, nil
}

func (u *promotionBannerUsecaseImpl) CreatePromotionBanner(promotionBannerReqDTO dto.UpsertPromotionBannerReqDTO) (*dto.PromotionBannerResDTO, error) {
	promotionBanner := entity.PromotionBanner{
		Name:        promotionBannerReqDTO.Name,
		Description: promotionBannerReqDTO.Description,
	}

	if promotionBannerReqDTO.Image != nil {
		bannerUrl, err := u.mediaUsecase.UploadFileForBinding(*promotionBannerReqDTO.Image, promotionBannerReqDTO.Image.Filename)
		if err != nil {
			return nil, err
		}
		promotionBanner.ImageUrl = bannerUrl
	}

	createdPromotionBanner, err := u.promotionBannerRepo.CreatePromotionBanner(promotionBanner)
	if err != nil {
		return nil, err
	}

	return &dto.PromotionBannerResDTO{
		ID:          createdPromotionBanner.ID,
		Name:        createdPromotionBanner.Name,
		Description: createdPromotionBanner.Description,
		ImageUrl:    createdPromotionBanner.ImageUrl,
	}, nil
}

func (u *promotionBannerUsecaseImpl) UpdatePromotionBanner(id uint, promotionBannerReqDTO dto.UpsertPromotionBannerReqDTO) (*dto.PromotionBannerResDTO, error) {
	newPromotionBanner := entity.PromotionBanner{
		ID:          id,
		Name:        promotionBannerReqDTO.Name,
		Description: promotionBannerReqDTO.Description,
	}

	if promotionBannerReqDTO.Image != nil {
		bannerUrl, err := u.mediaUsecase.UploadFileForBinding(*promotionBannerReqDTO.Image, promotionBannerReqDTO.Image.Filename)
		if err != nil {
			return nil, err
		}
		newPromotionBanner.ImageUrl = bannerUrl
	}

	updatedPromotionBanner, err := u.promotionBannerRepo.UpdatePromotionBanner(newPromotionBanner)
	if err != nil {
		return nil, err
	}

	return &dto.PromotionBannerResDTO{
		ID:          updatedPromotionBanner.ID,
		Name:        updatedPromotionBanner.Name,
		Description: updatedPromotionBanner.Description,
		ImageUrl:    updatedPromotionBanner.ImageUrl,
	}, nil
}

func (u *promotionBannerUsecaseImpl) DeletePromotionBanner(id uint) (*dto.PromotionBannerResDTO, error) {
	promotionBanner, err := u.promotionBannerRepo.GetPromotionBannerByID(id)
	if err != nil {
		return nil, err
	}

	promotionBanner, err = u.promotionBannerRepo.DeletePromotionBanner(*promotionBanner)
	if err != nil {
		return nil, err
	}

	return &dto.PromotionBannerResDTO{
		ID:          promotionBanner.ID,
		Name:        promotionBanner.Name,
		Description: promotionBanner.Description,
		ImageUrl:    promotionBanner.ImageUrl,
	}, nil
}
