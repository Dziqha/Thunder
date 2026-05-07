# Contributing to Thunder

Thanks for contributing.

## Development Setup

1. Install Go (same major/minor as CI).
2. Fork and clone the repository.
3. Create a feature branch from `master`.

## Local Validation

Run these checks before opening a PR:

```bash
go test ./...
go test -race ./...
go build ./...
thunder release-check
```

## Pull Request Guidelines

- Keep PR scope focused.
- Add/update tests for behavior changes.
- Update docs (`README.md`, `docs/*`) when config or CLI behavior changes.
- Add notable user-facing changes to `CHANGELOG.md` under `Unreleased`.

## Commit Messages

Use concise messages that reflect intent, for example:

- `feat: add profile event filtering`
- `fix: prevent dependency wait deadlock`
- `docs: update orchestration config examples`

## Reporting Bugs

Use the bug report template and include:

- OS + Go version
- Thunder version
- minimal `thunder.toml`
- exact command run
- relevant output logs
