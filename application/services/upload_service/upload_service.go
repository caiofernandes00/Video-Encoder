package upload_service

import (
	"cloud.google.com/go/storage"
	"context"
	"encoder/application/utils"
	"io"
	"os"
	"strings"
)

type UploadUseCase interface {
	Execute(objectPath string, client *storage.Client, ctx context.Context) error
}

type UploadService struct {
	OutputBucket string
}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (up *UploadService) Execute(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv(utils.LocalStoragePath)+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(up.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}
