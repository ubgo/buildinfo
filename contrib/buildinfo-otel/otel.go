// Package buildinfootel exposes build metadata as OpenTelemetry resource
// attributes.
package buildinfootel

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/ubgo/buildinfo"
)

// Attribute key constants. Stable across releases.
const (
	KeyVersion   = "build.version"
	KeyCommit    = "build.commit"
	KeyBranch    = "build.branch"
	KeyTime      = "build.time"
	KeyGoVersion = "build.go_version"
	KeyGOOS      = "build.goos"
	KeyGOARCH    = "build.goarch"
	KeyModified  = "build.modified"
)

// Attributes returns OpenTelemetry attribute.KeyValue values describing the
// current build, suitable for resource construction:
//
//	res, _ := resource.New(ctx,
//	    resource.WithAttributes(attribute.String("service.name", "myapi")),
//	    resource.WithAttributes(buildinfootel.Attributes()...),
//	)
//
// Any extra attributes are appended after the build attributes so caller
// values override on duplicates.
func Attributes(extra ...attribute.KeyValue) []attribute.KeyValue {
	info := buildinfo.Get()
	return append([]attribute.KeyValue{
		attribute.String(KeyVersion, info.Version),
		attribute.String(KeyCommit, info.Commit),
		attribute.String(KeyBranch, info.Branch),
		attribute.String(KeyTime, info.BuildTime),
		attribute.String(KeyGoVersion, info.GoVersion),
		attribute.String(KeyGOOS, info.GOOS),
		attribute.String(KeyGOARCH, info.GOARCH),
		attribute.Bool(KeyModified, info.Modified),
	}, extra...)
}
