/*
 *
 * Copyright 2019 codestation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"strconv"
	"time"
)

var (
	// Version indicates the application version
	Version = "dev"
	// Commit indicates the git commit of the build
	Commit = "unknown"
	// BuildTime indicates the date when the binary was built (set by -ldflags)
	BuildTime = "unknown"
)

func init() {
	if BuildTime != "unknown" {
		i, err := strconv.ParseInt(BuildTime, 10, 64)
		if err == nil {
			tm := time.Unix(i, 0)
			BuildTime = tm.Format("Mon Jan _2 15:04:05 2006")
		} else {
			BuildTime = "unknown"
		}
	}
}
