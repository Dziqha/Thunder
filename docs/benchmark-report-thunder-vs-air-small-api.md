# Thunder vs Air Benchmark Report

## Summary

Short conclusion:

- Thunder average startup-to-health latency: `6037.69 ms`
- Air average startup-to-health latency: `6366.00 ms`
- Main takeaway: Thunder showed lower average startup latency than Air on this small API fixture in this environment.

## Environment

- Machine: local Windows workstation
- OS: Windows
- Go version: Go 1.24.x toolchain in current environment
- Thunder version: local current source build
- Air version: `v1.61.1`
- Date: 2026-05-07

## Fixture

- Fixture: `benchmarks/fixture-small-api`
- Health endpoint: `http://127.0.0.1:18080/health`
- Scenario: startup-to-health latency
- Iterations: `3`

## Commands

Thunder:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool thunder -Fixture benchmarks/fixture-small-api -Iterations 3 -OutJson docs/thunder-startup.json -OutMd docs/thunder-startup.md
```

Air:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool air -Fixture benchmarks/fixture-small-api -Iterations 3 -OutJson docs/air-startup.json -OutMd docs/air-startup.md
```

## Results

### Startup latency

| Tool | Avg ms | Min ms | Max ms | Iterations |
|------|--------|--------|--------|------------|
| Thunder | 6037.69 | 5780.04 | 6365.95 | 3 |
| Air | 6366.00 | 5754.61 | 6994.10 | 3 |

## Raw Reports

- Thunder JSON: `docs/thunder-startup.json`
- Thunder Markdown: `docs/thunder-startup.md`
- Air JSON: `docs/air-startup.json`
- Air Markdown: `docs/air-startup.md`

## Notes

- This run focuses only on startup-to-health latency.
- This is a small single-service fixture, not a full orchestration comparison.
- Iteration count is still low (`3`), so treat this as an initial baseline rather than a final public performance claim.
- No outliers were removed from the recorded samples.

## Claim Wording

Safe wording for public use:

- "Thunder showed lower average startup-to-health latency than Air on the small API fixture in this Windows test environment."

Avoid broader claims like:

- "Thunder is always faster than Air in all cases."
