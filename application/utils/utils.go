package utils

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/storage"
)

func PrintOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s", string(out))
	}
}

func GetClientStorage() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, err
}

func IsJson(s string) error {
	var js struct{}

	if err := json.Unmarshal([]byte(s), &js); err != nil {
		return err
	}

	return nil
}
