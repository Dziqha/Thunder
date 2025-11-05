// internal/cli/cli.go
package cli

import (
	"fmt"
	"os"

	"github.com/Dziqha/Thunder/internal/commands"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
)

func Execute() error {
	if len(os.Args) < 2 {
		showHelp()
		return nil
	}

	command := os.Args[1]

	switch command {
	case "init":
		return commands.Init()
	case "run":
		return commands.Run()
	case "version", "-v", "--version":
		showVersion()
		return nil
	case "help", "-h", "--help":
		showHelp()
		return nil
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		showHelp()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func showVersion() {
	fmt.Printf("%s⚡ Thunder v1.0.0%s\n", colorBlue, colorReset)
	fmt.Println("Ultra-fast hot reload for Go")
}

func showHelp() {
	banner := fmt.Sprintf(`
%s╔════════════════════════════════════╗
║     ⚡ THUNDER HOT RELOAD ⚡       ║
║   Faster than Air, Lighter than   ║
║         Lightning Strike!          ║
╚════════════════════════════════════╝%s

Usage:
  Thunder <command> [arguments]

Commands:
  init        Initialize Thunder in current directory
  run         Run your app with hot reload
  version     Show Thunder version
  help        Show this help message

Examples:
  Thunder init              # Initialize Thunder
  Thunder run               # Run with hot reload (uses main.go)
  Thunder run cmd/api       # Run specific package

Installation:
  go install github.com/Dziqha/thunder/cmd/thunder@latest

Learn more: https://github.com/Dziqha/thunder
`, colorBlue, colorReset)

	fmt.Println(banner)
}