# Config Reference

Thunder reads `thunder.toml` from the current working directory.

## Top-level (single-app mode)

- `build_path` (string): output binary path
- `main_file` (string): package/file build target for `thunder run`
- `watch_dirs` (array): directories to watch
- `exclude_dirs` (array): directory names to skip
- `watch_exts` (array): watched file extensions (with or without dot)
- `watch_files` (array): explicit filenames to watch
- `build_args` (array): extra args for `go build`
- `run_args` (array): args passed to built binary
- `debounce` (int ms): debounce delay

## `[project]`

- `name` (string)
- `default_profile` (string)
- `log_format` (`text` | `json`)

## `[services.<name>]`

- `type` (`go` | `process`)
- `package` (required for `go`)
- `build_path` (optional for `go`)
- `command` (required for `process`)
- `working_dir` (optional)
- `depends_on` (array)
- `env_files` (array)
- `env` (table)
- `restart_policy` (`always` | `on-failure` | `never`)
- `max_restarts` (int >= 0)
- `depends_timeout_ms` (int >= 0)

### `[services.<name>.restart_backoff]`

- `strategy` (`fixed` | `linear` | `exponential`)
- `base_ms` (int >= 0)
- `max_ms` (int >= 0)
- `jitter_pct` (0-100)

### `[services.<name>.healthcheck]`

- `type` (`none` | `http` | `tcp`)
- `url` (required for `http`)
- `host`, `port` (required for `tcp`)
- `interval_ms`
- `timeout_ms`
- `retries`
- `delay_ms`

## `[profiles.<name>]`

- `services` (array of service names)
