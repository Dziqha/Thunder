package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
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
	lastEvent     string
	config        *config.Config
	cancelFunc    context.CancelFunc
	rootCtx       context.Context
}

func New(cfg *config.Config) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		config:  cfg,
		rootCtx: context.Background(),
	}, nil
}

func (w *Watcher) Start() error {
	if err := w.config.NormalizeAndValidate(); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(w.config.BuildPath), 0755); err != nil {
		return err
	}

	for _, dir := range w.config.WatchDirs {
		if err := w.addRecursive(dir); err != nil {
			log.Printf("%s⚠ Warning: Could not watch directory %s: %v%s\n", colorYellow, dir, err, colorReset)
		}
	}

	if err := w.rebuild(); err != nil {
		log.Printf("%s✗ Initial build failed: %v%s\n", colorRed, err, colorReset)
	}

	go w.watch()

	sigCtx, stop := signal.NotifyContext(w.rootCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-sigCtx.Done()

	log.Printf("%s⏹ Shutting down Thunder...%s\n", colorBlue, colorReset)
	w.Close()
	return nil
}

func (w *Watcher) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			for _, exclude := range w.config.ExcludeDirs {
				if info.Name() == exclude {
					return filepath.SkipDir
				}
			}
			if err := w.watcher.Add(path); err != nil {
				return err
			}
			return nil
		}

		return nil
	})
}

func (w *Watcher) shouldWatchFile(name string) bool {
	base := filepath.Base(name)
	for _, f := range w.config.WatchFiles {
		if base == filepath.Base(f) {
			return true
		}
	}

	ext := strings.ToLower(filepath.Ext(base))
	for _, allowed := range w.config.WatchExts {
		if ext == allowed {
			return true
		}
	}

	return false
}

func (w *Watcher) isExcludedDir(path string) bool {
	name := filepath.Base(path)
	for _, exclude := range w.config.ExcludeDirs {
		if name == exclude {
			return true
		}
	}
	return false
}

func (w *Watcher) handleDirCreate(path string) {
	if w.isExcludedDir(path) {
		return
	}
	if err := w.addRecursive(path); err != nil {
		log.Printf("%s⚠ Failed to watch new directory %s: %v%s\n", colorYellow, path, err, colorReset)
	}
}

func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					w.handleDirCreate(event.Name)
					continue
				}
			}

			if !w.shouldWatchFile(event.Name) {
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

	w.lastEvent = event.Name

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(w.config.DebounceD, func() {
		log.Printf("%s⚡ File changed: %s%s\n", colorYellow, filepath.Base(w.lastEvent), colorReset)
		if err := w.rebuild(); err != nil {
			log.Printf("%s✗ Build failed: %v%s\n", colorRed, err, colorReset)
		}
	})
}

func (w *Watcher) rebuild() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.stopLocked()

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

	return w.runLocked()
}

func (w *Watcher) runLocked() error {
	ctx, cancel := context.WithCancel(w.rootCtx)
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

	go func(cmd *exec.Cmd) {
		if err := cmd.Wait(); err != nil {
			if ctx.Err() == nil {
				log.Printf("%s✗ Application exited with error: %v%s\n", colorRed, err, colorReset)
			}
		}
	}(w.cmd)

	return nil
}

func (w *Watcher) stopLocked() {
	if w.cancelFunc != nil {
		w.cancelFunc()
		w.cancelFunc = nil
	}

	if w.cmd != nil && w.cmd.Process != nil {
		if runtime.GOOS != "windows" {
			_ = w.cmd.Process.Signal(syscall.SIGTERM)
		}
		time.Sleep(50 * time.Millisecond)
		_ = w.cmd.Process.Kill()
		w.cmd = nil
	}
}

func (w *Watcher) stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.stopLocked()
}

func (w *Watcher) Close() {
	w.stop()
	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}
	_ = w.watcher.Close()
}
