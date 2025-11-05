# ====================
# Makefile
# ====================
.PHONY: dev build install clean test run help

# Build thunder binary
build:
	@echo "🔨 Building Thunder..."
	@go build -o thunder thunder.go
	@echo "✅ Thunder built successfully!"

# Install thunder to GOPATH
install:
	@echo "📦 Installing Thunder..."
	@go install thunder.go
	@echo "✅ Thunder installed! Use 'thunder' command anywhere"

# Run with hot reload (development mode)
dev:
	@echo "⚡ Starting Thunder hot reload..."
	@go run thunder.go

# Run specific file
run:
	@go run thunder.go $(FILE)

# Build your app
build-app:
	@echo "🔨 Building application..."
	@go build -o bin/app main.go
	@echo "✅ Application built to bin/app"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf tmp/
	@rm -f thunder
	@rm -f main
	@rm -rf bin/
	@echo "✅ Cleaned!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Format code
fmt:
	@echo "💅 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies ready!"

# Initialize new project
init:
	@echo "🎯 Initializing Thunder project..."
	@go mod init myapp
	@go get github.com/fsnotify/fsnotify
	@echo "✅ Project initialized!"

# Show help
help:
	@echo "⚡ Thunder - Ultra Fast Hot Reload"
	@echo ""
	@echo "Available commands:"
	@echo "  make dev        - Start hot reload (default: main.go)"
	@echo "  make run FILE=  - Run specific file (e.g., make run FILE=cmd/api/main.go)"
	@echo "  make build      - Build thunder binary"
	@echo "  make install    - Install thunder globally"
	@echo "  make build-app  - Build your application"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make fmt        - Format code"
	@echo "  make deps       - Download dependencies"
	@echo "  make init       - Initialize new project"
	@echo ""
	@echo "Quick start:"
	@echo "  1. make init    # First time setup"
	@echo "  2. make dev     # Start developing"
