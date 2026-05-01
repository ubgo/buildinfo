// Package buildinfofiber exposes build metadata via a Fiber-native handler
// and a Mount helper for typical /version routing with optional middleware.
package buildinfofiber

import (
	"github.com/gofiber/fiber/v2"

	"github.com/ubgo/buildinfo"
)

// Handler returns a fiber.Handler that responds with buildinfo.Get() as JSON.
func Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(buildinfo.Get())
	}
}

// MountOption configures Mount.
type MountOption func(*mountConfig)

type mountConfig struct {
	path       string
	middleware []fiber.Handler
}

// WithPath overrides the default route ("/version").
func WithPath(p string) MountOption {
	return func(c *mountConfig) { c.path = p }
}

// WithMiddleware applies user middleware to the version route.
func WithMiddleware(mw ...fiber.Handler) MountOption {
	return func(c *mountConfig) { c.middleware = append(c.middleware, mw...) }
}

// Mount registers Handler on r at /version (or the path overridden via
// WithPath), wrapping it with any middleware supplied via WithMiddleware.
//
// Accepts any fiber.Router so mounting on the root *fiber.App or on a
// route group both work.
func Mount(r fiber.Router, opts ...MountOption) {
	cfg := &mountConfig{path: "/version"}
	for _, o := range opts {
		o(cfg)
	}
	chain := append([]fiber.Handler{}, cfg.middleware...)
	chain = append(chain, Handler())
	r.Get(cfg.path, chain...)
}
