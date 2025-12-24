package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, options minio.PutObjectOptions) error
	Download(ctx context.Context, key string, options minio.GetObjectOptions) (io.ReadCloser, error)
	Delete(ctx context.Context, key string, options minio.RemoveObjectOptions) error
}
