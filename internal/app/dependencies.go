package app

import "github.com/Rumm1/eduhub-backend/internal/config"

type Dependencies struct {
	Config config.Config
}

func NewDependencies(cfg config.Config) (*Dependencies, error) {
	return &Dependencies{Config: cfg}, nil
}

func (d *Dependencies) Close() error {
	return nil
}
