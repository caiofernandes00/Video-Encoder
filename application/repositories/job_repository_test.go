package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_JobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := generateNewValidVideo()
	videoRepo := repositories.NewVideoRepository(db)
	videoRepo.Insert(video)

	newJob := generateValidJob(video)
	jobRepo := repositories.NewJobRepository(db)
	jobRepo.Insert(newJob)

	job, err := jobRepo.Find(newJob.ID)

	require.NotEmpty(t, job.ID)
	require.Nil(t, err)
	require.Equal(t, job.ID, newJob.ID)
	require.Equal(t, job.VideoID, video.ID, newJob.VideoID)
}

func Test_JobRepositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := generateNewValidVideo()
	videoRepo := repositories.NewVideoRepository(db)
	videoRepo.Insert(video)

	newJob := generateValidJob(video)
	jobRepo := repositories.NewJobRepository(db)
	jobRepo.Insert(newJob)

	jobToUpdate := generateValidJob2(video)
	jobRepo.Update(jobToUpdate)

	job, err := jobRepo.Find(newJob.ID)

	require.NotEmpty(t, job.ID)
	require.Nil(t, err)
	require.Equal(t, "another-output", jobToUpdate.OutputBucketPath)
	require.Equal(t, "another-status", jobToUpdate.Status)
}

func generateValidJob(video *domain.Video) *domain.Job {
	job, _ := domain.NewJob("output", "status", video)
	return job
}

func generateValidJob2(video *domain.Video) *domain.Job {
	job, _ := domain.NewJob("another-output", "another-status", video)
	return job
}
