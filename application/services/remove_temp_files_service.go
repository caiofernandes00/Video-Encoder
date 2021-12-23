package services

import (
	"encoder/application/utils"
	"encoder/domain"
	"log"
	"os"
)

type RemoveTempFilesUseCase interface {
	Execute() error
}

type RemoveTempFilesService struct {
	Video *domain.Video
}

func NewRemoveTempFilesService(v *domain.Video) *RemoveTempFilesService {
	return &RemoveTempFilesService{Video: v}
}

func (v *RemoveTempFilesService) Execute() error {

	err := os.Remove(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + utils.Mp4Format)
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+utils.Mp4Format, " with error: ", err)
		return err
	}

	err = os.Remove(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + utils.FragFile)
	if err != nil {
		log.Println("error removing mp4", v.Video.ID+utils.FragFile, " with error: ", err)
		return err
	}

	err = os.RemoveAll(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID)
	if err != nil {
		log.Println("error removing", v.Video.ID, "folder with error: ", err)
		return err
	}

	log.Println("files have been removed: ", v.Video.ID)
	return nil
}
