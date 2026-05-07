package commands

import (
	"fmt"
	"os"

	"github.com/Dziqha/Thunder/internal/config"
	"github.com/Dziqha/Thunder/internal/diagnostics"
)

func Doctor() error {
	cfg, err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			cfg = config.Default()
		} else {
			return fmt.Errorf("failed to load config: %v", err)
		}
	}

	issues := diagnostics.Doctor(cfg)
	fmt.Printf("%s🩺 Thunder Doctor%s\n", colorCyan, colorReset)
	for _, issue := range issues {
		switch issue.Level {
		case "ok":
			fmt.Printf("%s✓ %s%s\n", colorGreen, issue.Message, colorReset)
		case "warn":
			fmt.Printf("%s⚠ %s%s\n", colorYellow, issue.Message, colorReset)
		default:
			fmt.Printf("%s✗ %s%s\n", colorRed, issue.Message, colorReset)
		}
	}

	return nil
}
