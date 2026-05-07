package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Dziqha/Thunder/internal/config"
)

func Inspect() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	profile := ""
	if len(os.Args) > 2 {
		profile = os.Args[2]
	}

	services, err := cfg.ServicesForProfile(profile)
	if err != nil {
		return err
	}
	sort.Strings(services)

	fmt.Printf("%s🔎 Thunder Inspect%s\n", colorCyan, colorReset)
	fmt.Printf("Project: %s\n", cfg.Project.Name)
	fmt.Printf("Default profile: %s\n", cfg.Project.DefaultProfile)
	if profile != "" {
		fmt.Printf("Selected profile: %s\n", profile)
	}

	fmt.Println("Services:")
	for _, name := range services {
		svc := cfg.Services[name]
		t := strings.TrimSpace(svc.Type)
		if t == "" {
			t = "go"
		}
		fmt.Printf("  - %s (%s)", name, t)
		if svc.Package != "" {
			fmt.Printf(" pkg=%s", svc.Package)
		}
		if len(svc.DependsOn) > 0 {
			fmt.Printf(" depends_on=%v", svc.DependsOn)
		}
		fmt.Println()
	}

	return nil
}
