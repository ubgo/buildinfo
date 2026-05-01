package buildinfo

import (
	"encoding/json"
	"runtime"
	"runtime/debug"
	"sync"
)

// Info contains build metadata for a Go binary.
//
// String fields default to "dev" (Version) or "unknown" (Commit, BuildTime,
// Branch) when neither -ldflags nor runtime/debug VCS data populate them, so
// callers can render them safely without nil checks.
type Info struct {
	Version   string   `json:"version"`
	Commit    string   `json:"commit"`
	BuildTime string   `json:"build_time"`
	Branch    string   `json:"branch"`
	GoVersion string   `json:"go_version"`
	GOOS      string   `json:"goos"`
	GOARCH    string   `json:"goarch"`
	Modified  bool     `json:"modified"`
	Modules   []Module `json:"modules,omitempty"`
}

// Module describes a single dependency module entry from runtime/debug.
type Module struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Sum     string `json:"sum,omitempty"`
}

var (
	cached Info
	once   sync.Once

	// readBuildInfo is the seam for tests to inject a synthetic *debug.BuildInfo.
	readBuildInfo = debug.ReadBuildInfo
)

// Get returns the populated Info struct, cached after the first call.
//
// Population precedence (highest first):
//
//  1. -ldflags overrides set at build time.
//  2. runtime/debug.ReadBuildInfo VCS data (vcs.revision, vcs.time, vcs.modified).
//  3. Sentinel defaults ("dev" for Version, "unknown" for Commit / BuildTime / Branch).
func Get() Info {
	once.Do(load)
	return cached
}

// Map returns Info as a flat string-only map for legacy or string-typed
// consumers (e.g. simple key-value renderers).
func Map() map[string]string {
	i := Get()
	return map[string]string{
		"version":    i.Version,
		"commit":     i.Commit,
		"build_time": i.BuildTime,
		"branch":     i.Branch,
		"go_version": i.GoVersion,
		"goos":       i.GOOS,
		"goarch":     i.GOARCH,
	}
}

// JSON returns the Info marshalled as JSON bytes.
func (i Info) JSON() ([]byte, error) {
	return json.Marshal(i)
}

func load() {
	info := Info{
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	}

	if bi, ok := readBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				info.Commit = s.Value
			case "vcs.time":
				info.BuildTime = s.Value
			case "vcs.modified":
				info.Modified = s.Value == "true"
			}
		}
		info.Modules = collectModules(bi)
	}

	if Version != "" {
		info.Version = Version
	}
	if Commit != "" {
		info.Commit = Commit
	}
	if BuildTime != "" {
		info.BuildTime = BuildTime
	}
	if Branch != "" {
		info.Branch = Branch
	}

	if info.Version == "" {
		info.Version = "dev"
	}
	if info.Commit == "" {
		info.Commit = "unknown"
	}
	if info.BuildTime == "" {
		info.BuildTime = "unknown"
	}
	if info.Branch == "" {
		info.Branch = "unknown"
	}

	cached = info
}

func collectModules(bi *debug.BuildInfo) []Module {
	mods := make([]Module, 0, len(bi.Deps))
	for _, d := range bi.Deps {
		for d.Replace != nil {
			d = d.Replace
		}
		mods = append(mods, Module{
			Path:    d.Path,
			Version: d.Version,
			Sum:     d.Sum,
		})
	}
	return mods
}
