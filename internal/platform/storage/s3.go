package storage

import (
	"context"
	"errors"
	"io"
)

type S3Config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type S3Store struct {
	config S3Config
}

func NewS3Store(config S3Config) *S3Store {
	return &S3Store{config: config}
}

func (s *S3Store) Put(_ context.Context, _ Object, _ io.Reader) error {
	return errors.New("s3 storage is not configured")
}

func (s *S3Store) Get(_ context.Context, _ string) (io.ReadCloser, Object, error) {
	return nil, Object{}, errors.New("s3 storage is not configured")
}

func (s *S3Store) Delete(_ context.Context, _ string) error {
	return errors.New("s3 storage is not configured")
}
