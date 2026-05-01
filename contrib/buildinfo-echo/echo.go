// Package buildinfoecho exposes build metadata via an Echo-native handler
// and a Mount helper for typical /version routing with optional middleware.
package buildinfoecho

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ubgo/buildinfo"
)

// Handler returns an echo.HandlerFunc that responds with buildinfo.Get() as JSON.
func Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, buildinfo.Get())
	}
}

// MountOption configures Mount.
type MountOption func(*mountConfig)

type mountConfig struct {
	path       string
	middleware []echo.MiddlewareFunc
}

// WithPath overrides the default route ("/version").
func WithPath(p string) MountOption {
	return func(c *mountConfig) { c.path = p }
}

// WithMiddleware applies user middleware to the version route.
func WithMiddleware(mw ...echo.MiddlewareFunc) MountOption {
	return func(c *mountConfig) { c.middleware = append(c.middleware, mw...) }
}

// Mount registers Handler on e at /version (or the path overridden via
// WithPath), wrapping it with any middleware supplied via WithMiddleware.
func Mount(e *echo.Echo, opts ...MountOption) {
	cfg := &mountConfig{path: "/version"}
	for _, o := range opts {
		o(cfg)
	}
	grp := e.Group("", cfg.middleware...)
	grp.GET(cfg.path, Handler())
}
