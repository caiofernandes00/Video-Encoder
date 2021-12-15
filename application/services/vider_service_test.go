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

func Test_VideoServiceDownload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService(video, repo)
	err := videoService.Download(BucketName)
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
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
