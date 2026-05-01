package buildinfoslog

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestAttrs_ContainsExpectedKeys(t *testing.T) {
	attrs := Attrs()

	want := []string{"build_version", "build_commit", "build_branch", "build_time", "build_goversion"}
	got := make(map[string]bool, len(attrs))
	for _, a := range attrs {
		got[a.Key] = true
	}
	for _, k := range want {
		if !got[k] {
			t.Errorf("Attrs missing key %q", k)
		}
	}
}

func TestAttrs_LoggerEmitsBuildFields(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	for _, a := range Attrs() {
		logger = logger.With(a)
	}
	logger.Info("test")

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal log line: %v body=%s", err, buf.String())
	}
	for _, k := range []string{"build_version", "build_commit", "build_branch", "build_time", "build_goversion"} {
		if _, ok := out[k]; !ok {
			t.Errorf("log output missing key %q: %s", k, buf.String())
		}
	}
}

func TestGroup_NestedUnderBuildKey(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil)).With(Group())
	logger.Info("test")

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal log line: %v body=%s", err, buf.String())
	}
	build, ok := out["build"].(map[string]any)
	if !ok {
		t.Fatalf("expected build group, got %v", out)
	}
	for _, k := range []string{"version", "commit", "branch", "time", "goversion"} {
		if _, ok := build[k]; !ok {
			t.Errorf("group missing key %q: %v", k, build)
		}
	}
}
