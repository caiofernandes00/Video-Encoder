package domain_test

import (
	"encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidateValidVideo(t *testing.T) {
	video, _ := domain.NewVideo("fake-id", "fake-path")

	err := video.Validate()
	require.Nil(t, err)
}

func Test_ValidateEmptyVideo(t *testing.T) {
	_, err := domain.NewVideo("", "")

	require.Error(t, err)
}
