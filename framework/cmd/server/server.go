package server

import (
	"encoder/application/repositories"
	"encoder/application/services/job_service"
	"encoder/application/utils"
	"encoder/framework/database"
	"encoder/framework/queue"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv(utils.AutoMigrateDb))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	debug, err := strconv.ParseBool(os.Getenv(utils.Debug))
	if err != nil {
		log.Fatalf("Error parsing boolean env var")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv(utils.DsnTest)
	db.Dsn = os.Getenv(utils.Dsn)
	db.DbTypeTest = os.Getenv(utils.DbTypeTest)
	db.DbType = os.Getenv(utils.DbType)
	db.Env = os.Getenv(utils.Env)
}

func main() {

	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan job_service.JobWorkerService)

	dbConnection, err := db.Connect()

	if err != nil {
		log.Fatalf("error connecting to DB")
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	videoRepository := repositories.NewVideoRepository(dbConnection)
	jobRepository := repositories.NewJobRepository(dbConnection)

	jobManager := job_service.NewJobManagerService(dbConnection, rabbitMQ, jobReturnChannel, messageChannel, videoRepository, jobRepository)
	jobManager.Start(ch)

}
