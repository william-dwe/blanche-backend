package dto

import "mime/multipart"

type UploadImageReqDTO struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}

type UploadImageResDTO struct {
	ImageURLs string `json:"image_url"`
}
