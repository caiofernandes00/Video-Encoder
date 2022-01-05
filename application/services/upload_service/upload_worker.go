package upload_service

import (
	"cloud.google.com/go/storage"
	"encoder/application/utils"
	"golang.org/x/net/context"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
)

type UploadWorkersUseCase interface {
	Execute(concurrency int, doneUpload chan string, client *storage.Client, ctx context.Context) error
}

type UploadWorkersService struct {
	FilePaths     []string
	VideoPathDir  string
	Errors        []string
	UploadUseCase UploadUseCase
}

func NewUploadWorkersService(uploadUseCase *UploadService, videoPathDir string) *UploadWorkersService {
	return &UploadWorkersService{
		VideoPathDir: videoPathDir, UploadUseCase: uploadUseCase,
	}
}

func (uw *UploadWorkersService) Execute(concurrency int, doneUpload chan string, client *storage.Client, ctx context.Context) error {

	in := make(chan int, runtime.NumCPU())
	returnChan := make(chan string)

	err := uw.loadPaths()
	if err != nil {
		return err
	}

	go func() {
		for x := 0; x < len(uw.FilePaths); x++ {
			in <- x
		}
		close(in)
	}()

	for process := 0; process < concurrency; process++ {
		go uw.uploadWorker(in, returnChan, client, ctx)
	}

	uw.verifyReturnChan(returnChan, doneUpload)

	return nil
}

func (uw *UploadWorkersService) verifyReturnChan(returnChan chan string, doneUpload chan string) {
	for r := range returnChan {
		if r != "" {
			doneUpload <- r
			break
		}
	}
}

func (uw *UploadWorkersService) uploadWorker(in chan int, returnChan chan string, client *storage.Client, ctx context.Context) {
	for x := range in {
		err := uw.UploadUseCase.Execute(uw.FilePaths[x], utils.OutputBucketName, client, ctx)

		if err != nil {
			uw.Errors = append(uw.Errors, uw.FilePaths[x])
			log.Printf("error during the upload: %v. Error: %v", uw.FilePaths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- utils.UploadCompleted
}

func (uw *UploadWorkersService) loadPaths() error {
	err := filepath.Walk(uw.VideoPathDir, func(path string, info fs.FileInfo, err error) error {

		if !info.IsDir() {
			uw.FilePaths = append(uw.FilePaths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
