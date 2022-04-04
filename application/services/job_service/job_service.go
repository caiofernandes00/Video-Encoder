package job_service

import (
	"context"
	"encoder/application/repositories"
	"encoder/application/services/download_service"
	"encoder/application/services/upload_service"
	"encoder/application/utils"
	"encoder/domain"
	"errors"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
)

type JobUseCase interface {
	Start() error
	Insert() error
}

type JobService struct {
	Job             *domain.Job
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
	JobRepository   repositories.JobRepository

	DownloadUseCase        download_service.DownloadUseCase
	FragmentUseCase        download_service.FragmentUseCase
	EncodeUseCase          download_service.EncodeUseCase
	RemoveTempFilesUseCase download_service.RemoveTempFilesUseCase
	UploadWorkersUseCase   upload_service.UploadWorkersUseCase
}

func NewJobService(
	job *domain.Job,
	video *domain.Video,
	VideoRepository repositories.VideoRepository,
	JobRepository repositories.JobRepository,
) *JobService {
	return &JobService{
		Job:             job,
		Video:           video,
		VideoRepository: VideoRepository,
		JobRepository:   JobRepository,
	}
}

func (j *JobService) Start() error {
	j.getServices()

	client, ctx, err := utils.GetClientStorage()
	if err != nil {
		return err
	}

	mp4TargetFile := os.Getenv(utils.LocalStoragePath) + "/" + j.Video.ID + utils.Mp4Format
	encodeTargetDir := os.Getenv(utils.LocalStoragePath) + "/" + j.Video.ID
	fragTargetFile := os.Getenv(utils.LocalStoragePath) + "/" + j.Video.ID + utils.FragmentCommand

	err = j.changeJobStatus(domain.StatusDownloading)
	if err != nil {
		return j.failJob(err)
	}
	err = j.DownloadUseCase.Execute(os.Getenv(utils.InputBucketName), mp4TargetFile, client, ctx)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.StatusFragmenting)
	if err != nil {
		return j.failJob(err)
	}
	err = j.FragmentUseCase.Execute(mp4TargetFile, fragTargetFile)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.StatusEncoding)
	if err != nil {
		return j.failJob(err)
	}
	err = j.EncodeUseCase.Execute(encodeTargetDir)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.StatusUploading)
	if err != nil {
		return j.failJob(err)
	}
	err = j.performUpload(encodeTargetDir, client, ctx)
	if err != nil {
		j.failJob(err)
	}

	err = j.changeJobStatus(domain.StatusRemovingFiles)
	if err != nil {
		return j.failJob(err)
	}
	err = j.RemoveTempFilesUseCase.Execute()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus(domain.StatusFinished)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) Insert() error {
	_, err := j.JobRepository.Insert(j.Job)

	if err != nil {
		return err
	}

	return nil
}

func (j JobService) changeJobStatus(status string) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = domain.StatusFailed
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobService) performUpload(videoPath string, client *storage.Client, ctx context.Context) error {

	err := j.changeJobStatus(domain.StatusUploading)
	if err != nil {
		return j.failJob(err)
	}

	concurrency, _ := strconv.Atoi(os.Getenv(utils.ConcurrencyUpload))
	doneUpload := make(chan string)

	uploadService := upload_service.NewUploadService(j.Video)
	uploadWorkers := upload_service.NewUploadWorkersService(uploadService, videoPath)

	go uploadWorkers.Execute(concurrency, doneUpload, client, ctx)

	var uploadResult string
	uploadResult = <-doneUpload

	if uploadResult != utils.UploadCompleted {
		return j.failJob(errors.New(uploadResult))
	}

	return err
}

func (j *JobService) getServices() {
	j.DownloadUseCase = download_service.NewDownloadService(j.Video)
	j.FragmentUseCase = download_service.NewFragmentService(j.Video)
	j.EncodeUseCase = download_service.NewEncodeService(j.Video)
	j.RemoveTempFilesUseCase = download_service.NewRemoveTempFilesService(j.Video)
}
