package sonar

import "runtime/debug"

// buildUserAgent returns the default User-Agent string.
func buildUserAgent() string {
	return "sonarqube-client-go/" + resolveVersion()
}

// resolveVersion returns the most specific version string available.
func resolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return vcsRevision(info)
}

// vcsRevision extracts the short VCS revision from build info.
func vcsRevision(info *debug.BuildInfo) string {
	var revision string

	var dirty bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) > 7 {//nolint:mnd // 7 is the conventional short-hash length
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