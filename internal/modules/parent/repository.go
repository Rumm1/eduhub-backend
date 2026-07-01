package parent

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("entity not found")

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) List(ctx context.Context) ([]Entity, error) {
	return []Entity{}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (Entity, error) {
	if id == "" {
		return Entity{}, ErrNotFound
	}
	return Entity{ID: id, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (r *Repository) Create(ctx context.Context, request CreateRequest) (Entity, error) {
	now := time.Now()
	return Entity{ID: now.Format("20060102150405.000000000"), Name: request.Name, CreatedAt: now, UpdatedAt: now}, nil
}

func (r *Repository) Update(ctx context.Context, id string, request UpdateRequest) (Entity, error) {
	if id == "" {
		return Entity{}, ErrNotFound
	}
	now := time.Now()
	return Entity{ID: id, Name: request.Name, CreatedAt: now, UpdatedAt: now}, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrNotFound
	}
	return nil
}
