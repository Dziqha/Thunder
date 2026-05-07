package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/config"
	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/repository"
	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/response"
	"github.com/Dziqha/Thunder/benchmarks/fixture-medium-service/internal/service"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	cfg := config.Load()
	repo := repository.NewServiceRepository()
	healthService := service.NewHealthService(cfg, repo)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h := healthService.Check()
		_ = json.NewEncoder(w).Encode(response.Health{
			Status:      h.Status,
			Service:     h.Service,
			Environment: h.Environment,
			Region:      h.Region,
			BuildNumber: h.BuildNumber,
			Dependency: map[string]any{
				"name":   h.Dependency.Name,
				"status": h.Dependency.Status,
				"region": h.Dependency.Region,
			},
			Timestamp: h.Timestamp,
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fixture-medium-service"))
	})

	return mux
}
