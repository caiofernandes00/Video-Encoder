package domain_test

import (
	"encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidateValidVideo(t *testing.T) {
	video, _ := domain.NewVideo("fake-id", "fake-path", uuid.NewV4().String())

	err := video.Validate()
	require.Nil(t, err)
}

func Test_ValidateEmptyVideo(t *testing.T) {
	_, err := domain.NewVideo("", "", "")

	require.Error(t, err)
}

func Test_ValidateInvalidVideoID(t *testing.T) {
	_, err := domain.NewVideo("fake-id", "fake-path", "invalid-id")

	require.Error(t, err)
}
