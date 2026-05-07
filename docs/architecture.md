# Architecture

Thunder is split into a few core layers:

- `internal/cli`: command routing and user-facing output
- `internal/config`: config loading, normalization, and validation
- `internal/watcher`: hot reload watcher for single-app mode
- `internal/orchestrator`: multi-service lifecycle manager
- `internal/commands`: command handlers (`run`, `dev`, `doctor`, `inspect`, `events`)
- `internal/diagnostics`: runtime and config checks

## Runtime Modes

## 1) `thunder run`

Single application hot reload loop:

1. watch filesystem changes
2. debounce events
3. rebuild binary
4. restart process

## 2) `thunder dev`

Service orchestration flow:

1. resolve selected profile
2. topologically sort by `depends_on`
3. wait for dependency health
4. start service process
5. run health check
6. mark healthy
7. restart based on policy when service exits

## Event Stream

Orchestrator emits runtime events through an internal channel.
`thunder events` can consume and serialize them as text or JSON.

## Security Notes

- log output performs best-effort redaction for common secret key patterns
- command execution uses explicit argument arrays
- config validation rejects unsupported policy values
