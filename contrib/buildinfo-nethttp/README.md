# buildinfo-nethttp

> **Role: HTTP renderer.** This adapter reads `buildinfo.Get()` and exposes the result over HTTP. It does no I/O itself — see the [system diagram](https://github.com/ubgo/buildinfo#how-the-pieces-fit-together) for how all eight adapters consume the same Info struct.

Stdlib `net/http` adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata via an `http.Handler` and a `Mount` helper that registers the route with optional middleware.

Zero third-party dependencies. Standard library only.

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
                       │  │ buildinfo-nethttp│ reads Info     │
                       │  │  (HTTP RENDERER) │                │
                       │  └────────┬─────────┘                │
                       │           │ http.Handler             │
                       │           ▼                          │
                       │  ┌──────────────────┐                │
                       │  │ http.ServeMux    │                │
                       │  │   /version       │                │
                       │  └────────┬─────────┘                │
                       └───────────┼──────────────────────────┘
                                   ▼
                          curl / dashboards / k8s describe
```

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-nethttp
```

## Quick start

```go
package main

import (
    "net/http"

    binethttp "github.com/ubgo/buildinfo/contrib/buildinfo-nethttp"
)

func main() {
    mux := http.NewServeMux()
    binethttp.Mount(mux)                    // registers GET /version
    http.ListenAndServe(":8080", mux)
}
```

```sh
$ curl http://localhost:8080/version
{"version":"dev","commit":"...","build_time":"...","branch":"unknown","go_version":"go1.24.0",...}
```

## Custom path

```go
binethttp.Mount(mux, binethttp.WithPath("/api/v1/version"))
```

## With middleware (auth, logging, rate-limit, …)

Middleware uses the standard stdlib shape `func(http.Handler) http.Handler`. Apply via `WithMiddleware` and they wrap in declaration order.

```go
import "crypto/subtle"

func internalKeyAuth(expected string) binethttp.Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if subtle.ConstantTimeCompare(
                []byte(r.Header.Get("X-Internal-Key")),
                []byte(expected),
            ) != 1 {
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

binethttp.Mount(mux,
    binethttp.WithPath("/internal/version"),
    binethttp.WithMiddleware(internalKeyAuth("secret")),
)
```

## API

| Symbol | Purpose |
|--------|---------|
| `Handler() http.Handler` | The version handler in isolation; mount it manually if you need finer control. |
| `Mount(mux *http.ServeMux, opts ...MountOption)` | Registers Handler on the mux with default path `/version`. |
| `WithPath(p string) MountOption` | Override the route. |
| `WithMiddleware(mw ...Middleware) MountOption` | Apply user middleware in declaration order. |
| `Middleware = func(http.Handler) http.Handler` | The stdlib-compatible middleware shape. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
