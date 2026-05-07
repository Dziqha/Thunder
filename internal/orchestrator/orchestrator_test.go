package orchestrator

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Dziqha/Thunder/internal/config"
)

func TestResolveOrder(t *testing.T) {
	o := &Orchestrator{cfg: &config.Config{Services: config.ServicesMap{
		"redis": {Type: "process", Command: []string{"redis-server"}},
		"api":   {Type: "go", Package: "./cmd/api", DependsOn: []string{"redis"}},
	}}}

	ordered, err := o.resolveOrder([]string{"api", "redis"})
	if err != nil {
		t.Fatalf("resolveOrder error: %v", err)
	}

	if len(ordered) != 2 {
		t.Fatalf("expected 2 services, got %d", len(ordered))
	}
	if ordered[0] != "redis" || ordered[1] != "api" {
		t.Fatalf("unexpected order: %v", ordered)
	}
}

func TestRestartDelayRange(t *testing.T) {
	svc := config.ServiceConfig{
		RestartBackoff: config.RestartBackoffConfig{Strategy: "fixed", BaseMS: 500, MaxMS: 500, JitterPct: 0},
	}
	d := restartDelay(svc, 3)
	if d != 500*time.Millisecond {
		t.Fatalf("expected 500ms, got %v", d)
	}
}

func TestRedactSecrets(t *testing.T) {
	line := "user=alice token=abc123 PASSWORD=qwerty normal=value"
	out := redactSecrets(line)
	if out == line {
		t.Fatal("expected redaction to modify line")
	}
	if strings.Contains(out, "abc123") || strings.Contains(out, "qwerty") {
		t.Fatalf("secret leaked in output: %s", out)
	}
}

func TestWaitForDependenciesHealthyTimeout(t *testing.T) {
	o := &Orchestrator{
		cfg: &config.Config{Services: config.ServicesMap{
			"redis": {Type: "process", Command: []string{"redis-server"}},
			"api":   {Type: "go", Package: "./cmd/api", DependsOn: []string{"redis"}, DependsTimeoutMS: 50},
		}},
		healthy: make(map[string]bool),
	}

	err := o.waitForDependenciesHealthy(context.Background(), "api")
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestEmitNonBlocking(t *testing.T) {
	o := &Orchestrator{events: make(chan Event, 1)}
	for i := 0; i < 10; i++ {
		o.emit("api", "test", "msg")
	}
}

func TestWaitForHealthHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	err := waitForHealth(context.Background(), config.HealthcheckConfig{
		Type:       "http",
		URL:        srv.URL,
		IntervalMS: 10,
		TimeoutMS:  500,
		Retries:    2,
	})
	if err != nil {
		t.Fatalf("expected healthcheck success, got: %v", err)
	}
}

func TestWaitForHealthTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	err = waitForHealth(context.Background(), config.HealthcheckConfig{
		Type:       "tcp",
		Host:       "127.0.0.1",
		Port:       addr.Port,
		IntervalMS: 10,
		TimeoutMS:  500,
		Retries:    2,
	})
	if err != nil {
		t.Fatalf("expected tcp healthcheck success, got: %v", err)
	}
}

func TestStartProfileProcessStartFailure(t *testing.T) {
	cfg := &config.Config{
		Services: config.ServicesMap{
			"bad": {
				Type:          "process",
				Command:       []string{"this-command-does-not-exist-xyz"},
				RestartPolicy: "never",
			},
		},
		Profiles: config.ProfilesMap{
			"dev": {Services: []string{"bad"}},
		},
		Project:  config.ProjectConfig{DefaultProfile: "dev"},
		Debounce: 100,
	}
	if err := cfg.NormalizeAndValidate(); err != nil {
		t.Fatalf("config normalize failed: %v", err)
	}

	o, err := New(cfg)
	if err != nil {
		t.Fatalf("new orchestrator failed: %v", err)
	}

	err = o.StartProfile(context.Background(), "dev")
	if err == nil {
		t.Fatal("expected start profile failure")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "start failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
