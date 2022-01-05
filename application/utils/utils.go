package utils

import (
	"cloud.google.com/go/storage"
	"context"
	"log"
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
