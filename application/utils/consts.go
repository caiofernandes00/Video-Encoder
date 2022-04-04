package utils

const (
	LocalStoragePath = "LOCAL_STORAGE_PATH"
	FragmentCommand  = "mp4fragment"
	EncodeCommand    = "mp4dash"
	Mp4Format        = ".mp4"
	FragFile         = ".frag"
	UploadCompleted  = "upload completed"
	ContentTypeJson  = "application/json"

	// Env variables
	InputBucketName                = "INPUT_BUCKET_NAME"
	OutputBucketName               = "OUTPUT_BUCKET_NAME"
	ConcurrencyUpload              = "CONCURRENCY_UPLOAD"
	ConcurrencyWorkers             = "CONCURRENCY_WORKERS"
	RabbitMQNotificationEx         = "RABBITMQ_NOTIFICATION_EX"
	RabbitMQNotificationRoutingKey = "RABBITMQ_NOTIFICATION_ROUTING_KEY"
	AutoMigrateDb                  = "AUTO_MIGRATE_DB"
	Debug                          = "DEBUG"
	Env                            = "ENV"
	DbType                         = "DB_TYPE"
	DbTypeTest                     = "DB_TYPE_TEST"
	Dsn                            = "DSN"
	DsnTest                        = "DSN_TEST"
)
