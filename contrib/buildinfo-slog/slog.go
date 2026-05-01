// Package buildinfoslog exposes build metadata as []slog.Attr for stdlib slog.
package buildinfoslog

import (
	"log/slog"

	"github.com/ubgo/buildinfo"
)

// Attrs returns slog.Attr values describing the current build, suitable for
// logger initialisation:
//
//	logger := slog.Default()
//	for _, a := range buildinfoslog.Attrs() {
//	    logger = logger.With(a)
//	}
func Attrs() []slog.Attr {
	info := buildinfo.Get()
	return []slog.Attr{
		slog.String("build_version", info.Version),
		slog.String("build_commit", info.Commit),
		slog.String("build_branch", info.Branch),
		slog.String("build_time", info.BuildTime),
		slog.String("build_goversion", info.GoVersion),
	}
}

// Group returns the build information as a single slog.Attr group named "build":
//
//	logger := slog.Default().With(buildinfoslog.Group())
//	logger.Info("starting")
//	// → "build":{"version":"...", "commit":"...", ...}
func Group() slog.Attr {
	info := buildinfo.Get()
	return slog.Group("build",
		slog.String("version", info.Version),
		slog.String("commit", info.Commit),
		slog.String("branch", info.Branch),
		slog.String("time", info.BuildTime),
		slog.String("goversion", info.GoVersion),
	)
}
