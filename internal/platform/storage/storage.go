package storage

import (
	"context"
	"io"
)

type Object struct {
	Key         string
	ContentType string
	Size        int64
}

type Store interface {
	Put(ctx context.Context, object Object, body io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, Object, error)
	Delete(ctx context.Context, key string) error
}
