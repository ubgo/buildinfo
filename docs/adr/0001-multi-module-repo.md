# ADR 0001 — Multi-module repository layout

**Status:** Accepted
**Date:** 2026-04-26

## Context

`buildinfo` ships a small, dependency-free core plus a growing set of adapters that bridge to specific HTTP frameworks (Gin, Chi, Echo, Fiber, stdlib net/http) and observability libraries (OpenTelemetry, Zap, slog). Each adapter brings its own third-party dependencies.

## Decision

The repository contains **multiple Go modules**, not one:

- The core module lives at the repository root: `github.com/ubgo/buildinfo`.
- Each adapter lives under `contrib/<adapter-name>/` with its own `go.mod`: e.g. `github.com/ubgo/buildinfo/contrib/buildinfo-gin`.

A consumer picks exactly the modules they need:

```go
import (
    "github.com/ubgo/buildinfo"                                      // core only
    bgin "github.com/ubgo/buildinfo/contrib/buildinfo-gin"           // optional Gin adapter
)
```

`go.work` at the repo root is used for **local development only** and is never published.

## Consequences

- Users only download the dependencies of the adapters they actually import.
- Each module versions independently: `v1.2.0` of `buildinfo` and `v0.5.0` of `buildinfo-gin` can coexist.
- CI runs the test suite per module in a matrix.
- More repositories of metadata files (LICENSE, README) only at the repo root, not per module — adapters reuse the repo root LICENSE.

## Alternatives considered

- **Single module with everything inside.** Forces every consumer to download Gin, Echo, OTEL, etc., even if they only want the core. Rejected.
- **Separate repository per adapter** (e.g. `ubgo/buildinfo`, `ubgo/buildinfo-gin`). Repo explosion for what is fundamentally a coherent product. Rejected.
