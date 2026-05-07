package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotatingWriterRotates(t *testing.T) {
	tmp := t.TempDir()
	logPath := filepath.Join(tmp, "events.log")

	w, err := newRotatingWriter(logPath, 64, 3)
	if err != nil {
		t.Fatalf("newRotatingWriter error: %v", err)
	}
	defer w.Close()

	line := []byte("abcdefghijklmnopqrstuvwxyz0123456789\n")
	for i := 0; i < 10; i++ {
		if _, err := w.Write(line); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}

	if _, err := os.Stat(logPath + ".1"); err != nil {
		t.Fatalf("expected rotated file .1, got error: %v", err)
	}
}
