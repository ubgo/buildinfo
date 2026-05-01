package buildinfo

import (
	"encoding/json"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
)

// resetCache clears the package-level cache so each test exercises load() afresh.
// Also resets ldflags variables to their zero values to keep tests independent.
func resetCache(t *testing.T) {
	t.Helper()
	cached = Info{}
	once = sync.Once{}
	Version = ""
	Commit = ""
	BuildTime = ""
	Branch = ""
}

func TestGet_PopulatesRuntimeFields(t *testing.T) {
	resetCache(t)
	info := Get()

	if info.GoVersion != runtime.Version() {
		t.Errorf("GoVersion: got %q, want %q", info.GoVersion, runtime.Version())
	}
	if info.GOOS != runtime.GOOS {
		t.Errorf("GOOS: got %q, want %q", info.GOOS, runtime.GOOS)
	}
	if info.GOARCH != runtime.GOARCH {
		t.Errorf("GOARCH: got %q, want %q", info.GOARCH, runtime.GOARCH)
	}
}

func TestGet_SentinelDefaultsWhenLdflagsEmpty(t *testing.T) {
	resetCache(t)
	info := Get()

	if info.Version != "dev" {
		t.Errorf("Version: got %q, want %q", info.Version, "dev")
	}
	if info.Branch != "unknown" {
		t.Errorf("Branch: got %q, want %q", info.Branch, "unknown")
	}
	// Commit and BuildTime may be populated by runtime/debug VCS data when
	// running `go test` inside a git repo, so we only assert the sentinel for
	// fields runtime/debug never sets (Branch).
}

func TestGet_LdflagsOverrideWin(t *testing.T) {
	resetCache(t)
	Version = "1.2.3"
	Commit = "abc123"
	BuildTime = "2026-04-26T12:00:00Z"
	Branch = "main"

	// Re-trigger the once.Do load with the new ldflags values.
	cached = Info{}
	once = sync.Once{}

	info := Get()

	if info.Version != "1.2.3" {
		t.Errorf("Version: got %q, want %q", info.Version, "1.2.3")
	}
	if info.Commit != "abc123" {
		t.Errorf("Commit: got %q, want %q", info.Commit, "abc123")
	}
	if info.BuildTime != "2026-04-26T12:00:00Z" {
		t.Errorf("BuildTime: got %q, want %q", info.BuildTime, "2026-04-26T12:00:00Z")
	}
	if info.Branch != "main" {
		t.Errorf("Branch: got %q, want %q", info.Branch, "main")
	}
}

func TestGet_Cached(t *testing.T) {
	resetCache(t)
	a := Get()
	b := Get()

	if a.GoVersion != b.GoVersion {
		t.Errorf("Get cache mismatch: a.GoVersion=%q b.GoVersion=%q", a.GoVersion, b.GoVersion)
	}
	if a.Version != b.Version {
		t.Errorf("Get cache mismatch: a.Version=%q b.Version=%q", a.Version, b.Version)
	}
}

func TestMap_ContainsExpectedKeys(t *testing.T) {
	resetCache(t)
	m := Map()

	want := []string{"version", "commit", "build_time", "branch", "go_version", "goos", "goarch"}
	for _, k := range want {
		if _, ok := m[k]; !ok {
			t.Errorf("Map missing key %q (got keys %v)", k, mapKeys(m))
		}
	}
}

func TestInfo_JSON_RoundTrip(t *testing.T) {
	resetCache(t)
	original := Get()

	b, err := original.JSON()
	if err != nil {
		t.Fatalf("JSON: %v", err)
	}
	if !strings.Contains(string(b), `"version":`) {
		t.Errorf("JSON missing version key: %s", b)
	}

	var got Info
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Version != original.Version {
		t.Errorf("RoundTrip Version: got %q, want %q", got.Version, original.Version)
	}
	if got.GoVersion != original.GoVersion {
		t.Errorf("RoundTrip GoVersion: got %q, want %q", got.GoVersion, original.GoVersion)
	}
}

func TestGet_ModulesFieldShape(t *testing.T) {
	resetCache(t)
	info := Get()

	// runtime/debug.ReadBuildInfo behaviour varies across build modes. We only
	// assert that when it does return modules, every entry has a Path set.
	for i, m := range info.Modules {
		if m.Path == "" {
			t.Errorf("Modules[%d].Path is empty: %+v", i, m)
		}
	}
}

// withFakeBuildInfo swaps the readBuildInfo seam to a synthetic value for the
// duration of the test, restoring the original on cleanup.
func withFakeBuildInfo(t *testing.T, bi *debug.BuildInfo, ok bool) {
	t.Helper()
	prev := readBuildInfo
	readBuildInfo = func() (*debug.BuildInfo, bool) { return bi, ok }
	t.Cleanup(func() { readBuildInfo = prev })
}

