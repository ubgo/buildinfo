package buildinfozap

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestFields_ContainsExpectedKeys(t *testing.T) {
	fields := Fields()

	want := map[string]bool{
		"build_version":   true,
		"build_commit":    true,
		"build_branch":    true,
		"build_time":      true,
		"build_goversion": true,
	}
	for _, f := range fields {
		delete(want, f.Key)
	}
	for k := range want {
		t.Errorf("Fields missing key %q", k)
	}
}

func TestFields_LoggerEmitsBuildKeys(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core).With(Fields()...)
	logger.Info("test")

	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	ctx := entries[0].ContextMap()
	for _, k := range []string{"build_version", "build_commit", "build_branch", "build_time", "build_goversion"} {
		if _, ok := ctx[k]; !ok {
			t.Errorf("log entry missing %q: %v", k, ctx)
		}
	}
}

func TestNamespace_NestsUnderBuildKey(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core).With(Namespace())
	logger.Info("test")

	entries := recorded.All()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	ctx := entries[0].ContextMap()
	build, ok := ctx["build"].(map[string]any)
	if !ok {
		t.Fatalf("expected build namespace as map, got %v (type %T)", ctx["build"], ctx["build"])
	}
	for _, k := range []string{"version", "commit", "branch", "time", "goversion"} {
		if _, ok := build[k]; !ok {
			t.Errorf("namespace missing %q: %v", k, build)
		}
	}
}
