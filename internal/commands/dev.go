package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dziqha/Thunder/internal/config"
	"github.com/Dziqha/Thunder/internal/orchestrator"
)

func Dev() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	orch, err := orchestrator.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator: %v", err)
	}

	profile := ""
	if len(os.Args) > 2 {
		profile = os.Args[2]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := orch.StartProfile(ctx, profile); err != nil {
		return err
	}

	fmt.Printf("%s⚡ Thunder dev orchestration running%s\n", colorYellow, colorReset)
	fmt.Printf("%s💡 Press Ctrl+C to stop all services%s\n\n", colorCyan, colorReset)

	<-ctx.Done()
	orch.StopAll()
	return nil
}
