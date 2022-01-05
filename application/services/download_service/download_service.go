package download_service

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/domain"
	"io/ioutil"
	"log"
	"os"
)

type DownloadUseCase interface {
	Execute(bucketName string, mp4TargetDir string, client *storage.Client, ctx context.Context) error
}

type DownloadService struct {
	Video *domain.Video
}

func NewDownloadService(v *domain.Video) *DownloadService {
	return &DownloadService{Video: v}
}

func (d *DownloadService) Execute(bucketName string, mp4TargetDir string, client *storage.Client, ctx context.Context) error {
	bkt := client.Bucket(bucketName)
	obj := bkt.Object(d.Video.FilePath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(mp4TargetDir)
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has bees stored", d.Video.ID)

	closeConnections(reader, file)
	return nil
}

func closeConnections(reader *storage.Reader, file *os.File) {
	defer reader.Close()
	defer file.Close()
}
