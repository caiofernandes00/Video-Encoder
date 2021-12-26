package services

import (
	"encoder/application/repositories"
	"encoder/application/services/download_service"
	"encoder/application/services/upload_service"
	"encoder/domain"
)

type JobUseCase interface {
	Start() error
}

type JobService struct {
	Job                    *domain.Job
	JobRepository          repositories.JobRepository
	DownloadUseCase        download_service.DownloadUseCase
	FragmentUseCase        download_service.FragmentUseCase
	EncodeUseCase          download_service.EncodeUseCase
	RemoveTempFilesUseCase download_service.RemoveTempFilesUseCase
	UploadUseCase          upload_service.UploadUseCase
	UploadWorkersUseCase   upload_service.UploadWorkersUseCase
}

const (
	FailedStatus = "FAILED"
	Downloading  = "DOWNLOADING"
)

func NewJobService(
	job *domain.Job,
	jobRepository repositories.JobRepository,
	downloadUseCase download_service.DownloadUseCase,
	fragmentUseCase download_service.FragmentUseCase,
	encodeUseCase download_service.EncodeUseCase,
	removeTempFilesUseCase download_service.RemoveTempFilesUseCase,
	uploadUseCase upload_service.UploadUseCase,
	uploadWorkersUseCase upload_service.UploadWorkersUseCase,
) *JobService {
	return &JobService{
		Job:                    job,
		JobRepository:          jobRepository,
		DownloadUseCase:        downloadUseCase,
		FragmentUseCase:        fragmentUseCase,
		EncodeUseCase:          encodeUseCase,
		RemoveTempFilesUseCase: removeTempFilesUseCase,
		UploadUseCase:          uploadUseCase,
		UploadWorkersUseCase:   uploadWorkersUseCase,
	}
}

func (j *JobService) Start() error {
	var err error

	err = j.changeJobStatus(Downloading)

	if err != nil {
		return j.failJob(err)
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
	j.Job.Status = FailedStatus
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return nil
}
