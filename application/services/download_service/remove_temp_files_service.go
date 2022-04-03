package download_service

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

func (v *RemoveTempFilesService) Execute(mp4File string, fragFile string, encodeDir string) error {

	err := v.removeFile(mp4File, fragFile)
	if err != nil {
		return err
	}

	err = v.removeDir(encodeDir)
	if err != nil {
		return err
	}

	log.Println("files have been removed: ", v.Video.ID)
	return nil
}

func (v *RemoveTempFilesService) removeFile(files ...string) error {
	var err error

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			log.Println("error removing mp4", v.Video.ID+utils.Mp4Format, " with error: ", err)
			break
		}
	}

	return err
}

func (v *RemoveTempFilesService) removeDir(dirs ...string) error {
	var err error

	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Println("error removing", v.Video.ID, "folder with error: ", err)
		}
	}

	return err
}
