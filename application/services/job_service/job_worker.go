package job_service

import (
	"encoder/application/repositories"
	"encoder/application/services/download_service"
	"encoder/application/services/upload_service"
	"encoder/application/services/video_service"
	"encoder/application/utils"
	"encoder/domain"
	"encoding/json"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerUseCase interface {
	Execute(messageChan chan amqp.Delivery, returnChan chan JobWorkerService)
}

type JobWorkerService struct {
	Job          *domain.Job
	Video        *domain.Video
	VideoUseCase video_service.VideoUseCase
	JobUseCase   JobUseCase
	Error        error

	Message *amqp.Delivery

	VideoRepository repositories.VideoRepository

	JobRepository          repositories.JobRepository
	DownloadUseCase        download_service.DownloadUseCase
	FragmentUseCase        download_service.FragmentUseCase
	EncodeUseCase          download_service.EncodeUseCase
	RemoveTempFilesUseCase download_service.RemoveTempFilesUseCase
	UploadWorkersUseCase   upload_service.UploadWorkersUseCase
}

func NewJobWorkerService(
	message *amqp.Delivery,

	videoRepository repositories.VideoRepository,

	jobRepository repositories.JobRepository,
	downloadUseCase download_service.DownloadUseCase,
	fragmentUseCase download_service.FragmentUseCase,
	encodeUseCase download_service.EncodeUseCase,
	removeTempFilesUseCase download_service.RemoveTempFilesUseCase,
	uploadWorkersUseCase upload_service.UploadWorkersUseCase,
) *JobWorkerService {
	return &JobWorkerService{
		Message: message,

		VideoRepository: videoRepository,

		JobRepository:          jobRepository,
		DownloadUseCase:        downloadUseCase,
		FragmentUseCase:        fragmentUseCase,
		EncodeUseCase:          encodeUseCase,
		RemoveTempFilesUseCase: removeTempFilesUseCase,
		UploadWorkersUseCase:   uploadWorkersUseCase,
	}
}

func (jw *JobWorkerService) Execute(messageChan chan amqp.Delivery, returnChan chan JobWorkerService) {
	for message := range messageChan {
		err := utils.IsJson(string(message.Body))

		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		err = jw.createVideo(message)
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		jw.createVideoService()
		err = jw.VideoUseCase.InsertVideo()
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		err = jw.createJob(message)
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		jw.createJobService()
		err = jw.JobUseCase.Insert()
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		err = jw.JobUseCase.Start()
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		returnChan <- returnJobResult(jw.Job, &message, nil)

	}
}

func (jw *JobWorkerService) createVideo(message amqp.Delivery) error {
	jw.Video = &domain.Video{}
	err := json.Unmarshal(message.Body, jw.Video)
	if err != nil {
		return err
	}
	jw.Video.ID = uuid.NewV4().String()

	err = jw.Video.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (jw *JobWorkerService) createJob(message amqp.Delivery) error {
	job, err := domain.NewJob(os.Getenv(utils.OutputBucketName), domain.StatusFinished, jw.Video)

	if err != nil {
		return err
	}

	jw.Job = job

	return nil
}

func (jw *JobWorkerService) createVideoService() {
	videoUseCase := video_service.NewVideoService(jw.Video, jw.VideoRepository)
	jw.VideoUseCase = videoUseCase
}

func (jw *JobWorkerService) createJobService() {
	jobUseCase := NewJobService(jw.Job, jw.Video, jw.JobRepository, jw.DownloadUseCase, jw.FragmentUseCase, jw.EncodeUseCase, jw.RemoveTempFilesUseCase, jw.UploadWorkersUseCase)
	jw.JobUseCase = jobUseCase
}

func returnJobResult(job *domain.Job, message *amqp.Delivery, err error) JobWorkerService {
	return JobWorkerService{
		Job:     job,
		Message: message,
		Error:   err,
	}
}
