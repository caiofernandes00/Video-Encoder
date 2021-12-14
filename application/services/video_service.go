package services

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/application/repositories"
	"encoder/application/utils"
	"encoder/domain"
	"io/ioutil"
	"log"
	"os"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketname string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketname)
	obj := bkt.Object(v.Video.FilePath)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv(utils.LocalStoragePath) + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has bees stored", v.Video.ID)

	closeConnections(reader, file)
	return nil
}

func closeConnections(reader *storage.Reader, file *os.File) {
	defer reader.Close()
	defer file.Close()
}
