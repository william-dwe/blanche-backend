package usecase

import (
	"mime/multipart"

	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/util"
)

type MediaUsecase interface {
	UploadFileForBinding(file multipart.FileHeader, object string) (string, error)
	DeleteFile(object string) error
}

type MediaUsecaseConfig struct {
	GCSUploader util.GCSUploader
}

type mediaUsecaseImpl struct {
	gcsUploader util.GCSUploader
}

func NewMediaUsecase(c MediaUsecaseConfig) MediaUsecase {
	return &mediaUsecaseImpl{
		gcsUploader: c.GCSUploader,
	}
}

func (u *mediaUsecaseImpl) UploadFileForBinding(file multipart.FileHeader, object string) (string, error) {
	url, err := u.gcsUploader.UploadFileFromFileHeader(file, object)
	if err != nil {
		return "", domain.ErrUploadFile
	}

	return url, nil
}

func (u *mediaUsecaseImpl) DeleteFile(object string) error {
	err := u.gcsUploader.DeleteFile(object)
	if err != nil {
		return domain.ErrDeleteFile
	}

	return nil
}
