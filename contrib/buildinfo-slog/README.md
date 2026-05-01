# buildinfo-slog

> **Role: Logger renderer.** This adapter reads `buildinfo.Get()` and returns logger fields ready to attach to your structured logger. It does no I/O — see the [system diagram](https://github.com/ubgo/buildinfo#how-the-pieces-fit-together) for how all eight adapters consume the same Info struct.

Stdlib `log/slog` adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata as `[]slog.Attr` (flat) or as a single grouped `slog.Attr` (nested under a `build` group).

Zero third-party dependencies. Standard library only.

## How it works

```
                       ┌──────────────────────────────────────┐
                       │            YOUR SERVICE              │
                       │                                      │
   -ldflags ────→      │  buildinfo.Get() → Info{Version, …}  │
   runtime/debug ──→   │             │                        │
                       │             ▼                        │
                       │  ┌──────────────────┐                │
                       │  │  buildinfo-slog  │                │
                       │  │  (LOG RENDERER)  │                │
                       │  │  Attrs() →       │                │
                       │  │   []slog.Attr    │                │
                       │  │  Group() →       │                │
                       │  │   slog.Attr      │                │
                       │  └────────┬─────────┘                │
                       │           │                          │
                       │           ▼                          │
                       │  slog.Logger.With(Attrs()…)          │
                       │           │                          │
                       │           ▼                          │
                       │  Every log line carries              │
                       │  build_version, build_commit, …      │
                       └───────────┬──────────────────────────┘
                                   ▼
                       Datadog / Loki / ELK / vendor log
                       aggregator
```

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-slog
```

## Quick start (flat attrs)

```go
package main

import (
    "log/slog"
    "os"

    bslog "github.com/ubgo/buildinfo/contrib/buildinfo-slog"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    for _, a := range bslog.Attrs() {
        logger = logger.With(a)
    }
    logger.Info("server starting")
}
```

Emits:

```json
{
  "time": "...",
  "level": "INFO",
  "msg": "server starting",
  "build_version": "dev",
  "build_commit": "abcdef0",
  "build_branch": "unknown",
  "build_time": "2026-04-26T...",
  "build_goversion": "go1.24.0"
}
```

## Quick start (grouped)

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(bslog.Group())
logger.Info("server starting")
```

Emits:

```json
{
  "time": "...",
  "level": "INFO",
  "msg": "server starting",
  "build": {
    "version": "dev",
    "commit": "abcdef0",
    "branch": "unknown",
    "time": "2026-04-26T...",
    "goversion": "go1.24.0"
  }
}
```

## Choosing between Attrs and Group

| Use | When |
|-----|------|
| `Attrs()` | Log-aggregation pipelines that expect flat `build_*` keys (Datadog, Loki, ELK with default mappings). |
| `Group()` | When you want a single `build` sub-object — keeps the rest of the log line uncluttered. |

## API

| Symbol | Purpose |
|--------|---------|
| `Attrs() []slog.Attr` | Flat `build_version`, `build_commit`, `build_branch`, `build_time`, `build_goversion`. |
| `Group() slog.Attr` | Single `build` group containing nested `version`, `commit`, `branch`, `time`, `goversion`. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
