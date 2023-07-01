package usecase

import (
	"strings"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type CategoryUsecase interface {
	GetCategoryTree(req dto.CategoryListReqParamDTO) ([]dto.CategoryResDTO, error)
	GetCategoryAncestorsBySlug(slug string) (*dto.CategoryResDTO, error)
	GetCategoryBySlug(slug string) (*dto.CategoryResDTO, error)
	GetCategoryList(req dto.PaginationRequest) (*dto.CategoryAdminListResDTO, error)

	CreateCategory(req dto.UpsertCategoryReqDTO) (*dto.UpsertCategoryResDTO, error)
	DeleteCategory(categoryId uint) (*dto.UpsertCategoryResDTO, error)
	UpdateCategory(categoryId uint, req dto.UpsertCategoryReqDTO) (*dto.UpsertCategoryResDTO, error)
	GetCategoryByID(id uint) (*dto.CategoryDetailResDTO, error)
}

type CategoryUsecaseConfig struct {
	CategoryRepository repository.CategoryRepository
	MediaUsecase       MediaUsecase
}

type categoryUsecaseImpl struct {
	categoryRepository repository.CategoryRepository
	mediaUsecase       MediaUsecase
}

func NewCategoryUsecase(c CategoryUsecaseConfig) CategoryUsecase {
	return &categoryUsecaseImpl{
		categoryRepository: c.CategoryRepository,
		mediaUsecase:       c.MediaUsecase,
	}
}

func (u *categoryUsecaseImpl) GetCategoryTree(req dto.CategoryListReqParamDTO) ([]dto.CategoryResDTO, error) {
	categories, err := u.categoryRepository.GetCategoryTree(req)
	if err != nil {
		return nil, err
	}

	var categoryTreeDTO []dto.CategoryResDTO
	for _, category := range categories {
		categoryTreeDTO = append(categoryTreeDTO, util.CategoryToDTOList(category))
	}

	return categoryTreeDTO, nil
}

func (u *categoryUsecaseImpl) GetCategoryList(req dto.PaginationRequest) (*dto.CategoryAdminListResDTO, error) {
	categories, total, err := u.categoryRepository.GetCategoryList(req)
	if err != nil {
		return nil, err
	}

	if categories == nil {
		return &dto.CategoryAdminListResDTO{
			PaginationResponse: dto.PaginationResponse{
				TotalData: int64(req.Page),
			},
			Categories: []dto.CategoryAdminResDTO{},
		}, nil
	}

	categoryListDTO := make([]dto.CategoryAdminResDTO, 0)
	for _, category := range categories {
		categoryListDTO = append(categoryListDTO, dto.CategoryAdminResDTO{
			Name: category.Name,
			Slug: category.Slug,
			ID:   category.ID,
		})
		if category.ParentId == 0 && category.GrandparentId == 0 {
			categoryListDTO[len(categoryListDTO)-1].Level = 1
		}
		if category.ParentId != 0 && category.GrandparentId == 0 {
			categoryListDTO[len(categoryListDTO)-1].Level = 2
		}
		if category.ParentId != 0 && category.GrandparentId != 0 {
			categoryListDTO[len(categoryListDTO)-1].Level = 3
		}
	}

	return &dto.CategoryAdminListResDTO{
		PaginationResponse: dto.PaginationResponse{
			TotalData:   total,
			TotalPage:   (total + int64(req.Limit) - 1) / int64(req.Limit),
			CurrentPage: req.Page,
		},
		Categories: categoryListDTO,
	}, nil
}

func (u *categoryUsecaseImpl) GetCategoryAncestorsBySlug(slug string) (*dto.CategoryResDTO, error) {
	categories, err := u.categoryRepository.GetCategoryAncestorsBySlug(slug)
	if err != nil {
		return nil, err
	}

	var categoryAncestorsDTO dto.CategoryResDTO
	var parent *dto.CategoryResDTO
	for _, category := range categories {
		if category.ParentId == 0 {
			categoryAncestorsDTO = util.CategoryToDTOList(category)
			parent = &categoryAncestorsDTO
		} else {
			if parent.ID == category.ParentId {
				parent.Children = append(parent.Children, util.CategoryToDTOList(category))
				parent = &parent.Children[len(parent.Children)-1]
			}
		}
	}
	return &categoryAncestorsDTO, nil
}

func (u *categoryUsecaseImpl) GetCategoryBySlug(slug string) (*dto.CategoryResDTO, error) {
	category, err := u.categoryRepository.GetCategoryBySlug(slug)
	if err != nil {
		return nil, err
	}

	categoryDTO := util.CategoryToDTOList(*category)

	return &categoryDTO, nil
}

func (u *categoryUsecaseImpl) CreateCategory(req dto.UpsertCategoryReqDTO) (*dto.UpsertCategoryResDTO, error) {
	categoryNew := entity.Category{
		Name: req.Name,
	}

	if req.ParentId != 0 {
		categoryParent, err := u.categoryRepository.GetCategoryByID(req.ParentId)
		if err != nil {
			return nil, err
		}
		if categoryParent != nil {
			categoryNew.ParentId = categoryParent.ID
			if categoryParent.ParentId != 0 && categoryParent.GrandparentId != 0 {
				return nil, domain.ErrCategoryExceedLimit
			}
			if categoryParent.ParentId != 0 && categoryParent.GrandparentId == 0 {
				categoryNew.GrandparentId = categoryParent.ParentId
			}
		}
	}

	if req.Image != nil {
		imageUrl, err := u.mediaUsecase.UploadFileForBinding(*req.Image, req.Image.Filename)
		if err != nil {
			return nil, err
		}

		categoryNew.ImageUrl = imageUrl
	}

	slug, err := u.generateCategorySlug(categoryNew.Name, categoryNew.ParentId, categoryNew.GrandparentId)
	if err != nil {
		return nil, err
	}
	categoryNew.Slug = slug

	category, err := u.categoryRepository.CreateCategory(&categoryNew)
	if err != nil {
		return nil, err
	}

	return &dto.UpsertCategoryResDTO{
		ID:            category.ID,
		Name:          category.Name,
		Slug:          category.Slug,
		ImageUrl:      category.ImageUrl,
		ParentId:      category.ParentId,
		GrandparentId: category.GrandparentId,
	}, nil
}

func (u *categoryUsecaseImpl) generateCategorySlug(categoryName string, parentId, grandparentId uint) (string, error) {
	slug := strings.Join(strings.Split(strings.ToLower(categoryName), " "), "-")
	if parentId == 0 || grandparentId == 0 {
		return slug, nil
	}

	parentCategory, err := u.categoryRepository.GetCategoryByID(parentId)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{parentCategory.Slug, slug}, "-"), nil
}

