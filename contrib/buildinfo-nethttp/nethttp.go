// Package buildinfonethttp exposes a stdlib net/http handler that returns
// buildinfo.Get() as JSON.
package buildinfonethttp

import (
	"net/http"

	"github.com/ubgo/buildinfo"
)

// Handler returns an http.Handler that responds with buildinfo.Get() as JSON.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		body, err := buildinfo.Get().JSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}

// Middleware is a stdlib net/http middleware shape.
type Middleware = func(http.Handler) http.Handler

// MountOption configures Mount.
type MountOption func(*mountConfig)

type mountConfig struct {
	path       string
	middleware []Middleware
}

// WithPath overrides the default route ("/version").
func WithPath(p string) MountOption {
	return func(c *mountConfig) { c.path = p }
}

// WithMiddleware applies user middleware to the handler in declaration order.
func WithMiddleware(mw ...Middleware) MountOption {
	return func(c *mountConfig) { c.middleware = append(c.middleware, mw...) }
}

// Mount registers Handler on mux at /version (or the path overridden via
// WithPath), wrapping it with any middleware supplied via WithMiddleware.
func Mount(mux *http.ServeMux, opts ...MountOption) {
	cfg := &mountConfig{path: "/version"}
	for _, o := range opts {
		o(cfg)
	}
	h := Handler()
	for i := len(cfg.middleware) - 1; i >= 0; i-- {
		h = cfg.middleware[i](h)
	}
	mux.Handle(cfg.path, h)
}
