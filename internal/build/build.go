package build

import (
	"fmt"
	"runtime/debug"
)

// Version is dynamically set by -ldflags or falls back to module info.
var Version = "DEV"

// Date is the build date in YYYY-MM-DD format, set by -ldflags.
var Date = ""

// CommitHash is the git commit hash at build time, set by -ldflags or vcs info.
var CommitHash = ""

// EnsureBuildInfo populates Version, CommitHash, and Date from Go build info
// when they were not set via ldflags.
func EnsureBuildInfo() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "DEV" && info.Main.Version != "(devel)" && info.Main.Version != "" {
			Version = info.Main.Version
		}
		if CommitHash == "" {
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					CommitHash = setting.Value
				case "vcs.time":
					if Date == "" {
						Date = setting.Value
					}
				}
			}
		}
	}
	if Version == "" {
		Version = "DEV"
	}
}

// ShortCommitHash returns the first 7 characters of the commit hash.
func ShortCommitHash() string {
	EnsureBuildInfo()
	if len(CommitHash) > 7 {
		return CommitHash[:7]
	}
	return CommitHash
}

// GetInfo returns a formatted version string including the version and commit hash.
func GetInfo() string {
	EnsureBuildInfo()
	res := Version
	if h := ShortCommitHash(); h != "" {
		res += fmt.Sprintf(" (%s)", h)
	}
	if Date != "" {
		res += fmt.Sprintf(" %s", Date)
	}
	return res
}

func init() {
	EnsureBuildInfo()
}
