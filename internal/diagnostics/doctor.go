package diagnostics

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Dziqha/Thunder/internal/config"
)

type Issue struct {
	Level   string
	Message string
}

func Doctor(cfg *config.Config) []Issue {
	issues := make([]Issue, 0)

	if cfg.Debounce < 50 {
		issues = append(issues, Issue{Level: "warn", Message: "debounce is too low; recommended >= 50ms"})
	}

	buildDir := filepath.Dir(cfg.BuildPath)
	if st, err := os.Stat(buildDir); err != nil || !st.IsDir() {
		issues = append(issues, Issue{Level: "warn", Message: fmt.Sprintf("build directory %q does not exist yet", buildDir)})
	}

	if len(cfg.WatchDirs) == 0 {
		issues = append(issues, Issue{Level: "error", Message: "watch_dirs is empty"})
	}

	if cfg.Project.LogFormat != "" && cfg.Project.LogFormat != "text" && cfg.Project.LogFormat != "json" {
		issues = append(issues, Issue{Level: "error", Message: "project.log_format must be text or json"})
	}

	for name, svc := range cfg.Services {
		for _, envFile := range svc.EnvFiles {
			if _, err := os.Stat(envFile); err != nil {
				issues = append(issues, Issue{Level: "warn", Message: fmt.Sprintf("service %q env file not found: %s", name, envFile)})
			}
		}

		if svc.Type == "process" && len(svc.Command) > 0 {
			if _, err := execLookPathSafe(svc.Command[0]); err != nil {
				issues = append(issues, Issue{Level: "warn", Message: fmt.Sprintf("service %q command not found in PATH: %s", name, svc.Command[0])})
			}
		}

		if svc.DependsTimeoutMS < 0 {
			issues = append(issues, Issue{Level: "error", Message: fmt.Sprintf("service %q depends_timeout_ms must be >= 0", name)})
		}

		ht := svc.Healthcheck
		switch ht.Type {
		case "", "none":
		case "http":
			if ht.URL == "" {
				issues = append(issues, Issue{Level: "error", Message: fmt.Sprintf("service %q healthcheck.type=http requires healthcheck.url", name)})
			}
		case "tcp":
			if ht.Host == "" || ht.Port <= 0 {
				issues = append(issues, Issue{Level: "error", Message: fmt.Sprintf("service %q healthcheck.type=tcp requires host and port", name)})
			}
		default:
			issues = append(issues, Issue{Level: "error", Message: fmt.Sprintf("service %q has unsupported healthcheck type: %s", name, ht.Type)})
		}
	}

	if len(issues) == 0 {
		issues = append(issues, Issue{Level: "ok", Message: "no critical issues found"})
	}

	return issues
}

func execLookPathSafe(bin string) (string, error) {
	return exec.LookPath(bin)
}
