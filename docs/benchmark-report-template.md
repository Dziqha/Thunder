# Thunder vs Air Benchmark Report

## Summary

Short conclusion:

- Thunder:
- Air:
- Main takeaway:

## Environment

- Machine:
- CPU:
- RAM:
- OS:
- Go version:
- Thunder version:
- Air version:
- Date:

## Fixture

- Fixture: `benchmarks/fixture-small-api`
- Health endpoint: `http://127.0.0.1:18080/health`
- Scenario: startup-to-health latency

## Commands

Thunder:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool thunder -Fixture benchmarks/fixture-small-api -Iterations 10 -OutJson docs/thunder-startup.json -OutMd docs/thunder-startup.md
```

Air:

```bash
pwsh ./scripts/bench-compare.ps1 -Tool air -Fixture benchmarks/fixture-small-api -Iterations 10 -OutJson docs/air-startup.json -OutMd docs/air-startup.md
```

## Results

### Startup latency

| Tool | Avg ms | Min ms | Max ms | Iterations |
|------|--------|--------|--------|------------|
| Thunder |  |  |  |  |
| Air |  |  |  |  |

## Raw Reports

- Thunder JSON: `docs/thunder-startup.json`
- Thunder Markdown: `docs/thunder-startup.md`
- Air JSON: `docs/air-startup.json`
- Air Markdown: `docs/air-startup.md`

## Notes

- Mention whether runs were cold or warm.
- Mention whether any outliers were removed.
- Mention whether other background processes were running.

## Claim Wording

Use precise language, for example:

- "Thunder showed lower average startup-to-health latency than Air on the small API fixture in this environment."

Avoid broad claims like:

- "Thunder is always faster than Air in all cases."
