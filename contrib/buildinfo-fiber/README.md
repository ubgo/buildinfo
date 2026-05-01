# buildinfo-fiber

> **Role: HTTP renderer.** This adapter reads `buildinfo.Get()` and exposes the result over HTTP. It does no I/O itself — see the [system diagram](https://github.com/ubgo/buildinfo#how-the-pieces-fit-together) for how all eight adapters consume the same Info struct.

Fiber adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata via a `fiber.Handler` and a `Mount` helper that registers the route with optional middleware.

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
                       │  │  buildinfo-fiber │ reads Info     │
                       │  │  (HTTP RENDERER) │                │
                       │  └────────┬─────────┘                │
                       │           │ fiber.Handler            │
                       │           ▼                          │
                       │  ┌──────────────────┐                │
                       │  │  fiber.App       │                │
                       │  │   /version       │                │
                       │  └────────┬─────────┘                │
                       └───────────┼──────────────────────────┘
                                   ▼
                          curl / dashboards / k8s describe
```

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-fiber
```

## Quick start

```go
package main

import (
    "github.com/gofiber/fiber/v2"

    bfiber "github.com/ubgo/buildinfo/contrib/buildinfo-fiber"
)

func main() {
    app := fiber.New()
    bfiber.Mount(app)                            // registers GET /version
    app.Listen(":8080")
}
```

```sh
$ curl http://localhost:8080/version
{"version":"dev","commit":"...","go_version":"go1.24.0",...}
```

## Custom path

```go
bfiber.Mount(app, bfiber.WithPath("/api/v1/version"))
```

## With middleware (auth, logging, rate-limit, …)

Middleware is `fiber.Handler`. Apply via `WithMiddleware`; they run in declaration order before the version handler.

```go
import (
    "crypto/subtle"
    "net/http"
)

func internalKeyAuth(expected string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if subtle.ConstantTimeCompare(
            []byte(c.Get("X-Internal-Key")),
            []byte(expected),
        ) != 1 {
            return c.SendStatus(http.StatusUnauthorized)
        }
        return c.Next()
    }
}

bfiber.Mount(app,
    bfiber.WithPath("/internal/version"),
    bfiber.WithMiddleware(internalKeyAuth("secret")),
)
```

## Mounting on a route group

`Mount` accepts any `fiber.Router`, so a group works exactly the same:

```go
api := app.Group("/api/v1", authMiddleware)
bfiber.Mount(api)            // → GET /api/v1/version
```

## API

| Symbol | Purpose |
|--------|---------|
| `Handler() fiber.Handler` | The version handler in isolation. |
| `Mount(r fiber.Router, opts ...MountOption)` | Registers Handler on the router with default path `/version`. |
| `WithPath(p string) MountOption` | Override the route. |
| `WithMiddleware(mw ...fiber.Handler) MountOption` | Apply user middleware in declaration order. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
