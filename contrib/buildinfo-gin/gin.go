// Package buildinfogin exposes build metadata via a Gin-native handler and
// a Mount helper for typical /version routing with optional middleware.
package buildinfogin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ubgo/buildinfo"
)

// Handler returns a gin.HandlerFunc that responds with buildinfo.Get() as JSON.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, buildinfo.Get())
	}
}

// MountOption configures Mount.
type MountOption func(*mountConfig)

type mountConfig struct {
	path       string
	middleware []gin.HandlerFunc
}

// WithPath overrides the default route ("/version").
func WithPath(p string) MountOption {
	return func(c *mountConfig) { c.path = p }
}

// WithMiddleware applies user middleware to the version route.
func WithMiddleware(mw ...gin.HandlerFunc) MountOption {
	return func(c *mountConfig) { c.middleware = append(c.middleware, mw...) }
}

// Mount registers Handler on r at /version (or the path overridden via
// WithPath), wrapping it with any middleware supplied via WithMiddleware.
func Mount(r gin.IRouter, opts ...MountOption) {
	cfg := &mountConfig{path: "/version"}
	for _, o := range opts {
		o(cfg)
	}
	grp := r.Group("/", cfg.middleware...)
	grp.GET(cfg.path, Handler())
}
