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
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
)

// GiteaConfig has the config options for the GiteaConfig service
type GiteaConfig struct {
	ConfigPath string
	DataPath   string
	SaveDir    string
}

// GiteaAppPath points to the gitea binary location
var GiteaAppPath = "/app/gitea/gitea"

func (g *GiteaConfig) newGiteaCmd() *CmdConfig {
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
		WorkDir:    g.SaveDir,
	}
}

// Backup generates a tarball of the GiteaConfig repositories and returns the path where is stored
func (g *GiteaConfig) Backup() (*BackupResults, error) {
	namePrefix := "gitea-dump"
	filename := generateFilename("", namePrefix) + ".zip"
	args := []string{"dump", "--skip-log", "--tempdir", g.SaveDir, "--file", filename}

	if g.ConfigPath != "" {
		args = append(args, "--config", g.ConfigPath)
	}

	app := g.newGiteaCmd()

	if err := os.MkdirAll(g.SaveDir, 0o755); err != nil {
		return nil, err
	}

	// fix folder permission so gitea can write the backup
	if err := os.Chown(g.SaveDir, int(app.Credential.Uid), int(app.Credential.Gid)); err != nil {
		return nil, err
	}

	if err := app.CmdRun(GiteaAppPath, args...); err != nil {
		return nil, fmt.Errorf("couldn't execute %s, %v", GiteaAppPath, err)
	}

	result := &BackupResults{[]BackupResult{{NamePrefix: namePrefix, Path: path.Join(g.SaveDir, filename)}}}

	return result, nil
}

// Restore takes a GiteaConfig backup and restores it to the service
func (g *GiteaConfig) Restore(_ string) error {
	return errors.New("not implemented")
}
