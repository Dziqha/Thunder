# Benchmark Fixture: Small API

This fixture is the first reproducible benchmark target for Thunder vs Air.

Location:

- `benchmarks/fixture-small-api`

## What It Contains

- `cmd/api` entrypoint
- `/health` endpoint on port `18080`
- small internal package layout
- environment-based config loading
- `thunder.toml` for Thunder
- `.air.toml` for Air

## Run With Thunder

```bash
cd benchmarks/fixture-small-api
thunder run ./cmd/api
```

or orchestration mode:

```bash
cd benchmarks/fixture-small-api
thunder dev
```

## Run With Air

```bash
cd benchmarks/fixture-small-api
air -c .air.toml
```

## Suggested Benchmark Scenarios

1. Startup latency
- start tool
- wait until `/health` returns `200`

2. Rebuild latency
- edit a file in `internal/config/config.go`
- measure time until `/health` returns `200` again

3. Repeated change storm
- modify the same file 5-10 times quickly
- measure total rebuild behavior and stabilization time

## Comparison Runner

You can run the startup benchmark runner from the repo root.

Thunder:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool thunder -Fixture benchmarks/fixture-small-api -Iterations 5 -OutJson docs/thunder-startup.json -OutMd docs/thunder-startup.md
```

Air:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool air -Fixture benchmarks/fixture-small-api -Iterations 5 -OutJson docs/air-startup.json -OutMd docs/air-startup.md
```

This runner currently focuses on startup-to-health latency.

If a run fails, inspect:

- `benchmarks/fixture-small-api/tmp/bench-thunder-stdout.log`
- `benchmarks/fixture-small-api/tmp/bench-thunder-stderr.log`
- `benchmarks/fixture-small-api/tmp/bench-air-stdout.log`
- `benchmarks/fixture-small-api/tmp/bench-air-stderr.log`

## Notes

- Keep the same machine and Go version for both tools.
- Run at least 10 iterations per scenario.
- Record raw timings, not just averages.
