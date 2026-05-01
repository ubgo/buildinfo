# buildinfo-echo

> **Role: HTTP renderer.** This adapter reads `buildinfo.Get()` and exposes the result over HTTP. It does no I/O itself — see the [system diagram](https://github.com/ubgo/buildinfo#how-the-pieces-fit-together) for how all eight adapters consume the same Info struct.

Echo adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata via an `echo.HandlerFunc` and a `Mount` helper that registers the route with optional middleware.

## How it works

```
                       ┌──────────────────────────────────────┐
                       │            YOUR SERVICE              │
                       │                                      │
   -ldflags ────→      │  buildinfo.Get() → Info{Version,     │
   runtime/debug ──→   │                          Commit, …}  │
                       │             │                        │
                       │             ▼                        │
                       │  ┌──────────────────┐                │
                       │  │  buildinfo-echo  │ reads Info     │
                       │  │  (HTTP RENDERER) │                │
                       │  └────────┬─────────┘                │
                       │           │ echo.HandlerFunc         │
                       │           ▼                          │
                       │  ┌──────────────────┐                │
                       │  │  echo.Echo       │                │
                       │  │   /version       │                │
                       │  └────────┬─────────┘                │
                       └───────────┼──────────────────────────┘
                                   ▼
                          curl / dashboards / k8s describe
```

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-echo
```

## Quick start

```go
package main

import (
    "github.com/labstack/echo/v4"

    becho "github.com/ubgo/buildinfo/contrib/buildinfo-echo"
)

func main() {
    e := echo.New()
    becho.Mount(e)                            // registers GET /version
    e.Logger.Fatal(e.Start(":8080"))
}
```

```sh
$ curl http://localhost:8080/version
{"version":"dev","commit":"...","go_version":"go1.24.0",...}
```

## Custom path

```go
becho.Mount(e, becho.WithPath("/api/v1/version"))
```

## With middleware (auth, logging, rate-limit, …)

Middleware is `echo.MiddlewareFunc`. Apply via `WithMiddleware`; they wrap in declaration order.

```go
import (
    "crypto/subtle"
    "net/http"
)

func internalKeyAuth(expected string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if subtle.ConstantTimeCompare(
                []byte(c.Request().Header.Get("X-Internal-Key")),
                []byte(expected),
            ) != 1 {
                return c.NoContent(http.StatusUnauthorized)
            }
            return next(c)
        }
    }
}

becho.Mount(e,
    becho.WithPath("/internal/version"),
    becho.WithMiddleware(internalKeyAuth("secret")),
)
```

## API

| Symbol | Purpose |
|--------|---------|
| `Handler() echo.HandlerFunc` | The version handler in isolation. |
| `Mount(e *echo.Echo, opts ...MountOption)` | Registers Handler on the Echo instance with default path `/version`. |
| `WithPath(p string) MountOption` | Override the route. |
| `WithMiddleware(mw ...echo.MiddlewareFunc) MountOption` | Apply user middleware in declaration order. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
