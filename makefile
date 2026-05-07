.PHONY: build install run dev test race fmt clean bench release-check help

BINARY=thunder
CMD=./cmd/thunder

build:
	@echo "Building Thunder..."
	@go build -o $(BINARY) $(CMD)

install:
	@echo "Installing Thunder..."
	@go install $(CMD)

run:
	@go run $(CMD) run

dev:
	@go run $(CMD) dev

test:
	@go test ./...

race:
	@go test -race ./...

fmt:
	@go fmt ./...

bench:
	@pwsh ./scripts/bench.ps1 -Iterations 20 -Target $(CMD)

release-check:
	@go run $(CMD) release-check

clean:
	@echo "Cleaning artifacts..."
	@rm -rf tmp/
	@rm -f $(BINARY)

help:
	@echo "Thunder Make targets"
	@echo "  make build         Build thunder binary"
	@echo "  make install       Install thunder to GOPATH/bin"
	@echo "  make run           Run single-app hot reload mode"
	@echo "  make dev           Run orchestration mode"
	@echo "  make test          Run tests"
	@echo "  make race          Run race tests"
	@echo "  make fmt           Format source"
	@echo "  make bench         Run benchmark harness"
	@echo "  make release-check Run release readiness checks"
	@echo "  make clean         Remove local artifacts"
