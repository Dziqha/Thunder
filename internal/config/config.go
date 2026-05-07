package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BuildPath   string        `toml:"build_path"`
	MainFile    string        `toml:"main_file"`
	Project     ProjectConfig `toml:"project"`
	Services    ServicesMap   `toml:"services"`
	Profiles    ProfilesMap   `toml:"profiles"`
	WatchExts   []string      `toml:"watch_exts"`
	WatchFiles  []string      `toml:"watch_files"`
	WatchDirs   []string      `toml:"watch_dirs"`
	ExcludeDirs []string      `toml:"exclude_dirs"`
	BuildArgs   []string      `toml:"build_args"`
	RunArgs     []string      `toml:"run_args"`
	Debounce    int           `toml:"debounce"` // in milliseconds
	DebounceD   time.Duration `toml:"-"`
}

type ProjectConfig struct {
	Name           string `toml:"name"`
	DefaultProfile string `toml:"default_profile"`
	LogFormat      string `toml:"log_format"`
}

type ServiceConfig struct {
	Type             string               `toml:"type"`
	Package          string               `toml:"package"`
	BuildPath        string               `toml:"build_path"`
	Command          []string             `toml:"command"`
	WorkingDir       string               `toml:"working_dir"`
	DependsOn        []string             `toml:"depends_on"`
	EnvFiles         []string             `toml:"env_files"`
	Env              map[string]string    `toml:"env"`
	RestartPolicy    string               `toml:"restart_policy"`
	MaxRestarts      int                  `toml:"max_restarts"`
	RestartBackoff   RestartBackoffConfig `toml:"restart_backoff"`
	DependsTimeoutMS int                  `toml:"depends_timeout_ms"`
	Healthcheck      HealthcheckConfig    `toml:"healthcheck"`
}

type RestartBackoffConfig struct {
	Strategy  string `toml:"strategy"`
	BaseMS    int    `toml:"base_ms"`
	MaxMS     int    `toml:"max_ms"`
	JitterPct int    `toml:"jitter_pct"`
}

type HealthcheckConfig struct {
	Type       string `toml:"type"`
	URL        string `toml:"url"`
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	IntervalMS int    `toml:"interval_ms"`
	TimeoutMS  int    `toml:"timeout_ms"`
	Retries    int    `toml:"retries"`
	DelayMS    int    `toml:"delay_ms"`
}

type ServicesMap map[string]ServiceConfig

type ProfileConfig struct {
	Services []string `toml:"services"`
}

type ProfilesMap map[string]ProfileConfig

