# buildinfo-gin

Gin adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata via a `gin.HandlerFunc` and a `Mount` helper that registers the route with optional middleware.

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-gin
```

## Quick start

```go
package main

import (
    "github.com/gin-gonic/gin"

    bgin "github.com/ubgo/buildinfo/contrib/buildinfo-gin"
)

func main() {
    r := gin.Default()
    bgin.Mount(r)                            // registers GET /version
    r.Run(":8080")
}
```

```sh
$ curl http://localhost:8080/version
{"version":"dev","commit":"...","go_version":"go1.24.0",...}
```

## Custom path

```go
bgin.Mount(r, bgin.WithPath("/api/v1/version"))
```

## With middleware (auth, logging, rate-limit, …)

Middleware is `gin.HandlerFunc`. Apply via `WithMiddleware`; they run in declaration order before the version handler.

```go
import (
    "crypto/subtle"
    "net/http"
)

func internalKeyAuth(expected string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if subtle.ConstantTimeCompare(
            []byte(c.GetHeader("X-Internal-Key")),
            []byte(expected),
        ) != 1 {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}

bgin.Mount(r,
    bgin.WithPath("/internal/version"),
    bgin.WithMiddleware(internalKeyAuth("secret")),
)
```

## Mounting on a route group

`Mount` accepts any `gin.IRouter`, so a group works exactly the same:

```go
api := r.Group("/api/v1", authMiddleware())
bgin.Mount(api, bgin.WithPath("/version"))
// → GET /api/v1/version, protected by authMiddleware
```

## API

| Symbol | Purpose |
|--------|---------|
| `Handler() gin.HandlerFunc` | The version handler in isolation. |
| `Mount(r gin.IRouter, opts ...MountOption)` | Registers Handler on the router with default path `/version`. |
| `WithPath(p string) MountOption` | Override the route. |
| `WithMiddleware(mw ...gin.HandlerFunc) MountOption` | Apply user middleware in declaration order. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
