/*
Copyright 2025 codestation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
