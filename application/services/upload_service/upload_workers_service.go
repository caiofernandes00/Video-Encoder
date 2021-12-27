package upload_service

import (
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
)

type UploadWorkersUseCase interface {
	Execute(concurrency int, doneUpload chan string) error
}

type UploadWorkersService struct {
	Paths         []string
	VideoPath     string
	Errors        []string
	UploadUseCase UploadUseCase
}

func NewUploadWorkersService(uploadUseCase *UploadService, paths []string, videoPath string) *UploadWorkersService {
	return &UploadWorkersService{
		Paths: paths, VideoPath: videoPath, UploadUseCase: uploadUseCase,
	}
}

func (uw *UploadWorkersService) Execute(concurrency int, doneUpload chan string) error {

	in := make(chan int, runtime.NumCPU())
	returnChan := make(chan string)

	err := uw.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go uw.uploadWorker(in, returnChan, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(uw.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChan {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (uw *UploadWorkersService) uploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for x := range in {
		err := uw.UploadUseCase.Execute(uw.Paths[x], uploadClient, ctx)

		if err != nil {
			uw.Errors = append(uw.Errors, uw.Paths[x])
			log.Printf("error during the upload: %v. Error: %v", uw.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "uploaded completed"
}

func (uw *UploadWorkersService) loadPaths() error {
	err := filepath.Walk(uw.VideoPath, func(path string, info fs.FileInfo, err error) error {

		if !info.IsDir() {
			uw.Paths = append(uw.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, err
}
