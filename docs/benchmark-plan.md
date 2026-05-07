# Benchmark Plan

This document defines a fair and reproducible way to compare Thunder with other Go dev runners, especially Air.

## Goals

- measure local development feedback speed
- measure idle overhead
- compare behavior under repeated file changes
- avoid misleading one-off numbers

## Tools Under Comparison

- Thunder
- Air

Optional later:
- Fresh
- CompileDaemon

## Principles

1. Same machine
- run all tools on the same machine
- keep CPU mode and power profile consistent

2. Same project
- benchmark the exact same Go project layout
- same dependencies
- same entrypoint

3. Same scenario
- same file edit pattern
- same warm/cold state assumptions
- same number of iterations

4. Separate cold and warm runs
- cold: first run after cleanup
- warm: repeated runs without dependency changes

5. Publish raw numbers
- keep average, min, max, and all samples
- avoid only reporting best-case values

## Benchmark Scenarios

### 1. Startup latency
Measure time from command start to app ready.

Metrics:
- startup total time
- first health-ready time

### 2. Rebuild latency after code change
Measure time from file save to healthy restarted process.

Metrics:
- change detection latency
- build duration
- restart duration
- ready duration
- end-to-end latency

### 3. Repeated edit storm
Simulate rapid changes to the same file and to different files.

Metrics:
- rebuild count
- dropped/merged event behavior
- final healthy restart time

### 4. Idle overhead
Run tool without file changes for fixed duration.

Metrics:
- average CPU usage
- peak CPU usage
- memory footprint

### 5. Multi-service orchestration
Thunder-specific comparison should be explicit here.
Do not compare orchestration features to Air as if they were the same feature set.

Metrics:
- service startup order correctness
- dependency health gating delay
- restart isolation per service

## Benchmark Fixtures

Use at least these projects:

### Fixture A: Small API
- one `cmd/api`
- small dependency graph
- simple `/health` endpoint

### Fixture B: Medium service
- one `cmd/api`
- internal packages
- config loading
- moderate compile workload

### Fixture C: Multi-service stack
- `api`
- `worker`
- optional `redis` dependency

## Execution Rules

- run at least 10 iterations per scenario
- report average, min, max, median if available
- discard obviously invalid runs only with explanation
- record machine specs in report
- record Go version and tool versions in report

## Suggested Output Format

For each scenario report:

- environment
- fixture
- command used
- iteration count
- raw samples
- average/min/max
- notes about anomalies

## Example Report Structure

```md
# Thunder vs Air Benchmark

## Environment
- CPU:
- RAM:
- OS:
- Go version:

## Fixture A: Small API

### Startup latency
| Tool | Avg ms | Min ms | Max ms |
|------|--------|--------|--------|
| Thunder | ... | ... | ... |
| Air | ... | ... | ... |

### Rebuild latency
| Tool | Avg ms | Min ms | Max ms |
|------|--------|--------|--------|
| Thunder | ... | ... | ... |
| Air | ... | ... | ... |
```

## Reproducibility Checklist

- [ ] same fixture committed for all tools
- [ ] same Go version
- [ ] same machine and OS
- [ ] same health endpoint semantics
- [ ] at least 10 iterations per scenario
- [ ] raw sample files stored in repo or attached to release/docs

## Notes on Claims

Do not claim Thunder is faster than Air globally unless the report clearly shows:

- which scenario
- which fixture
- which environment
- how many iterations

Preferred wording:

- "Thunder showed lower average rebuild latency than Air in Fixture A on Windows 11 with Go 1.24.1."

Avoid vague wording like:

- "Thunder is always faster than Air."