func TestLoad_VCSSettingsPopulateFields(t *testing.T) {
	resetCache(t)
	withFakeBuildInfo(t, &debug.BuildInfo{
		GoVersion: "go1.24.0",
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "deadbeef"},
			{Key: "vcs.time", Value: "2026-01-15T10:00:00Z"},
			{Key: "vcs.modified", Value: "true"},
		},
	}, true)

	info := Get()

	if info.Commit != "deadbeef" {
		t.Errorf("Commit: got %q, want %q", info.Commit, "deadbeef")
	}
	if info.BuildTime != "2026-01-15T10:00:00Z" {
		t.Errorf("BuildTime: got %q, want %q", info.BuildTime, "2026-01-15T10:00:00Z")
	}
	if !info.Modified {
		t.Errorf("Modified: got false, want true")
	}
}

func TestLoad_VCSModifiedFalseWhenSettingFalse(t *testing.T) {
	resetCache(t)
	withFakeBuildInfo(t, &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs.modified", Value: "false"},
		},
	}, true)

	info := Get()
	if info.Modified {
		t.Errorf("Modified: got true, want false")
	}
}

func TestLoad_LdflagsOverrideRuntimeDebug(t *testing.T) {
	resetCache(t)
	Version = "1.0.0"
	Commit = "fromldflags"
	BuildTime = "2026-04-26T00:00:00Z"
	Branch = "release"

	withFakeBuildInfo(t, &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "fromvcs"},
			{Key: "vcs.time", Value: "2026-01-15T10:00:00Z"},
		},
	}, true)

	info := Get()

	if info.Version != "1.0.0" {
		t.Errorf("Version: got %q, want %q (ldflags should win)", info.Version, "1.0.0")
	}
	if info.Commit != "fromldflags" {
		t.Errorf("Commit: got %q, want %q (ldflags should win)", info.Commit, "fromldflags")
	}
	if info.BuildTime != "2026-04-26T00:00:00Z" {
		t.Errorf("BuildTime: got %q, want %q (ldflags should win)", info.BuildTime, "2026-04-26T00:00:00Z")
	}
}

func TestLoad_NoBuildInfoFallsBackToSentinels(t *testing.T) {
	resetCache(t)
	withFakeBuildInfo(t, nil, false)

	info := Get()

	if info.Version != "dev" {
		t.Errorf("Version: got %q, want %q", info.Version, "dev")
	}
	if info.Commit != "unknown" {
		t.Errorf("Commit: got %q, want %q", info.Commit, "unknown")
	}
	if info.BuildTime != "unknown" {
		t.Errorf("BuildTime: got %q, want %q", info.BuildTime, "unknown")
	}
	if info.Branch != "unknown" {
		t.Errorf("Branch: got %q, want %q", info.Branch, "unknown")
	}
}

func TestCollectModules_ReplaceResolvesToTarget(t *testing.T) {
	bi := &debug.BuildInfo{
		Deps: []*debug.Module{
			{
				Path:    "example.com/origin",
				Version: "v1.0.0",
				Replace: &debug.Module{
					Path:    "example.com/replacement",
					Version: "v1.2.3",
					Sum:     "h1:abcdef",
				},
			},
			{
				Path:    "example.com/plain",
				Version: "v0.5.0",
				Sum:     "h1:plain",
			},
		},
	}

	mods := collectModules(bi)
	if len(mods) != 2 {
		t.Fatalf("len(mods): got %d, want 2", len(mods))
	}

	// Replaced entry should report the replacement path/version/sum.
	if mods[0].Path != "example.com/replacement" {
		t.Errorf("mods[0].Path: got %q, want %q", mods[0].Path, "example.com/replacement")
	}
	if mods[0].Version != "v1.2.3" {
		t.Errorf("mods[0].Version: got %q, want %q", mods[0].Version, "v1.2.3")
	}
	if mods[0].Sum != "h1:abcdef" {
		t.Errorf("mods[0].Sum: got %q, want %q", mods[0].Sum, "h1:abcdef")
	}

	// Non-replaced entry passes through.
	if mods[1].Path != "example.com/plain" {
		t.Errorf("mods[1].Path: got %q, want %q", mods[1].Path, "example.com/plain")
	}
}

func TestCollectModules_EmptyDepsReturnsEmptySlice(t *testing.T) {
	mods := collectModules(&debug.BuildInfo{})
	if mods == nil {
		t.Errorf("collectModules nil; want non-nil empty slice")
	}
	if len(mods) != 0 {
		t.Errorf("len(mods): got %d, want 0", len(mods))
	}
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
