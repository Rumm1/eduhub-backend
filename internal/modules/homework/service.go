package homework

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	if repository == nil {
		repository = NewRepository()
	}
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]Entity, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (Entity, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (Entity, error) {
	return s.repository.Create(ctx, request)
}

func (s *Service) Update(ctx context.Context, id string, request UpdateRequest) (Entity, error) {
	return s.repository.Update(ctx, id, request)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