func (u *categoryUsecaseImpl) GetCategoryByID(id uint) (*dto.CategoryDetailResDTO, error) {
	category, err := u.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	categoryResDTO := &dto.CategoryDetailResDTO{
		ID:       category.ID,
		Name:     category.Name,
		Slug:     category.Slug,
		ImageUrl: category.ImageUrl,
		ParentId: category.ParentId,
	}
	if category.ParentId == 0 && category.GrandparentId == 0 {
		categoryResDTO.Level = 1
	}
	if category.ParentId != 0 && category.GrandparentId == 0 {
		categoryResDTO.Level = 2
	}
	if category.ParentId != 0 && category.GrandparentId != 0 {
		categoryResDTO.Level = 3
	}

	return categoryResDTO, nil
}

func (u *categoryUsecaseImpl) UpdateCategory(categoryId uint, req dto.UpsertCategoryReqDTO) (*dto.UpsertCategoryResDTO, error) {
	category, err := u.categoryRepository.GetCategoryByID(categoryId)
	if err != nil {
		return nil, err
	}

	if categoryId != category.ID {
		return nil, domain.ErrCategoryNotAuthorized
	}

	categoryNew := entity.Category{
		ID:   category.ID,
		Name: req.Name,
	}

	if req.Image != nil {
		imageUrl, err := u.mediaUsecase.UploadFileForBinding(*req.Image, req.Image.Filename)
		if err != nil {
			return nil, err
		}

		categoryNew.ImageUrl = imageUrl
	}

	if req.ParentId != 0 {
		categoryParent, err := u.categoryRepository.GetCategoryByID(req.ParentId)
		if err != nil {
			return nil, err
		}
		if categoryParent != nil {
			categoryNew.ParentId = categoryParent.ID
			if categoryParent.ParentId != 0 && categoryParent.GrandparentId != 0 {
				return nil, domain.ErrCategoryExceedLimit
			}
			if categoryParent.ParentId != 0 && categoryParent.GrandparentId == 0 {
				categoryNew.GrandparentId = categoryParent.ParentId
			}
		}
	}

	slug, err := u.generateCategorySlug(categoryNew.Name, categoryNew.ParentId, categoryNew.GrandparentId)
	if err != nil {
		return nil, err
	}
	categoryNew.Slug = slug

	ok, err := u.categoryRepository.CheckCategoryInProduct(categoryId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrCategoryInUse
	}

	category, err = u.categoryRepository.UpdateCategory(&categoryNew)
	if err != nil {
		return nil, err
	}

	return &dto.UpsertCategoryResDTO{
		ID:            category.ID,
		Name:          category.Name,
		Slug:          category.Slug,
		ImageUrl:      category.ImageUrl,
		ParentId:      category.ParentId,
		GrandparentId: category.GrandparentId,
	}, nil
}

func (u *categoryUsecaseImpl) DeleteCategory(categoryId uint) (*dto.UpsertCategoryResDTO, error) {
	ok, err := u.categoryRepository.CheckCategoryInProduct(categoryId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrCategoryInUse
	}

	category, err := u.categoryRepository.GetCategoryByID(categoryId)
	if err != nil {
		return nil, err
	}

	if categoryId != category.ID {
		return nil, domain.ErrCategoryNotAuthorized
	}

	category, err = u.categoryRepository.DeleteCategory(category)
	if err != nil {
		return nil, err
	}

	return &dto.UpsertCategoryResDTO{
		ID:            category.ID,
		Name:          category.Name,
		Slug:          category.Slug,
		ImageUrl:      category.ImageUrl,
		ParentId:      category.ParentId,
		GrandparentId: category.GrandparentId,
	}, nil
}
