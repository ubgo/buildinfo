// Package buildinfochi exposes build metadata via a Chi-native Mount helper.
//
// Chi accepts any http.Handler natively, so the handler is constructed inline
// rather than depending on the buildinfo-nethttp adapter, keeping each
// adapter standalone.
package buildinfochi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

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

// Middleware is the chi-compatible middleware shape.
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

// WithMiddleware applies user middleware to the version route.
func WithMiddleware(mw ...Middleware) MountOption {
	return func(c *mountConfig) { c.middleware = append(c.middleware, mw...) }
}

// Mount registers Handler on r at /version (or the path overridden via
// WithPath), wrapping it with any middleware supplied via WithMiddleware.
func Mount(r chi.Router, opts ...MountOption) {
	cfg := &mountConfig{path: "/version"}
	for _, o := range opts {
		o(cfg)
	}
	r.Group(func(sub chi.Router) {
		for _, mw := range cfg.middleware {
			sub.Use(mw)
		}
		sub.Method(http.MethodGet, cfg.path, Handler())
	})
}
