# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial implementation of the `buildinfo` core module: `Info`, `Module`, `Get`, `Map`, `Info.JSON`.
- `-ldflags` overrides for `Version`, `Commit`, `BuildTime`, `Branch`.
- Auto-population of `GoVersion`, `GOOS`, `GOARCH`, `Modified`, `Modules` via `runtime/debug.ReadBuildInfo`.
- Sentinel default values (`"dev"`, `"unknown"`) for empty fields.
- Test suite with 100% statement coverage covering ldflags precedence, VCS fallback, replace-resolution, JSON round-trip, and Map shape.
- Eight adapter modules under `contrib/`:
  - `buildinfo-nethttp` — stdlib `net/http` `/version` handler + Mount helper.
  - `buildinfo-gin` — Gin handler + Mount helper.
  - `buildinfo-chi` — Chi Mount helper using stdlib `http.Handler`.
  - `buildinfo-echo` — Echo handler + Mount helper.
  - `buildinfo-fiber` — Fiber handler + Mount helper.
  - `buildinfo-otel` — OpenTelemetry resource attributes via `Attributes()`.
  - `buildinfo-zap` — Zap log fields via `Fields()` and `Namespace()`.
  - `buildinfo-slog` — stdlib `slog` Attrs via `Attrs()` and `Group()`.
- Each adapter ships its own Go module with its own `go.mod`, so consumers only download dependencies of the adapters they import.
- Every adapter exposes a uniform API: `Handler`, `Mount(router, ...)` with `WithPath` and `WithMiddleware` options for HTTP renderers; `Fields`/`Attrs`/`Attributes` for logging and observability adapters.
- Taskfile, CI workflows, README, CONTRIBUTING, NOTICE.
- Licensed under Apache License 2.0.

[Unreleased]: https://github.com/ubgo/buildinfo/compare/v0.0.0...HEAD
