package watcher

import (
	"testing"

	"github.com/Dziqha/Thunder/internal/config"
)

func TestShouldWatchFile(t *testing.T) {
	w := &Watcher{config: &config.Config{WatchExts: []string{".go", ".mod"}, WatchFiles: []string{"go.sum", ".env"}}}

	if !w.shouldWatchFile("internal/main.go") {
		t.Fatal("expected .go to be watched")
	}
	if !w.shouldWatchFile("go.sum") {
		t.Fatal("expected go.sum to be watched")
	}
	if w.shouldWatchFile("README.md") {
		t.Fatal("expected README.md not to be watched")
	}
}

func TestIsExcludedDir(t *testing.T) {
	w := &Watcher{config: &config.Config{ExcludeDirs: []string{"tmp", "vendor", ".git"}}}

	if !w.isExcludedDir("tmp") {
		t.Fatal("expected tmp to be excluded")
	}
	if w.isExcludedDir("internal") {
		t.Fatal("expected internal not to be excluded")
	}
}
