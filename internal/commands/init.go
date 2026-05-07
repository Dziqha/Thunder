package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Dziqha/Thunder/internal/config"
	"github.com/Dziqha/Thunder/internal/utils"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func Init() error {
	fmt.Printf("%s⚡ Initializing Thunder...%s\n\n", colorCyan, colorReset)

	if _, err := os.Stat("thunder.toml"); err == nil {
		fmt.Printf("%s⚠ Thunder is already initialized%s\n", colorYellow, colorReset)
		return nil
	}

	mainPath, err := utils.DetectMainFile(".")
	if err != nil {
		mainPath = "main.go"
		fmt.Printf("%s⚠ No main.go found automatically, using default: %s%s\n", colorYellow, mainPath, colorReset)
	} else {
		rel, _ := filepath.Rel(".", mainPath)
		mainPath = filepath.ToSlash(rel)
		fmt.Printf("%s✓ Detected entry point: %s%s\n", colorGreen, mainPath, colorReset)
	}

	defaults := config.Default()

	thunderConfig := fmt.Sprintf(`# Thunder Configuration
# Fast hot reload for Go applications

# Build settings
build_path = "%s"
main_file = "%s"

# Watch settings
watch_dirs = ["."]
exclude_dirs = ["tmp", "vendor", ".git", "node_modules", ".idea", "bin"]
watch_exts = [".go", ".mod", ".sum", ".env", ".tmpl", ".tpl", ".yaml", ".yml", ".json"]
watch_files = ["go.mod", "go.sum", "thunder.toml", ".env"]

# Build and run args
build_args = []
run_args = []

# Debounce time
debounce = 100

# Optional orchestration settings
[project]
name = "myapp"
default_profile = "dev"
log_format = "text"

# Example service definitions
# [services.api]
# type = "go"
# package = "./cmd/api"
# depends_on = ["redis"]
# env_files = [".env"]
# restart_policy = "always"
# max_restarts = 3
# depends_timeout_ms = 15000
# [services.api.restart_backoff]
# strategy = "exponential"
# base_ms = 500
# max_ms = 10000
# jitter_pct = 20
# [services.api.healthcheck]
# type = "http"
# url = "http://localhost:8080/health"
# interval_ms = 500
# timeout_ms = 3000
# retries = 20

# [services.redis]
# type = "process"
# command = ["docker", "compose", "up", "redis"]

# [profiles.dev]
# services = ["api", "redis"]
`, defaults.BuildPath, mainPath)

	if err := os.WriteFile("thunder.toml", []byte(thunderConfig), 0644); err != nil {
		return fmt.Errorf("failed to create thunder.toml: %v", err)
	}

	fmt.Printf("%s✓ Created thunder.toml%s\n", colorGreen, colorReset)

	if err := os.MkdirAll("tmp", 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %v", err)
	}
	fmt.Printf("%s✓ Created tmp/ directory%s\n", colorGreen, colorReset)

	fmt.Printf("\n%s🎉 Thunder initialized successfully!%s\n\n", colorGreen, colorReset)
	fmt.Printf("Next steps:\n  1. go mod init yourmodule (if not done)\n  2. thunder run\n\n")

	return nil
}
