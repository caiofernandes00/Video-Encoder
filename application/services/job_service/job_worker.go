package job_service

import (
	"encoder/application/repositories"
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
	Job     *domain.Job
	Video   *domain.Video
	Error   error
	Message *amqp.Delivery

	VideoRepository repositories.VideoRepository
	JobRepository   repositories.JobRepository
}

func NewJobWorkerService(videoRepository repositories.VideoRepository, jobRepository repositories.JobRepository) *JobWorkerService {
	return &JobWorkerService{
		VideoRepository: videoRepository,
		JobRepository:   jobRepository,
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

		videoUseCase := jw.createVideoService()
		err = videoUseCase.InsertVideo()
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		err = jw.createJob(message)
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		jobUseCase := jw.createJobService()
		err = jobUseCase.Insert()
		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

		err = jobUseCase.Start()
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

func (jw *JobWorkerService) createVideoService() video_service.VideoUseCase {
	return video_service.NewVideoService(jw.Video, jw.VideoRepository)
}

func (jw *JobWorkerService) createJobService() JobUseCase {
	return NewJobService(jw.Job, jw.Video, jw.VideoRepository, jw.JobRepository)
}

func returnJobResult(job *domain.Job, message *amqp.Delivery, err error) JobWorkerService {
	return JobWorkerService{
		Job:     job,
		Message: message,
		Error:   err,
	}
}
