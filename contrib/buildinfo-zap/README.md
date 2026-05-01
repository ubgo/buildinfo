# buildinfo-zap

Zap adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata as `[]zap.Field` (flat) or as a single grouped `zap.Field` (nested under a `build` namespace).

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-zap
```

## Quick start (flat fields)

```go
package main

import (
    "go.uber.org/zap"

    bzap "github.com/ubgo/buildinfo/contrib/buildinfo-zap"
)

func main() {
    logger, _ := zap.NewProduction()
    logger = logger.With(bzap.Fields()...)

    logger.Info("server starting")
}
```

Emits:

```json
{
  "level": "info",
  "msg": "server starting",
  "build_version": "dev",
  "build_commit": "abcdef0",
  "build_branch": "unknown",
  "build_time": "2026-04-26T...",
  "build_goversion": "go1.24.0"
}
```

## Quick start (grouped namespace)

If you'd rather have build metadata nested under a single `build` key:

```go
logger = logger.With(bzap.Namespace())
logger.Info("server starting")
```

Emits:

```json
{
  "level": "info",
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

## Choosing between Fields and Namespace

| Use | When |
|-----|------|
| `Fields()` | Log-aggregation pipelines that expect flat `build_*` keys (Datadog, Loki, ELK with default mappings). |
| `Namespace()` | When you want a single `build.{...}` sub-object — keeps the rest of the log line uncluttered. |

## API

| Symbol | Purpose |
|--------|---------|
| `Fields() []zap.Field` | Flat `build_version`, `build_commit`, `build_branch`, `build_time`, `build_goversion`. |
| `Namespace() zap.Field` | Single `build` group containing nested `version`, `commit`, `branch`, `time`, `goversion`. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
