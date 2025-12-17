package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucketName, key, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *MinioStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
}

func (s *MinioStorage) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
}
