package job_service

import (
	"encoder/application/utils"
	"encoder/domain"
	"github.com/streadway/amqp"
)

type JobWorkerUseCase interface {
}

type JobWorkerService struct {
	Job        *domain.Job
	Message    *amqp.Delivery
	Error      error
	JobUseCase JobUseCase
}

func NewJobWorkerService(job *domain.Job, jobUseCase JobUseCase, message *amqp.Delivery) *JobWorkerService {
	return &JobWorkerService{Job: job, JobUseCase: jobUseCase, Message: message}
}

func JobWorker(messageChan chan amqp.Delivery, returnChan chan JobWorkerService) {
	for message := range messageChan {
		err := utils.IsJson(string(message.Body))

		if err != nil {
			returnChan <- returnJobResult(&domain.Job{}, &message, err)
			continue
		}

	}
}

func returnJobResult(job *domain.Job, message *amqp.Delivery, err error) JobWorkerService {
	return JobWorkerService{
		Job:     job,
		Message: message,
		Error:   err,
	}
}
