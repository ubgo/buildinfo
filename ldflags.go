package buildinfo

// These variables can be overridden via -ldflags at build time:
//
//	-ldflags="-X github.com/ubgo/buildinfo.Version=1.2.3 \
//	          -X github.com/ubgo/buildinfo.Commit=abc123 \
//	          -X github.com/ubgo/buildinfo.BuildTime=2026-04-26T12:00:00Z \
//	          -X github.com/ubgo/buildinfo.Branch=main"
//
// When unset, values fall back to runtime/debug.ReadBuildInfo (Go 1.18+ VCS
// data) where applicable. ldflags overrides always win.
var (
	// Version is the semver string for this build (e.g. "1.2.3").
	Version string

	// Commit is the VCS commit hash.
	Commit string

	// BuildTime is the build timestamp in RFC3339 format.
	BuildTime string

	// Branch is the VCS branch name.
	Branch string
)
