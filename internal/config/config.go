package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BuildPath   string        `toml:"build_path"`
	MainFile    string        `toml:"main_file"`
	WatchDirs   []string      `toml:"watch_dirs"`
	ExcludeDirs []string      `toml:"exclude_dirs"`
	BuildArgs   []string      `toml:"build_args"`
	RunArgs     []string      `toml:"run_args"`
	Debounce    int           `toml:"debounce"` // in milliseconds
	DebounceD   time.Duration `toml:"-"`
}

func Default() *Config {
	return &Config{
		BuildPath:   "./tmp/main",
		MainFile:    "main.go",
		WatchDirs:   []string{"."},
		ExcludeDirs: []string{"tmp", "vendor", ".git", "node_modules", ".idea", "bin"},
		BuildArgs:   []string{},
		RunArgs:     []string{},
		Debounce:    100,
		DebounceD:   100 * time.Millisecond,
	}
}

func Load() (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile("thunder.toml")
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.DebounceD = time.Duration(cfg.Debounce) * time.Millisecond

	if len(cfg.WatchDirs) == 0 {
		cfg.WatchDirs = []string{"."}
	}

	return cfg, nil
}