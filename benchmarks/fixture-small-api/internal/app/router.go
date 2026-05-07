package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dziqha/Thunder/benchmarks/fixture-small-api/internal/config"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	cfg := config.Load()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "ok",
			"service":   cfg.ServiceName,
			"env":       cfg.Environment,
			"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		})
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fixture-small-api"))
	})

	return mux
}
