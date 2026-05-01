// Package buildinfozap exposes build metadata as []zap.Field for go.uber.org/zap.
package buildinfozap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ubgo/buildinfo"
)

// Fields returns zap.Field values describing the current build, suitable for
// logger initialisation:
//
//	logger := zap.NewProduction()
//	logger = logger.With(buildinfozap.Fields()...)
//	logger.Info("starting")
func Fields() []zap.Field {
	info := buildinfo.Get()
	return []zap.Field{
		zap.String("build_version", info.Version),
		zap.String("build_commit", info.Commit),
		zap.String("build_branch", info.Branch),
		zap.String("build_time", info.BuildTime),
		zap.String("build_goversion", info.GoVersion),
	}
}

// Namespace returns the build information as a single zap.Field namespace
// named "build", which groups all build keys under a nested object.
//
//	logger := zap.NewProduction().With(buildinfozap.Namespace())
//	logger.Info("starting")
//	// → "build":{"version":"...", "commit":"...", ...}
func Namespace() zap.Field {
	return zap.Object("build", buildMarshaler{info: buildinfo.Get()})
}

type buildMarshaler struct {
	info buildinfo.Info
}

func (b buildMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("version", b.info.Version)
	enc.AddString("commit", b.info.Commit)
	enc.AddString("branch", b.info.Branch)
	enc.AddString("time", b.info.BuildTime)
	enc.AddString("goversion", b.info.GoVersion)
	return nil
}
