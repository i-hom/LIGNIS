package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinio(config *Config) (*MinioStorage, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Username, config.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &MinioStorage{
		client: client,
		bucket: config.Bucket,
	}, nil
}
