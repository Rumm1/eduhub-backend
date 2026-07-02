package platformdashboard

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetDashboard(ctx context.Context) (DashboardResponse, error) {
	return s.repository.GetDashboard(ctx)
}
