/*
Copyright 2018 codestation

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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func gogsBackup() (string, error) {
	filename := fmt.Sprintf("gogs-backup-%s.zip", time.Now().Format("20060102150405"))
	cmd := exec.Command("gosu", "git", appPath, "backup", "--target", saveDir, "--archive-name", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "USER=git")
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", appPath, err)
	}

	return filename, nil
}
