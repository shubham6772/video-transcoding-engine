package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinIOStorage struct {
	Client *minio.Client
}

func (m *MinIOStorage) PutObject(bucket, key string, reader io.Reader, size int64) error {
	_, err := m.Client.PutObject(
		context.Background(),
		bucket,
		key,
		reader,
		size,
		minio.PutObjectOptions{},
	)
	return err
}

func (m *MinIOStorage) GetObject(bucket, key string) (io.ReadCloser, error) {
	return m.Client.GetObject(
		context.Background(),
		bucket,
		key,
		minio.GetObjectOptions{},
	)
}

func (m *MinIOStorage) ListObjects(bucket, prefix string) ([]ObjectInfo, error) {
	ctx := context.Background()

	ch := m.Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var res []ObjectInfo

	for obj := range ch {
		if obj.Err != nil {
			return nil, obj.Err
		}

		res = append(res, ObjectInfo{
			Key:  obj.Key,
			Size: obj.Size,
		})
	}

	return res, nil
}

func (m *MinIOStorage) DeleteObject(bucket, key string) error {
	return m.Client.RemoveObject(
		context.Background(),
		bucket,
		key,
		minio.RemoveObjectOptions{},
	)
}