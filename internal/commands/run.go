// internal/commands/run.go
package commands

import (
	"fmt"
	"os"

	"github.com/Dziqha/thunder/internal/config"
	"github.com/Dziqha/thunder/internal/watcher"
)

func Run() error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		// If no config found, use defaults
		if os.IsNotExist(err) {
			fmt.Printf("%s⚠ No thunder.toml found, using default configuration%s\n", colorYellow, colorReset)
			fmt.Printf("%s💡 Run 'thunder init' to create a configuration file%s\n\n", colorYellow, colorReset)
			cfg = config.Default()
		} else {
			return fmt.Errorf("failed to load config: %v", err)
		}
	}

	// Override main file if specified in args
	if len(os.Args) > 2 {
		cfg.MainFile = os.Args[2]
	}

	// Create and start watcher
	w, err := watcher.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}
	defer w.Close()

	fmt.Printf("%s⚡ Thunder is watching for changes...%s\n", colorYellow, colorReset)
	fmt.Printf("%s💡 Press Ctrl+C to stop%s\n\n", colorCyan, colorReset)

	return w.Start()
}