package util

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
)

func CategoryToDTOList(category entity.Category) dto.CategoryResDTO {
	returnCategory := dto.CategoryResDTO{
		ID:       category.ID,
		Name:     category.Name,
		Slug:     category.Slug,
		ImageUrl: category.ImageUrl,
	}
	if len(category.Children) > 0 {
		for _, child := range category.Children {
			returnCategory.Children = append(returnCategory.Children, CategoryToDTOList(child))
		}
	}

	return returnCategory
}
