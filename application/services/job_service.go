package services

import (
	"encoder/application/repositories"
	"encoder/application/services/download_service"
	"encoder/application/services/upload_service"
	"encoder/domain"
	"os"
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
	UploadWorkersUseCase   upload_service.UploadWorkersUseCase
}

const (
	FailedStatus        = "FAILED"
	DownloadingStatus   = "DOWNLOADING"
	FragmentingStatus   = "FRAGMENTING"
	EncodingStatus      = "ENCODING"
	RemovingFilesStatus = "REMOVING_REMAINING_STATUS"
	UploadingStatus     = "UPLOADING"
)

func NewJobService(
	job *domain.Job,
	jobRepository repositories.JobRepository,
	downloadUseCase download_service.DownloadUseCase,
	fragmentUseCase download_service.FragmentUseCase,
	encodeUseCase download_service.EncodeUseCase,
	removeTempFilesUseCase download_service.RemoveTempFilesUseCase,
	uploadWorkersUseCase upload_service.UploadWorkersUseCase,
) *JobService {
	return &JobService{
		Job:                    job,
		JobRepository:          jobRepository,
		DownloadUseCase:        downloadUseCase,
		FragmentUseCase:        fragmentUseCase,
		EncodeUseCase:          encodeUseCase,
		RemoveTempFilesUseCase: removeTempFilesUseCase,
		UploadWorkersUseCase:   uploadWorkersUseCase,
	}
}

func (j *JobService) Start() error {
	var err error

	err = j.changeJobStatus(DownloadingStatus)

	if err != nil {
		return j.failJob(err)
	}

	err = j.DownloadUseCase.Execute(os.Getenv("BUCKET_NAME"))
	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus(FragmentingStatus)
	if err != nil {
		return j.failJob(err)
	}

	err = j.FragmentUseCase.Execute()
	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus(FragmentingStatus)
	if err != nil {
		return j.failJob(err)
	}

	err = j.EncodeUseCase.Execute()
	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus(EncodingStatus)
	if err != nil {
		return j.failJob(err)
	}

	err = j.RemoveTempFilesUseCase.Execute()
	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus(RemovingFilesStatus)
	if err != nil {
		return j.failJob(err)
	}

	//err = j.UploadWorkersUseCase.Execute()
	//if err != nil {
	//	return j.failJob(err)
	//}
	//err = j.changeJobStatus(UploadingStatus)
	//if err != nil {
	//	return j.failJob(err)
	//}

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

func (j *JobService) performUpload() error {

	err := j.changeJobStatus(UploadingStatus)
	if err != nil {
		return j.failJob(err)
	}

}
