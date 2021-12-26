package services

import uploadService "encoder/application/services/upload_service"

type UploadManagerFacadeUseCase interface {
	Execute() error
}

type UploadManagerFacadeService struct {
	UploadService *uploadService.UploadService
}

func NewVideoUpload(uploadService *uploadService.UploadService) *UploadManagerFacadeService {
	return &UploadManagerFacadeService{UploadService: uploadService}
}

func (u *UploadManagerFacadeService) Execute() error {
	return nil
}
