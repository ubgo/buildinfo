# buildinfo

Build metadata for Go binaries — version, commit, build time, branch, Go version, OS / arch, dirty flag, and dependency module list — populated automatically from `runtime/debug.ReadBuildInfo` (Go 1.18+) with `-ldflags` overrides for CI-stamped builds.

Zero third-party dependencies in the core. HTTP, OTEL, Zap, and slog integrations live in separate adapter modules under `contrib/`.

## How the pieces fit together

`buildinfo` populates a single `Info` struct at process start from two sources, then every contrib **consumes** that struct and surfaces it somewhere — an HTTP route, OTEL resource attributes, Zap fields, slog attrs. There is no registry, no observer pattern, no async work; it's a build-time fact set, exposed many ways.

```
               ┌────────────────────────────────────────────────┐
               │                YOUR SERVICE                    │
               │                                                │
   -ldflags ─→ │  ┌──────────────────────┐                      │
   runtime    │  │   buildinfo.Get()    │  cached on first call │
   /debug ─→  │  │   ↓                  │                       │
               │  │   Info{Version,      │                       │
               │  │     Commit,          │                       │
               │  │     BuildTime,       │                       │
               │  │     Branch,          │                       │
               │  │     GoVersion,       │                       │
               │  │     GOOS, GOARCH,    │                       │
               │  │     Modified,        │                       │
               │  │     Modules[]}       │                       │
               │  └──────────┬───────────┘                       │
               │             │                                   │
               │             ├────────────────────┐              │
               │             │                    │              │
               │             ▼                    ▼              │
               │  ┌──────────────────┐   ┌──────────────────┐    │
               │  │ HTTP ADAPTERS    │   │ LOGGER + OTEL    │    │
               │  │  buildinfo-      │   │  buildinfo-otel  │    │
               │  │   nethttp / gin /│   │  buildinfo-zap   │    │
               │  │   chi / echo /   │   │  buildinfo-slog  │    │
               │  │   fiber          │   │                  │    │
               │  └────────┬─────────┘   └────────┬─────────┘    │
               │           │ /version JSON        │ Attrs/Fields │
               └───────────┼──────────────────────┼──────────────┘
                           ▼                      ▼
                      curl / k8s          attached to every span,
                      release dashboard    metric, and log line
```

Every adapter is read-only against `buildinfo.Info`. None of them perform any I/O on their own; they just make the same struct available in different output formats.

## Install

```sh
go get github.com/ubgo/buildinfo
```

## Quick start

```go
package main

import (
    "log"

    "github.com/ubgo/buildinfo"
)

func main() {
    info := buildinfo.Get()
    log.Printf("starting %s commit=%s go=%s",
        info.Version, info.Commit, info.GoVersion)
}
```

Without any build configuration, you'll see something like:

```
starting dev commit=abcdef0 go=go1.24.0
```

`Commit` was filled by `runtime/debug.ReadBuildInfo` reading the VCS metadata Go embeds since 1.18. `Version` defaults to `"dev"` until you stamp it via `-ldflags`.

## CI version stamping

Stamp release values at build time via `-ldflags`:

```sh
go build -ldflags="\
  -X github.com/ubgo/buildinfo.Version=$(git describe --tags --always) \
  -X github.com/ubgo/buildinfo.Commit=$(git rev-parse HEAD) \
  -X github.com/ubgo/buildinfo.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/ubgo/buildinfo.Branch=$(git rev-parse --abbrev-ref HEAD)"
```

`-ldflags` overrides win over `runtime/debug` data.

## API

```go
// Auto-populated, cached after first call.
info := buildinfo.Get()

// Flat string-only map for simple renderers.
m := buildinfo.Map()

// JSON-marshalled bytes for HTTP / log payloads.
b, _ := info.JSON()
```

`Info` fields:

| Field | Source | Default if empty |
|-------|--------|------------------|
| `Version` | `-ldflags` | `"dev"` |
| `Commit` | `-ldflags` → `vcs.revision` | `"unknown"` |
| `BuildTime` | `-ldflags` → `vcs.time` | `"unknown"` |
| `Branch` | `-ldflags` only | `"unknown"` |
| `GoVersion` | `runtime.Version()` | — |
| `GOOS` | `runtime.GOOS` | — |
| `GOARCH` | `runtime.GOARCH` | — |
| `Modified` | `vcs.modified` | `false` |
| `Modules` | `runtime/debug.BuildInfo.Deps` | empty slice |

## Adapters

Adapter modules ship as separate Go modules under `contrib/`. Import only the ones you use; each pulls in its own dependencies.

| Adapter | Module path | Purpose |
|---------|-------------|---------|
| [`buildinfo-nethttp`](./contrib/buildinfo-nethttp) | `github.com/ubgo/buildinfo/contrib/buildinfo-nethttp` | stdlib `/version` handler + Mount |
| [`buildinfo-gin`](./contrib/buildinfo-gin) | `github.com/ubgo/buildinfo/contrib/buildinfo-gin` | Gin `/version` handler + Mount |
| [`buildinfo-chi`](./contrib/buildinfo-chi) | `github.com/ubgo/buildinfo/contrib/buildinfo-chi` | Chi `/version` Mount helper |
| [`buildinfo-echo`](./contrib/buildinfo-echo) | `github.com/ubgo/buildinfo/contrib/buildinfo-echo` | Echo `/version` handler + Mount |
| [`buildinfo-fiber`](./contrib/buildinfo-fiber) | `github.com/ubgo/buildinfo/contrib/buildinfo-fiber` | Fiber `/version` handler + Mount |
| [`buildinfo-otel`](./contrib/buildinfo-otel) | `github.com/ubgo/buildinfo/contrib/buildinfo-otel` | OpenTelemetry resource attributes |
| [`buildinfo-zap`](./contrib/buildinfo-zap) | `github.com/ubgo/buildinfo/contrib/buildinfo-zap` | Zap log fields (`Fields` + `Namespace`) |
| [`buildinfo-slog`](./contrib/buildinfo-slog) | `github.com/ubgo/buildinfo/contrib/buildinfo-slog` | stdlib `slog` Attrs (`Attrs` + `Group`) |

Click any adapter for its dedicated README with install, quick start, middleware, and API tables.

All eight adapters ship in v0.1.0. Each is a separate Go module under `contrib/<adapter>/` and pulls only its own dependencies.

## Compatibility

Requires Go 1.24 or later.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE) and [`NOTICE`](./NOTICE).
