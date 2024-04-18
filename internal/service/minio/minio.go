package minio

import "github.com/minio/minio-go/v7"

type MinioStorage struct {
	client *minio.Client
	bucket string
}

func NewMinio(config *Config) (*MinioStorage, error) {
	client, err := minio.New()
	if err != nil {
		return nil, err
	}
	return &MinioStorage{
		client: client,
		bucket: config.Bucket,
	}, nil
}
