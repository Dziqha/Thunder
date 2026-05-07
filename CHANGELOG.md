# Changelog

All notable changes to this project will be documented in this file.

The format is inspired by Keep a Changelog, and this project follows semantic versioning as closely as practical during active development.

## [Unreleased]

## [v0.2.0] - 2026-05-07

### Added
- Multi-service orchestration mode (`thunder dev`)
- Runtime diagnostics (`thunder doctor`)
- Config inspection (`thunder inspect`)
- Event streaming (`thunder events`) with filtering and optional file output
- Event file rotation options (`--max-size-mb`, `--max-files`)
- Release readiness check command (`thunder release-check`)
- CI and release workflows

### Changed
- Hot reload lifecycle and watcher stability improvements
- Config model expanded for services, profiles, health checks, and restart policies

### Security
- Log redaction for common secret-like keys
