package service

import (
	"time"

	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/config"
	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/domain"
	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/repository"
)

type HealthService struct {
	cfg  config.Config
	repo repository.ServiceRepository
}

type HealthResponse struct {
	Status      string         `json:"status"`
	Service     string         `json:"service"`
	Environment string         `json:"environment"`
	Region      string         `json:"region"`
	BuildNumber int            `json:"build_number"`
	Dependency  domain.Service `json:"dependency"`
	Timestamp   string         `json:"timestamp"`
}

func NewHealthService(cfg config.Config, repo repository.ServiceRepository) HealthService {
	return HealthService{cfg: cfg, repo: repo}
}

func (s HealthService) Check() HealthResponse {
	dep := s.repo.Get()

	return HealthResponse{
		Status:      "ok",
		Service:     s.cfg.ServiceName,
		Environment: s.cfg.Environment,
		Region:      s.cfg.Region,
		BuildNumber: s.cfg.BuildNumber,
		Dependency:  dep,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
	}
}
