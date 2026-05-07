package repository

import "github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/domain"

type ServiceRepository struct{}

func NewServiceRepository() ServiceRepository {
	return ServiceRepository{}
}

func (r ServiceRepository) Get() domain.Service {
	return domain.Service{
		Name:   "fixture-medium-service",
		Status: "ok",
		Region: "local",
	}
}
