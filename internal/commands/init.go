// internal/commands/init.go
package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Dziqha/Thunder/internal/utils"
)

const (
	colorReset  = "\033[0m"
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
		mainPath = rel
		fmt.Printf("%s✓ Detected entry point: %s%s\n", colorGreen, mainPath, colorReset)
	}

	thunderConfig := fmt.Sprintf(`# Thunder Configuration
# Fast hot reload for Go applications

# Build settings
build_path = "./tmp/main.exe"
main_file = "%s"

# Watch settings
watch_dirs = ["."]
exclude_dirs = ["tmp", "vendor", ".git", "node_modules", ".idea", "bin"]

# Debounce time
debounce = 100
`, mainPath)

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