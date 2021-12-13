package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewVideRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := generateNewValidVideo()

	repo := repositories.NewVideoRepository(db)
	repo.Insert(video)
	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}

func generateNewValidVideo() *domain.Video {
	v, _ := domain.NewVideo("resource", "path", uuid.NewV4().String())
	return v
}
