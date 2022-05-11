package version

import (
	"runtime/debug"
	"time"
)

var (
	// Tag indicates the commit tag
	Tag = "none"
	// Revision indicates the git commit of the build
	Revision = "unknown"
	// LastCommit indicates the date of the commit
	LastCommit time.Time
	// Modified indicates if the binary was built from a unmodified source code
	Modified = true
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				Revision = setting.Value
			case "vcs.time":
				LastCommit, _ = time.Parse(time.RFC3339, setting.Value)
			case "vcs.modified":
				Modified = setting.Value == "true"
			}
		}
	}
}
