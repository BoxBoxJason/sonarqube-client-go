package cli

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// version and buildTime are set at build time via -ldflags when using `make build`.
// When a binary is installed with `go install`, these remain empty and version
// resolution falls back to runtime/debug.ReadBuildInfo instead.
//
//nolint:gochecknoglobals // intentional ldflags injection targets
var (
	version   string // set via -X github.com/boxboxjason/sonarqube-client-go/internal/cli.version=x.y.z
	buildTime string // set via -X github.com/boxboxjason/sonarqube-client-go/internal/cli.buildTime=2006-01-02T15:04:05Z
)

// versionInfo returns a human-readable version string of the form:
//
//	v1.2.3 (go1.25.7, built: 2026-02-22T20:00:00Z)
//
// Version resolution order:
//  1. ldflags-injected version  (make build / make build version=x.y.z)
//  2. Module version from debug.ReadBuildInfo  (go install @vX.Y.Z)
//  3. VCS commit hash from build settings  (local go build / go install @latest)
//  4. "dev" as final fallback
//
// Build time resolution order:
//  1. ldflags-injected buildTime  (make build)
//  2. vcs.time from build settings
//  3. "unknown"
func versionInfo() string {
	ver := resolveVersion()
	bt := resolveBuildTime()

	return fmt.Sprintf("%s (go: %s, built: %s)", ver, runtime.Version(), bt)
}

// resolveVersion returns the most specific version string available.
func resolveVersion() string {
	// ldflags injection wins — used by `make build` and CI releases.
	if version != "" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	// `go install @vX.Y.Z` populates Main.Version with the module tag.
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	// Local `go build` / `go install @latest` — fall back to VCS revision.
	return vcsRevision(info)
}

// vcsRevision extracts the short commit hash (and a "-dirty" suffix when the
// working tree has uncommitted changes) from the VCS build settings.
func vcsRevision(info *debug.BuildInfo) string {
	var revision string

	var dirty bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) > 7 { //nolint:mnd // 7 is the conventional short-hash length
				revision = setting.Value[:7]
			} else {
				revision = setting.Value
			}
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}

	if revision == "" {
		return "dev"
	}

	if dirty {
		return revision + "-dirty"
	}

	return revision
}

// resolveBuildTime returns the build timestamp, falling back through ldflags →
// vcs.time build setting → "unknown".
func resolveBuildTime() string {
	if buildTime != "" {
		return buildTime
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.time" {
			return setting.Value
		}
	}

	return "unknown"
}
