package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	BucketName = "video-encoder-test-bkt"
)

func Test_DownloadManagerFacadeService(t *testing.T) {
	video, repo := prepare()

	downloadManagerFacadeService := services.NewDownloadManagerFacadeService()
	err = downloadManagerFacadeService.Execute()
	require.Nil(t, err)
}

func prepare() (*domain.Video, *repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()

	video := generateVideoObject()
	videoRepo := repositories.NewVideoRepository(db)
	videoRepo.Insert(video)

	return video, videoRepo
}

func generateVideoObject() *domain.Video {
	return &domain.Video{
		ID:        uuid.NewV4().String(),
		FilePath:  "convite.mp4",
		CreatedAt: time.Now(),
	}
}
