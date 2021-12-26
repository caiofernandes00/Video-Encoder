package services

import (
	downloadService "encoder/application/services/download_service"
	"github.com/joho/godotenv"
	"log"
)

type DownloadManagerFacadeUseCase interface {
	Execute() error
}

type DownloadManagerFacadeService struct {
	DownloadUseCase        downloadService.DownloadUseCase
	FragmentUseCase        downloadService.FragmentUseCase
	EncodeUseCase          downloadService.EncodeUseCase
	RemoveTempFilesUseCase downloadService.RemoveTempFilesUseCase
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
}

func NewDownloadManagerFacadeService(
	downloadUseCase downloadService.DownloadUseCase,
	fragmentUseCase downloadService.FragmentUseCase,
	encodeUseCase downloadService.EncodeUseCase,
	removeTempFilesUseCase downloadService.RemoveTempFilesUseCase,
) DownloadManagerFacadeService {
	return DownloadManagerFacadeService{
		DownloadUseCase:        downloadUseCase,
		FragmentUseCase:        fragmentUseCase,
		EncodeUseCase:          encodeUseCase,
		RemoveTempFilesUseCase: removeTempFilesUseCase,
	}
}

func (v *DownloadManagerFacadeService) Execute(bucketName string) error {
	var err error

	err = v.DownloadUseCase.Execute(bucketName)
	if err != nil {
		return err
	}

	err = v.FragmentUseCase.Execute()
	if err != nil {
		return err
	}

	err = v.EncodeUseCase.Execute()
	if err != nil {
		return err
	}

	err = v.RemoveTempFilesUseCase.Execute()
	if err != nil {
		return err
	}

	return nil
}
