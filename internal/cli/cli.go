package cli

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Dziqha/Thunder/internal/commands"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorYellow = "\033[33m"
)

var version = "dev"

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
	case "dev":
		return commands.Dev()
	case "doctor":
		return commands.Doctor()
	case "inspect":
		return commands.Inspect()
	case "events":
		return commands.Events()
	case "release-check":
		return commands.ReleaseCheck()
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
	fmt.Printf("%s⚡ Thunder %s%s\n", colorBlue, resolveVersion(), colorReset)
	fmt.Println("Ultra-fast hot reload for Go")
}

func resolveVersion() string {
	if version != "" && version != "dev" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	for _, dep := range info.Deps {
		if dep.Path == "github.com/Dziqha/Thunder" && dep.Version != "" && dep.Version != "(devel)" {
			return dep.Version
		}
	}

	return "dev"
}

func showHelp() {
	banner := fmt.Sprintf(`
%s╔════════════════════════════════════╗
║     ⚡ THUNDER HOT RELOAD ⚡       ║
║   Faster than Air, Lighter than   ║
║         Lightning Strike!          ║
╚════════════════════════════════════╝%s

Usage:
  thunder <command> [arguments]

Commands:
  init        Initialize Thunder in current directory
  run         Run your app with hot reload
  dev         Run multi-service orchestration from config
  doctor      Validate config, paths, and runtime readiness
  inspect     Show resolved project/profile/service config
  events      Stream runtime events (json/text)
  release-check Validate release readiness
  version     Show Thunder version
  help        Show this help message

Examples:
  thunder init              # Initialize Thunder
  thunder run               # Run with hot reload (uses main.go)
  thunder run ./cmd/api     # Run specific package
  thunder dev               # Run default profile services
  thunder dev backend       # Run specific profile
  thunder doctor            # Run diagnostics
  thunder inspect backend   # Show resolved profile services
  thunder events dev --format=json --service=api --type=restart.
  thunder events dev --format=json --out=events.log --also-stdout
  thunder release-check

Installation:
  go install github.com/Dziqha/thunder/cmd/thunder@latest

Learn more: https://github.com/Dziqha/thunder
`, colorBlue, colorReset)

	fmt.Println(banner)
}
