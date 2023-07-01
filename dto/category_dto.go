package dto

import "mime/multipart"

type CategoryListReqParamDTO struct {
	Level int `form:"level,default=3"`
}

type CategoryResDTO struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	Slug     string           `json:"slug"`
	ImageUrl string           `json:"image_url"`
	Children []CategoryResDTO `json:"children"`
}

type UpsertCategoryReqDTO struct {
	Name     string                `form:"name" binding:"required"`
	ParentId uint                  `form:"parent_id"`
	Image    *multipart.FileHeader `form:"image"`
}

type UpsertCategoryResDTO struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	ImageUrl      string `json:"image_url"`
	ParentId      uint   `json:"parent_id"`
	GrandparentId uint   `json:"grandparent_id"`
}

type CategoryDetailResDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ImageUrl string `json:"image_url"`
	ParentId uint   `json:"parent_id"`
	Level    int    `json:"level"`
}

type CategoryAdminResDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Level int    `json:"level"`
}

type CategoryAdminListResDTO struct {
	PaginationResponse
	Categories []CategoryAdminResDTO `json:"categories"`
}
