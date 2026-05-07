package commands

import (
	"fmt"
	"os"

	"github.com/Dziqha/Thunder/internal/config"
)

func ReleaseCheck() error {
	failures := 0

	if _, err := os.Stat(".github/workflows/ci.yml"); err != nil {
		fmt.Printf("%s✗ Missing CI workflow (.github/workflows/ci.yml)%s\n", colorRed, colorReset)
		failures++
	} else {
		fmt.Printf("%s✓ CI workflow found%s\n", colorGreen, colorReset)
	}

	if _, err := os.Stat(".github/workflows/release.yml"); err != nil {
		fmt.Printf("%s✗ Missing release workflow (.github/workflows/release.yml)%s\n", colorRed, colorReset)
		failures++
	} else {
		fmt.Printf("%s✓ Release workflow found%s\n", colorGreen, colorReset)
	}

	if _, err := os.Stat("README.md"); err != nil {
		fmt.Printf("%s✗ Missing README.md%s\n", colorRed, colorReset)
		failures++
	} else {
		fmt.Printf("%s✓ README.md found%s\n", colorGreen, colorReset)
	}

	if _, err := os.Stat("LICENSE"); err != nil {
		fmt.Printf("%s✗ Missing LICENSE%s\n", colorRed, colorReset)
		failures++
	} else {
		fmt.Printf("%s✓ LICENSE found%s\n", colorGreen, colorReset)
	}

	if _, err := os.Stat("CHANGELOG.md"); err != nil {
		fmt.Printf("%s✗ Missing CHANGELOG.md%s\n", colorRed, colorReset)
		failures++
	} else {
		fmt.Printf("%s✓ CHANGELOG.md found%s\n", colorGreen, colorReset)
	}

	if _, err := os.Stat("thunder.toml"); err == nil {
		if _, err := config.Load(); err != nil {
			fmt.Printf("%s✗ thunder.toml invalid: %v%s\n", colorRed, err, colorReset)
			failures++
		} else {
			fmt.Printf("%s✓ thunder.toml loads successfully%s\n", colorGreen, colorReset)
		}
	} else {
		fmt.Printf("%s⚠ thunder.toml not found (ok for library repo)%s\n", colorYellow, colorReset)
	}

	if failures > 0 {
		return fmt.Errorf("release check failed with %d issue(s)", failures)
	}

	fmt.Printf("%s🎉 Release check passed%s\n", colorGreen, colorReset)
	return nil
}
