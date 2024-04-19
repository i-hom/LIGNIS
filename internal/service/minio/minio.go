package minio

import (
	"context"
	"io"

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
	isExist, _ := client.BucketExists(context.Background(), config.Bucket)

	if !isExist {
		_ = client.MakeBucket(context.Background(), config.Bucket, minio.MakeBucketOptions{})
	}

	return &MinioStorage{
		client: client,
		bucket: config.Bucket,
	}, nil
}

func (storage MinioStorage) Upload(ctx context.Context, objName string, reader io.Reader, objSize int64) (minio.UploadInfo, error) {
	return storage.client.PutObject(ctx, storage.bucket, objName, reader, objSize, minio.PutObjectOptions{})
}

func (storage MinioStorage) Download(ctx context.Context, objName string) (*minio.Object, error) {
	return storage.client.GetObject(ctx, storage.bucket, objName, minio.GetObjectOptions{})
}
