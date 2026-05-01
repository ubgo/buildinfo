# buildinfo-chi

Chi adapter for [`github.com/ubgo/buildinfo`](https://github.com/ubgo/buildinfo) — exposes build metadata via a stdlib-compatible `http.Handler` and a Chi-native `Mount` helper that registers the route with optional middleware.

## Install

```sh
go get github.com/ubgo/buildinfo
go get github.com/ubgo/buildinfo/contrib/buildinfo-chi
```

## Quick start

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    bchi "github.com/ubgo/buildinfo/contrib/buildinfo-chi"
)

func main() {
    r := chi.NewRouter()
    bchi.Mount(r)                            // registers GET /version
    http.ListenAndServe(":8080", r)
}
```

```sh
$ curl http://localhost:8080/version
{"version":"dev","commit":"...","go_version":"go1.24.0",...}
```

## Custom path

```go
bchi.Mount(r, bchi.WithPath("/api/v1/version"))
```

## With middleware (auth, logging, rate-limit, …)

Middleware uses the standard stdlib shape `func(http.Handler) http.Handler` (the same shape Chi uses natively). Apply via `WithMiddleware`; they run in declaration order before the version handler.

```go
import "crypto/subtle"

func internalKeyAuth(expected string) bchi.Middleware {
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

bchi.Mount(r,
    bchi.WithPath("/internal/version"),
    bchi.WithMiddleware(internalKeyAuth("secret")),
)
```

## Mounting under a sub-router

`Mount` accepts any `chi.Router`, so a sub-router works the same:

```go
r.Route("/api/v1", func(api chi.Router) {
    api.Use(authMiddleware)
    bchi.Mount(api)              // → GET /api/v1/version
})
```

## API

| Symbol | Purpose |
|--------|---------|
| `Handler() http.Handler` | The version handler as a stdlib `http.Handler`. Usable on any router that accepts `http.Handler`. |
| `Mount(r chi.Router, opts ...MountOption)` | Registers Handler on the router with default path `/version`. |
| `WithPath(p string) MountOption` | Override the route. |
| `WithMiddleware(mw ...Middleware) MountOption` | Apply user middleware in declaration order. |
| `Middleware = func(http.Handler) http.Handler` | The Chi-compatible middleware shape. |

## License

Apache-2.0 — see [`LICENSE`](../../LICENSE) at the repository root.
