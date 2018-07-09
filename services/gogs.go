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

package services

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"
)

type Gogs struct {
	ConfigPath string
}

var GogsAppPath = "/app/gogs/gogs"

func (g *Gogs) Backup() (string, error) {
	filepath := fmt.Sprintf("gogs-backup-%s.zip", time.Now().Format("20060102150405"))

	args := []string{
		"git",
		GogsAppPath, "backup",
		"--target", SaveDir,
		"--archive-name", filepath,
	}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	cmd := exec.Command("gosu", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	cmd.Env = append(env, "USER=git")

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", GogsAppPath, err)
	}

	return path.Join(SaveDir, filepath), nil
}

func (g *Gogs) Restore(filepath string) error {
	args := []string{
		"git",
		GogsAppPath, "restore",
		"--from", filepath,
	}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	cmd := exec.Command("gosu", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	env := os.Environ()
	cmd.Env = append(env, "USER=git")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("couldn't execute gogs restore, %v", err)
	}

	return nil
}
