// internal/watcher/watcher.go
package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/Dziqha/Thunder/internal/config"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

type Watcher struct {
	watcher       *fsnotify.Watcher
	cmd           *exec.Cmd
	mutex         sync.Mutex
	debounceTimer *time.Timer
	config        *config.Config
	cancelFunc    context.CancelFunc
}

func New(cfg *config.Config) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		config:  cfg,
	}, nil
}

func (w *Watcher) Start() error {
	// Create tmp directory
	if err := os.MkdirAll(filepath.Dir(w.config.BuildPath), 0755); err != nil {
		return err
	}

	// Add directories to watcher
	for _, dir := range w.config.WatchDirs {
		if err := w.addRecursive(dir); err != nil {
			log.Printf("%s⚠ Warning: Could not watch directory %s: %v%s\n", colorYellow, dir, err, colorReset)
		}
	}

	// Initial build and run
	if err := w.rebuild(); err != nil {
		log.Printf("%s✗ Initial build failed: %v%s\n", colorRed, err, colorReset)
	}

	// Watch for changes
	go w.watch()

	// Keep running
	select {}
}

func (w *Watcher) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if info.IsDir() {
			for _, exclude := range w.config.ExcludeDirs {
				if info.Name() == exclude {
					return filepath.SkipDir
				}
			}
			return w.watcher.Add(path)
		}
		return nil
	})
}

func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only watch .go files
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}

			// Ignore certain operations
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			w.scheduleRebuild(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("%s✗ Watcher error: %v%s\n", colorRed, err, colorReset)
		}
	}
}

func (w *Watcher) scheduleRebuild(event fsnotify.Event) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(w.config.DebounceD, func() {
		log.Printf("%s⚡ File changed: %s%s\n", colorYellow, filepath.Base(event.Name), colorReset)
		if err := w.rebuild(); err != nil {
			log.Printf("%s✗ Build failed: %v%s\n", colorRed, err, colorReset)
		}
	})
}

func (w *Watcher) rebuild() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Stop current process
	w.stop()

	// Build
	start := time.Now()
	log.Printf("%s⚙ Building...%s\n", colorCyan, colorReset)

	args := append([]string{"build", "-o", w.config.BuildPath}, w.config.BuildArgs...)
	args = append(args, w.config.MainFile)

	buildCmd := exec.Command("go", args...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return err
	}

	buildTime := time.Since(start)
	log.Printf("%s✓ Build completed in %dms%s\n", colorGreen, buildTime.Milliseconds(), colorReset)

	// Run
	return w.run()
}

func (w *Watcher) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	w.cancelFunc = cancel

	w.cmd = exec.CommandContext(ctx, w.config.BuildPath, w.config.RunArgs...)
	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr
	w.cmd.Stdin = os.Stdin

	log.Printf("%s▶ Starting application...%s\n", colorPurple, colorReset)
	fmt.Println(strings.Repeat("─", 50))

	if err := w.cmd.Start(); err != nil {
		return err
	}

	// Wait in goroutine
	go func() {
		if err := w.cmd.Wait(); err != nil {
			// Only log if not killed by us
			if ctx.Err() == nil {
				log.Printf("%s✗ Application exited with error: %v%s\n", colorRed, err, colorReset)
			}
		}
	}()

	return nil
}

func (w *Watcher) stop() {
	if w.cancelFunc != nil {
		w.cancelFunc()
		w.cancelFunc = nil
	}

	if w.cmd != nil && w.cmd.Process != nil {
		// Give it a moment to cleanup
		time.Sleep(50 * time.Millisecond)
		w.cmd = nil
	}
}

func (w *Watcher) Close() {
	w.stop()
	w.watcher.Close()
}