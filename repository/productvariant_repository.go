package repository

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/db"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	GetVariantItemById(uint) (*entity.VariantItem, error)
	GetVariantItemsByProductId(uint) ([]entity.VariantItem, error)
}

type ProductVariantRepositoryConfig struct {
	DB *gorm.DB
}

type productVariantRepositoryImpl struct {
	db *gorm.DB
}

func NewProductVariantRepository(c ProductVariantRepositoryConfig) ProductVariantRepository {
	return &productVariantRepositoryImpl{
		db: c.DB,
	}
}

func (r *productVariantRepositoryImpl) GetVariantItemById(variantItemId uint) (*entity.VariantItem, error) {
	var variantItem entity.VariantItem
	err := db.Get().Debug().
		Preload("VariantSpecs").
		Preload("VariantSpecs.VariationGroup").
		Where("id = ?", variantItemId).
		Find(&variantItem).Error

	if err != nil {
		return nil, domain.ErrGetProductVariant
	}

	return &variantItem, nil
}

func (r *productVariantRepositoryImpl) GetVariantItemsByProductId(productId uint) ([]entity.VariantItem, error) {
	var variantItems []entity.VariantItem
	err := db.Get().Debug().
		Preload("VariantSpecs").
		Preload("VariantSpecs.VariationGroup").
		Where("product_id = ?", productId).
		Order("id ASC").Find(&variantItems).Error

	if err != nil {
		return nil, domain.ErrGetProductVariant
	}

	return variantItems, nil
}
