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
	"path"
	"syscall"
)

// GogsConfig has the config options for the GogsConfig service
type GogsConfig struct {
	ConfigPath string
	DataPath   string
	SaveDir    string
}

// GogsAppPath points to the gogs binary location
var GogsAppPath = "/app/gogs/gogs"

func (g *GogsConfig) newGogsCmd() *CmdConfig {
	uid := uint32(getEnvInt("PUID", 1000))
	gid := uint32(getEnvInt("PGID", 1000))
	creds := &syscall.Credential{Uid: uid, Gid: gid}

	env := os.Environ()
	home := fmt.Sprintf("HOME=%s", path.Join(g.DataPath, "git"))
	env = append(env, "USER=git", home)

	return &CmdConfig{
		OutputFile: os.Stdout,
		Env:        env,
		Credential: creds,
	}
}

// Backup generates a tarball of the GogsConfig repositories and returns the path where is stored
func (g *GogsConfig) Backup() (string, error) {
	filename := generateFilename("", "gogs-backup") + ".zip"
	args := []string{"backup", "--target", g.SaveDir, "--archive-name", filename}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	app := g.newGogsCmd()

	if err := app.CmdRun(GogsAppPath, args...); err != nil {
		return "", fmt.Errorf("couldn't execute %s, %v", GogsAppPath, err)
	}

	return path.Join(g.SaveDir, filename), nil
}

// Restore takes a GogsConfig backup and restores it to the service
func (g *GogsConfig) Restore(filepath string) error {
	args := []string{"restore", "--from", filepath, "--tempdir", g.SaveDir}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	app := g.newGogsCmd()

	if err := app.CmdRun(GogsAppPath, args...); err != nil {
		return fmt.Errorf("couldn't execute gogs restore, %v", err)
	}

	return nil
}
