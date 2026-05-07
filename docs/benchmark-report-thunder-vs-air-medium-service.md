# Thunder vs Air Benchmark Report (Medium Service)

## Summary

- Thunder average startup-to-health latency: `6375.64 ms`
- Air average startup-to-health latency: `7450.04 ms`
- Main takeaway: Thunder showed lower average startup latency than Air on the medium service fixture in this Windows test environment.

## Environment

- Machine: local Windows workstation
- OS: Windows
- Go version: Go 1.24.x toolchain in current environment
- Thunder version: local current source build
- Air version: `v1.61.1`
- Date: 2026-05-07

## Fixture

- Fixture: `benchmarks/fixture-medium-service`
- Health endpoint: `http://127.0.0.1:18180/health`
- Scenario: startup-to-health latency
- Iterations: `5`

## Commands

Thunder:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool thunder -Fixture benchmarks/fixture-medium-service -HealthUrl http://127.0.0.1:18180/health -Iterations 5 -OutJson docs/thunder-medium-startup.json -OutMd docs/thunder-medium-startup.md
```

Air:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool air -Fixture benchmarks/fixture-medium-service -HealthUrl http://127.0.0.1:18180/health -Iterations 5 -OutJson docs/air-medium-startup.json -OutMd docs/air-medium-startup.md
```

## Results

### Startup latency

| Tool | Avg ms | Min ms | Max ms | Iterations |
|------|--------|--------|--------|------------|
| Thunder | 6375.64 | 4751.53 | 7603.31 | 5 |
| Air | 7450.04 | 6373.38 | 9245.32 | 5 |

## Raw Reports

- Thunder JSON: `docs/thunder-medium-startup.json`
- Thunder Markdown: `docs/thunder-medium-startup.md`
- Air JSON: `docs/air-medium-startup.json`
- Air Markdown: `docs/air-medium-startup.md`

## Notes

- This run measures startup-to-health latency only.
- The fixture is intentionally heavier than the small API fixture and uses more internal package layering.
- Iteration count is still limited to `5`, but the gap is consistent enough to use as a stronger baseline than the small fixture.
