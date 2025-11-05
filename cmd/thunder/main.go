// cmd/thunder/main.go
package main

import (
	"fmt"
	"os"

	"github.com/Dziqha/thunder/internal/cli"
)

const version = "1.0.0"

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}