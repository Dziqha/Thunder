// internal/commands/init.go
package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

const thunderConfig = `# Thunder Configuration
# Fast hot reload for Go applications

# Build settings
build_path = "./tmp/main"
main_file = "main.go"

# Watch settings
watch_dirs = ["."]
exclude_dirs = ["tmp", "vendor", ".git", "node_modules", ".idea", "bin"]

# Build arguments (optional)
# build_args = ["-tags=dev", "-race"]

# Run arguments (optional)
# run_args = ["-port=8080"]

# Debounce time in milliseconds (default: 100)
debounce = 100
`

const exampleMainGo = `package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Thunder! 🚀\n")
		fmt.Fprintf(w, "Time: %s\n", time.Now().Format("15:04:05"))
		fmt.Fprintf(w, "\nTry editing this file and watch Thunder reload instantly!\n")
	})
	
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"status\": \"running\", \"message\": \"Thunder is awesome!\"}")
	})

	addr := ":8080"
	log.Printf("⚡ Server starting on http://localhost%s\n", addr)
	log.Printf("💡 Edit any .go file to see Thunder in action!\n")
	
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
`

const gitignore = `# Thunder
tmp/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test
*.test
*.out

# Build
bin/
dist/

# Dependencies
vendor/

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`

const readme = `# My Thunder App

⚡ Fast Go application with hot reload powered by Thunder

## Quick Start

` + "```bash" + `
# Start development server with hot reload
thunder run

# Visit http://localhost:8080
` + "```" + `

## Commands

- ` + "`thunder run`" + ` - Start app with hot reload
- ` + "`thunder init`" + ` - Initialize Thunder (already done!)

## Project Structure

` + "```" + `
.
├── main.go           # Your application entry point
├── thunder.toml      # Thunder configuration
├── .gitignore        # Git ignore file
└── tmp/              # Build output (auto-created)
` + "```" + `

## Configuration

Edit ` + "`thunder.toml`" + ` to customize:

- Build path and main file
- Watch directories
- Excluded directories
- Build and run arguments
- Debounce time

## Learn More

- [Thunder Documentation](https://github.com/yourusername/thunder)
- [Go Documentation](https://go.dev/doc/)

---

Built with ⚡ by Thunder
`

func Init() error {
	fmt.Printf("%s⚡ Initializing Thunder...%s\n\n", colorCyan, colorReset)

	// Check if already initialized
	if _, err := os.Stat("thunder.toml"); err == nil {
		fmt.Printf("%s⚠ Thunder is already initialized in this directory%s\n", colorYellow, colorReset)
		fmt.Println("thunder.toml already exists")
		return nil
	}

	// Create thunder.toml
	if err := os.WriteFile("thunder.toml", []byte(thunderConfig), 0644); err != nil {
		return fmt.Errorf("failed to create thunder.toml: %v", err)
	}
	fmt.Printf("%s✓ Created thunder.toml%s\n", colorGreen, colorReset)

	// Create .gitignore
	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		if err := os.WriteFile(".gitignore", []byte(gitignore), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %v", err)
		}
		fmt.Printf("%s✓ Created .gitignore%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("%s→ .gitignore already exists (skipped)%s\n", colorYellow, colorReset)
	}

	// Create main.go if not exists
	if _, err := os.Stat("main.go"); os.IsNotExist(err) {
		if err := os.WriteFile("main.go", []byte(exampleMainGo), 0644); err != nil {
			return fmt.Errorf("failed to create main.go: %v", err)
		}
		fmt.Printf("%s✓ Created main.go (example app)%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("%s→ main.go already exists (skipped)%s\n", colorYellow, colorReset)
	}

	// Create README.md if not exists
	if _, err := os.Stat("README.md"); os.IsNotExist(err) {
		if err := os.WriteFile("README.md", []byte(readme), 0644); err != nil {
			return fmt.Errorf("failed to create README.md: %v", err)
		}
		fmt.Printf("%s✓ Created README.md%s\n", colorGreen, colorReset)
	} else {
		fmt.Printf("%s→ README.md already exists (skipped)%s\n", colorYellow, colorReset)
	}

	// Create tmp directory
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %v", err)
	}
	fmt.Printf("%s✓ Created tmp/ directory%s\n", colorGreen, colorReset)

	// Initialize go.mod if not exists
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		dirName := filepath.Base(getCurrentDir())
		fmt.Printf("\n%s📦 Initializing Go module...%s\n", colorCyan, colorReset)
		fmt.Printf("Module name: %s\n", dirName)
		// Note: In real implementation, we would execute: go mod init
		fmt.Printf("%s💡 Run: go mod init %s%s\n", colorYellow, dirName, colorReset)
	}

	fmt.Printf("\n%s🎉 Thunder initialized successfully!%s\n\n", colorGreen, colorReset)
	fmt.Println("Next steps:")
	fmt.Printf("  1. %sgo mod init yourmodule%s (if not done)\n", colorCyan, colorReset)
	fmt.Printf("  2. %sthunder run%s - Start development with hot reload\n", colorCyan, colorReset)
	fmt.Printf("  3. Edit main.go and watch Thunder reload instantly! ⚡\n\n")

	return nil
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "myapp"
	}
	return dir
}