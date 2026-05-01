# buildinfo-otel

OpenTelemetry adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata as `[]attribute.KeyValue` suitable for use in an OTEL `resource.Resource`.

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-otel
```

## Quick start

```go
package main

import (
    "context"

    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/sdk/resource"
    semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

    botel "github.com/ubgo/buildinfo/contrib/buildinfo-otel"
)

func newResource(ctx context.Context) (*resource.Resource, error) {
    return resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName("myapi"),
        ),
        resource.WithAttributes(botel.Attributes()...),
    )
}
```

The resulting resource carries:

```
service.name      = "myapi"
build.version     = "dev"
build.commit      = "abcdef0..."
build.branch      = "unknown"
build.time        = "2026-04-26T..."
build.go_version  = "go1.24.0"
build.goos        = "darwin"
build.goarch      = "arm64"
build.modified    = false
```

…which become resource attributes on every span and metric your service emits.

## With extra attributes (override / append)

`Attributes(extra ...attribute.KeyValue)` appends `extra` after the build attributes. `resource.New` then deduplicates with later wins, so callers can override individual build attributes:

```go
botel.Attributes(
    attribute.String(botel.KeyVersion, "1.2.3-canary"),  // overrides build.version
    attribute.String("deployment.environment", "production"),
)
```

## Stable attribute key constants

The package exposes the attribute keys as constants so callers can refer to them safely:

```go
botel.KeyVersion    // "build.version"
botel.KeyCommit     // "build.commit"
botel.KeyBranch     // "build.branch"
botel.KeyTime       // "build.time"
botel.KeyGoVersion  // "build.go_version"
botel.KeyGOOS       // "build.goos"
botel.KeyGOARCH     // "build.goarch"
botel.KeyModified   // "build.modified"
```

These keys are part of the public API and follow semver rules — they will not change in a minor release.

## API

| Symbol | Purpose |
|--------|---------|
| `Attributes(extra ...attribute.KeyValue) []attribute.KeyValue` | Returns build metadata as OTEL attributes, with optional caller-supplied extras appended. |
| `Key*` constants | Stable attribute key strings. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
