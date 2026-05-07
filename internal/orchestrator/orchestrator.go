package orchestrator

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Dziqha/Thunder/internal/config"
)

type Orchestrator struct {
	cfg      *config.Config
	cmds     map[string]*exec.Cmd
	cancels  map[string]context.CancelFunc
	healthy  map[string]bool
	mu       sync.Mutex
	basePath string
	logJSON  bool
	events   chan Event
}

type Event struct {
	TS      time.Time `json:"ts"`
	Service string    `json:"service"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

func New(cfg *config.Config) (*Orchestrator, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Orchestrator{
		cfg:      cfg,
		cmds:     make(map[string]*exec.Cmd),
		cancels:  make(map[string]context.CancelFunc),
		healthy:  make(map[string]bool),
		basePath: wd,
		logJSON:  strings.EqualFold(cfg.Project.LogFormat, "json"),
		events:   make(chan Event, 256),
	}, nil
}

func (o *Orchestrator) Events() <-chan Event {
	return o.events
}

func (o *Orchestrator) StartProfile(ctx context.Context, profile string) error {
	serviceNames, err := o.cfg.ServicesForProfile(profile)
	if err != nil {
		return err
	}
	if len(serviceNames) == 0 {
		return fmt.Errorf("no services configured")
	}

	ordered, err := o.resolveOrder(serviceNames)
	if err != nil {
		return err
	}

	for _, name := range ordered {
		o.emit(name, "dependency.wait", "waiting for dependencies")
		if err := o.waitForDependenciesHealthy(ctx, name); err != nil {
			o.emit(name, "dependency.timeout", err.Error())
			o.StopAll()
			return err
		}
		o.emit(name, "dependency.ready", "dependencies healthy")
		if err := o.startService(ctx, name, o.cfg.Services[name]); err != nil {
			o.emit(name, "service.error", err.Error())
			o.StopAll()
			return err
		}
		o.emit(name, "service.started", "service started")
	}

	return nil
}

func (o *Orchestrator) StopAll() {
	o.mu.Lock()
	defer o.mu.Unlock()

	for name, cancel := range o.cancels {
		cancel()
		delete(o.cancels, name)
	}

	for name, cmd := range o.cmds {
		if cmd != nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
			o.emit(name, "service.stop", "service stopped")
		}
		delete(o.cmds, name)
		delete(o.healthy, name)
	}
}

func (o *Orchestrator) startService(rootCtx context.Context, name string, svc config.ServiceConfig) error {
	t := strings.ToLower(strings.TrimSpace(svc.Type))
	if t == "" {
		t = "go"
	}

	ctx, cancel := context.WithCancel(rootCtx)

	var cmd *exec.Cmd
	switch t {
	case "go":
		buildPath := svc.BuildPath
		if strings.TrimSpace(buildPath) == "" {
			buildPath = filepath.Join("tmp", name+binarySuffix())
		}
		if err := os.MkdirAll(filepath.Dir(buildPath), 0755); err != nil {
			cancel()
			return fmt.Errorf("service %s: create build dir: %w", name, err)
		}

		buildArgs := []string{"build", "-o", buildPath, svc.Package}
		buildCmd := exec.Command("go", buildArgs...)
		buildCmd.Stdout = o.prefixedWriter(name, os.Stdout)
		buildCmd.Stderr = o.prefixedWriter(name, os.Stderr)
		if err := buildCmd.Run(); err != nil {
			cancel()
			return fmt.Errorf("service %s: build failed: %w", name, err)
		}

		cmd = exec.CommandContext(ctx, buildPath)
	case "process":
		cmd = exec.CommandContext(ctx, svc.Command[0], svc.Command[1:]...)
	default:
		cancel()
		return fmt.Errorf("service %s: unsupported type %q", name, svc.Type)
	}

	if svc.WorkingDir != "" {
		cmd.Dir = filepath.Join(o.basePath, svc.WorkingDir)
	}

	env, err := o.buildServiceEnv(svc)
	if err != nil {
		cancel()
		return fmt.Errorf("service %s: env setup failed: %w", name, err)
	}
	cmd.Env = env

	cmd.Stdout = o.prefixedWriter(name, os.Stdout)
	cmd.Stderr = o.prefixedWriter(name, os.Stderr)
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("service %s: start failed: %w", name, err)
	}

	if err := waitForHealth(ctx, svc.Healthcheck); err != nil {
		cancel()
		_ = cmd.Process.Kill()
		o.emit(name, "health.fail", err.Error())
		return fmt.Errorf("service %s: healthcheck failed: %w", name, err)
	}
	o.emit(name, "health.ok", "service is healthy")

	o.mu.Lock()
	o.cmds[name] = cmd
	o.cancels[name] = cancel
	o.healthy[name] = true
	o.mu.Unlock()

	go func(serviceName string, c *exec.Cmd, serviceConfig config.ServiceConfig) {
		err := c.Wait()
		if err != nil {
			o.emit(serviceName, "service.exit.error", err.Error())
		} else {
			o.emit(serviceName, "service.exit", "service exited")
		}
		o.mu.Lock()
		delete(o.cmds, serviceName)
		delete(o.cancels, serviceName)
		delete(o.healthy, serviceName)
		o.mu.Unlock()

		if shouldRestart(serviceConfig, err) {
			for attempt := 1; attempt <= maxRestarts(serviceConfig); attempt++ {
				if rootCtx.Err() != nil {
					return
				}
				delay := restartDelay(serviceConfig, attempt)
				o.emit(serviceName, "restart.wait", fmt.Sprintf("attempt=%d delay=%s", attempt, delay))
				time.Sleep(delay)
				if restartErr := o.startService(rootCtx, serviceName, serviceConfig); restartErr == nil {
					fmt.Fprintf(os.Stderr, "[%s] restarted (%d/%d)\n", serviceName, attempt, maxRestarts(serviceConfig))
					o.emit(serviceName, "restart.ok", fmt.Sprintf("attempt=%d", attempt))
					return
				}
				o.emit(serviceName, "restart.fail", fmt.Sprintf("attempt=%d", attempt))
			}
		}
	}(name, cmd, svc)

	return nil
}

func (o *Orchestrator) emit(service, kind, message string) {
	e := Event{TS: time.Now().UTC(), Service: service, Type: kind, Message: message}
	select {
	case o.events <- e:
	default:
	}
}

func (o *Orchestrator) waitForDependenciesHealthy(ctx context.Context, name string) error {
	svc, ok := o.cfg.Services[name]
	if !ok {
		return fmt.Errorf("unknown service %q", name)
	}
	timeoutMS := svc.DependsTimeoutMS
	if timeoutMS <= 0 {
		timeoutMS = 15_000
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutMS)*time.Millisecond)
	defer cancel()

	for _, dep := range svc.DependsOn {
		for {
			select {
			case <-timeoutCtx.Done():
				return fmt.Errorf("service %q timed out waiting dependency %q health", name, dep)
			default:
			}

			o.mu.Lock()
			healthy := o.healthy[dep]
			o.mu.Unlock()
			if healthy {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}

func (o *Orchestrator) buildServiceEnv(svc config.ServiceConfig) ([]string, error) {
	merged := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			merged[parts[0]] = parts[1]
		}
	}

	for _, envFile := range svc.EnvFiles {
		if err := applyEnvFile(merged, envFile); err != nil {
			return nil, err
		}
	}

	for k, v := range svc.Env {
		merged[k] = v
	}

	out := make([]string, 0, len(merged))
	for k, v := range merged {
		out = append(out, k+"="+v)
	}
	sort.Strings(out)
	return out, nil
}

func applyEnvFile(dst map[string]string, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
		if k != "" {
			dst[k] = v
		}
	}
	return s.Err()
}

func shouldRestart(svc config.ServiceConfig, err error) bool {
	policy := strings.ToLower(strings.TrimSpace(svc.RestartPolicy))
	if policy == "" {
		policy = "always"
	}

	switch policy {
	case "never":
		return false
	case "on-failure":
		return err != nil
	default:
		return true
	}
}

func maxRestarts(svc config.ServiceConfig) int {
	if svc.MaxRestarts <= 0 {
		return 3
	}
	return svc.MaxRestarts
}

func restartDelay(svc config.ServiceConfig, attempt int) time.Duration {
	b := svc.RestartBackoff
	strategy := strings.ToLower(strings.TrimSpace(b.Strategy))
	if strategy == "" {
		strategy = "exponential"
	}
	base := b.BaseMS
	if base <= 0 {
		base = 500
	}
	max := b.MaxMS
	if max <= 0 {
		max = 10_000
	}

	ms := base
	switch strategy {
	case "fixed":
		ms = base
	case "linear":
		ms = base * attempt
	default:
		if attempt <= 1 {
			ms = base
		} else {
			ms = base * (1 << (attempt - 1))
		}
	}
	if ms > max {
		ms = max
	}

	jitterPct := b.JitterPct
	if jitterPct < 0 {
		jitterPct = 0
	}
	jitter := int(float64(ms) * (float64(jitterPct) / 100.0))
	if jitter > 0 {
		delta := rand.Intn((2*jitter)+1) - jitter
		ms += delta
	}
	if ms < 0 {
		ms = 0
	}
	return time.Duration(ms) * time.Millisecond
}

func waitForHealth(ctx context.Context, hc config.HealthcheckConfig) error {
	t := strings.ToLower(strings.TrimSpace(hc.Type))
	if t == "" || t == "none" {
		return nil
	}

	interval := durationOrDefault(hc.IntervalMS, 500*time.Millisecond)
	timeout := durationOrDefault(hc.TimeoutMS, 3*time.Second)
	retries := hc.Retries
	if retries <= 0 {
		retries = 20
	}
	if hc.DelayMS > 0 {
		select {
		case <-time.After(time.Duration(hc.DelayMS) * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	var lastErr error
	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		switch t {
		case "http":
			lastErr = checkHTTPHealth(hc.URL, timeout)
		case "tcp":
			lastErr = checkTCPHealth(hc.Host, hc.Port, timeout)
		default:
			return fmt.Errorf("unsupported healthcheck type %q", hc.Type)
		}

		if lastErr == nil {
			return nil
		}
		time.Sleep(interval)
	}

	if lastErr == nil {
		lastErr = errors.New("unknown healthcheck error")
	}
	return lastErr
}

func checkHTTPHealth(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}
	return nil
}

func checkTCPHealth(host string, port int, timeout time.Duration) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func durationOrDefault(ms int, fallback time.Duration) time.Duration {
	if ms <= 0 {
		return fallback
	}
	return time.Duration(ms) * time.Millisecond
}

func (o *Orchestrator) resolveOrder(selected []string) ([]string, error) {
	selectedSet := make(map[string]struct{}, len(selected))
	for _, s := range selected {
		selectedSet[s] = struct{}{}
	}

	visited := make(map[string]int)
	order := make([]string, 0, len(selected))

	var visit func(string) error
	visit = func(name string) error {
		state := visited[name]
		if state == 1 {
			return fmt.Errorf("cyclic dependency detected at service %q", name)
		}
		if state == 2 {
			return nil
		}

		svc, ok := o.cfg.Services[name]
		if !ok {
			return fmt.Errorf("unknown service %q", name)
		}

		visited[name] = 1
		for _, dep := range svc.DependsOn {
			if _, ok := selectedSet[dep]; !ok {
				continue
			}
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[name] = 2
		order = append(order, name)
		return nil
	}

	sortedSelected := append([]string(nil), selected...)
	sort.Strings(sortedSelected)
	for _, name := range sortedSelected {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

type prefixWriter struct {
	name string
	to   io.Writer
	json bool
}

func (p *prefixWriter) Write(data []byte) (int, error) {
	trimmed := strings.TrimSuffix(string(data), "\n")
	if trimmed == "" {
		return len(data), nil
	}
	redacted := redactSecrets(trimmed)
	if p.json {
		payload := map[string]string{
			"service": p.name,
			"message": redacted,
			"ts":      time.Now().Format(time.RFC3339Nano),
		}
		b, err := json.Marshal(payload)
		if err != nil {
			return 0, err
		}
		_, err = fmt.Fprintf(p.to, "%s\n", string(b))
		if err != nil {
			return 0, err
		}
		return len(data), nil
	}

	_, err := fmt.Fprintf(p.to, "[%s] %s\n", p.name, redacted)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func redactSecrets(line string) string {
	parts := strings.Fields(line)
	for i := range parts {
		kv := strings.SplitN(parts[i], "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(kv[0]))
		if strings.Contains(key, "TOKEN") || strings.Contains(key, "SECRET") || strings.Contains(key, "PASSWORD") || strings.Contains(key, "API_KEY") {
			parts[i] = kv[0] + "=***REDACTED***"
		}
	}
	return strings.Join(parts, " ")
}

func (o *Orchestrator) prefixedWriter(name string, to io.Writer) io.Writer {
	return &prefixWriter{name: name, to: to, json: o.logJSON}
}

func binarySuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
