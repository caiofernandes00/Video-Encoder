package upload_service

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/domain"
	"io"
	"os"
)

type UploadUseCase interface {
	Execute(dirPath string, outputBucket string, client *storage.Client, ctx context.Context) error
}

type UploadService struct {
	Video *domain.Video
}

func NewUploadService(video *domain.Video) *UploadService {
	return &UploadService{Video: video}
}

func (up *UploadService) Execute(dirPath string, outputBucket string, client *storage.Client, ctx context.Context) error {

	f, err := os.Open(dirPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(outputBucket).Object(up.Video.ID).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
