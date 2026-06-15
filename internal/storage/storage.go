package storage

import "io"

type ObjectInfo struct {
	Key  string
	Size int64
}

type Storage interface {
	PutObject(bucket, key string, reader io.Reader, size int64) error
	GetObject(bucket, key string) (io.ReadCloser, error)
	ListObjects(bucket, prefix string) ([]ObjectInfo, error)
	DeleteObject(bucket, key string) error
}