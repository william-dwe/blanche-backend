package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoryTree(req dto.CategoryListReqParamDTO) ([]entity.Category, error)
	GetCategoryList(req dto.PaginationRequest) ([]entity.Category, int64, error)
	GetCategoryAncestorsBySlug(slug string) ([]entity.Category, error)
	GetCategoryTreeByListId(categoryIds []uint) ([]entity.Category, error)
	GetCategoryBySlug(slug string) (*entity.Category, error)
	GetCategoryByID(id uint) (*entity.Category, error)

	CreateCategory(category *entity.Category) (*entity.Category, error)
	UpdateCategory(category *entity.Category) (*entity.Category, error)
	DeleteCategory(category *entity.Category) (*entity.Category, error)
	CheckCategoryInProduct(categoryId uint) (bool, error)
}

type CategoryRepositoryConfig struct {
	DB *gorm.DB
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(c CategoryRepositoryConfig) CategoryRepository {
	return &categoryRepositoryImpl{
		db: c.DB,
	}
}

func (r *categoryRepositoryImpl) GetCategoryTree(req dto.CategoryListReqParamDTO) ([]entity.Category, error) {
	var categories []entity.Category
	var query *gorm.DB
	switch req.Level {
	case 1:
		query = r.db.Where("parent_id IS NULL")
	case 2:
		query = r.db.Where("grandparent_id IS NULL").Preload("Children")
	default:
		query = r.db.Where("parent_id IS NULL").Preload("Children.Children")
	}
	err := query.Order("name asc").Find(&categories).Error
	if err != nil {
		return nil, domain.ErrGetCategories
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) GetCategoryList(req dto.PaginationRequest) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64
	pageOffset := req.Limit * (req.Page - 1)
	err := r.db.Offset(pageOffset).Limit(req.Limit).Order("created_at desc").Find(&categories).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return nil, total, domain.ErrGetCategories
	}

	return categories, total, nil
}

func (r *categoryRepositoryImpl) GetCategoryAncestorsBySlug(slug string) ([]entity.Category, error) {
	var categoryAncestor []entity.Category
	err := r.db.Raw(
		`WITH RECURSIVE ancestors AS (
			SELECT * FROM categories WHERE slug = ?
			UNION ALL
			SELECT c.* FROM categories c
			JOIN ancestors a ON c.id = a.parent_id
		)
		SELECT * FROM ancestors
		ORDER BY ROW_NUMBER () OVER () DESC;
		`, slug).Scan(&categoryAncestor).Error
	if err != nil {
		return nil, domain.ErrGetCategories
	}

	return categoryAncestor, nil
}

func (r *categoryRepositoryImpl) GetCategoryTreeByListId(categoryIds []uint) ([]entity.Category, error) {
	var categoryAncestor []entity.Category
	err := r.db.Preload("Parent").
		Preload("Grandparent").
		Where("id IN (?)", categoryIds).
		Find(&categoryAncestor).Error

	if err != nil {
		return nil, domain.ErrGetCategories
	}

	return categoryAncestor, nil
}

func (r *categoryRepositoryImpl) GetCategoryBySlug(slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("slug = ?", slug).Preload("Children.Children").First(&category).Error
	if err != nil {
		return nil, domain.ErrGetCategories
	}

	return &category, nil
}

func (r *categoryRepositoryImpl) GetCategoryByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("id = ?", id).Preload("Children.Children").First(&category).Error
	if err != nil {
		return nil, domain.ErrGetCategories
	}

	return &category, nil
}

func (r *categoryRepositoryImpl) CreateCategory(category *entity.Category) (*entity.Category, error) {
	err := r.db.Create(category).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"categories_un_slug": domain.ErrCategorySlugAlreadyExist,
			},
			domain.ErrCreateCategory,
		)

		return nil, maskedErr
	}

	return category, nil
}

func (r *categoryRepositoryImpl) UpdateCategory(category *entity.Category) (*entity.Category, error) {
	err := r.db.Model(category).Updates(category).Error
	if err != nil {
		maskedErr := util.PgConsErrMasker(
			err,
			entity.ConstraintErrMaskerMap{
				"categories_un_slug": domain.ErrCategorySlugAlreadyExist,
			},
			domain.ErrCreateCategory,
		)

		return nil, maskedErr
	}

	return category, nil
}

func (r *categoryRepositoryImpl) DeleteCategory(category *entity.Category) (*entity.Category, error) {
	err := r.db.Model(category).Delete(category).Error
	if err != nil {
		return nil, domain.ErrDeleteCategory
	}

	return category, nil
}

func (r *categoryRepositoryImpl) CheckCategoryInProduct(categoryId uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Product{}).Where("category_id = ?", categoryId).Count(&count).Error
	if err != nil {
		return false, domain.ErrGetCategories
	}

	return count == 0, nil
}
