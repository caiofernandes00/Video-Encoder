package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ValidJson(t *testing.T) {

	json := `{
				"id": "3e542365-7840-4185-8af6-a26f71e5b51d",
				"file_path": "test.mp4",
				"status": "pending"
			}`

	err := IsJson(json)
	require.Nil(t, err)
}

func Test_InvalidJson(t *testing.T) {

	json := `Invalid Json`

	err := IsJson(json)
	require.Error(t, err)
}
