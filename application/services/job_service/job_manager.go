package job_service

import (
	"encoder/application/repositories"
	"encoder/application/utils"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManagerUseCase interface {
}

type JobManagerService struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerService
	RabbitMQ         *queue.RabbitMQ
	VideoRepository  repositories.VideoRepository
	JobRepository    repositories.JobRepository
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManagerService(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerService, messageChannel chan amqp.Delivery, videoRepository repositories.VideoRepository, jobRepository repositories.JobRepository) *JobManagerService {
	return &JobManagerService{
		Db:               db,
		Domain:           domain.Job{},
		RabbitMQ:         rabbitMQ,
		JobReturnChannel: jobReturnChannel,
		MessageChannel:   messageChannel,
		VideoRepository:  videoRepository,
		JobRepository:    jobRepository,
	}
}

func (j *JobManagerService) Start(ch *amqp.Channel) {

	concurrency, err := strconv.Atoi(os.Getenv(utils.ConcurrencyWorkers))
	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS.")
	}

	jobWorker := NewJobWorkerService(j.VideoRepository, j.JobRepository)
	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go jobWorker.Execute(j.MessageChannel, j.JobReturnChannel)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManagerService) notifySuccess(jobResult JobWorkerService, ch *amqp.Channel) error {

	// Mutex.Lock()
	jobJson, err := json.Marshal(jobResult.Job)
	// Mutex.Unlock()

	if err != nil {
		return err
	}

	err = j.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobManagerService) checkParseErrors(jobResult JobWorkerService) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error)
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManagerService) notify(jobJson []byte) error {

	err := j.RabbitMQ.Notify(
		string(jobJson),
		utils.ContentTypeJson,
		os.Getenv(utils.RabbitMQNotificationEx),
		os.Getenv(utils.RabbitMQNotificationRoutingKey),
	)

	if err != nil {
		return err
	}

	return nil
}