func Default() *Config {
	return &Config{
		BuildPath:   defaultBuildPath(),
		MainFile:    "main.go",
		WatchExts:   []string{".go", ".mod", ".sum", ".env", ".tmpl", ".tpl", ".yaml", ".yml", ".json"},
		WatchFiles:  []string{"go.mod", "go.sum", "thunder.toml", ".env"},
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

	if err := cfg.NormalizeAndValidate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) NormalizeAndValidate() error {
	if c.BuildPath == "" {
		c.BuildPath = defaultBuildPath()
	}

	if c.MainFile == "" {
		c.MainFile = "main.go"
	}

	if c.Debounce < 50 {
		c.Debounce = 50
	}
	c.DebounceD = time.Duration(c.Debounce) * time.Millisecond

	if len(c.WatchExts) == 0 {
		c.WatchExts = Default().WatchExts
	}
	for i := range c.WatchExts {
		ext := strings.TrimSpace(strings.ToLower(c.WatchExts[i]))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		c.WatchExts[i] = ext
	}

	if len(c.WatchFiles) == 0 {
		c.WatchFiles = Default().WatchFiles
	}
	for i := range c.WatchFiles {
		c.WatchFiles[i] = filepath.Clean(c.WatchFiles[i])
	}

	if len(c.WatchDirs) == 0 {
		c.WatchDirs = []string{"."}
	}
	for i := range c.WatchDirs {
		c.WatchDirs[i] = filepath.Clean(c.WatchDirs[i])
	}

	if c.MainFile == "" {
		return fmt.Errorf("main_file cannot be empty")
	}

	lf := strings.TrimSpace(strings.ToLower(c.Project.LogFormat))
	if lf == "" {
		c.Project.LogFormat = "text"
	} else if lf != "text" && lf != "json" {
		return fmt.Errorf("project.log_format must be text or json")
	}

	if c.Services != nil {
		for name, svc := range c.Services {
			if err := validateService(name, svc); err != nil {
				return err
			}
		}
	}

	if c.Profiles != nil {
		for profileName, profile := range c.Profiles {
			for _, service := range profile.Services {
				if _, ok := c.Services[service]; !ok {
					return fmt.Errorf("profile %q references unknown service %q", profileName, service)
				}
			}
		}
	}

	return nil
}

func validateService(name string, svc ServiceConfig) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	typeName := strings.TrimSpace(strings.ToLower(svc.Type))
	if typeName == "" {
		typeName = "go"
	}

	switch typeName {
	case "go":
		if strings.TrimSpace(svc.Package) == "" {
			return fmt.Errorf("service %q: package is required for type=go", name)
		}
	case "process":
		if len(svc.Command) == 0 {
			return fmt.Errorf("service %q: command is required for type=process", name)
		}
	default:
		return fmt.Errorf("service %q: unsupported type %q", name, svc.Type)
	}

	rp := strings.TrimSpace(strings.ToLower(svc.RestartPolicy))
	if rp == "" {
		rp = "always"
	}
	if rp != "always" && rp != "on-failure" && rp != "never" {
		return fmt.Errorf("service %q: unsupported restart_policy %q", name, svc.RestartPolicy)
	}

	if svc.MaxRestarts < 0 {
		return fmt.Errorf("service %q: max_restarts must be >= 0", name)
	}
	if svc.DependsTimeoutMS < 0 {
		return fmt.Errorf("service %q: depends_timeout_ms must be >= 0", name)
	}

	bs := strings.TrimSpace(strings.ToLower(svc.RestartBackoff.Strategy))
	if bs == "" {
		bs = "exponential"
	}
	if bs != "fixed" && bs != "linear" && bs != "exponential" {
		return fmt.Errorf("service %q: unsupported restart_backoff.strategy %q", name, svc.RestartBackoff.Strategy)
	}
	if svc.RestartBackoff.BaseMS < 0 || svc.RestartBackoff.MaxMS < 0 {
		return fmt.Errorf("service %q: restart_backoff base_ms/max_ms must be >= 0", name)
	}
	if svc.RestartBackoff.JitterPct < 0 || svc.RestartBackoff.JitterPct > 100 {
		return fmt.Errorf("service %q: restart_backoff.jitter_pct must be 0-100", name)
	}

	ht := strings.TrimSpace(strings.ToLower(svc.Healthcheck.Type))
	if ht != "" {
		switch ht {
		case "http":
			if strings.TrimSpace(svc.Healthcheck.URL) == "" {
				return fmt.Errorf("service %q: healthcheck.url is required for http type", name)
			}
		case "tcp":
			if strings.TrimSpace(svc.Healthcheck.Host) == "" || svc.Healthcheck.Port <= 0 {
				return fmt.Errorf("service %q: healthcheck.host and healthcheck.port are required for tcp type", name)
			}
		case "none":
		default:
			return fmt.Errorf("service %q: unsupported healthcheck.type %q", name, svc.Healthcheck.Type)
		}
	}

	for _, dep := range svc.DependsOn {
		if strings.TrimSpace(dep) == "" {
			return fmt.Errorf("service %q: depends_on contains empty name", name)
		}
	}

	return nil
}

func (c *Config) ServicesForProfile(profile string) ([]string, error) {
	if len(c.Services) == 0 {
		return nil, nil
	}

	resolved := strings.TrimSpace(profile)
	if resolved == "" {
		resolved = strings.TrimSpace(c.Project.DefaultProfile)
	}

	if resolved == "" {
		all := make([]string, 0, len(c.Services))
		for name := range c.Services {
			all = append(all, name)
		}
		return all, nil
	}

	p, ok := c.Profiles[resolved]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", resolved)
	}

	return append([]string(nil), p.Services...), nil
}

func defaultBuildPath() string {
	if runtime.GOOS == "windows" {
		return ".\\tmp\\main.exe"
	}
	return "./tmp/main"
}
