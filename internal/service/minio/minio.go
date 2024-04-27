package minio

import (
	"bytes"
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

func (storage MinioStorage) Upload(ctx context.Context, objName string, file io.Reader) (minio.UploadInfo, error) {
	buffer, err := io.ReadAll(file)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return storage.client.PutObject(ctx, storage.bucket, objName, bytes.NewReader(buffer), int64(len(buffer)), minio.PutObjectOptions{})
}

func (storage MinioStorage) Download(ctx context.Context, objName string) (*minio.Object, error) {
	return storage.client.GetObject(ctx, storage.bucket, objName, minio.GetObjectOptions{})
}

func (storage MinioStorage) Delete(ctx context.Context, objName string) error {
	return storage.client.RemoveObject(ctx, storage.bucket, objName, minio.RemoveObjectOptions{})
}
