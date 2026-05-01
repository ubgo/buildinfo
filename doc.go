// Package buildinfo exposes build-time and runtime metadata for a Go binary —
// version, commit, build time, branch, Go version, OS/arch, dirty flag, and
// the list of dependency modules.
//
// Values are populated from two sources, in priority order:
//
//  1. -ldflags overrides set at build time (highest precedence).
//  2. runtime/debug.ReadBuildInfo VCS data (Go 1.18+).
//
// The package has zero third-party dependencies. HTTP, OTEL, Zap, and slog
// integrations live in separate adapter modules under contrib/.
//
// Typical use:
//
//	info := buildinfo.Get()
//	log.Printf("starting %s commit=%s", info.Version, info.Commit)
//
// Build with version stamping:
//
//	go build -ldflags="-X github.com/ubgo/buildinfo.Version=1.2.3"
package buildinfo
