package domain_test

import (
	"encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidateValidJob(t *testing.T) {
	video := prepareValidVideo()

	job, err := domain.NewJob(
		"fake-output",
		"fake-status",
		video,
	)

	require.NotNil(t, job)
	require.Nil(t, err)
}

func Test_ValidateEmptyJob(t *testing.T) {
	video := prepareValidVideo()

	_, err := domain.NewJob(
		"",
		"",
		video,
	)

	require.Error(t, err)
}

func prepareValidVideo() *domain.Video {
	video, _ := domain.NewVideo("fake-resource", "fake-filepath")
	return video
}
