package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/Dziqha/Thunder/internal/config"
	"github.com/Dziqha/Thunder/internal/orchestrator"
)

func Events() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	profile := ""
	format := "json"
	serviceFilter := ""
	typePrefixFilter := ""
	outFile := ""
	alsoStdout := false
	maxSizeMB := 0
	maxFiles := 5

	for _, arg := range os.Args[2:] {
		switch {
		case strings.HasPrefix(arg, "--format="):
			format = strings.TrimPrefix(arg, "--format=")
		case strings.HasPrefix(arg, "--service="):
			serviceFilter = strings.TrimPrefix(arg, "--service=")
		case strings.HasPrefix(arg, "--type="):
			typePrefixFilter = strings.TrimPrefix(arg, "--type=")
		case strings.HasPrefix(arg, "--out="):
			outFile = strings.TrimPrefix(arg, "--out=")
		case arg == "--also-stdout":
			alsoStdout = true
		case strings.HasPrefix(arg, "--max-size-mb="):
			v := strings.TrimPrefix(arg, "--max-size-mb=")
			if n, err := strconv.Atoi(v); err == nil {
				maxSizeMB = n
			}
		case strings.HasPrefix(arg, "--max-files="):
			v := strings.TrimPrefix(arg, "--max-files=")
			if n, err := strconv.Atoi(v); err == nil {
				maxFiles = n
			}
		case strings.HasPrefix(arg, "--"):
			continue
		case profile == "":
			profile = arg
		}
	}

	var writers []io.Writer
	if outFile == "" || alsoStdout {
		writers = append(writers, os.Stdout)
	}
	if outFile != "" {
		if maxSizeMB > 0 {
			rot, err := newRotatingWriter(outFile, int64(maxSizeMB)*1024*1024, maxFiles)
			if err != nil {
				return fmt.Errorf("failed to create rotating output file: %w", err)
			}
			defer rot.Close()
			writers = append(writers, rot)
		} else {
			f, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open output file: %w", err)
			}
			defer f.Close()
			writers = append(writers, f)
		}
	}
	output := io.MultiWriter(writers...)

	orch, err := orchestrator.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := orch.StartProfile(ctx, profile); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			orch.StopAll()
			return nil
		case ev := <-orch.Events():
			if serviceFilter != "" && ev.Service != serviceFilter {
				continue
			}
			if typePrefixFilter != "" && !strings.HasPrefix(ev.Type, typePrefixFilter) {
				continue
			}
			if format == "text" {
				_, _ = fmt.Fprintf(output, "[%s] %s %s\n", ev.Service, ev.Type, ev.Message)
				continue
			}
			b, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			_, _ = fmt.Fprintf(output, "%s\n", string(b))
		}
	}
}

type rotatingWriter struct {
	path     string
	maxBytes int64
	maxFiles int
	file     *os.File
	mu       sync.Mutex
}

func newRotatingWriter(path string, maxBytes int64, maxFiles int) (*rotatingWriter, error) {
	if maxFiles < 2 {
		maxFiles = 2
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil && filepath.Dir(path) != "." {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &rotatingWriter{path: path, maxBytes: maxBytes, maxFiles: maxFiles, file: f}, nil
}

func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(int64(len(p))); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *rotatingWriter) rotateIfNeeded(incoming int64) error {
	if w.maxBytes <= 0 {
		return nil
	}
	st, err := w.file.Stat()
	if err != nil {
		return err
	}
	if st.Size()+incoming <= w.maxBytes {
		return nil
	}

	if err := w.file.Close(); err != nil {
		return err
	}

	for i := w.maxFiles - 1; i >= 1; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.path, i)
		newPath := fmt.Sprintf("%s.%d", w.path, i+1)
		if i == w.maxFiles-1 {
			_ = os.Remove(oldPath)
			continue
		}
		if _, err := os.Stat(oldPath); err == nil {
			_ = os.Rename(oldPath, newPath)
		}
	}
	if _, err := os.Stat(w.path); err == nil {
		_ = os.Rename(w.path, w.path+".1")
	}

	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	w.file = f
	return nil
}

func (w *rotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	return w.file.Close()
}
